package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tfauto/internal/config"
	"tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var destroyPath string
var destroyYes bool

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Run terraform destroy with confirmation",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if destroyPath == "" {
			destroyPath = "."
		}
		configResult, err := config.LoadForPath(destroyPath)
		if err != nil {
			return err
		}
		if configResult.Found && configResult.Config.Terraform.ProtectDestroy && !destroyYes {
			return fmt.Errorf("destroy is protected by %s; pass --yes only after reviewing the plan and policy", configResult.Path)
		}
		if destroyYes {
			return nil
		}
		if isNonInteractive() {
			return fmt.Errorf("destroy requires --yes in non-interactive mode")
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to destroy resources in %s? Type yes to confirm: ", destroyPath)
		text, _ := reader.ReadString('\n')

		if strings.TrimSpace(text) != "yes" {
			return fmt.Errorf("aborted")
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
		fmt.Println("Destroy completed")
	},
}

func init() {
	destroyCmd.Flags().StringVar(&destroyPath, "path", ".", "Path to Terraform project")
	destroyCmd.Flags().BoolVar(&destroyYes, "yes", false, "Skip interactive confirmation (non-interactive)")
	rootCmd.AddCommand(destroyCmd)
}

func isNonInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (info.Mode() & os.ModeCharDevice) == 0
}
