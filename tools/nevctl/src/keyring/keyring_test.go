package keyring

import (
	"errors"
	"testing"
)

func TestSystemKeyring(t *testing.T) {
	kr := NewSystemKeyring()
	service := "nevctl-test"
	key := "test-key"
	val := "secret123"

	// Cleanup just in case
	_ = kr.Delete(service, key)

	// Set
	if err := kr.Set(service, key, val); err != nil {
		t.Skipf("Keyring not available or permission denied: %v", err)
		return
	}

	// Get
	got, err := kr.Get(service, key)
	if err != nil {
		t.Fatalf("Failed to get secret: %v", err)
	}
	if got != val {
		t.Fatalf("Expected %s, got %s", val, got)
	}

	// List
	keys, err := kr.List(service)
	if err != nil {
		t.Fatalf("Failed to list secrets: %v", err)
	}
	found := false
	for _, k := range keys {
		if k == key {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Expected key %s in list %v", key, keys)
	}

	// Delete
	if err := kr.Delete(service, key); err != nil {
		t.Fatalf("Failed to delete secret: %v", err)
	}

	// Get again should fail
	_, err = kr.Get(service, key)
	if err == nil || !errors.Is(err, ErrNotFound) {
		t.Fatalf("Expected ErrNotFound, got %v", err)
	}
}
