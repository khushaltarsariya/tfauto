package cmd

import (
	"fmt"

	"tfauto/internal/config"
	"tfauto/internal/doctor"

	"github.com/spf13/cobra"
)

var doctorPath string
var doctorJSON bool

type doctorJSONOutput struct {
	Command string                `json:"command"`
	Path    string                `json:"path"`
	Results []doctor.Result       `json:"results"`
	Config  doctorConfigJSONState `json:"config"`
	OK      bool                  `json:"ok"`
}

type doctorConfigJSONState struct {
	Found            bool     `json:"found"`
	Path             string   `json:"path,omitempty"`
	RequirePlanFile  bool     `json:"require_plan_file,omitempty"`
	ProtectDestroy   bool     `json:"protect_destroy,omitempty"`
	AllowedTemplates []string `json:"allowed_templates,omitempty"`
	RequiredTags     []string `json:"required_tags,omitempty"`
	Error            string   `json:"error,omitempty"`
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run environment diagnostics for a Terraform project",
	Args:  cobra.NoArgs,
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

		if doctorJSON {
			output := doctorJSONOutput{
				Command: "doctor",
				Path:    doctorPath,
				Results: report.Results,
				OK:      configErr == nil && !report.HasFailures(),
			}
			if configErr != nil {
				output.Config = doctorConfigJSONState{
					Found: false,
					Error: configErr.Error(),
				}
			} else {
				output.Config = doctorConfigJSONState{
					Found:            configResult.Found,
					Path:             configResult.Path,
					RequirePlanFile:  configResult.Config.Terraform.RequirePlanFile,
					ProtectDestroy:   configResult.Config.Terraform.ProtectDestroy,
					AllowedTemplates: configResult.Config.Templates.Allowed,
					RequiredTags:     configResult.Config.Policy.RequireTags,
				}
			}

			if err := writeJSON(cmd.OutOrStdout(), output); err != nil {
				return err
			}

			if configErr != nil || report.HasFailures() {
				return fmt.Errorf("tfauto doctor: one or more blocking issues found")
			}
			return nil
		}

		fmt.Println("tfauto: doctor")
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
			return fmt.Errorf("tfauto doctor: one or more blocking issues found")
		}
		if report.HasFailures() {
			return fmt.Errorf("tfauto doctor: one or more blocking issues found")
		}

		fmt.Println(tfautoMessage("doctor", "completed. no blocking issues found"))
		return nil
	},
}

func init() {
	doctorCmd.Flags().StringVar(&doctorPath, "path", ".", "Path to Terraform project")
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Output diagnostics as JSON")
	rootCmd.AddCommand(doctorCmd)
}
