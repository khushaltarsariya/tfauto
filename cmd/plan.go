package cmd

import (
	"fmt"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var pathFlag string

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan in a path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run terraform plan in ", pathFlag)
		if err := terraform.Init(pathFlag); err != nil {
			fmt.Println("terraform init failed:", err)
			return
		}

		if err := terraform.Plan(pathFlag); err != nil {
			fmt.Println("terraform plan failed", err)
			return
		}
	},
}

func init() {
	planCmd.Flags().StringVar(&pathFlag, "path", "./tf-project", "Path to terraform project")
	rootCmd.AddCommand(planCmd)
}
