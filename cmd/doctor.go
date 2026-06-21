package cmd

import (
	"fmt"

	"tfauto/internal/config"
	"tfauto/internal/doctor"

	"github.com/spf13/cobra"
)

var doctorPath string

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run environment diagnostics for a Terraform project",
	Long: `Run environment diagnostics before Terraform operations.

Checks include:
- target path validity
- Terraform installation and version
- presence of Terraform files
- project initialization state
- backend and workspace hints
- variable prompt risk detection
- AWS region resolution
- AWS CLI and caller identity when available`,
	RunE: func(cmd *cobra.Command, args []string) error {
		report := doctor.Run(cmd.Context(), doctorPath)
		configResult, configErr := config.LoadForPath(doctorPath)

		fmt.Println("tfauto doctor")
		fmt.Println()

		for _, result := range report.Results {
			fmt.Printf("[%s] %s\n", result.Status, result.Name)
			for _, detail := range result.Details {
				fmt.Printf("  - %s\n", detail)
			}
			fmt.Println()
		}

		if configErr != nil {
			fmt.Println("[FAIL] tfauto config")
			fmt.Printf("  - %s\n\n", configErr)
		} else if configResult.Found {
			fmt.Println("[PASS] tfauto config")
			fmt.Printf("  - Found %s\n", configResult.Path)
			fmt.Printf("  - require_plan_file: %t\n", configResult.Config.Terraform.RequirePlanFile)
			fmt.Printf("  - protect_destroy: %t\n", configResult.Config.Terraform.ProtectDestroy)
			if len(configResult.Config.Policy.RequireTags) > 0 {
				fmt.Printf("  - required tags: %v\n", configResult.Config.Policy.RequireTags)
			}
			fmt.Println()
		} else {
			fmt.Println("[WARN] tfauto config")
			fmt.Println("  - No .tfauto.yaml found")
			fmt.Println("  - Command defaults will be used")
			fmt.Println()
		}

		if configErr != nil {
			return fmt.Errorf("doctor found one or more blocking issues")
		}
		if report.HasFailures() {
			return fmt.Errorf("doctor found one or more blocking issues")
		}

		fmt.Println("Doctor completed. No blocking issues found.")
		return nil
	},
}

func init() {
	doctorCmd.Flags().StringVar(&doctorPath, "path", ".", "Path to Terraform project")
	rootCmd.AddCommand(doctorCmd)
}
