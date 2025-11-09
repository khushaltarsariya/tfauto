package cmd

import (
	"fmt"
	"os"
	"tfauto/internal/generator"

	"github.com/spf13/cobra"
)

var templateName string
var targetDir string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a terraform template into a directory",

	Run: func(cmd *cobra.Command, args []string) {
		if templateName == "" {
			fmt.Println("Please specify --template (e.g., aws-basic)")
			return
		}

		if targetDir == "" {
			targetDir = "./tf-project"
		}

		if _, err := os.Stat(targetDir); err == nil {
			fmt.Printf("Target directory %s allreday exists. Aborting ", targetDir)
			return
		}

		if err := generator.CopyTemplate(templateName, targetDir); err != nil {
			fmt.Println("Error coping template", err)
			return
		}

		fmt.Println("Template copy to", targetDir)
	},
}

func init() {
	initCmd.Flags().StringVar(&templateName, "template", "aws-basic", "template name in templates")
	initCmd.Flags().StringVar(&targetDir, "target", "./tf-project", "Target directory to copy templates")
	rootCmd.AddCommand(initCmd)
}
