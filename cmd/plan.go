package cmd

import (
	"fmt"

	"github.com/khushaltarsariya/tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var pathFlag string
var planOut string
var planDetailedExitCode bool

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate a Terraform plan for a project directory",
	Long: `Run terraform plan in a project directory.

Examples:
  tfauto plan --path ./app
  tfauto plan --path ./app --out tfplan --detailed-exitcode`,
	Example: `  tfauto plan --path ./app
  tfauto plan --path ./app --out tfplan --detailed-exitcode`,
	Args: cobra.NoArgs,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if pathFlag == "" {
			pathFlag = "."
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if err := terraform.Init(ctx, pathFlag); err != nil {
			return fmt.Errorf("tfauto plan: terraform init failed: %w", err)

		}

		if planOut != "" {
			result, err := terraform.PlanOut(ctx, pathFlag, planOut, planDetailedExitCode)
			if err != nil {
				return fmt.Errorf("tfauto plan: terraform plan -out: %w", err)
			}
			fmt.Println(tfautoMessage("plan", "wrote plan to %s", planOut))
			if result.HasChanges {
				fmt.Println(tfautoMessage("plan", "Terraform plan detected pending changes"))
				return ExitError{Code: 2}
			}
			fmt.Println(tfautoMessage("plan", "completed successfully"))
			return nil
		}
		result, err := terraform.Plan(ctx, pathFlag, planDetailedExitCode)
		if err != nil {
			return fmt.Errorf("tfauto plan: terraform plan failed: %w", err)
		}
		if result.HasChanges {
			fmt.Println(tfautoMessage("plan", "Terraform plan detected pending changes"))
			return ExitError{Code: 2}
		}
		fmt.Println(tfautoMessage("plan", "completed successfully"))
		return nil
	},
}

func init() {
	planCmd.Flags().StringVar(&pathFlag, "path", ".", "Path to Terraform project")
	planCmd.Flags().StringVar(&planOut, "out", "", "Write the plan to a file")
	planCmd.Flags().BoolVar(&planDetailedExitCode, "detailed-exitcode", false, "Return exit code 2 when the plan contains changes")
	rootCmd.AddCommand(planCmd)
}
