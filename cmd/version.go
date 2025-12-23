package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show terraform version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tfauto version", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
