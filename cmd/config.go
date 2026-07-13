package cmd

import (
	"fmt"

	"github.com/khushaltarsariya/tfauto/internal/config"

	"github.com/spf13/cobra"
)

var configPath string
var configJSON bool

type configJSONOutput struct {
	Command string           `json:"command"`
	Path    string           `json:"path"`
	Found   bool             `json:"found"`
	OK      bool             `json:"ok"`
	Details configJSONDetail `json:"details"`
	Error   string           `json:"error,omitempty"`
}

type configJSONDetail struct {
	ConfigPath       string   `json:"config_path,omitempty"`
	Project          string   `json:"project,omitempty"`
	Environment      string   `json:"environment,omitempty"`
	RequirePlanFile  bool     `json:"require_plan_file,omitempty"`
	ProtectDestroy   bool     `json:"protect_destroy,omitempty"`
	AllowedTemplates []string `json:"allowed_templates,omitempty"`
	RequiredTags     []string `json:"required_tags,omitempty"`
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Inspect tfauto project configuration",
	Long: `Inspect the active tfauto project configuration.

Use this command to confirm whether a .tfauto.yaml file is being applied.`,
	Example: `  tfauto config check --path ./app
  tfauto config check --path ./app --json`,
	Args: cobra.NoArgs,
}

var configCheckCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"inspect"},
	Short:   "Check the active .tfauto.yaml configuration",
	Long: `Check the active .tfauto.yaml configuration and show the resolved project rules.

Use this command before init, apply, or destroy when you want to confirm policy behavior.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := config.LoadForPath(configPath)
		if configJSON {
			output := configJSONOutput{
				Command: "config check",
				Path:    configPath,
				OK:      err == nil,
			}
			if err != nil {
				output.Error = err.Error()
			} else {
				output.Found = result.Found
				if result.Found {
					output.Details = configJSONDetail{
						ConfigPath:       result.Path,
						Project:          result.Config.Project,
						Environment:      result.Config.Environment,
						RequirePlanFile:  result.Config.Terraform.RequirePlanFile,
						ProtectDestroy:   result.Config.Terraform.ProtectDestroy,
						AllowedTemplates: result.Config.Templates.Allowed,
						RequiredTags:     result.Config.Policy.RequireTags,
					}
				}
			}

			if encodeErr := writeJSON(cmd.OutOrStdout(), output); encodeErr != nil {
				return encodeErr
			}

			return err
		}

		if err != nil {
			return err
		}
		if !result.Found {
			fmt.Println("tfauto config: no .tfauto.yaml found")
			fmt.Println("tfauto config: command defaults will be used")
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
	configCheckCmd.Flags().BoolVar(&configJSON, "json", false, "Output configuration as JSON")
	configCmd.AddCommand(configCheckCmd)
	rootCmd.AddCommand(configCmd)
}
