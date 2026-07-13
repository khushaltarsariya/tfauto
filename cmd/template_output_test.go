package cmd

import (
	"bytes"
	"strings"
	"testing"

	tplfs "github.com/khushaltarsariya/tfauto/templates"
)

func TestRenderTemplatesOutput(t *testing.T) {
	metadata := []tplfs.TemplateMetadata{
		{
			Name:                     "aws-basic",
			Version:                  "1.0.0",
			CloudProvider:            "aws",
			RequiredTerraformVersion: ">= 1.7.0",
			Category:                 "starter",
			EstimatedMonthlyCost:     "$0-5/mo",
			Description:              "Simple AWS EC2 infrastructure.",
		},
	}

	var buf bytes.Buffer
	renderTemplates(&buf, metadata)

	out := buf.String()
	for _, want := range []string{
		"tfauto: templates",
		"NAME",
		"aws-basic",
		"$0-5/mo",
		"Next steps:",
		"tfauto init --template <name> --target ./my-project",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}

func TestRenderTemplateOutput(t *testing.T) {
	metadata := tplfs.TemplateMetadata{
		Name:                     "aws-basic",
		Version:                  "1.0.0",
		Author:                   "tfauto",
		Category:                 "starter",
		CloudProvider:            "aws",
		EstimatedMonthlyCost:     "$0-5/mo",
		RequiredTerraformVersion: ">= 1.7.0",
		RequiredProviders:        []string{"hashicorp/aws"},
		Tags:                     []string{"aws", "starter"},
		MetadataSource:           "template.json",
		Description:              "Simple AWS EC2 infrastructure.",
	}

	var buf bytes.Buffer
	renderTemplate(&buf, metadata, []string{"main.tf", "variables.tf"}, "aws-basic")

	out := buf.String()
	for _, want := range []string{
		"tfauto: template aws-basic",
		"Metadata:",
		"Required providers: hashicorp/aws",
		"Files:",
		"main.tf",
		"Next step:",
		"tfauto init --template aws-basic --target ./my-project",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q:\n%s", want, out)
		}
	}
}
