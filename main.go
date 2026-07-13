package main

import (
	"fmt"
	"os"

	"github.com/khushaltarsariya/tfauto/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		exitCode := 1
		if exitCoder, ok := err.(interface{ ExitCode() int }); ok {
			exitCode = exitCoder.ExitCode()
		}
		if err.Error() != "" {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(exitCode)
	}
}
