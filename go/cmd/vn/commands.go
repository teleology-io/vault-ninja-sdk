package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/teleology-io/vault-ninja-sdk/go/cmd/vn/internal/config"
	"github.com/teleology-io/vault-ninja-sdk/go/cmd/vn/internal/creds"
)

func cmdLogin(ctx context.Context, args []string) error {
	var g globals
	if err := parse(newFlags("login", &g), args); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Paste your API key (vn_org_2_…) and press Enter:")
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	key, err := creds.Parse(line)
	if err != nil {
		return err
	}
	// Check it works before storing it.
	c := clientFor(key)
	secrets, err := c.List(ctx)
	if err != nil {
		return fmt.Errorf("key rejected: %w", err)
	}
	if err := creds.Save(line); err != nil {
		return fmt.Errorf("store in keychain: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Saved. This key can read %d secret(s).\n", len(secrets))
	return nil
}

func cmdLogout(args []string) error {
	var g globals
	if err := parse(newFlags("logout", &g), args); err != nil {
		return err
	}
	if err := creds.Delete(); err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Removed the stored key.")
	return nil
}

func cmdWhoami(ctx context.Context, args []string) error {
	var g globals
	if err := parse(newFlags("whoami", &g), args); err != nil {
		return err
	}
	key, source, err := creds.Resolve(g.tokenFile)
	if err != nil {
		return err
	}
	secrets, err := clientFor(key).List(ctx)
	if err != nil {
		return err
	}
	kid := strings.SplitN(strings.TrimPrefix(key.Auth, "vn_org_2_"), "_", 2)[0]
	fmt.Printf("key:     vn_org_2_%s…\nsource:  %s\nsecrets: %d readable\n", kid[:8], source, len(secrets))
	return nil
}

func cmdList(ctx context.Context, args []string) error {
	var g globals
	fs := newFlags("ls", &g)
	asJSON := fs.Bool("json", false, "JSON output")
	if err := parse(fs, args); err != nil {
		return err
	}
	c, err := client(&g, "cli")
	if err != nil {
		return err
	}
	secrets, err := c.List(ctx)
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(secrets)
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ID\tTITLE\tTAGS\tUPDATED")
	for _, s := range secrets {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.ID, s.Title, strings.Join(s.Tags, ","), s.UpdatedAt.Local().Format("2006-01-02"))
	}
	return tw.Flush()
}

func cmdGet(ctx context.Context, args []string) error {
	var g globals
	fs := newFlags("get", &g)
	asJSON := fs.Bool("json", false, "print the whole secret as JSON")
	if err := parse(fs, args); err != nil {
		return err
	}
	if fs.NArg() < 1 || fs.NArg() > 2 {
		return errors.New("usage: vn get <secret> [field] [--json]")
	}
	c, err := client(&g, "cli")
	if err != nil {
		return err
	}
	s, err := c.Resolve(ctx, fs.Arg(0))
	if err != nil {
		return err
	}
	if *asJSON {
		return writeJSON(s)
	}
	if fs.NArg() == 1 && len(s.Fields) != 1 {
		// Don't dump every value; show what's there.
		tw := tabwriter.NewWriter(os.Stderr, 0, 4, 2, ' ', 0)
		fmt.Fprintf(tw, "%s has %d fields — pick one: vn get %q <field>\n", s.Title, len(s.Fields), fs.Arg(0))
		for _, f := range s.Fields {
			fmt.Fprintf(tw, "  %s\t%s\n", f.Label, f.Type)
		}
		tw.Flush()
		return exitCode(2)
	}
	v, err := s.Value(fs.Arg(1))
	if err != nil {
		return err
	}
	fmt.Print(v)
	if isTerminal(os.Stdout) {
		fmt.Println()
	}
	return nil
}

func cmdEnv(ctx context.Context, args []string) error {
	var g globals
	flags := envFlags{}
	fs := newFlags("env", &g)
	format := fs.String("format", "dotenv", "dotenv, json or shell")
	fs.Var(flags, "env", "NAME=<secret>[:<field>] (repeatable)")
	if err := parse(fs, args); err != nil {
		return err
	}
	refs, _, err := mergedRefs(&g, flags)
	if err != nil {
		return err
	}
	if len(refs) == 0 {
		return config.ErrNoEnv
	}
	c, err := client(&g, "cli")
	if err != nil {
		return err
	}
	vals, err := resolveEnv(ctx, c, refs)
	if err != nil {
		return err
	}
	switch *format {
	case "json":
		return writeJSON(vals)
	case "dotenv":
		for _, n := range config.Names(refs) {
			fmt.Printf("%s=%s\n", n, dotenvQuote(vals[n]))
		}
	case "shell":
		for _, n := range config.Names(refs) {
			fmt.Printf("export %s=%s\n", n, shellQuote(vals[n]))
		}
	default:
		return fmt.Errorf("unknown --format %q (dotenv, json, shell)", *format)
	}
	return nil
}

func cmdRun(ctx context.Context, args []string) error {
	var g globals
	flags := envFlags{}
	fs := newFlags("run", &g)
	fs.Var(flags, "env", "NAME=<secret>[:<field>] (repeatable)")
	if err := parseUntilCommand(fs, args); err != nil {
		return err
	}
	argv := fs.Args()
	if len(argv) == 0 {
		return errors.New("usage: vn run [--env NAME=<secret>[:<field>]]... -- <command> [args...]")
	}
	refs, _, err := mergedRefs(&g, flags)
	if err != nil {
		return err
	}
	if len(refs) == 0 {
		return config.ErrNoEnv
	}
	c, err := client(&g, "cli")
	if err != nil {
		return err
	}
	vals, err := resolveEnv(ctx, c, refs)
	if err != nil {
		return err
	}

	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = withoutCredentials(os.Environ())
	for n, v := range vals {
		cmd.Env = append(cmd.Env, n+"="+v)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	// Ctrl-C reaches the child through the process group; just wait.
	err = cmd.Wait()
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return exitCode(ee.ExitCode())
	}
	return err
}

// withoutCredentials strips vn's own key from a child's environment — the
// command needs the secrets it asked for, not the key that reads them all.
func withoutCredentials(env []string) []string {
	out := env[:0:0]
	for _, kv := range env {
		if strings.HasPrefix(kv, "VN_TOKEN=") || strings.HasPrefix(kv, "VN_API_KEY=") {
			continue
		}
		out = append(out, kv)
	}
	return out
}

func cmdFile(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] != "get" {
		return errors.New("usage: vn file get <secret> <file> [-o path|-]")
	}
	var g globals
	fs := newFlags("file get", &g)
	outPath := fs.String("o", "", "output path (default: the file's name; - for stdout)")
	if err := parse(fs, args[1:]); err != nil {
		return err
	}
	if fs.NArg() != 2 {
		return errors.New("usage: vn file get <secret> <file> [-o path|-]")
	}
	c, err := client(&g, "cli")
	if err != nil {
		return err
	}
	s, err := c.Resolve(ctx, fs.Arg(0))
	if err != nil {
		return err
	}
	f, err := s.FindFile(fs.Arg(1))
	if err != nil {
		return err
	}
	if *outPath == "-" {
		return c.DownloadFile(ctx, s, f.ID, os.Stdout)
	}
	path := *outPath
	if path == "" {
		path = filepath.Base(f.Name)
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".vn-download-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if err := c.DownloadFile(ctx, s, f.ID, tmp); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), path); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Saved %s (%d bytes)\n", path, f.Size)
	return nil
}

func writeJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func dotenvQuote(v string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "$", `\$`)
	return `"` + r.Replace(v) + `"`
}

func shellQuote(v string) string {
	return "'" + strings.ReplaceAll(v, "'", `'\''`) + "'"
}

func isTerminal(f *os.File) bool {
	st, err := f.Stat()
	return err == nil && st.Mode()&os.ModeCharDevice != 0
}
