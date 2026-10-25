// Package redact scrubs secret values out of command output before it's
// returned to an AI model. It defends against *accidental* leakage
// (a test printing its config, a stack trace with a DSN); it can't stop a
// command deliberately transforming a value (e.g. reversing it) — the real
// controls are a narrowly scoped, expiring, read-only key.
package redact

import (
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

// minLen: shorter values (e.g. "1", "true") would shred unrelated output.
const minLen = 4

type Redactor struct {
	r *strings.Replacer
}

// New builds a redactor for values (map of display name → value). Each
// value is also matched in common encodings.
func New(values map[string]string) *Redactor {
	type pair struct{ needle, label string }
	var pairs []pair
	seen := map[string]bool{}
	add := func(needle, label string) {
		if len(needle) >= minLen && !seen[needle] {
			seen[needle] = true
			pairs = append(pairs, pair{needle, label})
		}
	}
	for name, v := range values {
		label := "[REDACTED:" + name + "]"
		add(v, label)
		add(base64.StdEncoding.EncodeToString([]byte(v)), label)
		add(base64.RawStdEncoding.EncodeToString([]byte(v)), label)
		add(base64.URLEncoding.EncodeToString([]byte(v)), label)
		add(base64.RawURLEncoding.EncodeToString([]byte(v)), label)
		add(hex.EncodeToString([]byte(v)), label)
		add(url.QueryEscape(v), label)
		add(url.PathEscape(v), label)
		for _, line := range strings.Split(v, "\n") { // multi-line values (PEM keys) leak line by line
			add(strings.TrimSpace(line), label)
		}
	}
	// Longest first, so a value containing another is replaced whole.
	sort.Slice(pairs, func(i, j int) bool { return len(pairs[i].needle) > len(pairs[j].needle) })
	args := make([]string, 0, len(pairs)*2)
	for _, p := range pairs {
		args = append(args, p.needle, p.label)
	}
	return &Redactor{r: strings.NewReplacer(args...)}
}

func (r *Redactor) String(s string) string { return r.r.Replace(s) }
