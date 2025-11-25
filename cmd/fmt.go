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
	Short: "terraform formate command {fmt}",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if fmtPath == "" {
			fmtPath = "."
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		fmt.Println("terraform formate check", fmtPath, fmtCheck)

		if err := terraform.Fmt(ctx, fmtPath, fmtCheck); err != nil {
			return fmt.Errorf("terraform fmt: %w", err)
		}
		if fmtCheck {
			fmt.Println("fmt:check complete")
		} else {
			fmt.Println("fmt:formatted succesfully")
		}
		return nil
	},
}

func init() {
	fmtCmd.Flags().StringVar(&fmtPath, "path", ".", "path to terraform project")
	fmtCmd.Flags().BoolVar(&fmtCheck, "check", false, "check if files are formated")
	rootCmd.AddCommand(fmtCmd)

}
