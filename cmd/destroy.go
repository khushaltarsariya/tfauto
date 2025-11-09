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

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Run terraform destroy in a path(required confirmation)",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to destroy resource in %s? Type yes to confirm:", destroyPath)
		text, _ := reader.ReadString('\n')

		if strings.TrimSpace(text) != "yes" {
			fmt.Println("Aborted")
			return
		}

		if err := terraform.Init(destroyPath); err != nil {
			fmt.Println("terraform init failed", err)
			return
		}

		if err := terraform.Destroy(destroyPath); err != nil {
			fmt.Println("terraform destroy failed", err)
			return
		}
		fmt.Println("Destroy complete")
	},
}

func init() {
	destroyCmd.Flags().StringVar(&destroyPath, "path", "./tf-project", "Path to Terraform project")
	rootCmd.AddCommand(destroyCmd)
}
