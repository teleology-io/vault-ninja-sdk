package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/teleology-io/vault-ninja-sdk/go/cmd/vn/internal/api"
	"github.com/teleology-io/vault-ninja-sdk/go/cmd/vn/internal/config"
	"github.com/teleology-io/vault-ninja-sdk/go/cmd/vn/internal/creds"
)

// globals are accepted by every subcommand.
type globals struct {
	tokenFile  string
	configPath string
}

func newFlags(name string, g *globals) *flag.FlagSet {
	fs := flag.NewFlagSet("vn "+name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&g.tokenFile, "token-file", "", "file containing the API key")
	fs.StringVar(&g.configPath, "config", "", "path to "+config.FileName)
	return fs
}

// parse accepts flags before, between or after positional args (so
// `vn get "Prod DB" url --json` works); "--" ends flag parsing. Positionals
// end up in fs.Args().
func parse(fs *flag.FlagSet, args []string) error {
	var positional []string
	for {
		if err := fs.Parse(args); err != nil {
			return fmt.Errorf("%s: %w", fs.Name(), err)
		}
		rest := fs.Args()
		if len(rest) == 0 {
			break
		}
		// fs.Parse consumed a "--" if rest is shorter than what preceded it.
		if len(args) > 0 && len(rest) < len(args) && args[len(args)-len(rest)-1] == "--" {
			positional = append(positional, rest...)
			break
		}
		positional = append(positional, rest[0])
		args = rest[1:]
	}
	return fs.Parse(append([]string{"--"}, positional...))
}

// parseUntilCommand stops at the first positional arg: everything from
// there on belongs to the child command (vn run).
func parseUntilCommand(fs *flag.FlagSet, args []string) error {
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("%s: %w", fs.Name(), err)
	}
	return nil
}

func client(g *globals, kind string) (*api.Client, error) {
	key, _, err := creds.Resolve(g.tokenFile)
	if err != nil {
		return nil, err
	}
	return api.New(key, kind), nil
}

// envFlags collects repeated --env NAME=<secret>[:<field>].
type envFlags map[string]config.Ref

func (e envFlags) String() string { return "" }
func (e envFlags) Set(s string) error {
	name, ref, err := config.ParseEnvFlag(s)
	if err != nil {
		return err
	}
	e[name] = ref
	return nil
}

// resolveEnv decrypts every reference, fetching each secret once.
func resolveEnv(ctx context.Context, c *api.Client, refs map[string]config.Ref) (map[string]string, error) {
	cache := map[string]*api.Secret{}
	out := make(map[string]string, len(refs))
	for _, name := range config.Names(refs) {
		ref := refs[name]
		s, ok := cache[strings.ToLower(ref.Secret)]
		if !ok {
			var err error
			if s, err = c.Resolve(ctx, ref.Secret); err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			cache[strings.ToLower(ref.Secret)] = s
		}
		v, err := s.Value(ref.Field)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		out[name] = v
	}
	return out, nil
}

// mergedRefs is the project file's env: section overlaid with --env flags.
func mergedRefs(g *globals, flags envFlags) (map[string]config.Ref, *config.Config, error) {
	cfg, err := config.Load(g.configPath)
	if err != nil {
		return nil, nil, err
	}
	refs := map[string]config.Ref{}
	for n, r := range cfg.Env {
		refs[n] = r
	}
	for n, r := range flags {
		refs[n] = r
	}
	return refs, cfg, nil
}

func clientFor(key *creds.Key) *api.Client { return api.New(key, "cli") }
