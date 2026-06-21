package cmd

import (
	"fmt"

	"tfauto/internal/config"

	"github.com/spf13/cobra"
)

var configPath string

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Inspect tfauto project configuration",
}

var configCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the active .tfauto.yaml configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := config.LoadForPath(configPath)
		if err != nil {
			return err
		}
		if !result.Found {
			fmt.Println("No .tfauto.yaml found.")
			fmt.Println("tfauto will use command defaults.")
			return nil
		}

		printConfigSummary(result)
		return nil
	},
}

func printConfigSummary(result config.LoadResult) {
	cfg := result.Config

	fmt.Println("tfauto config")
	fmt.Println()
	fmt.Println("Path:", result.Path)
	if cfg.Project != "" {
		fmt.Println("Project:", cfg.Project)
	}
	if cfg.Environment != "" {
		fmt.Println("Environment:", cfg.Environment)
	}
	fmt.Println()
	fmt.Println("Terraform:")
	fmt.Println("  require_plan_file:", cfg.Terraform.RequirePlanFile)
	fmt.Println("  protect_destroy:", cfg.Terraform.ProtectDestroy)

	if len(cfg.Templates.Allowed) > 0 {
		fmt.Println()
		fmt.Println("Allowed templates:")
		for _, name := range cfg.Templates.Allowed {
			fmt.Println("  -", name)
		}
	}

	if len(cfg.Policy.RequireTags) > 0 {
		fmt.Println()
		fmt.Println("Required tags:")
		for _, name := range cfg.Policy.RequireTags {
			fmt.Println("  -", name)
		}
	}
}

func init() {
	configCheckCmd.Flags().StringVar(&configPath, "path", ".", "Path to Terraform project")
	configCmd.AddCommand(configCheckCmd)
	rootCmd.AddCommand(configCmd)
}
