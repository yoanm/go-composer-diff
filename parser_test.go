package compdiff_test

import (
	"testing"

	compdiff "github.com/yoanm/go-composer-diff"
)

func TestParseLock_SupportTypeMismatch(t *testing.T) {
	t.Parallel()

	// When support is a string instead of object, unmarshal should fail
	data := []byte(`{"packages": [{"name":"vendor/bad","version":"1.0.0","support":"https://wiki.example"}]}`)
	if _, err := compdiff.ParseLock(data); err == nil {
		t.Fatalf("expected ParseLock to fail on support type mismatch, but it succeeded")
	}
}

func TestParseLock_WithoutOptionalFields(t *testing.T) {
	t.Parallel()

	// Packages without optional fields should parse fine
	data := []byte(`{
  "packages": [
    {"name":"vendor/minimal","version":"1.2.3"}
  ]
}`)

	lock, err := compdiff.ParseLock(data)
	if err != nil {
		t.Fatalf("ParseLock failed: %v", err)
	}

	if len(lock.Packages) != 1 {
		t.Fatalf("expected 1 package, got %d", len(lock.Packages))
	}
}
