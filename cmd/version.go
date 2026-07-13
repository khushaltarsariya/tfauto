package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver", "v"},
	Short:   "Show build information for tfauto",
	Long: `Show tfauto version.

Examples:
  tfauto version
  tfauto version --json`,
	Example: `  tfauto version
  tfauto version --json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if jsonRequested(cmd) {
			return writeJSON(cmd.OutOrStdout(), map[string]any{
				"command": "version",
				"ok":      true,
				"version": version,
				"commit":  commit,
				"date":    date,
			})
		}

		fmt.Printf("tfauto: version %s\n", version)
		fmt.Printf("tfauto: commit %s\n", commit)
		fmt.Printf("tfauto: built %s\n", date)
		return nil
	},
}

func init() {
	versionCmd.Flags().Bool("json", false, "Output version information as JSON")
	rootCmd.AddCommand(versionCmd)
}
