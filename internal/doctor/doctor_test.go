package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReportHasFailures(t *testing.T) {
	t.Parallel()

	report := Report{
		Results: []Result{
			{Name: "ok", Status: StatusPass},
			{Name: "warn", Status: StatusWarn},
		},
	}

	if report.HasFailures() {
		t.Fatal("HasFailures returned true for report without failures")
	}

	report.Results = append(report.Results, Result{Name: "fail", Status: StatusFail})
	if !report.HasFailures() {
		t.Fatal("HasFailures returned false for report with a failure")
	}
}

func TestReportHasWarnings(t *testing.T) {
	t.Parallel()

	report := Report{
		Results: []Result{
			{Name: "warn", Status: StatusWarn},
		},
	}

	if !report.HasWarnings() {
		t.Fatal("HasWarnings returned false for report with a warning")
	}
}

func TestCheckPath(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	result, absPath := checkPath(tempDir)
	if result.Status != StatusPass {
		t.Fatalf("checkPath status = %s, want %s", result.Status, StatusPass)
	}
	if absPath == "" {
		t.Fatal("checkPath returned empty absolute path")
	}

	filePath := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	result, _ = checkPath(filePath)
	if result.Status != StatusFail {
		t.Fatalf("checkPath on file status = %s, want %s", result.Status, StatusFail)
	}
}

func TestCheckTerraformFiles(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	result := checkTerraformFiles(tempDir)
	if result.Status != StatusFail {
		t.Fatalf("checkTerraformFiles status = %s, want %s for empty dir", result.Status, StatusFail)
	}

	tfPath := filepath.Join(tempDir, "main.tf")
	if err := os.WriteFile(tfPath, []byte(`terraform {}`), 0o644); err != nil {
		t.Fatalf("write terraform file: %v", err)
	}

	result = checkTerraformFiles(tempDir)
	if result.Status != StatusPass {
		t.Fatalf("checkTerraformFiles status = %s, want %s", result.Status, StatusPass)
	}

	found := false
	for _, detail := range result.Details {
		if strings.Contains(detail, "main.tf") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected details to mention main.tf, got %v", result.Details)
	}
}

func TestCheckTerraformInitialized(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	result := checkTerraformInitialized(tempDir)
	if result.Status != StatusWarn {
		t.Fatalf("status without init = %s, want %s", result.Status, StatusWarn)
	}

	lockPath := filepath.Join(tempDir, ".terraform.lock.hcl")
	if err := os.WriteFile(lockPath, []byte("provider"), 0o644); err != nil {
		t.Fatalf("write lock file: %v", err)
	}

	result = checkTerraformInitialized(tempDir)
	if result.Status != StatusWarn {
		t.Fatalf("status with partial init = %s, want %s", result.Status, StatusWarn)
	}

	if err := os.Mkdir(filepath.Join(tempDir, ".terraform"), 0o755); err != nil {
		t.Fatalf("create .terraform dir: %v", err)
	}

	result = checkTerraformInitialized(tempDir)
	if result.Status != StatusPass {
		t.Fatalf("status with full init = %s, want %s", result.Status, StatusPass)
	}
}

func TestAWSEnvironmentDetails(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("USERPROFILE", homeDir)
	t.Setenv("AWS_PROFILE", "test-profile")
	t.Setenv("AWS_REGION", "us-east-1")

	awsDir := filepath.Join(homeDir, ".aws")
	if err := os.MkdirAll(awsDir, 0o755); err != nil {
		t.Fatalf("create aws dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(awsDir, "config"), []byte("[default]"), 0o644); err != nil {
		t.Fatalf("write aws config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(awsDir, "credentials"), []byte("[default]"), 0o644); err != nil {
		t.Fatalf("write aws credentials: %v", err)
	}

	details := awsEnvironmentDetails(t.TempDir())
	joined := strings.Join(details, "\n")

	for _, expected := range []string{
		"AWS_PROFILE=test-profile",
		"AWS_REGION=us-east-1",
		"Shared config:",
		"Shared credentials:",
	} {
		if !strings.Contains(joined, expected) {
			t.Fatalf("awsEnvironmentDetails missing %q in %q", expected, joined)
		}
	}
}

func TestCheckWritePermissions(t *testing.T) {
	t.Parallel()

	result := checkWritePermissions(t.TempDir())
	if result.Status != StatusPass {
		t.Fatalf("checkWritePermissions status = %s, want %s", result.Status, StatusPass)
	}
}

func TestCheckCurrentDirectory(t *testing.T) {
	t.Parallel()

	result := checkCurrentDirectory(".")
	if result.Status != StatusPass {
		t.Fatalf("checkCurrentDirectory status = %s, want %s", result.Status, StatusPass)
	}
}
