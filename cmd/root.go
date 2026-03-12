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
	Short: "Terraform automation cli",
	Long:  "tfauto simplifies Terraform workflows: init, plan, apply, destroy, and update variables.",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		startedAt = time.Now()

		//handle --chdir

		if globalChdir != "" {
			if err := os.Chdir(globalChdir); err != nil {
				return fmt.Errorf("chdir to %q:%w", globalChdir, err)
			}
			abs, _ := filepath.Abs(".")
			fmt.Println("Working directory:", abs)
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
			return fmt.Errorf("terraform not found in PATH: %w", err)
		}

		if globalDebug {
			_ = os.Setenv("TF_LOG", "DEBUG")
			fmt.Println("TF_LOG=DEBUG")
		}

		return nil
	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if v := cmd.Context().Value(ctxCancelKey{}); v != nil {
			if cancle, ok := v.(context.CancelFunc); ok {
				cancle()
			}
		}
		fmt.Printf("Done in %s\n", time.Since(startedAt).Round(time.Millisecond))
	},
}

type ctxCancelKey struct{}

func requiresTerraform(cmd *cobra.Command) bool {
	if cmd.Parent() == nil {
		return false
	}

	for current := cmd; current != nil; current = current.Parent() {
		switch current.Name() {
		case "init", "templates", "template", "version", "doctor", "completion", "help":
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
