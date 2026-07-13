package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
)

func TestPlanJSONReturnsExitCodeTwoOnChanges(t *testing.T) {
	setupPlanFakeTerraform(t)

	tmpDir := t.TempDir()
	var out bytes.Buffer

	oldPathFlag := pathFlag
	oldPlanOut := planOut
	oldPlanDetailedExitCode := planDetailedExitCode
	oldPlanJSON := planJSON
	defer func() {
		pathFlag = oldPathFlag
		planOut = oldPlanOut
		planDetailedExitCode = oldPlanDetailedExitCode
		planJSON = oldPlanJSON
	}()

	pathFlag = tmpDir
	planOut = ""
	planDetailedExitCode = true
	planJSON = true

	testCmd := &cobra.Command{}
	testCmd.SetContext(context.Background())
	testCmd.SetOut(&out)

	err := planCmd.RunE(testCmd, nil)
	if err == nil {
		t.Fatal("expected exit code error")
	}

	exitErr, ok := err.(ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.ExitCode() != 2 {
		t.Fatalf("ExitCode() = %d, want 2", exitErr.ExitCode())
	}

	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal JSON: %v", err)
	}
	if got := payload["exit_code"]; got != float64(2) {
		t.Fatalf("exit_code = %v, want 2", got)
	}
	if got := payload["has_changes"]; got != true {
		t.Fatalf("has_changes = %v, want true", got)
	}
}

func setupPlanFakeTerraform(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()
	execPath, err := os.Executable()
	if err != nil {
		t.Fatalf("get executable: %v", err)
	}

	helperBinary := filepath.Join(tempDir, helperPlanBinaryName())
	if err := copyPlanTestBinary(execPath, helperBinary); err != nil {
		t.Fatalf("copy helper binary: %v", err)
	}

	if runtime.GOOS == "windows" {
		helperScript := filepath.Join(tempDir, "terraform.bat")
		content := "@echo off\r\nset GO_WANT_PLAN_HELPER_PROCESS=1\r\n\"" + helperBinary + "\" -test.run=TestPlanHelperProcess -- %*\r\n"
		if err := os.WriteFile(helperScript, []byte(content), 0o755); err != nil {
			t.Fatalf("write helper script: %v", err)
		}
	} else {
		helperScript := filepath.Join(tempDir, "terraform")
		content := "#!/usr/bin/env bash\nexport GO_WANT_PLAN_HELPER_PROCESS=1\nexec \"" + helperBinary + "\" -test.run=TestPlanHelperProcess -- \"$@\"\n"
		if err := os.WriteFile(helperScript, []byte(content), 0o755); err != nil {
			t.Fatalf("write helper script: %v", err)
		}
	}

	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func helperPlanBinaryName() string {
	if runtime.GOOS == "windows" {
		return "terraform-plan-helper.exe"
	}
	return "terraform-plan-helper"
}

func copyPlanTestBinary(src string, dst string) error {
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

func TestPlanHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_PLAN_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	for i, arg := range os.Args {
		if arg == "--" {
			args = os.Args[i+1:]
			break
		}
	}

	switch {
	case len(args) > 0 && args[0] == "init":
		os.Exit(0)
	case len(args) > 0 && args[0] == "plan":
		os.Exit(2)
	default:
		os.Exit(0)
	}
}
