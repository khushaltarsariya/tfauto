package terraform

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func runInDir(ctx context.Context, dir string, name string, args ...string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("%s not found in PATH %w", name, err)
	}

	if dir == "" {
		dir = "."
	}

	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("path %s does not exist %w ", dir, err)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func Init(ctx context.Context, path string) error {
	fmt.Println("Terraform Init")
	return runInDir(ctx, path, "terraform", "init", "-input=false")
}

func InitForValidation(ctx context.Context, path string) error {
	fmt.Println("terraform init (validation mode)")
	return runInDir(ctx, path, "terraform", "init", "-backend=false", "-input=false")
}

func Plan(ctx context.Context, path string) error {
	fmt.Println("terraform plan")
	return runInDir(ctx, path, "terraform", "plan", "-input=false")
}

func PlanOut(ctx context.Context, path string, outFile string) error {
	outPath := outFile
	if filepath.IsAbs(outFile) {
		// outPath = filepath.Join(path, outFile)
	} else {
		outPath = outFile
	}
	fmt.Println("terraform plan -out", outPath)
	return runInDir(ctx, path, "terraform", "plan", "-input=false", "-out", outPath)
}

func Apply(ctx context.Context, path string) error {
	fmt.Println("terraform apply")
	return runInDir(ctx, path, "terraform", "apply", "-auto-approve", "-input=false")
}

func Validate(ctx context.Context, path string) error {
	fmt.Println("terraform validate")
	return runInDir(ctx, path, "terraform", "validate")
}

func Fmt(ctx context.Context, path string, check bool) error {
	fmt.Println("terraform fmt")
	args := []string{"fmt"}
	if check {
		args = append(args, "-check")
	}
	err := runInDir(ctx, path, "terraform", args...)

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			if code == 2 || code == 3 {
				return fmt.Errorf("unformatted files detected")
			}
		}
		return err
	}
	return nil
}

func Destroy(ctx context.Context, path string) error {
	fmt.Println("terraform destroy")
	return runInDir(ctx, path, "terraform", "destroy", "-auto-approve", "-input=false")
}
