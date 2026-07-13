package tplfs

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

//go:embed *
var builtins embed.FS

type TemplateMetadata struct {
	Name                     string   `json:"name"`
	Description              string   `json:"description"`
	Version                  string   `json:"version,omitempty"`
	Author                   string   `json:"author,omitempty"`
	Category                 string   `json:"category,omitempty"`
	CloudProvider            string   `json:"cloud_provider,omitempty"`
	EstimatedMonthlyCost     string   `json:"estimated_monthly_cost,omitempty"`
	RequiredTerraformVersion string   `json:"required_terraform_version,omitempty"`
	RequiredProviders        []string `json:"required_providers,omitempty"`
	Tags                     []string `json:"tags,omitempty"`
	HasManifest              bool     `json:"has_manifest"`
	ManifestFile             string   `json:"manifest_file,omitempty"`
	MetadataSource           string   `json:"metadata_source,omitempty"`
}

func List() ([]string, error) {
	entries, err := fs.ReadDir(builtins, ".")
	if err != nil {
		return nil, fmt.Errorf("read embedded templates: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	sort.Strings(names)
	return names, nil
}

func ListMetadata() ([]TemplateMetadata, error) {
	names, err := List()
	if err != nil {
		return nil, err
	}

	metadata := make([]TemplateMetadata, 0, len(names))
	for _, name := range names {
		item, err := Metadata(name)
		if err != nil {
			return nil, err
		}
		metadata = append(metadata, item)
	}

	sort.SliceStable(metadata, func(i, j int) bool {
		return metadata[i].Name < metadata[j].Name
	})

	return metadata, nil
}

func Exists(name string) bool {
	info, err := fs.Stat(builtins, name)
	return err == nil && info.IsDir()
}

func Metadata(name string) (TemplateMetadata, error) {
	if !Exists(name) {
		return TemplateMetadata{}, fmt.Errorf("template %q does not exist", name)
	}

	if manifest, manifestFile, found, err := readManifest(name); err != nil {
		return TemplateMetadata{}, err
	} else if found {
		manifest = normalizeManifest(name, manifest)
		manifest.HasManifest = true
		manifest.ManifestFile = manifestFile
		manifest.MetadataSource = filepath.Base(manifestFile)
		if strings.TrimSpace(manifest.Description) == "" {
			if description, err := readDescription(name); err == nil {
				manifest.Description = description
			}
		}
		return manifest, nil
	}

	description, err := readDescription(name)
	if err != nil {
		return TemplateMetadata{}, err
	}

	legacy := inferLegacyMetadata(name, description)
	return legacy, nil
}

func Description(name string) (string, error) {
	metadata, err := Metadata(name)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(metadata.Description), nil
}

func Files(name string) ([]string, error) {
	if !Exists(name) {
		return nil, fmt.Errorf("template %q does not exist", name)
	}

	var files []string
	err := fs.WalkDir(builtins, name, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(name, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if isMetadataFile(rel) || rel == "DESCRIPTION.md" {
			return nil
		}

		files = append(files, rel)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk embedded template files: %w", err)
	}

	sort.Strings(files)
	return files, nil
}

func Copy(name string, targetDir string) error {
	if !Exists(name) {
		return fmt.Errorf("template %q does not exist", name)
	}

	return fs.WalkDir(builtins, name, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(name, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, rel)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}

		data, err := fs.ReadFile(builtins, path)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0o644)
	})
}

func readManifest(name string) (TemplateMetadata, string, bool, error) {
	for _, manifestFile := range []string{"template.json", "template.yaml", "template.yml"} {
		path := filepath.ToSlash(filepath.Join(name, manifestFile))
		data, err := fs.ReadFile(builtins, path)
		if err != nil {
			continue
		}

		var metadata TemplateMetadata
		switch filepath.Ext(manifestFile) {
		case ".json":
			if err := json.Unmarshal(data, &metadata); err != nil {
				return TemplateMetadata{}, "", false, fmt.Errorf("parse %s: %w", path, err)
			}
		default:
			parsed, err := parseYAMLManifest(data)
			if err != nil {
				return TemplateMetadata{}, "", false, fmt.Errorf("parse %s: %w", path, err)
			}
			metadata = parsed
		}

		return metadata, path, true, nil
	}

	return TemplateMetadata{}, "", false, nil
}

func parseYAMLManifest(data []byte) (TemplateMetadata, error) {
	var metadata TemplateMetadata
	var listKey string

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := stripComment(scanner.Text())
		if strings.TrimSpace(line) == "" {
			continue
		}

		indent := leadingSpaces(line)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") {
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			switch listKey {
			case "required_providers":
				metadata.RequiredProviders = append(metadata.RequiredProviders, value)
			case "tags":
				metadata.Tags = append(metadata.Tags, value)
			default:
				return TemplateMetadata{}, fmt.Errorf("template manifest line %d: unsupported list item", lineNumber)
			}
			continue
		}

		key, value, ok := splitKeyValue(trimmed)
		if !ok {
			return TemplateMetadata{}, fmt.Errorf("template manifest line %d: invalid line %q", lineNumber, trimmed)
		}

		if indent != 0 {
			return TemplateMetadata{}, fmt.Errorf("template manifest line %d: nested keys are not supported", lineNumber)
		}

		listKey = ""
		switch key {
		case "name":
			metadata.Name = value
		case "description":
			metadata.Description = value
		case "version":
			metadata.Version = value
		case "author":
			metadata.Author = value
		case "category":
			metadata.Category = value
		case "cloud_provider":
			metadata.CloudProvider = value
		case "estimated_monthly_cost":
			metadata.EstimatedMonthlyCost = value
		case "required_terraform_version":
			metadata.RequiredTerraformVersion = value
		case "required_providers", "tags":
			if value != "" {
				return TemplateMetadata{}, fmt.Errorf("template manifest line %d: %s must be a list", lineNumber, key)
			}
			listKey = key
		default:
			return TemplateMetadata{}, fmt.Errorf("template manifest line %d: unsupported key %q", lineNumber, key)
		}
	}

	if err := scanner.Err(); err != nil {
		return TemplateMetadata{}, fmt.Errorf("read template manifest: %w", err)
	}

	return metadata, nil
}

func normalizeManifest(name string, metadata TemplateMetadata) TemplateMetadata {
	metadata.Name = strings.TrimSpace(metadata.Name)
	if metadata.Name == "" {
		metadata.Name = name
	}

	if strings.TrimSpace(metadata.CloudProvider) == "" {
		metadata.CloudProvider = inferCloudProvider(name)
	}
	if strings.TrimSpace(metadata.EstimatedMonthlyCost) == "" {
		metadata.EstimatedMonthlyCost = inferEstimatedMonthlyCost(name)
	}
	if strings.TrimSpace(metadata.RequiredTerraformVersion) == "" {
		metadata.RequiredTerraformVersion = inferRequiredTerraformVersion(name)
	}
	if len(metadata.RequiredProviders) == 0 {
		metadata.RequiredProviders = inferRequiredProviders(name)
	}
	if len(metadata.Tags) == 0 {
		metadata.Tags = inferTags(name)
	}
	if strings.TrimSpace(metadata.Category) == "" {
		metadata.Category = inferCategory(name)
	}
	if strings.TrimSpace(metadata.Version) == "" {
		metadata.Version = "1.0.0"
	}
	if strings.TrimSpace(metadata.Author) == "" {
		metadata.Author = "tfauto"
	}

	return metadata
}

func inferLegacyMetadata(name, description string) TemplateMetadata {
	return TemplateMetadata{
		Name:                     name,
		Description:              description,
		Version:                  "legacy",
		Author:                   "tfauto",
		Category:                 inferCategory(name),
		CloudProvider:            inferCloudProvider(name),
		EstimatedMonthlyCost:     inferEstimatedMonthlyCost(name),
		RequiredTerraformVersion: inferRequiredTerraformVersion(name),
		RequiredProviders:        inferRequiredProviders(name),
		Tags:                     inferTags(name),
		MetadataSource:           "legacy",
	}
}

func readDescription(name string) (string, error) {
	data, err := fs.ReadFile(builtins, filepath.ToSlash(filepath.Join(name, "DESCRIPTION.md")))
	if err != nil {
		if os.IsNotExist(err) || strings.Contains(err.Error(), "file does not exist") {
			return "", nil
		}
		return "", fmt.Errorf("read embedded description: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

func inferCloudProvider(name string) string {
	if strings.HasPrefix(name, "aws-") || strings.Contains(name, "aws") {
		return "aws"
	}
	return "unknown"
}

func inferEstimatedMonthlyCost(name string) string {
	switch {
	case strings.Contains(name, "cloudfront") || strings.Contains(name, "s3-static-site"):
		return "$0-10/mo"
	case strings.Contains(name, "lambda"):
		return "$0-20/mo"
	case strings.Contains(name, "rds"):
		return "$50-150/mo"
	case strings.Contains(name, "ecs"):
		return "$25-100/mo"
	case strings.Contains(name, "three-tier"):
		return "$20-60/mo"
	case strings.Contains(name, "alb") || strings.Contains(name, "asg"):
		return "$20-80/mo"
	case strings.Contains(name, "vpc"):
		return "$0-5/mo"
	default:
		return "varies"
	}
}

func inferRequiredTerraformVersion(name string) string {
	content, err := fs.ReadFile(builtins, filepath.ToSlash(filepath.Join(name, "provider.tf")))
	if err != nil {
		return ""
	}

	matches := regexp.MustCompile(`required_version\s*=\s*"([^"]+)"`).FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func inferRequiredProviders(name string) []string {
	content, err := fs.ReadFile(builtins, filepath.ToSlash(filepath.Join(name, "provider.tf")))
	if err != nil {
		return nil
	}

	blockPattern := regexp.MustCompile(`(?s)([A-Za-z0-9_]+)\s*=\s*\{(.*?)\}`)
	sourcePattern := regexp.MustCompile(`source\s*=\s*"([^"]+)"`)
	matches := blockPattern.FindAllStringSubmatch(string(content), -1)
	if len(matches) == 0 {
		return nil
	}

	providers := make([]string, 0, len(matches))
	for _, match := range matches {
		source := sourcePattern.FindStringSubmatch(match[2])
		if len(source) > 1 {
			providers = append(providers, source[1])
			continue
		}
		providers = append(providers, match[1])
	}
	sort.Strings(providers)
	return providers
}

func inferTags(name string) []string {
	tags := []string{"tfauto", "terraform", inferCloudProvider(name)}
	for _, part := range strings.Split(name, "-") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		tags = append(tags, part)
	}

	seen := make(map[string]struct{}, len(tags))
	deduped := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		deduped = append(deduped, tag)
	}

	sort.Strings(deduped)
	return deduped
}

func inferCategory(name string) string {
	switch {
	case strings.Contains(name, "rds"):
		return "database"
	case strings.Contains(name, "ecs") || strings.Contains(name, "alb") || strings.Contains(name, "asg"):
		return "application"
	case strings.Contains(name, "lambda"):
		return "serverless"
	case strings.Contains(name, "cloudfront") || strings.Contains(name, "static-site"):
		return "frontend"
	case strings.Contains(name, "vpc"):
		return "networking"
	default:
		return "starter"
	}
}

func isMetadataFile(path string) bool {
	switch filepath.Base(path) {
	case "template.json", "template.yaml", "template.yml":
		return true
	default:
		return false
	}
}

func stripComment(line string) string {
	if index := strings.Index(line, "#"); index >= 0 {
		return line[:index]
	}
	return line
}

func leadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

func splitKeyValue(line string) (string, string, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	if key == "" {
		return "", "", false
	}

	return key, value, true
}
