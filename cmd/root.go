package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tfauto",
	Short: "Terraform automation cli",
	Long:  "tfauto simplifies Terraform workflows: init, plan, apply, destroy, and update variables.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
