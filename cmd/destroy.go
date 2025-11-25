package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var destroyPath string
var destroyYes bool

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Run terraform destroy in a path(required confirmation)",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if destroyPath == "" {
			destroyPath = "."
		}
		if destroyYes {
			return nil
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to destroy resource in %s? Type yes to confirm:", destroyPath)
		text, _ := reader.ReadString('\n')

		if strings.TrimSpace(text) != "yes" {
			return fmt.Errorf("Aborted")

		}

		return nil

	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if err := terraform.Init(ctx, destroyPath); err != nil {
			return fmt.Errorf("terraform init failed: %w", err)

		}

		if err := terraform.Destroy(ctx, destroyPath); err != nil {
			return fmt.Errorf("terraform destroy failed: %w", err)

		}
		fmt.Println("Destroy complete")

		return nil
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("destroy: complete")
	},
}

func init() {
	destroyCmd.Flags().StringVar(&destroyPath, "path", "./tf-project", "Path to Terraform project")
	destroyCmd.Flags().BoolVar(&destroyYes, "yes", false, "Skip interactive confirmation (non-interactive)")
	rootCmd.AddCommand(destroyCmd)
}
