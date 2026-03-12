package cmd

import (
	"fmt"
	"tfauto/internal/terraform"
	"time"

	"github.com/spf13/cobra"
)

var pathFlag string
var planOut string
var planDetailedExitCode bool

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
			result, err := terraform.PlanOut(ctx, pathFlag, planOut, planDetailedExitCode)
			if err != nil {
				return fmt.Errorf("terraform plan -out: %w ", err)
			}
			fmt.Println("Plan written to", planOut)
			if result.HasChanges {
				fmt.Println("Terraform plan detected pending changes.")
				return ExitError{Code: 2}
			}
			return nil
		}
		result, err := terraform.Plan(ctx, pathFlag, planDetailedExitCode)
		if err != nil {
			return fmt.Errorf("terraform plan failed: %w", err)
		}
		if result.HasChanges {
			fmt.Println("Terraform plan detected pending changes.")
			return ExitError{Code: 2}
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
	planCmd.Flags().BoolVar(&planDetailedExitCode, "detailed-exitcode", false, "Return exit code 2 when the plan contains changes")
	rootCmd.AddCommand(planCmd)
}
