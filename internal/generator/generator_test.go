package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyFile(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	src := filepath.Join(tempDir, "source.txt")
	dst := filepath.Join(tempDir, "nested", "dest.txt")

	if err := os.WriteFile(src, []byte("tfauto"), 0o644); err != nil {
		t.Fatalf("write source file: %v", err)
	}

	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile returned error: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read copied file: %v", err)
	}

	if got := string(data); got != "tfauto" {
		t.Fatalf("copied file contents = %q, want %q", got, "tfauto")
	}
}

func TestCopyDir(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	dstDir := filepath.Join(tempDir, "dst")

	if err := os.MkdirAll(filepath.Join(srcDir, "nested"), 0o755); err != nil {
		t.Fatalf("create source directory: %v", err)
	}

	files := map[string]string{
		filepath.Join(srcDir, "main.tf"):           `resource "null_resource" "example" {}`,
		filepath.Join(srcDir, "nested", "vars.tf"): `variable "name" { type = string }`,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write source file %s: %v", path, err)
		}
	}

	if err := copyDir(srcDir, dstDir); err != nil {
		t.Fatalf("copyDir returned error: %v", err)
	}

	for originalPath, wantContent := range files {
		relPath, err := filepath.Rel(srcDir, originalPath)
		if err != nil {
			t.Fatalf("relative path: %v", err)
		}

		copiedPath := filepath.Join(dstDir, relPath)
		data, err := os.ReadFile(copiedPath)
		if err != nil {
			t.Fatalf("read copied file %s: %v", copiedPath, err)
		}

		if got := string(data); got != wantContent {
			t.Fatalf("copied content for %s = %q, want %q", copiedPath, got, wantContent)
		}
	}
}

func TestCopyTemplateMissingTemplate(t *testing.T) {
	tempDir := t.TempDir()
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}

	err = CopyTemplate("missing-template", filepath.Join(tempDir, "out"))
	if err == nil {
		t.Fatal("CopyTemplate returned nil error for a missing template")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("CopyTemplate error = %q, want message about missing template", err)
	}
}
