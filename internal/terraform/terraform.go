package terraform

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type PlanResult struct {
	HasChanges bool
}

func runInDirWithExitCode(ctx context.Context, dir string, name string, args ...string) (int, error) {
	if _, err := exec.LookPath(name); err != nil {
		return 0, fmt.Errorf("%s not found in PATH %w", name, err)
	}

	if dir == "" {
		dir = "."
	}

	if _, err := os.Stat(dir); err != nil {
		return 0, fmt.Errorf("path %s does not exist %w ", dir, err)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err == nil {
		return 0, nil
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), err
	}

	return 0, err
}

func runInDir(ctx context.Context, dir string, name string, args ...string) error {
	_, err := runInDirWithExitCode(ctx, dir, name, args...)
	return err
}

func Init(ctx context.Context, path string) error {
	fmt.Println("Terraform Init")
	return runInDir(ctx, path, "terraform", "init", "-input=false")
}

func InitForValidation(ctx context.Context, path string) error {
	fmt.Println("terraform init (validation mode)")
	return runInDir(ctx, path, "terraform", "init", "-backend=false", "-input=false")
}

func Plan(ctx context.Context, path string, detailedExitCode bool) (PlanResult, error) {
	fmt.Println("terraform plan")
	args := []string{"plan", "-input=false"}
	if detailedExitCode {
		args = append(args, "-detailed-exitcode")
	}

	code, err := runInDirWithExitCode(ctx, path, "terraform", args...)
	if err != nil {
		if detailedExitCode && code == 2 {
			return PlanResult{HasChanges: true}, nil
		}
		return PlanResult{}, err
	}

	return PlanResult{}, nil
}

func PlanOut(ctx context.Context, path string, outFile string, detailedExitCode bool) (PlanResult, error) {
	outPath := outFile
	if filepath.IsAbs(outFile) {
		// outPath = filepath.Join(path, outFile)
	} else {
		outPath = outFile
	}
	fmt.Println("terraform plan -out", outPath)
	args := []string{"plan", "-input=false", "-out", outPath}
	if detailedExitCode {
		args = append(args, "-detailed-exitcode")
	}

	code, err := runInDirWithExitCode(ctx, path, "terraform", args...)
	if err != nil {
		if detailedExitCode && code == 2 {
			return PlanResult{HasChanges: true}, nil
		}
		return PlanResult{}, err
	}

	return PlanResult{}, nil
}

func Apply(ctx context.Context, path string) error {
	fmt.Println("terraform apply")
	return runInDir(ctx, path, "terraform", "apply", "-auto-approve", "-input=false")
}

func ApplyPlan(ctx context.Context, path string, planFile string) error {
	fmt.Println("terraform apply saved plan", planFile)
	return runInDir(ctx, path, "terraform", "apply", "-input=false", planFile)
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
