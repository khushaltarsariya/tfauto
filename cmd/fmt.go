package cmd

import (
	"fmt"

	"github.com/khushaltarsariya/tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var fmtPath string
var fmtCheck bool

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format Terraform files in a project directory",
	Long: `Format Terraform files in a project directory.

Examples:
  tfauto fmt --path ./app
  tfauto fmt --path ./app --check`,
	Example: `  tfauto fmt --path ./app
  tfauto fmt --path ./app --check`,
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
			return fmt.Errorf("tfauto fmt: terraform fmt failed: %w", err)
		}
		if fmtCheck {
			fmt.Println(tfautoMessage("fmt", "check completed successfully"))
		} else {
			fmt.Println(tfautoMessage("fmt", "completed successfully"))
		}
		return nil
	},
}

func init() {
	fmtCmd.Flags().StringVar(&fmtPath, "path", ".", "Path to Terraform project")
	fmtCmd.Flags().BoolVar(&fmtCheck, "check", false, "Check whether files are already formatted")
	rootCmd.AddCommand(fmtCmd)

}
