package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var doctorPath string

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run basic diagnostics for the terraform environmnrt",
	Long: `Checking comman issues before running terraform:
	- Terraform installation
	- Terraform version
	- Aws credentials
	- Presence of .tf files`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if doctorPath == "" {
			doctorPath = "."
		}

		fmt.Println("tfauto doctor checks")

		//check terraform binary
		fmt.Println("Checking terraform installation")
		if _, err := exec.LookPath("terraform"); err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("terraform is not install or not in a path")
		}
		fmt.Println("OK")

		//check terraform version
		fmt.Println("Checking terraform version")
		out, err := exec.CommandContext(cmd.Context(), "terraform", "version").Output()
		if err != nil {
			fmt.Println("FAILED")
			return fmt.Errorf("failed to get terraform version")
		}

		fmt.Println("OK")
		fmt.Println(string(out))

		//check AWS credentials
		fmt.Println("Checking AWS Credentials")
		if os.Getenv("AWS_ACCESS_KEY_ID") == "" && os.Getenv("AWS_PROFILE") == "" {
			fmt.Println("WARNING")
			fmt.Println("  AWS credentials not found in environment variables")
		} else {
			fmt.Println("OK")
		}

		//check .tf files
		fmt.Println("Checking terraform file in the path")
		found := false
		err = filepath.WalkDir(doctorPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(path) == ".tf" {
				found = true
				return fs.SkipAll
			}
			return nil
		})
		if err != nil && err != fs.SkipAll {
			return err
		}

		if !found {
			fmt.Println("FAILED")
			return fmt.Errorf("no .tf file found in %s", doctorPath)
		}

		fmt.Println("Environmnet looks good.you can run terrform safely")
		return nil
	},
}

func init() {
	doctorCmd.Flags().StringVar(&doctorPath, "path", ".", "Path to terraform project")
	rootCmd.AddCommand(doctorCmd)
}
