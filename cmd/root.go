package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	globalChdir    string
	globalDebug    bool
	gloabalTimeout time.Duration
	startedAt      time.Time
)

var rootCmd = &cobra.Command{
	Use:   "tfauto",
	Short: "Terraform workflow automation CLI",
	Long:  "tfauto standardizes Terraform workflows with safer commands, reusable templates, and consistent project setup.",
	Example: `  tfauto version
  tfauto templates
  tfauto init --template aws-basic --target ./my-project
  tfauto plan --path ./my-project
  tfauto doctor --path ./my-project`,
	SilenceUsage:  true,
	SilenceErrors: true,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		startedAt = time.Now()
		jsonOutput := jsonRequested(cmd)

		// Handle --chdir before running any command.

		if globalChdir != "" {
			if err := os.Chdir(globalChdir); err != nil {
				return fmt.Errorf("tfauto: chdir to %q: %w", globalChdir, err)
			}
			if !jsonOutput && shouldShowRuntimeBanner(cmd) {
				abs, _ := filepath.Abs(".")
				fmt.Printf("tfauto: working directory %s\n", abs)
			}
		}

		if gloabalTimeout <= 0 {
			gloabalTimeout = 15 * time.Minute
		}

		cxt, cancle := context.WithTimeout(context.Background(), gloabalTimeout)
		ctx := context.WithValue(cxt, ctxCancelKey{}, cancle)
		cmd.SetContext(ctx)

		if !requiresTerraform(cmd) {
			return nil
		}

		if _, err := exec.LookPath("terraform"); err != nil {
			return fmt.Errorf("tfauto: terraform not found in PATH: %w", err)
		}

		if globalDebug && !jsonOutput {
			_ = os.Setenv("TF_LOG", "DEBUG")
			fmt.Println("tfauto: TF_LOG=DEBUG")
		} else if globalDebug {
			_ = os.Setenv("TF_LOG", "DEBUG")
		}

		return nil
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if v := cmd.Context().Value(ctxCancelKey{}); v != nil {
			if cancle, ok := v.(context.CancelFunc); ok {
				cancle()
			}
		}
		if !jsonRequested(cmd) && shouldShowRuntimeBanner(cmd) {
			fmt.Printf("tfauto: completed in %s\n", time.Since(startedAt).Round(time.Millisecond))
		}
	},
}

type ctxCancelKey struct{}

func requiresTerraform(cmd *cobra.Command) bool {
	if cmd.Parent() == nil {
		return false
	}

	for current := cmd; current != nil; current = current.Parent() {
		switch current.Name() {
		case "init", "templates", "template", "version", "doctor", "config", "completion", "help":
			return false
		}
	}

	return true
}

func shouldShowRuntimeBanner(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}

	for current := cmd; current != nil; current = current.Parent() {
		switch current.Name() {
		case "version", "templates", "template", "config", "doctor", "completion", "help":
			return false
		}
	}

	return true
}

func init() {
	rootCmd.PersistentFlags().DurationVar(&gloabalTimeout, "timeout", 15*time.Minute, "Global timeout for terraform operations")
	rootCmd.PersistentFlags().BoolVar(&globalDebug, "debug", false, "Enable debug logs (sets TF_LOG=DEBUG)")
	rootCmd.PersistentFlags().StringVar(&globalChdir, "chdir", "", "Change working directory before running (like terraform -chdir)")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
