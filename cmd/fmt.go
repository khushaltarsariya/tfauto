package cmd

import (
	"fmt"

	"github.com/khushaltarsariya/tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var fmtPath string
var fmtCheck bool
var fmtJSON bool

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format Terraform files",
	Long: `Format Terraform files in a project directory.

Examples:
  tfauto fmt --path ./app
  tfauto fmt --path ./app --check
  tfauto fmt --path ./app --json`,
	Example: `  tfauto fmt --path ./app
  tfauto fmt --path ./app --check
  tfauto fmt --path ./app --json`,
	Args: cobra.NoArgs,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if fmtPath == "" {
			fmtPath = "."
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if err := terraform.Fmt(ctx, fmtPath, fmtCheck); err != nil {
			return fmt.Errorf("tfauto fmt: terraform fmt: %w", err)
		}
		if fmtJSON || jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command": "fmt",
				"ok":      true,
				"path":    fmtPath,
				"check":   fmtCheck,
			})
		}
		if fmtCheck {
			fmt.Println(tfautoMessage("fmt", "completed successfully (check only)"))
		} else {
			fmt.Println(tfautoMessage("fmt", "completed successfully"))
		}
		return nil
	},
}

func init() {
	fmtCmd.Flags().StringVar(&fmtPath, "path", ".", "Path to Terraform project")
	fmtCmd.Flags().BoolVar(&fmtCheck, "check", false, "Check whether files are already formatted")
	fmtCmd.Flags().BoolVar(&fmtJSON, "json", false, "Output formatting details as JSON")
	rootCmd.AddCommand(fmtCmd)

}
