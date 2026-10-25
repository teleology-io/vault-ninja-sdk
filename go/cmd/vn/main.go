// Command vn is the zero-knowledge Vault Ninja CLI and MCP server. It reads
// org secrets through /api/sdk/v2 and decrypts them on this machine: the
// API key's private half never leaves it, and the server only ever sees
// ciphertext.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/teleology-io/vault-ninja-sdk/go/cmd/vn/internal/api"
)

const usage = `vn — Vault Ninja secrets, decrypted locally

Usage:
  vn login                         store an API key in the OS keychain (reads it from stdin)
  vn logout                        remove the stored key
  vn whoami                        show which key is in use and how many secrets it can read
  vn ls [--json]                   list secrets shared with the key
  vn get <secret> [field] [--json] print a field value (or the whole secret with --json)
  vn env [--format dotenv|json|shell] [--env NAME=<secret>[:<field>]]...
                                   print variables from .vaultninja.yaml / --env (stdout only)
  vn run [--env NAME=<secret>[:<field>]]... -- <command> [args...]
                                   run a command with secrets injected as env vars
  vn file get <secret> <file> [-o path|-]
                                   download and decrypt a file
  vn mcp [--allow-reveal]          run the MCP server on stdio (for AI tools)
  vn version

Global flags (any command): --token-file <path>, --config <path>

Credentials, first match wins: $VN_TOKEN, $VN_API_KEY, --token-file, keychain (vn login).
The key must be a vn_org_2_… key (Settings → API Keys). $VN_API_URL overrides the server.
`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	var err error
	switch cmd {
	case "login":
		err = cmdLogin(ctx, args)
	case "logout":
		err = cmdLogout(args)
	case "whoami":
		err = cmdWhoami(ctx, args)
	case "ls", "list":
		err = cmdList(ctx, args)
	case "get":
		err = cmdGet(ctx, args)
	case "env":
		err = cmdEnv(ctx, args)
	case "run":
		err = cmdRun(ctx, args)
	case "file":
		err = cmdFile(ctx, args)
	case "mcp":
		err = cmdMCP(ctx, args)
	case "version", "--version", "-v":
		fmt.Println("vn", api.Version)
	case "help", "--help", "-h":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "vn: unknown command %q\n\n%s", cmd, usage)
		os.Exit(2)
	}
	if err != nil {
		var exit exitCode
		if errors.As(err, &exit) {
			os.Exit(int(exit))
		}
		fmt.Fprintln(os.Stderr, "vn:", err)
		os.Exit(1)
	}
}

// exitCode lets `vn run` pass the child's exit status through.
type exitCode int

func (e exitCode) Error() string { return fmt.Sprintf("exit status %d", int(e)) }
