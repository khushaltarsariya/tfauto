package cmd

import (
	"fmt"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var fmtPath string
var fmtCheck bool

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format Terraform files in a project directory",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if fmtPath == "" {
			fmtPath = "."
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		fmt.Println("Running terraform fmt in", fmtPath)

		if err := terraform.Fmt(ctx, fmtPath, fmtCheck); err != nil {
			return fmt.Errorf("terraform fmt: %w", err)
		}
		if fmtCheck {
			fmt.Println("fmt check completed")
		} else {
			fmt.Println("Terraform files formatted successfully")
		}
		return nil
	},
}

func init() {
	fmtCmd.Flags().StringVar(&fmtPath, "path", ".", "Path to Terraform project")
	fmtCmd.Flags().BoolVar(&fmtCheck, "check", false, "Check whether files are already formatted")
	rootCmd.AddCommand(fmtCmd)

}
