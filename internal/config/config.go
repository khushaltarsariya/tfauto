package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const FileName = ".tfauto.yaml"

type Config struct {
	Project     string
	Environment string
	Terraform   TerraformConfig
	Templates   TemplatesConfig
	Policy      PolicyConfig
}

type TerraformConfig struct {
	RequirePlanFile bool
	ProtectDestroy  bool
}

type TemplatesConfig struct {
	Allowed []string
}

type PolicyConfig struct {
	RequireTags []string
}

type LoadResult struct {
	Config Config
	Path   string
	Found  bool
}

func LoadForPath(startPath string) (LoadResult, error) {
	if startPath == "" {
		startPath = "."
	}

	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return LoadResult{}, fmt.Errorf("resolve path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil && os.IsNotExist(err) {
		absPath = nearestExistingParent(absPath)
		info, err = os.Stat(absPath)
	}
	if err != nil {
		return LoadResult{}, fmt.Errorf("stat path: %w", err)
	}
	if !info.IsDir() {
		absPath = filepath.Dir(absPath)
	}

	configPath, found := findConfig(absPath)
	if !found {
		return LoadResult{}, nil
	}

	cfg, err := ParseFile(configPath)
	if err != nil {
		return LoadResult{}, err
	}

	return LoadResult{
		Config: cfg,
		Path:   configPath,
		Found:  true,
	}, nil
}

func ParseFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open %s: %w", FileName, err)
	}
	defer file.Close()

	var cfg Config
	var section string
	var listKey string

	scanner := bufio.NewScanner(file)
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
			switch section + "." + listKey {
			case "templates.allowed":
				cfg.Templates.Allowed = append(cfg.Templates.Allowed, value)
			case "policy.require_tags":
				cfg.Policy.RequireTags = append(cfg.Policy.RequireTags, value)
			default:
				return Config{}, fmt.Errorf("%s:%d unsupported list item", FileName, lineNumber)
			}
			continue
		}

		key, value, ok := splitKeyValue(trimmed)
		if !ok {
			return Config{}, fmt.Errorf("%s:%d invalid line %q", FileName, lineNumber, trimmed)
		}

		if indent == 0 {
			listKey = ""
			switch key {
			case "project":
				cfg.Project = value
			case "environment":
				cfg.Environment = value
			case "terraform", "templates", "policy":
				if value != "" {
					return Config{}, fmt.Errorf("%s:%d section %q cannot have inline value", FileName, lineNumber, key)
				}
				section = key
			default:
				return Config{}, fmt.Errorf("%s:%d unsupported key %q", FileName, lineNumber, key)
			}
			continue
		}

		if section == "" {
			return Config{}, fmt.Errorf("%s:%d nested key without section", FileName, lineNumber)
		}

		switch section {
		case "terraform":
			boolValue, err := parseBool(value)
			if err != nil {
				return Config{}, fmt.Errorf("%s:%d %w", FileName, lineNumber, err)
			}
			switch key {
			case "require_plan_file":
				cfg.Terraform.RequirePlanFile = boolValue
			case "protect_destroy":
				cfg.Terraform.ProtectDestroy = boolValue
			default:
				return Config{}, fmt.Errorf("%s:%d unsupported terraform key %q", FileName, lineNumber, key)
			}
		case "templates":
			if key != "allowed" || value != "" {
				return Config{}, fmt.Errorf("%s:%d unsupported templates key %q", FileName, lineNumber, key)
			}
			listKey = key
		case "policy":
			if key != "require_tags" || value != "" {
				return Config{}, fmt.Errorf("%s:%d unsupported policy key %q", FileName, lineNumber, key)
			}
			listKey = key
		}
	}

	if err := scanner.Err(); err != nil {
		return Config{}, fmt.Errorf("read %s: %w", FileName, err)
	}

	return cfg, nil
}

func findConfig(startDir string) (string, bool) {
	current := startDir
	for {
		candidate := filepath.Join(current, FileName)
		if fileExists(candidate) {
			return candidate, true
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
		current = parent
	}
}

func nearestExistingParent(path string) string {
	current := path
	for {
		parent := filepath.Dir(current)
		if parent == current {
			return path
		}
		if info, err := os.Stat(parent); err == nil && info.IsDir() {
			return parent
		}
		current = parent
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

func parseBool(value string) (bool, error) {
	switch strings.ToLower(value) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("expected true or false, got %q", value)
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
