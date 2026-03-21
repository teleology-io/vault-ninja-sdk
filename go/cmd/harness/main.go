// Vault Ninja SDK — Go integration test harness.
//
// Usage (from the go/ directory):
//
//	go run ./cmd/harness -api-key=vn_org_... -endpoint=http://...
//	VN_API_KEY=vn_org_... VAULT_API_URL=http://... go run ./cmd/harness
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/teleology-io/vault-ninja-sdk/go/vaultninja"
)

var pass, fail int

func ok(label string) {
	fmt.Printf("  ✓ %s\n", label)
	pass++
}

func bad(label string) {
	fmt.Printf("  ✗ %s\n", label)
	fail++
}

func main() {
	apiKey := flag.String("api-key", "", "API key (or set VN_API_KEY)")
	endpoint := flag.String("endpoint", "", "API base URL (or set VAULT_API_URL)")
	flag.Parse()

	if *apiKey == "" {
		*apiKey = os.Getenv("VN_API_KEY")
	}
	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: -api-key or VN_API_KEY required")
		os.Exit(1)
	}

	fmt.Println("[go SDK]")

	var opts []vaultninja.Option
	if *endpoint != "" {
		opts = append(opts, vaultninja.WithBaseURL(*endpoint))
	}
	// VAULT_API_URL is picked up automatically by the client when endpoint flag is empty.

	client := vaultninja.New(*apiKey, opts...)
	ctx := context.Background()

	// 1. listSecrets
	secrets, err := client.ListSecrets(ctx)
	if err != nil {
		bad(fmt.Sprintf("listSecrets — %v", err))
		printResults()
		os.Exit(1)
	}
	if len(secrets) == 0 {
		bad("listSecrets — returned 0 secrets (need at least 1 to continue)")
		printResults()
		os.Exit(1)
	}
	ok(fmt.Sprintf("listSecrets — found %d secret(s)", len(secrets)))

	secretID := secrets[0].GetId()

	// 2. getSecret
	secret, err := client.GetSecret(ctx, secretID)
	if err != nil {
		bad(fmt.Sprintf("getSecret(%s) — %v", secretID, err))
		printResults()
		os.Exit(1)
	}
	ok(fmt.Sprintf("getSecret(%s) — %q", secretID, secret.GetTitle()))

	// 3. getField
	fields := secret.GetFields()
	if len(fields) == 0 {
		fmt.Println("  - getField — skipped (no fields on secret)")
	} else {
		fieldID := fields[0].GetId()
		field, err := client.GetField(ctx, secretID, fieldID)
		if err != nil {
			bad(fmt.Sprintf("getField(%s, %s) — %v", secretID, fieldID, err))
		} else {
			ok(fmt.Sprintf("getField(%s, %s) — label: %q", secretID, fieldID, field.GetLabel()))
		}
	}

	// 4. getFile
	files := secret.GetFiles()
	if len(files) == 0 {
		fmt.Println("  - getFile — skipped (no files on secret)")
	} else {
		fileID := files[0].GetId()
		sizeBytes := files[0].GetSizeBytes()
		_, err := client.GetFile(ctx, secretID, fileID)
		if err != nil {
			bad(fmt.Sprintf("getFile(%s, %s) — %v", secretID, fileID, err))
		} else {
			ok(fmt.Sprintf("getFile(%s, %s) — %d bytes", secretID, fileID, sizeBytes))
		}
	}

	printResults()
	if fail > 0 {
		os.Exit(1)
	}
}

func printResults() {
	fmt.Printf("\nResults: %d passed, %d failed\n", pass, fail)
}
