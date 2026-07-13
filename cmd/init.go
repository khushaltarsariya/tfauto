package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/khushaltarsariya/tfauto/internal/generator"

	"github.com/spf13/cobra"
)

var templateName string
var targetDir string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a Terraform project from a built-in template",
	Long: `Create a Terraform project from a built-in template.

Examples:
  tfauto init --template aws-basic --target ./my-project`,
	Example: `  tfauto init --template aws-basic --target ./my-project`,
	Args:    cobra.NoArgs,

	RunE: func(cmd *cobra.Command, args []string) error {
		if templateName == "" {
			return fmt.Errorf("tfauto init: missing required flag --template (for example: aws-basic)")
		}

		if targetDir == "" {
			targetDir = "./tf-project"
		}

		if _, err := os.Stat(targetDir); err == nil {
			return fmt.Errorf("tfauto init: target directory %q already exists", targetDir)
		}

		configResult, err := loadConfigForPath(filepath.Dir(targetDir))
		if err != nil {
			return tfautoError("init", err)
		}
		if configResult.Found && len(configResult.Config.Templates.Allowed) > 0 && !templateAllowed(templateName, configResult.Config.Templates.Allowed) {
			return fmt.Errorf("tfauto init: template %q is not allowed by %q", templateName, configResult.Path)
		}

		if err := generator.CopyTemplate(templateName, targetDir); err != nil {
			return fmt.Errorf("tfauto init: copy template: %w", err)
		}

		if jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command":   "init",
				"ok":        true,
				"template":  templateName,
				"target":    targetDir,
				"config":    configResult.Path,
				"generated": true,
				"message":   "template copied successfully",
			})
		}

		fmt.Println(tfautoMessage("init", "template copied to %s", targetDir))
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&templateName, "template", "aws-basic", "Template name (see `tfauto templates`)")
	initCmd.Flags().StringVar(&targetDir, "target", "./tf-project", "Target directory to copy templates")
	initCmd.Flags().Bool("json", false, "Output initialization details as JSON")
	rootCmd.AddCommand(initCmd)
}

func templateAllowed(name string, allowed []string) bool {
	for _, allowedName := range allowed {
		if name == allowedName {
			return true
		}
	}
	return false
}
