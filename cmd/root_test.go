package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestShouldShowRuntimeBanner(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
		want bool
	}{
		{
			name: "plan shows banner",
			cmd:  &cobra.Command{Use: "plan"},
			want: true,
		},
		{
			name: "version hides banner",
			cmd:  &cobra.Command{Use: "version"},
			want: false,
		},
		{
			name: "templates hides banner",
			cmd:  &cobra.Command{Use: "templates"},
			want: false,
		},
		{
			name: "nested config command hides banner",
			cmd: func() *cobra.Command {
				root := &cobra.Command{Use: "tfauto"}
				parent := &cobra.Command{Use: "config"}
				child := &cobra.Command{Use: "check"}
				root.AddCommand(parent)
				parent.AddCommand(child)
				return child
			}(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldShowRuntimeBanner(tt.cmd); got != tt.want {
				t.Fatalf("shouldShowRuntimeBanner() = %v, want %v", got, tt.want)
			}
		})
	}
}

