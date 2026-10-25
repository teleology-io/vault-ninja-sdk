// Package config reads an optional project file, .vaultninja.yaml, found in
// the current directory or any parent:
//
//	env:
//	  DATABASE_URL: { secret: "Prod DB", field: url }
//	  DB_PASSWORD:  { secret: "Prod DB", field: "Login#password" }
//	mcp:
//	  allow_reveal: ["Prod DB"]          # secrets get_secret_value may return (needs --allow-reveal)
//	  allowed_commands: [npm, go, make]  # if set, run_with_secrets may only run these
//
// The file holds references, never secret values, so it's safe to commit.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const FileName = ".vaultninja.yaml"

type Ref struct {
	Secret string `yaml:"secret" json:"secret" jsonschema:"secret title or id"`
	Field  string `yaml:"field" json:"field,omitempty" jsonschema:"field label, id or type; append #login/#password/#number/#expiry/#pin for part of a login or card field; optional if the secret has one field"`
}

type Config struct {
	Env map[string]Ref `yaml:"env"`
	MCP struct {
		AllowReveal     []string `yaml:"allow_reveal"`
		AllowedCommands []string `yaml:"allowed_commands"`
	} `yaml:"mcp"`

	Path string `yaml:"-"`
}

var envNameRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidEnvName reports whether name is a usable environment variable name.
func ValidEnvName(name string) bool { return envNameRe.MatchString(name) }

// Load reads path, or searches upward from the working directory when path
// is empty. A missing file isn't an error (empty config).
func Load(path string) (*Config, error) {
	if path == "" {
		path = find()
		if path == "" {
			return &Config{Env: map[string]Ref{}}, nil
		}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if c.Env == nil {
		c.Env = map[string]Ref{}
	}
	for name, ref := range c.Env {
		if !ValidEnvName(name) {
			return nil, fmt.Errorf("%s: %q isn't a valid environment variable name", path, name)
		}
		if ref.Secret == "" {
			return nil, fmt.Errorf("%s: env %s has no secret", path, name)
		}
	}
	c.Path = path
	return &c, nil
}

func find() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		p := filepath.Join(dir, FileName)
		if _, err := os.Stat(p); err == nil {
			return p
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// ParseEnvFlag parses NAME=<secret>[:<field>] (split on the last ":", so
// titles may contain colons).
func ParseEnvFlag(s string) (string, Ref, error) {
	name, ref, ok := strings.Cut(s, "=")
	if !ok || !ValidEnvName(name) || ref == "" {
		return "", Ref{}, fmt.Errorf("--env %q: want NAME=<secret>[:<field>]", s)
	}
	if i := strings.LastIndex(ref, ":"); i > 0 {
		return name, Ref{Secret: ref[:i], Field: ref[i+1:]}, nil
	}
	return name, Ref{Secret: ref}, nil
}

// Names returns the env names in a stable order.
func Names(m map[string]Ref) []string {
	names := make([]string, 0, len(m))
	for n := range m {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

var ErrNoEnv = errors.New("no variables to inject — add an env: section to " + FileName + " or pass --env NAME=<secret>[:<field>]")
