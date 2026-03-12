package cmd

import (
	"fmt"
	"tfauto/internal/terraform"
	"time"

	"github.com/spf13/cobra"
)

var pathFlag string
var planOut string

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Run terraform plan in a project directory",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if pathFlag == "" {
			pathFlag = "."
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		fmt.Println("Running terraform plan in", pathFlag)

		if err := terraform.Init(ctx, pathFlag); err != nil {
			return fmt.Errorf("terraform init failed: %w", err)

		}

		if planOut != "" {
			if err := terraform.PlanOut(ctx, pathFlag, planOut); err != nil {
				return fmt.Errorf("terraform plan -out: %w ", err)
			}
			fmt.Println("Plan written to", planOut)
			return nil
		}
		if err := terraform.Plan(ctx, pathFlag); err != nil {
			return fmt.Errorf("terraform plan failed: %w", err)

		}
		return nil
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("plan: completed at", time.Now().Format(time.RFC3339))
	},
}

func init() {
	planCmd.Flags().StringVar(&pathFlag, "path", ".", "Path to Terraform project")
	planCmd.Flags().StringVar(&planOut, "out", "", "Write the plan to a file")
	rootCmd.AddCommand(planCmd)
}
