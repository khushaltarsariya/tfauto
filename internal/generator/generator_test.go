package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCopyTemplateCopiesEmbeddedTemplate(t *testing.T) {
	t.Parallel()

	targetDir := filepath.Join(t.TempDir(), "out")
	if err := CopyTemplate("aws-basic", targetDir); err != nil {
		t.Fatalf("CopyTemplate returned error: %v", err)
	}

	for _, path := range []string{
		filepath.Join(targetDir, "main.tf"),
		filepath.Join(targetDir, "variables.tf"),
		filepath.Join(targetDir, "outputs.tf"),
		filepath.Join(targetDir, "provider.tf"),
		filepath.Join(targetDir, "DESCRIPTION.md"),
		filepath.Join(targetDir, "template.json"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected copied template file %s: %v", path, err)
		}
	}
}

func TestCopyTemplateMissingTemplate(t *testing.T) {
	t.Parallel()

	err := CopyTemplate("missing-template", filepath.Join(t.TempDir(), "out"))
	if err == nil {
		t.Fatal("CopyTemplate returned nil error for a missing template")
	}

	if !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("CopyTemplate error = %q, want message about missing template", err)
	}
}
