package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigFileName = "nevctl.yaml"
	DefaultServiceName    = "nev-home"
	DefaultTofuBinary     = "tofu"
)

// TemplateMapping maps a template source file to an active destination file.
type TemplateMapping struct {
	Source    string `yaml:"source"`
	Dest      string `yaml:"dest"`
	Ephemeral bool   `yaml:"ephemeral"`
	Mode      string `yaml:"mode"`
}

// Config represents the nevctl configuration.
type Config struct {
	Service    string            `yaml:"service"`
	TofuBinary string            `yaml:"tofu_binary"`
	Templates  []TemplateMapping `yaml:"templates"`

	// RootDir is the directory containing the loaded nevctl.yaml file.
	RootDir string `yaml:"-"`
}

// FindAndLoadConfig locates nevctl.yaml (starting at startDir or explicit path) and loads it.
func FindAndLoadConfig(explicitPath, startDir string) (*Config, error) {
	configPath := explicitPath
	var rootDir string

	if configPath != "" {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, fmt.Errorf("resolving config path: %w", err)
		}
		configPath = absPath
		rootDir = filepath.Dir(configPath)
	} else {
		found, dir, err := findConfigFile(startDir)
		if err != nil {
			return nil, err
		}
		configPath = found
		rootDir = dir
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", configPath, err)
	}

	cfg := &Config{
		Service:    DefaultServiceName,
		TofuBinary: DefaultTofuBinary,
		RootDir:    rootDir,
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", configPath, err)
	}

	if envBin := os.Getenv("TOFU_BIN"); envBin != "" {
		cfg.TofuBinary = envBin
	}

	if cfg.Service == "" {
		cfg.Service = DefaultServiceName
	}
	if cfg.TofuBinary == "" {
		cfg.TofuBinary = DefaultTofuBinary
	}

	// Make template paths absolute relative to RootDir
	for i := range cfg.Templates {
		t := &cfg.Templates[i]
		if !filepath.IsAbs(t.Source) {
			t.Source = filepath.Join(rootDir, t.Source)
		}
		if !filepath.IsAbs(t.Dest) {
			t.Dest = filepath.Join(rootDir, t.Dest)
		}
	}

	return cfg, nil
}

func findConfigFile(startDir string) (string, string, error) {
	if startDir == "" {
		var err error
		startDir, err = os.Getwd()
		if err != nil {
			return "", "", fmt.Errorf("getting working directory: %w", err)
		}
	}

	curr := startDir
	for {
		candidate := filepath.Join(curr, DefaultConfigFileName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, curr, nil
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	return "", "", fmt.Errorf("could not find %s in %s or any parent directory", DefaultConfigFileName, startDir)
}
