package doctor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Status string

const (
	StatusPass Status = "PASS"
	StatusWarn Status = "WARN"
	StatusFail Status = "FAIL"
)

type Result struct {
	Name    string
	Status  Status
	Details []string
}

type Report struct {
	Path    string
	Results []Result
}

func (r Report) HasFailures() bool {
	for _, result := range r.Results {
		if result.Status == StatusFail {
			return true
		}
	}

	return false
}

func Run(ctx context.Context, path string) Report {
	if path == "" {
		path = "."
	}

	report := Report{Path: path}

	pathResult, absPath := checkPath(path)
	report.Results = append(report.Results, pathResult)
	if pathResult.Status == StatusFail {
		return report
	}

	report.Results = append(report.Results, checkTerraformBinary())
	report.Results = append(report.Results, checkTerraformVersion(ctx))
	report.Results = append(report.Results, checkTerraformFiles(absPath))
	report.Results = append(report.Results, checkTerraformInitialized(absPath))
	report.Results = append(report.Results, checkBackendConfig(absPath))
	report.Results = append(report.Results, checkWorkspace(absPath))
	report.Results = append(report.Results, checkVariablePromptRisk(absPath))
	report.Results = append(report.Results, checkAWSRegion(absPath))
	report.Results = append(report.Results, checkAWSConfig(ctx, absPath))

	return report
}

func checkPath(path string) (Result, string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Result{
			Name:    "Target path",
			Status:  StatusFail,
			Details: []string{fmt.Sprintf("Unable to resolve path %q: %v", path, err)},
		}, ""
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return Result{
			Name:    "Target path",
			Status:  StatusFail,
			Details: []string{fmt.Sprintf("Path does not exist: %s", absPath)},
		}, ""
	}

	if !info.IsDir() {
		return Result{
			Name:    "Target path",
			Status:  StatusFail,
			Details: []string{fmt.Sprintf("Path is not a directory: %s", absPath)},
		}, ""
	}

	return Result{
		Name:    "Target path",
		Status:  StatusPass,
		Details: []string{absPath},
	}, absPath
}

func checkTerraformBinary() Result {
	path, err := exec.LookPath("terraform")
	if err != nil {
		return Result{
			Name:    "Terraform binary",
			Status:  StatusFail,
			Details: []string{"terraform was not found in PATH"},
		}
	}

	return Result{
		Name:    "Terraform binary",
		Status:  StatusPass,
		Details: []string{fmt.Sprintf("Found at %s", path)},
	}
}

func checkTerraformVersion(ctx context.Context) Result {
	out, err := exec.CommandContext(ctx, "terraform", "version", "-json").Output()
	if err != nil {
		return Result{
			Name:    "Terraform version",
			Status:  StatusWarn,
			Details: []string{"Unable to read terraform version in JSON mode"},
		}
	}

	var payload struct {
		TerraformVersion string `json:"terraform_version"`
	}
	if err := json.Unmarshal(out, &payload); err != nil {
		return Result{
			Name:    "Terraform version",
			Status:  StatusWarn,
			Details: []string{"Terraform version output could not be parsed"},
		}
	}

	return Result{
		Name:    "Terraform version",
		Status:  StatusPass,
		Details: []string{payload.TerraformVersion},
	}
}

func checkTerraformFiles(path string) Result {
	var files []string

	err := filepath.WalkDir(path, func(current string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".terraform" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(current) == ".tf" {
			rel, relErr := filepath.Rel(path, current)
			if relErr != nil {
				rel = current
			}
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		return Result{
			Name:    "Terraform files",
			Status:  StatusFail,
			Details: []string{fmt.Sprintf("Unable to scan path: %v", err)},
		}
	}

	if len(files) == 0 {
		return Result{
			Name:    "Terraform files",
			Status:  StatusFail,
			Details: []string{"No .tf files found in the target path"},
		}
	}

	details := []string{fmt.Sprintf("Found %d Terraform file(s)", len(files))}
	if len(files) <= 5 {
		details = append(details, files...)
	}

	return Result{
		Name:    "Terraform files",
		Status:  StatusPass,
		Details: details,
	}
}

func checkTerraformInitialized(path string) Result {
	lockFile := filepath.Join(path, ".terraform.lock.hcl")
	pluginDir := filepath.Join(path, ".terraform")

	hasLock := fileExists(lockFile)
	hasTerraformDir := dirExists(pluginDir)

	switch {
	case hasLock && hasTerraformDir:
		return Result{
			Name:    "Terraform initialization",
			Status:  StatusPass,
			Details: []string{"Project appears to be initialized"},
		}
	case hasLock || hasTerraformDir:
		return Result{
			Name:   "Terraform initialization",
			Status: StatusWarn,
			Details: []string{
				"Project looks partially initialized",
				"Run `tfauto validate --path <dir>` or `terraform init` to refresh initialization state",
			},
		}
	default:
		return Result{
			Name:   "Terraform initialization",
			Status: StatusWarn,
			Details: []string{
				"Project is not initialized yet",
				"Run `tfauto validate --path <dir>` or `terraform init` before plan/apply",
			},
		}
	}
}

func checkBackendConfig(path string) Result {
	files, err := terraformFiles(path)
	if err != nil {
		return Result{
			Name:    "Terraform backend configuration",
			Status:  StatusWarn,
			Details: []string{fmt.Sprintf("Unable to inspect backend configuration: %v", err)},
		}
	}

	for _, file := range files {
		content, readErr := os.ReadFile(file)
		if readErr != nil {
			continue
		}
		if strings.Contains(string(content), "backend \"") {
			rel := relativeOrAbsolute(path, file)
			return Result{
				Name:   "Terraform backend configuration",
				Status: StatusPass,
				Details: []string{
					fmt.Sprintf("Backend block detected in %s", rel),
					"Ensure backend credentials and remote state access are configured",
				},
			}
		}
	}

	return Result{
		Name:   "Terraform backend configuration",
		Status: StatusWarn,
		Details: []string{
			"No backend block detected",
			"Terraform will use local state unless you configure a backend",
		},
	}
}

func checkWorkspace(path string) Result {
	if workspace := os.Getenv("TF_WORKSPACE"); workspace != "" {
		return Result{
			Name:    "Terraform workspace",
			Status:  StatusPass,
			Details: []string{fmt.Sprintf("TF_WORKSPACE=%s", workspace)},
		}
	}

	workspaceFile := filepath.Join(path, ".terraform", "environment")
	if data, err := os.ReadFile(workspaceFile); err == nil {
		workspace := strings.TrimSpace(string(data))
		if workspace == "" {
			workspace = "default"
		}
		return Result{
			Name:    "Terraform workspace",
			Status:  StatusPass,
			Details: []string{fmt.Sprintf("Current workspace: %s", workspace)},
		}
	}

	return Result{
		Name:   "Terraform workspace",
		Status: StatusWarn,
		Details: []string{
			"Workspace could not be determined",
			"Default workspace is likely in use unless TF_WORKSPACE or .terraform/environment says otherwise",
		},
	}
}

func checkVariablePromptRisk(path string) Result {
	requiredVariables, err := requiredVariablesWithoutDefaults(path)
	if err != nil {
		return Result{
			Name:    "Variable prompt risk",
			Status:  StatusWarn,
			Details: []string{fmt.Sprintf("Unable to inspect variables: %v", err)},
		}
	}

	tfvarsFiles := matchingFiles(path, func(name string) bool {
		return name == "terraform.tfvars" ||
			name == "terraform.tfvars.json" ||
			strings.HasSuffix(name, ".auto.tfvars") ||
			strings.HasSuffix(name, ".auto.tfvars.json")
	})

	if len(requiredVariables) == 0 {
		return Result{
			Name:    "Variable prompt risk",
			Status:  StatusPass,
			Details: []string{"No required variables without defaults detected"},
		}
	}

	details := []string{
		fmt.Sprintf("Required variables without defaults: %s", strings.Join(requiredVariables, ", ")),
	}
	if len(tfvarsFiles) == 0 {
		details = append(details, "No terraform.tfvars or .auto.tfvars files detected")
		return Result{
			Name:    "Variable prompt risk",
			Status:  StatusWarn,
			Details: details,
		}
	}

	relFiles := make([]string, 0, len(tfvarsFiles))
	for _, file := range tfvarsFiles {
		relFiles = append(relFiles, relativeOrAbsolute(path, file))
	}
	details = append(details, fmt.Sprintf("Variable files detected: %s", strings.Join(relFiles, ", ")))

	return Result{
		Name:    "Variable prompt risk",
		Status:  StatusWarn,
		Details: details,
	}
}

func checkAWSRegion(path string) Result {
	if region := os.Getenv("AWS_REGION"); region != "" {
		return Result{
			Name:    "AWS region resolution",
			Status:  StatusPass,
			Details: []string{fmt.Sprintf("Using AWS_REGION=%s", region)},
		}
	}
	if region := os.Getenv("AWS_DEFAULT_REGION"); region != "" {
		return Result{
			Name:    "AWS region resolution",
			Status:  StatusPass,
			Details: []string{fmt.Sprintf("Using AWS_DEFAULT_REGION=%s", region)},
		}
	}

	files, err := terraformFiles(path)
	if err != nil {
		return Result{
			Name:    "AWS region resolution",
			Status:  StatusWarn,
			Details: []string{fmt.Sprintf("Unable to inspect provider configuration: %v", err)},
		}
	}

	for _, file := range files {
		content, readErr := os.ReadFile(file)
		if readErr != nil {
			continue
		}
		if strings.Contains(string(content), "provider \"aws\"") && strings.Contains(string(content), "region") {
			return Result{
				Name:   "AWS region resolution",
				Status: StatusPass,
				Details: []string{
					fmt.Sprintf("AWS region appears to be configured in %s", relativeOrAbsolute(path, file)),
				},
			}
		}
	}

	return Result{
		Name:   "AWS region resolution",
		Status: StatusWarn,
		Details: []string{
			"No AWS region detected in environment or provider configuration",
			"Set AWS_REGION or configure region in the provider block",
		},
	}
}

func checkAWSConfig(ctx context.Context, path string) Result {
	awsPath, err := exec.LookPath("aws")
	if err != nil {
		return Result{
			Name:   "AWS CLI and identity",
			Status: StatusWarn,
			Details: []string{
				"AWS CLI was not found in PATH",
				"Install AWS CLI if you want identity verification from tfauto doctor",
			},
		}
	}

	envDetails := awsEnvironmentDetails(path)
	cmd := exec.CommandContext(ctx, "aws", "sts", "get-caller-identity", "--output", "json")
	out, err := cmd.Output()
	if err != nil {
		details := []string{fmt.Sprintf("AWS CLI found at %s", awsPath)}
		details = append(details, envDetails...)
		details = append(details, "Unable to resolve AWS caller identity")
		return Result{
			Name:    "AWS CLI and identity",
			Status:  StatusWarn,
			Details: details,
		}
	}

	var payload struct {
		Account string `json:"Account"`
		Arn     string `json:"Arn"`
	}
	if err := json.Unmarshal(out, &payload); err != nil {
		details := []string{fmt.Sprintf("AWS CLI found at %s", awsPath)}
		details = append(details, envDetails...)
		details = append(details, "Caller identity response could not be parsed")
		return Result{
			Name:    "AWS CLI and identity",
			Status:  StatusWarn,
			Details: details,
		}
	}

	details := []string{fmt.Sprintf("AWS CLI found at %s", awsPath)}
	details = append(details, envDetails...)
	if payload.Account != "" {
		details = append(details, fmt.Sprintf("Account: %s", payload.Account))
	}
	if payload.Arn != "" {
		details = append(details, fmt.Sprintf("ARN: %s", payload.Arn))
	}

	return Result{
		Name:    "AWS CLI and identity",
		Status:  StatusPass,
		Details: details,
	}
}

func awsEnvironmentDetails(path string) []string {
	var details []string

	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		details = append(details, fmt.Sprintf("AWS_PROFILE=%s", profile))
	}
	if region := os.Getenv("AWS_REGION"); region != "" {
		details = append(details, fmt.Sprintf("AWS_REGION=%s", region))
	} else if region := os.Getenv("AWS_DEFAULT_REGION"); region != "" {
		details = append(details, fmt.Sprintf("AWS_DEFAULT_REGION=%s", region))
	}

	sharedConfig := filepath.Join(os.Getenv("USERPROFILE"), ".aws", "config")
	sharedCreds := filepath.Join(os.Getenv("USERPROFILE"), ".aws", "credentials")
	if home, err := os.UserHomeDir(); err == nil {
		sharedConfig = filepath.Join(home, ".aws", "config")
		sharedCreds = filepath.Join(home, ".aws", "credentials")
	}

	if fileExists(sharedConfig) {
		details = append(details, fmt.Sprintf("Shared config: %s", sharedConfig))
	}
	if fileExists(sharedCreds) {
		details = append(details, fmt.Sprintf("Shared credentials: %s", sharedCreds))
	}

	if len(details) == 0 {
		details = append(details, fmt.Sprintf("No AWS environment hints detected for %s", strings.TrimSpace(path)))
	}

	return details
}

func terraformFiles(path string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(path, func(current string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".terraform" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(current) == ".tf" {
			files = append(files, current)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func matchingFiles(path string, match func(name string) bool) []string {
	var files []string

	_ = filepath.WalkDir(path, func(current string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".terraform" {
				return filepath.SkipDir
			}
			return nil
		}
		if match(filepath.Base(current)) {
			files = append(files, current)
		}
		return nil
	})

	return files
}

func requiredVariablesWithoutDefaults(path string) ([]string, error) {
	files, err := terraformFiles(path)
	if err != nil {
		return nil, err
	}

	blockPattern := regexp.MustCompile(`(?s)variable\s+"([^"]+)"\s*\{(.*?)\}`)
	defaultPattern := regexp.MustCompile(`(?m)^\s*default\s*=`)
	seen := map[string]struct{}{}
	var required []string

	for _, file := range files {
		content, readErr := os.ReadFile(file)
		if readErr != nil {
			continue
		}

		matches := blockPattern.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			name := match[1]
			body := match[2]
			if defaultPattern.MatchString(body) {
				continue
			}
			if _, exists := seen[name]; exists {
				continue
			}
			seen[name] = struct{}{}
			required = append(required, name)
		}
	}

	return required, nil
}

func relativeOrAbsolute(base string, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return rel
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
