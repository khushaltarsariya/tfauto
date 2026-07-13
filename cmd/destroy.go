package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/khushaltarsariya/tfauto/internal/terraform"

	"github.com/spf13/cobra"
)

var destroyPath string
var destroyYes bool
var destroyJSON bool

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy Terraform-managed infrastructure",
	Long: `Run terraform destroy with confirmation.

Examples:
  tfauto destroy --path ./app --yes
  tfauto destroy --path ./app
  tfauto destroy --path ./app --json`,
	Example: `  tfauto destroy --path ./app
  tfauto destroy --path ./app --yes
  tfauto destroy --path ./app --json`,
	Args: cobra.NoArgs,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if destroyPath == "" {
			destroyPath = "."
		}
		configResult, err := loadConfigForPath(destroyPath)
		if err != nil {
			return tfautoError("destroy", err)
		}
		if configResult.Found && configResult.Config.Terraform.ProtectDestroy && !destroyYes {
			return fmt.Errorf("tfauto destroy: destroy is protected by %q; pass --yes only after reviewing the plan and policy", configResult.Path)
		}
		if destroyYes {
			return nil
		}
		if isNonInteractive() {
			return fmt.Errorf("tfauto destroy: non-interactive mode requires --yes")
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s", tfautoMessage("destroy", "are you sure you want to destroy resources in %s? Type yes to confirm: ", destroyPath))
		text, _ := reader.ReadString('\n')

		if strings.TrimSpace(text) != "yes" {
			return fmt.Errorf("tfauto destroy: aborted")
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if err := terraform.Init(ctx, destroyPath); err != nil {
			return fmt.Errorf("tfauto destroy: terraform init: %w", err)
		}

		if err := terraform.Destroy(ctx, destroyPath); err != nil {
			return fmt.Errorf("tfauto destroy: terraform destroy: %w", err)
		}
		if destroyJSON || jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command": "destroy",
				"ok":      true,
				"path":    destroyPath,
			})
		}
		fmt.Println(tfautoMessage("destroy", "completed successfully"))
		return nil
	},
}

func init() {
	destroyCmd.Flags().StringVar(&destroyPath, "path", ".", "Path to Terraform project")
	destroyCmd.Flags().BoolVar(&destroyYes, "yes", false, "Skip interactive confirmation (non-interactive)")
	destroyCmd.Flags().BoolVar(&destroyJSON, "json", false, "Output destroy details as JSON")
	rootCmd.AddCommand(destroyCmd)
}

func isNonInteractive() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (info.Mode() & os.ModeCharDevice) == 0
}
