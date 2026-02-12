package envcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	return p
}

func TestCheck_MissingAndExtra(t *testing.T) {
	dir := t.TempDir()

	example := writeFile(t, dir, ".env.example", `
# required
DATABASE_URL=
JWT_SECRET=
export REDIS_URL=

; comment style
`)

	env := writeFile(t, dir, ".env", `
DATABASE_URL=postgres://...
SOME_UNUSED=1
# JWT_SECRET missing
REDIS_URL=redis://...
`)

	res, err := Check(example, env)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if len(res.Missing) != 1 || res.Missing[0] != "JWT_SECRET" {
		t.Fatalf("missing mismatch: %#v", res.Missing)
	}

	// SOME_UNUSED should be extra
	found := false
	for _, k := range res.Extra {
		if k == "SOME_UNUSED" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected SOME_UNUSED in extra, got: %#v", res.Extra)
	}
}

func TestCheck_OK(t *testing.T) {
	dir := t.TempDir()

	example := writeFile(t, dir, ".env.example", `
DATABASE_URL=
JWT_SECRET=
`)

	env := writeFile(t, dir, ".env", `
DATABASE_URL=postgres://...
JWT_SECRET=abc
`)

	res, err := Check(example, env)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if len(res.Missing) != 0 {
		t.Fatalf("expected no missing, got: %#v", res.Missing)
	}
	if len(res.Extra) != 0 {
		t.Fatalf("expected no extra, got: %#v", res.Extra)
	}
}
