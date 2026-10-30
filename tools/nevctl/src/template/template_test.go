package template

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"nevctl/src/config"
	"nevctl/src/keyring"
)

type mockKeyring struct {
	secrets map[string]map[string]string
}

func newMockKeyring() *mockKeyring {
	return &mockKeyring{secrets: make(map[string]map[string]string)}
}

func (m *mockKeyring) Get(service, key string) (string, error) {
	if s, ok := m.secrets[service]; ok {
		if val, exists := s[key]; exists {
			return val, nil
		}
	}
	return "", keyring.ErrNotFound
}

func (m *mockKeyring) Set(service, key, value string) error {
	if _, ok := m.secrets[service]; !ok {
		m.secrets[service] = make(map[string]string)
	}
	m.secrets[service][key] = value
	return nil
}

func (m *mockKeyring) Delete(service, key string) error {
	if s, ok := m.secrets[service]; ok {
		delete(s, key)
	}
	return nil
}

func (m *mockKeyring) List(service string) ([]string, error) {
	var keys []string
	if s, ok := m.secrets[service]; ok {
		for k := range s {
			keys = append(keys, k)
		}
	}
	return keys, nil
}

func TestTemplateRender(t *testing.T) {
	tmpDir := t.TempDir()

	srcFile := filepath.Join(tmpDir, "test.conf.tmpl")
	destFile := filepath.Join(tmpDir, "test.conf")

	tmplContent := `token = {{ secret "my_token" | quote }}
user = {{ secret "db_user" }}
env_val = {{ env "TEST_ENV" }}
`
	if err := os.WriteFile(srcFile, []byte(tmplContent), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("TEST_ENV", "production")

	kr := newMockKeyring()
	_ = kr.Set("nev-home", "my_token", "secret123")
	_ = kr.Set("nev-home", "db_user", "admin")

	cfg := &config.Config{
		Service: "nev-home",
		Templates: []config.TemplateMapping{
			{
				Source:    srcFile,
				Dest:      destFile,
				Ephemeral: true,
			},
		},
	}

	engine := NewEngine(cfg, kr)
	if err := engine.RenderAll(); err != nil {
		t.Fatalf("unexpected error rendering: %v", err)
	}

	out, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("reading destination file: %v", err)
	}

	expected := "token = \"secret123\"\nuser = admin\nenv_val = production\n"
	if string(out) != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, string(out))
	}

	// Verify permissions (default 0644 gives read access to group and other)
	info, err := os.Stat(destFile)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Fatalf("expected 0644 permissions, got %o", info.Mode().Perm())
	}

	// Test Clean ephemeral
	if err := engine.Clean(true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(destFile); !os.IsNotExist(err) {
		t.Fatalf("expected file to be deleted, stat err: %v", err)
	}
}

func TestTemplateMissingSecrets(t *testing.T) {
	tmpDir := t.TempDir()

	srcFile := filepath.Join(tmpDir, "test.conf.tmpl")
	destFile := filepath.Join(tmpDir, "test.conf")

	tmplContent := `secret1 = {{ secret "missing1" }}
secret2 = {{ secret "missing2" }}
`
	if err := os.WriteFile(srcFile, []byte(tmplContent), 0644); err != nil {
		t.Fatal(err)
	}

	kr := newMockKeyring()
	cfg := &config.Config{
		Service: "nev-home",
		Templates: []config.TemplateMapping{
			{
				Source: srcFile,
				Dest:   destFile,
			},
		},
	}

	engine := NewEngine(cfg, kr)
	err := engine.RenderAll()
	if err == nil {
		t.Fatal("expected MissingSecretsError, got nil")
	}

	var missingErr *MissingSecretsError
	if !errors.As(err, &missingErr) {
		t.Fatalf("expected *MissingSecretsError, got %T: %v", err, err)
	}

	if len(missingErr.Keys) != 2 {
		t.Fatalf("expected 2 missing keys, got %d", len(missingErr.Keys))
	}
}
