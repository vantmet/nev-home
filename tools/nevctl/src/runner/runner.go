package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"nevctl/src/config"
	"nevctl/src/keyring"
	"nevctl/src/template"
)

// Runner manages the lifecycle of running OpenTofu with hydrated secrets.
type Runner struct {
	cfg *config.Config
	kr  keyring.Keyring
}

// NewRunner creates a new Runner.
func NewRunner(cfg *config.Config, kr keyring.Keyring) *Runner {
	return &Runner{
		cfg: cfg,
		kr:  kr,
	}
}

// Execute hydrates templates, runs OpenTofu with the specified args, and cleans ephemeral secrets.
func (r *Runner) Execute(ctx context.Context, args []string) error {
	binPath, err := exec.LookPath(r.cfg.TofuBinary)
	if err != nil {
		return fmt.Errorf("opentofu binary %q not found in PATH: %w", r.cfg.TofuBinary, err)
	}

	engine := template.NewEngine(r.cfg, r.kr)

	fmt.Println("==> Hydrating configuration templates...")
	if err := engine.RenderAll(); err != nil {
		return err
	}

	// Always wipe ephemeral files on exit (normal, error, or signal)
	defer func() {
		if cleanErr := engine.Clean(true); cleanErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to clean ephemeral files: %v\n", cleanErr)
		}
	}()

	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = r.cfg.RootDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	// Intercept signals to forward to child process
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting %s: %w", r.cfg.TofuBinary, err)
	}

	// Signal forwarding goroutine
	doneChan := make(chan struct{})
	go func() {
		select {
		case sig := <-sigChan:
			if cmd.Process != nil {
				_ = cmd.Process.Signal(sig)
			}
		case <-doneChan:
		}
	}()

	cmdErr := cmd.Wait()
	close(doneChan)

	if cmdErr != nil {
		return fmt.Errorf("%s failed: %w", r.cfg.TofuBinary, cmdErr)
	}

	return nil
}
