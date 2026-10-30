package template

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"nevctl/src/config"
	"nevctl/src/keyring"
)

// MissingSecretsError captures all secrets that were required by templates but missing from the keyring.
type MissingSecretsError struct {
	Service string
	Keys    []string
}

func (e *MissingSecretsError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("missing %d required secret(s) in keyring for service %q:\n", len(e.Keys), e.Service))
	for _, k := range e.Keys {
		sb.WriteString(fmt.Sprintf("  - %s (to set: nevctl secret set %s)\n", k, k))
	}
	return sb.String()
}

// Engine handles rendering templates using secrets from the keyring.
type Engine struct {
	cfg *config.Config
	kr  keyring.Keyring
}

// NewEngine creates a new template Engine.
func NewEngine(cfg *config.Config, kr keyring.Keyring) *Engine {
	return &Engine{
		cfg: cfg,
		kr:  kr,
	}
}

// RenderAll renders all configured templates to their destination paths.
func (e *Engine) RenderAll() error {
	missingKeys := make(map[string]struct{})
	renderedFiles := make(map[string][]byte)

	for _, tmpl := range e.cfg.Templates {
		content, err := os.ReadFile(tmpl.Source)
		if err != nil {
			return fmt.Errorf("reading template source %s: %w", tmpl.Source, err)
		}

		funcMap := template.FuncMap{
			"secret": func(key string) (string, error) {
				val, err := e.kr.Get(e.cfg.Service, key)
				if err != nil {
					missingKeys[key] = struct{}{}
					return "", nil
				}
				return val, nil
			},
			"quote": func(v any) string {
				return fmt.Sprintf("%q", fmt.Sprint(v))
			},
			"squote": func(v any) string {
				return fmt.Sprintf("'%s'", strings.ReplaceAll(fmt.Sprint(v), "'", `\'`))
			},
			"defaultVal": func(def string, v any) string {
				s := fmt.Sprint(v)
				if s == "" || s == "<nil>" {
					return def
				}
				return s
			},
			"env": func(name string) string {
				return os.Getenv(name)
			},
			"b64enc": func(s string) string {
				return base64.StdEncoding.EncodeToString([]byte(s))
			},
			"b64dec": func(s string) (string, error) {
				b, err := base64.StdEncoding.DecodeString(s)
				if err != nil {
					return "", err
				}
				return string(b), nil
			},
		}

		t, err := template.New(filepath.Base(tmpl.Source)).Funcs(funcMap).Parse(string(content))
		if err != nil {
			return fmt.Errorf("parsing template %s: %w", tmpl.Source, err)
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, nil); err != nil {
			return fmt.Errorf("executing template %s: %w", tmpl.Source, err)
		}

		renderedFiles[tmpl.Dest] = buf.Bytes()
	}

	if len(missingKeys) > 0 {
		var keys []string
		for k := range missingKeys {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return &MissingSecretsError{
			Service: e.cfg.Service,
			Keys:    keys,
		}
	}

	// Write all files with 0644 permissions (user rw, group r, other r) by default, or custom mode
	for _, tmpl := range e.cfg.Templates {
		data := renderedFiles[tmpl.Dest]
		mode := os.FileMode(0644)
		if tmpl.Mode != "" {
			if parsed, err := strconv.ParseUint(tmpl.Mode, 8, 32); err == nil {
				mode = os.FileMode(parsed)
			}
		}

		if err := os.MkdirAll(filepath.Dir(tmpl.Dest), 0755); err != nil {
			return fmt.Errorf("creating directory for %s: %w", tmpl.Dest, err)
		}
		if err := os.WriteFile(tmpl.Dest, data, mode); err != nil {
			return fmt.Errorf("writing rendered file %s: %w", tmpl.Dest, err)
		}
		if err := os.Chmod(tmpl.Dest, mode); err != nil {
			return fmt.Errorf("setting permissions on %s: %w", tmpl.Dest, err)
		}
	}

	return nil
}

// Clean removes rendered files. If ephemeralOnly is true, only files marked ephemeral are deleted.
func (e *Engine) Clean(ephemeralOnly bool) error {
	var errs []string
	for _, tmpl := range e.cfg.Templates {
		if ephemeralOnly && !tmpl.Ephemeral {
			continue
		}
		if err := os.Remove(tmpl.Dest); err != nil && !os.IsNotExist(err) {
			errs = append(errs, fmt.Sprintf("removing %s: %v", tmpl.Dest, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors during clean: %s", strings.Join(errs, "; "))
	}
	return nil
}
