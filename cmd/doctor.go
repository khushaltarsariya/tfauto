package cmd

import (
	"fmt"
	"strings"

	"github.com/khushaltarsariya/tfauto/internal/config"
	"github.com/khushaltarsariya/tfauto/internal/doctor"

	"github.com/spf13/cobra"
)

var doctorPath string
var doctorJSON bool

type doctorJSONOutput struct {
	Command          string                `json:"command"`
	Path             string                `json:"path"`
	CurrentDirectory string                `json:"current_directory"`
	Results          []doctor.Result       `json:"results"`
	Summary          doctor.Summary        `json:"summary"`
	Config           doctorConfigJSONState `json:"config"`
	OK               bool                  `json:"ok"`
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
	Use:     "doctor",
	Aliases: []string{"check", "diagnose"},
	Short:   "Inspect a Terraform project and its environment",
	Args:    cobra.NoArgs,
	Long: `Run environment diagnostics before Terraform operations.

Checks include:
- Filesystem: current path, terraform binary, version, current directory, write permissions, Terraform files
- Terraform: fmt, validate, backend configuration, providers, modules, initialization, workspace, variable prompts
- AWS: credentials, profile, caller identity, region
- Git: git binary and current branch`,
	Example: `  tfauto doctor --path ./app
  tfauto doctor --path ./app --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		report := doctor.Run(cmd.Context(), doctorPath)
		configResult, configErr := config.LoadForPath(doctorPath)

		if doctorJSON {
			output := doctorJSONOutput{
				Command:          "doctor",
				Path:             doctorPath,
				CurrentDirectory: report.CurrentDirectory,
				Results:          report.Results,
				Summary:          report.Summary,
				OK:               configErr == nil && !report.HasFailures() && !report.HasWarnings(),
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

			return doctorExit(report, configErr)
		}

		fmt.Println("tfauto: doctor")
		fmt.Println()
		renderDoctorReport(report)

		if configErr != nil {
			fmt.Println("Config")
			fmt.Println("------")
			fmt.Printf("[FAIL] %s\n\n", configErr)
		} else if configResult.Found {
			fmt.Println("Config")
			fmt.Println("------")
			fmt.Printf("[PASS] %s\n", configResult.Path)
			fmt.Printf("  - require_plan_file: %t\n", configResult.Config.Terraform.RequirePlanFile)
			fmt.Printf("  - protect_destroy: %t\n", configResult.Config.Terraform.ProtectDestroy)
			if len(configResult.Config.Policy.RequireTags) > 0 {
				fmt.Printf("  - required tags: %v\n", configResult.Config.Policy.RequireTags)
			}
			fmt.Println()
		} else {
			fmt.Println("Config")
			fmt.Println("------")
			fmt.Println("[WARN] No .tfauto.yaml found")
			fmt.Println("  - Command defaults will be used")
			fmt.Println()
		}

		return doctorExit(report, configErr)
	},
}

func init() {
	doctorCmd.Flags().StringVar(&doctorPath, "path", ".", "Path to Terraform project")
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Output diagnostics as JSON")
	rootCmd.AddCommand(doctorCmd)
}

func renderDoctorReport(report doctor.Report) {
	currentSection := ""
	for _, result := range report.Results {
		if result.Section != currentSection {
			currentSection = result.Section
			fmt.Println(currentSection)
			fmt.Println(strings.Repeat("-", len(currentSection)))
		}

		fmt.Printf("[%s] %s\n", result.Status, result.Name)
		for _, detail := range result.Details {
			fmt.Printf("  - %s\n", detail)
		}
		fmt.Println()
	}
	fmt.Printf("Summary: %d passed, %d warnings, %d failed\n\n", report.Summary.Pass, report.Summary.Warn, report.Summary.Fail)
}

func doctorExit(report doctor.Report, configErr error) error {
	if configErr != nil {
		return ExitError{Code: 1, Message: fmt.Sprintf("tfauto doctor: %v", configErr)}
	}
	if report.HasFailures() {
		return ExitError{Code: 1, Message: "tfauto doctor: one or more blocking issues found"}
	}
	if report.HasWarnings() {
		return ExitError{Code: 2, Message: "tfauto doctor: completed with warnings"}
	}
	return nil
}
