package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), FileName)
	content := `
project: demo
environment: prod

terraform:
  require_plan_file: true
  protect_destroy: true

templates:
  allowed:
    - aws-three-tier-vpc
    - aws-rds-postgres

policy:
  require_tags:
    - Owner
    - Environment
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile returned error: %v", err)
	}

	if cfg.Project != "demo" || cfg.Environment != "prod" {
		t.Fatalf("unexpected project/environment: %#v", cfg)
	}
	if !cfg.Terraform.RequirePlanFile || !cfg.Terraform.ProtectDestroy {
		t.Fatalf("terraform policy was not parsed: %#v", cfg.Terraform)
	}
	if got := len(cfg.Templates.Allowed); got != 2 {
		t.Fatalf("allowed templates len = %d, want 2", got)
	}
	if got := len(cfg.Policy.RequireTags); got != 2 {
		t.Fatalf("required tags len = %d, want 2", got)
	}
}

func TestLoadForPathSearchesParents(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	child := filepath.Join(root, "project", "nested")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("create child dir: %v", err)
	}

	configPath := filepath.Join(root, FileName)
	if err := os.WriteFile(configPath, []byte("project: demo\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	result, err := LoadForPath(child)
	if err != nil {
		t.Fatalf("LoadForPath returned error: %v", err)
	}
	if !result.Found {
		t.Fatal("expected config to be found")
	}
	if result.Path != configPath {
		t.Fatalf("config path = %q, want %q", result.Path, configPath)
	}
}

func TestLoadForPathUsesNearestExistingParent(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	configPath := filepath.Join(root, FileName)
	if err := os.WriteFile(configPath, []byte("project: demo\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	missingChild := filepath.Join(root, "missing", "project")
	result, err := LoadForPath(missingChild)
	if err != nil {
		t.Fatalf("LoadForPath returned error: %v", err)
	}
	if !result.Found {
		t.Fatal("expected config to be found")
	}
	if result.Path != configPath {
		t.Fatalf("config path = %q, want %q", result.Path, configPath)
	}
}
