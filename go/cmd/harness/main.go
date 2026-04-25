// Vault Ninja SDK — Go integration test harness.
//
// Usage (from the go/ directory):
//
//	go run ./cmd/harness -api-key=vn_org_... -endpoint=http://...
//	VN_API_KEY=vn_org_... VN_API_URL=http://... go run ./cmd/harness
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
	endpoint := flag.String("endpoint", "", "API base URL (or set VN_API_URL)")
	flag.Parse()

	if *apiKey == "" {
		*apiKey = os.Getenv("VN_API_KEY")
	}
	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: -api-key or VN_API_KEY required")
		os.Exit(1)
	}

	fmt.Println("[go SDK]")

	client := vaultninja.New(*apiKey, *endpoint)
	ctx := context.Background()

	// 1. ListSecrets
	secrets, err := client.ListSecrets(ctx)
	if err != nil {
		bad(fmt.Sprintf("ListSecrets — %v", err))
		printResults()
		os.Exit(1)
	}
	if len(secrets) == 0 {
		bad("ListSecrets — returned 0 secrets (need at least 1 to continue)")
		printResults()
		os.Exit(1)
	}
	ok(fmt.Sprintf("ListSecrets — found %d secret(s)", len(secrets)))

	secretID := secrets[0].ID

	// 2. GetSecret
	secret, err := client.GetSecret(ctx, secretID)
	if err != nil {
		bad(fmt.Sprintf("GetSecret(%s) — %v", secretID, err))
		printResults()
		os.Exit(1)
	}
	ok(fmt.Sprintf("GetSecret(%s) — %q", secretID, secret.Title))

	// 3. GetField
	if len(secret.Fields) == 0 {
		fmt.Println("  - GetField — skipped (no fields on secret)")
	} else {
		fieldID := secret.Fields[0].ID
		field, err := client.GetField(ctx, secretID, fieldID)
		if err != nil {
			bad(fmt.Sprintf("GetField(%s, %s) — %v", secretID, fieldID, err))
		} else {
			ok(fmt.Sprintf("GetField(%s, %s) — label: %q", secretID, fieldID, field.Label))
		}
	}

	// 4. GetFile
	if len(secret.Files) == 0 {
		fmt.Println("  - GetFile — skipped (no files on secret)")
	} else {
		fileID := secret.Files[0].ID
		sizeBytes := secret.Files[0].SizeBytes
		_, err := client.GetFile(ctx, secretID, fileID)
		if err != nil {
			bad(fmt.Sprintf("GetFile(%s, %s) — %v", secretID, fileID, err))
		} else {
			ok(fmt.Sprintf("GetFile(%s, %s) — %d bytes", secretID, fileID, sizeBytes))
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
