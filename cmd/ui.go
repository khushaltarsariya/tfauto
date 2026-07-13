package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"tfauto/internal/config"

	"github.com/spf13/cobra"
)

func tfautoError(command string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("tfauto %s: %w", command, err)
}

func tfautoMessage(command, format string, args ...any) string {
	prefixArgs := append([]any{command}, args...)
	return fmt.Sprintf("tfauto %s: "+format, prefixArgs...)
}

func writeJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func loadConfigForPath(path string) (config.LoadResult, error) {
	return config.LoadForPath(path)
}

func jsonRequested(cmd *cobra.Command) bool {
	for current := cmd; current != nil; current = current.Parent() {
		flag := current.Flags().Lookup("json")
		if flag == nil {
			continue
		}
		value, err := current.Flags().GetBool("json")
		return err == nil && value
	}

	return false
}
