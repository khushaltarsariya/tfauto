package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"tfauto/internal/config"
	"tfauto/internal/generator"

	"github.com/spf13/cobra"
)

var templateName string
var targetDir string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a Terraform project from a built-in template",

	RunE: func(cmd *cobra.Command, args []string) error {
		if templateName == "" {
			return fmt.Errorf("please specify --template (for example: aws-basic)")
		}

		if targetDir == "" {
			targetDir = "./tf-project"
		}

		if _, err := os.Stat(targetDir); err == nil {
			return fmt.Errorf("target directory %s already exists", targetDir)
		}

		configResult, err := config.LoadForPath(filepath.Dir(targetDir))
		if err != nil {
			return err
		}
		if configResult.Found && len(configResult.Config.Templates.Allowed) > 0 && !templateAllowed(templateName, configResult.Config.Templates.Allowed) {
			return fmt.Errorf("template %q is not allowed by %s", templateName, configResult.Path)
		}

		if err := generator.CopyTemplate(templateName, targetDir); err != nil {
			return fmt.Errorf("copy template: %w", err)
		}

		fmt.Println("Template copied to", targetDir)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&templateName, "template", "aws-basic", "Template name (see `tfauto templates`)")
	initCmd.Flags().StringVar(&targetDir, "target", "./tf-project", "Target directory to copy templates")
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
