package cmd

import (
	"fmt"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var applyPath string

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run terraform apply in a path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running terraform apply in ", applyPath)

		if err := terraform.Init(applyPath); err != nil {
			fmt.Println("terraform init failed", err)
			return
		}

		if err := terraform.Apply(applyPath); err != nil {
			fmt.Println("terraform apply failed", err)
			return
		}
	},
}

func init() {
	applyCmd.Flags().StringVar(&applyPath, "path", "./tf-project", "Path to Terraform project")
	applyCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(applyCmd)
}
