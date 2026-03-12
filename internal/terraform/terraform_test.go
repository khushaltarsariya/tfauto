package terraform

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunInDirWithExitCodeMissingPath(t *testing.T) {
	_, err := runInDirWithExitCode(context.Background(), filepath.Join(t.TempDir(), "missing"), "terraform", "plan")
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestFmtReturnsFriendlyErrorForFormattingExitCodes(t *testing.T) {
	setupFakeTerraform(t)
	t.Setenv("TF_HELPER_FMT_EXIT", "3")

	err := Fmt(context.Background(), t.TempDir(), true)
	if err == nil {
		t.Fatal("expected fmt error")
	}
	if !strings.Contains(err.Error(), "unformatted files detected") {
		t.Fatalf("unexpected fmt error: %v", err)
	}
}

func TestPlanDetailedExitCodeReturnsChangesWithoutError(t *testing.T) {
	setupFakeTerraform(t)
	t.Setenv("TF_HELPER_PLAN_EXIT", "2")

	result, err := Plan(context.Background(), t.TempDir(), true)
	if err != nil {
		t.Fatalf("Plan returned error: %v", err)
	}
	if !result.HasChanges {
		t.Fatal("expected plan result to report changes")
	}
}

func TestApplyPlanPassesPlanFileArgument(t *testing.T) {
	setupFakeTerraform(t)
	argsFile := filepath.Join(t.TempDir(), "terraform-args.txt")
	t.Setenv("TF_HELPER_ARGS_FILE", argsFile)

	workDir := t.TempDir()
	if err := ApplyPlan(context.Background(), workDir, "saved.tfplan"); err != nil {
		t.Fatalf("ApplyPlan returned error: %v", err)
	}

	data, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("read args file: %v", err)
	}
	got := string(data)
	if !strings.Contains(got, "apply") || !strings.Contains(got, "-input=false") || !strings.Contains(got, "saved.tfplan") {
		t.Fatalf("unexpected apply args: %q", got)
	}
}

func setupFakeTerraform(t *testing.T) {
	t.Helper()

	if runtime.GOOS == "windows" {
		t.Skip("helper-based terraform command tests are skipped on Windows in this environment")
	}

	tempDir := t.TempDir()
	execPath, err := os.Executable()
	if err != nil {
		t.Fatalf("get executable: %v", err)
	}

	helperBinary := filepath.Join(tempDir, helperBinaryName())
	if runtime.GOOS == "windows" {
		if err := copyTestBinary(execPath, helperBinary); err != nil {
			t.Fatalf("copy helper binary: %v", err)
		}
		helperScript := filepath.Join(tempDir, "terraform.bat")
		content := fmt.Sprintf("@echo off\r\nset GO_WANT_HELPER_PROCESS=1\r\n\"%s\" -test.run=TestTerraformHelperProcess -- %%*\r\n", helperBinary)
		if err := os.WriteFile(helperScript, []byte(content), 0o755); err != nil {
			t.Fatalf("write helper script: %v", err)
		}
	} else {
		if err := copyTestBinary(execPath, helperBinary); err != nil {
			t.Fatalf("copy helper binary: %v", err)
		}
		helperScript := filepath.Join(tempDir, "terraform")
		content := fmt.Sprintf("#!/usr/bin/env bash\nexport GO_WANT_HELPER_PROCESS=1\nexec \"%s\" -test.run=TestTerraformHelperProcess -- \"$@\"\n", helperBinary)
		if err := os.WriteFile(helperScript, []byte(content), 0o755); err != nil {
			t.Fatalf("write helper script: %v", err)
		}
	}

	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func helperBinaryName() string {
	if runtime.GOOS == "windows" {
		return "terraform-helper.exe"
	}
	return "terraform-helper"
}

func copyTestBinary(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Close()
}

func TestTerraformHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for i, arg := range os.Args {
		if arg == "--" {
			args = os.Args[i+1:]
			break
		}
	}

	if argsFile := os.Getenv("TF_HELPER_ARGS_FILE"); argsFile != "" {
		_ = os.WriteFile(argsFile, []byte(strings.Join(args, " ")), 0o644)
	}

	if len(args) > 0 {
		switch args[0] {
		case "fmt":
			exitWithCode(envInt("TF_HELPER_FMT_EXIT"))
		case "plan":
			exitWithCode(envInt("TF_HELPER_PLAN_EXIT"))
		case "apply":
			exitWithCode(envInt("TF_HELPER_APPLY_EXIT"))
		case "init":
			exitWithCode(envInt("TF_HELPER_INIT_EXIT"))
		}
	}

	exitWithCode(envInt("TF_HELPER_EXIT_CODE"))
}

func envInt(name string) int {
	value := os.Getenv(name)
	if value == "" {
		return 0
	}
	switch value {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	}
	return 0
}

func exitWithCode(code int) {
	os.Exit(code)
}
