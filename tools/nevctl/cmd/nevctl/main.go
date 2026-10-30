package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"nevctl/src/config"
	"nevctl/src/keyring"
	"nevctl/src/runner"
	"nevctl/src/template"
)

var (
	configPathFlag string
	serviceFlag    string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "nevctl",
		Short: "Command and control CLI for OpenTofu with secure keyring secrets",
		Long: `nevctl manages secrets in the Linux system keyring and hydrates configuration
templates before executing OpenTofu commands, guaranteeing sensitive credentials never
reach git.`,
	}

	rootCmd.PersistentFlags().StringVarP(&configPathFlag, "config", "c", "", "Path to nevctl.yaml configuration file")
	rootCmd.PersistentFlags().StringVarP(&serviceFlag, "service", "s", "", "Override keyring service name")

	// Helper to load config and keyring
	loadEnv := func() (*config.Config, keyring.Keyring, error) {
		cfg, err := config.FindAndLoadConfig(configPathFlag, "")
		if err != nil {
			return nil, nil, err
		}
		if serviceFlag != "" {
			cfg.Service = serviceFlag
		}
		kr := keyring.NewSystemKeyring()
		return cfg, kr, nil
	}

	// Wrapper commands for OpenTofu
	tofuSubcommands := []string{"plan", "apply", "init", "destroy"}
	for _, subcmd := range tofuSubcommands {
		cmdName := subcmd
		cmd := &cobra.Command{
			Use:                fmt.Sprintf("%s [flags]", cmdName),
			Short:              fmt.Sprintf("Run 'tofu %s' with hydrated secrets", cmdName),
			DisableFlagParsing: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, kr, err := loadEnv()
				if err != nil {
					return err
				}
				r := runner.NewRunner(cfg, kr)
				tofuArgs := append([]string{cmdName}, args...)
				return r.Execute(context.Background(), tofuArgs)
			},
		}
		rootCmd.AddCommand(cmd)
	}

	// Raw pass-through command: nevctl tofu [args...]
	tofuRawCmd := &cobra.Command{
		Use:                "tofu [args...]",
		Short:              "Execute arbitrary OpenTofu command with hydrated secrets",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, kr, err := loadEnv()
			if err != nil {
				return err
			}
			r := runner.NewRunner(cfg, kr)
			return r.Execute(context.Background(), args)
		},
	}
	rootCmd.AddCommand(tofuRawCmd)

	// Render command
	renderCmd := &cobra.Command{
		Use:   "render",
		Short: "Hydrate all template files from keyring secrets without running tofu",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, kr, err := loadEnv()
			if err != nil {
				return err
			}
			engine := template.NewEngine(cfg, kr)
			if err := engine.RenderAll(); err != nil {
				return err
			}
			fmt.Printf("Successfully rendered %d template(s).\n", len(cfg.Templates))
			for _, t := range cfg.Templates {
				ephemStr := ""
				if t.Ephemeral {
					ephemStr = " (ephemeral)"
				}
				fmt.Printf("  -> %s%s\n", t.Dest, ephemStr)
			}
			return nil
		},
	}
	rootCmd.AddCommand(renderCmd)

	// Clean command
	var ephemeralOnlyFlag bool
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Remove rendered files from disk",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, kr, err := loadEnv()
			if err != nil {
				return err
			}
			engine := template.NewEngine(cfg, kr)
			if err := engine.Clean(ephemeralOnlyFlag); err != nil {
				return err
			}
			if ephemeralOnlyFlag {
				fmt.Println("Cleaned ephemeral rendered files.")
			} else {
				fmt.Println("Cleaned all rendered configuration files.")
			}
			return nil
		},
	}
	cleanCmd.Flags().BoolVarP(&ephemeralOnlyFlag, "ephemeral-only", "e", false, "Only delete files designated as ephemeral")
	rootCmd.AddCommand(cleanCmd)

	// Secret subcommands
	secretCmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage secrets in the system keyring",
	}

	secretSetCmd := &cobra.Command{
		Use:   "set <key> [value]",
		Short: "Store a secret in the system keyring",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, kr, err := loadEnv()
			if err != nil {
				return err
			}
			key := args[0]
			var val string

			if len(args) == 2 {
				val = args[1]
			} else {
				// Interactive secure password read
				if term.IsTerminal(int(syscall.Stdin)) {
					fmt.Printf("Enter secret value for %q: ", key)
					byteVal, err := term.ReadPassword(int(syscall.Stdin))
					fmt.Println()
					if err != nil {
						return fmt.Errorf("reading password: %w", err)
					}
					val = string(byteVal)
				} else {
					reader := bufio.NewReader(os.Stdin)
					input, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("reading stdin: %w", err)
					}
					val = strings.TrimRight(input, "\r\n")
				}
			}

			if err := kr.Set(cfg.Service, key, val); err != nil {
				return err
			}
			fmt.Printf("Secret %q successfully saved in keyring (service: %s).\n", key, cfg.Service)
			return nil
		},
	}

	secretGetCmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Retrieve and display a secret from the keyring",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, kr, err := loadEnv()
			if err != nil {
				return err
			}
			key := args[0]
			val, err := kr.Get(cfg.Service, key)
			if err != nil {
				return err
			}
			fmt.Println(val)
			return nil
		},
	}

	secretDeleteCmd := &cobra.Command{
		Use:   "delete <key>",
		Short: "Delete a secret from the keyring",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, kr, err := loadEnv()
			if err != nil {
				return err
			}
			key := args[0]
			if err := kr.Delete(cfg.Service, key); err != nil {
				return err
			}
			fmt.Printf("Secret %q deleted from keyring (service: %s).\n", key, cfg.Service)
			return nil
		},
	}

	secretListCmd := &cobra.Command{
		Use:   "list",
		Short: "List secrets stored in the keyring for the current service",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, kr, err := loadEnv()
			if err != nil {
				return err
			}
			keys, err := kr.List(cfg.Service)
			if err != nil {
				return fmt.Errorf("listing keyring secrets: %w", err)
			}

			fmt.Printf("Secrets in keyring for service %q:\n", cfg.Service)
			if len(keys) == 0 {
				fmt.Println("  (no secrets found)")
				return nil
			}

			for _, k := range keys {
				fmt.Printf("  • %s\n", k)
			}
			return nil
		},
	}

	secretCmd.AddCommand(secretSetCmd, secretGetCmd, secretDeleteCmd, secretListCmd)
	rootCmd.AddCommand(secretCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
