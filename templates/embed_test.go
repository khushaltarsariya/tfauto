package tplfs

import (
	"strings"
	"testing"
)

func TestMetadataReadsManifest(t *testing.T) {
	t.Parallel()

	meta, err := Metadata("aws-lambda-apigateway")
	if err != nil {
		t.Fatalf("Metadata returned error: %v", err)
	}

	if !meta.HasManifest {
		t.Fatal("expected manifest-backed metadata")
	}
	if meta.ManifestFile != "aws-lambda-apigateway/template.json" {
		t.Fatalf("ManifestFile = %q, want template.json path", meta.ManifestFile)
	}
	if meta.Name != "aws-lambda-apigateway" {
		t.Fatalf("Name = %q, want aws-lambda-apigateway", meta.Name)
	}
	if meta.CloudProvider != "aws" {
		t.Fatalf("CloudProvider = %q, want aws", meta.CloudProvider)
	}
	if len(meta.RequiredProviders) != 2 {
		t.Fatalf("RequiredProviders = %v, want 2 providers", meta.RequiredProviders)
	}
}

func TestMetadataReadsStartupTemplate(t *testing.T) {
	t.Parallel()

	meta, err := Metadata("aws-startup-webapp")
	if err != nil {
		t.Fatalf("Metadata returned error: %v", err)
	}

	if !meta.HasManifest {
		t.Fatal("expected manifest-backed metadata")
	}
	if meta.Category != "application" {
		t.Fatalf("Category = %q, want application", meta.Category)
	}
	if meta.CloudProvider != "aws" {
		t.Fatalf("CloudProvider = %q, want aws", meta.CloudProvider)
	}
	if len(meta.Tags) == 0 {
		t.Fatal("expected tags to be populated")
	}
}

func TestMetadataReadsServerlessWebappTemplate(t *testing.T) {
	t.Parallel()

	meta, err := Metadata("aws-serverless-webapp")
	if err != nil {
		t.Fatalf("Metadata returned error: %v", err)
	}

	if meta.Category != "serverless" {
		t.Fatalf("Category = %q, want serverless", meta.Category)
	}
	if meta.CloudProvider != "aws" {
		t.Fatalf("CloudProvider = %q, want aws", meta.CloudProvider)
	}
	if !meta.HasManifest {
		t.Fatal("expected manifest-backed metadata")
	}
}

func TestFilesExcludeMetadataFiles(t *testing.T) {
	t.Parallel()

	files, err := Files("aws-basic")
	if err != nil {
		t.Fatalf("Files returned error: %v", err)
	}

	for _, file := range files {
		if strings.Contains(file, "template.json") || strings.Contains(file, "DESCRIPTION.md") {
			t.Fatalf("Files included metadata file %q", file)
		}
	}
}
