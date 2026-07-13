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

const (
	SectionFilesystem = "Filesystem"
	SectionTerraform  = "Terraform"
	SectionAWS        = "AWS"
	SectionGit        = "Git"
)

type Result struct {
	Section string   `json:"section"`
	Name    string   `json:"name"`
	Status  Status   `json:"status"`
	Details []string `json:"details"`
}

type Summary struct {
	Pass int `json:"pass"`
	Warn int `json:"warn"`
	Fail int `json:"fail"`
}

type Report struct {
	Path             string   `json:"path"`
	CurrentDirectory string   `json:"current_directory"`
	Results          []Result `json:"results"`
	Summary          Summary  `json:"summary"`
}

func (r Report) HasFailures() bool {
	if r.Summary.Fail > 0 {
		return true
	}
	for _, result := range r.Results {
		if result.Status == StatusFail {
			return true
		}
	}
	return false
}

func (r Report) HasWarnings() bool {
	if r.Summary.Warn > 0 {
		return true
	}
	for _, result := range r.Results {
		if result.Status == StatusWarn {
			return true
		}
	}
	return false
}

func (r Report) ExitCode() int {
	switch {
	case r.HasFailures():
		return 1
	case r.HasWarnings():
		return 2
	default:
		return 0
	}
}

func Run(ctx context.Context, path string) Report {
	if path == "" {
		path = "."
	}

	currentDir, _ := os.Getwd()
	report := Report{Path: path, CurrentDirectory: currentDir}

	pathResult, absPath := checkPath(path)
	report.add(pathResult)
	if pathResult.Status == StatusFail {
		return report
	}

	report.add(checkCurrentDirectory(currentDir))
	report.add(checkWritePermissions(absPath))
	report.add(checkTerraformBinary())
	report.add(checkTerraformVersion(ctx))
	report.add(checkTerraformFiles(absPath))
	report.add(checkTerraformFmt(ctx, absPath))
	report.add(checkTerraformValidate(ctx, absPath))
	report.add(checkTerraformInitialized(absPath))
	report.add(checkBackendConfig(absPath))
	report.add(checkProvidersInstalled(absPath))
	report.add(checkModulesDownloaded(absPath))
	report.add(checkAWSCredentials(ctx, absPath))
	report.add(checkAWSProfile())
	report.add(checkAWSCallerIdentity(ctx, absPath))
	report.add(checkAWSRegion(absPath))
	report.add(checkGitBinary())
	report.add(checkGitBranch(ctx, absPath))

	return report
}

func (r *Report) add(result Result) {
	r.Results = append(r.Results, result)
	switch result.Status {
	case StatusPass:
		r.Summary.Pass++
	case StatusWarn:
		r.Summary.Warn++
	case StatusFail:
		r.Summary.Fail++
	}
}

func newResult(section, name string, status Status, details ...string) Result {
	return Result{
		Section: section,
		Name:    name,
		Status:  status,
		Details: details,
	}
}

func countTerraformFiles(path string) int {
	count := 0
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
		if filepath.Ext(current) == ".tf" {
			count++
		}
		return nil
	})
	return count
}

func checkPath(path string) (Result, string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return newResult(SectionFilesystem, "current path", StatusFail, fmt.Sprintf("Unable to resolve path %q: %v", path, err)), ""
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return newResult(SectionFilesystem, "current path", StatusFail, fmt.Sprintf("Path does not exist: %s", absPath)), ""
	}

	if !info.IsDir() {
		return newResult(SectionFilesystem, "current path", StatusFail, fmt.Sprintf("Path is not a directory: %s", absPath)), ""
	}

	return newResult(SectionFilesystem, "current path", StatusPass, absPath), absPath
}

func checkTerraformBinary() Result {
	path, err := exec.LookPath("terraform")
	if err != nil {
		return newResult(SectionFilesystem, "terraform exists", StatusFail, "terraform was not found in PATH")
	}

	return newResult(SectionFilesystem, "terraform exists", StatusPass, fmt.Sprintf("Found at %s", path))
}

func checkTerraformVersion(ctx context.Context) Result {
	out, err := exec.CommandContext(ctx, "terraform", "version", "-json").Output()
	if err != nil {
		return newResult(SectionFilesystem, "terraform version", StatusWarn, "Unable to read terraform version in JSON mode")
	}

	var payload struct {
		TerraformVersion string `json:"terraform_version"`
	}
	if err := json.Unmarshal(out, &payload); err != nil {
		return newResult(SectionFilesystem, "terraform version", StatusWarn, "Terraform version output could not be parsed")
	}

	return newResult(SectionFilesystem, "terraform version", StatusPass, payload.TerraformVersion)
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
		return newResult(SectionFilesystem, "terraform files found", StatusFail, fmt.Sprintf("Unable to scan path: %v", err))
	}

	if len(files) == 0 {
		return newResult(SectionFilesystem, "terraform files found", StatusFail, "No .tf files found in the target path")
	}

	details := []string{fmt.Sprintf("Found %d Terraform file(s)", len(files))}
	if len(files) <= 5 {
		details = append(details, files...)
	}

	return newResult(SectionFilesystem, "terraform files found", StatusPass, details...)
}

func checkCurrentDirectory(currentDir string) Result {
	if currentDir == "" {
		return newResult(SectionFilesystem, "current directory", StatusWarn, "Unable to determine current working directory")
	}
	return newResult(SectionFilesystem, "current directory", StatusPass, currentDir)
}

func checkWritePermissions(path string) Result {
	file, err := os.CreateTemp(path, ".tfauto-write-check-*")
	if err != nil {
		return newResult(SectionFilesystem, "write permissions", StatusFail, fmt.Sprintf("Unable to write to %s: %v", path, err))
	}
	name := file.Name()
	if _, err := file.WriteString("tfauto"); err != nil {
		_ = file.Close()
		_ = os.Remove(name)
		return newResult(SectionFilesystem, "write permissions", StatusFail, fmt.Sprintf("Unable to write to %s: %v", path, err))
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(name)
		return newResult(SectionFilesystem, "write permissions", StatusFail, fmt.Sprintf("Unable to close write test file: %v", err))
	}
	_ = os.Remove(name)
	return newResult(SectionFilesystem, "write permissions", StatusPass, fmt.Sprintf("Writable: %s", path))
}

func checkTerraformInitialized(path string) Result {
	lockFile := filepath.Join(path, ".terraform.lock.hcl")
	pluginDir := filepath.Join(path, ".terraform")

	hasLock := fileExists(lockFile)
	hasTerraformDir := dirExists(pluginDir)

	switch {
	case hasLock && hasTerraformDir:
		return newResult(SectionTerraform, "project initialized", StatusPass, "Project appears to be initialized")
	case hasLock || hasTerraformDir:
		return newResult(SectionTerraform, "project initialized", StatusWarn, "Project looks partially initialized", "Run `tfauto validate --path <dir>` or `terraform init` to refresh initialization state")
	default:
		return newResult(SectionTerraform, "project initialized", StatusWarn, "Project is not initialized yet", "Run `tfauto validate --path <dir>` or `terraform init` before plan/apply")
	}
}

func checkBackendConfig(path string) Result {
	files, err := terraformFiles(path)
	if err != nil {
		return newResult(SectionTerraform, "backend configured", StatusWarn, fmt.Sprintf("Unable to inspect backend configuration: %v", err))
	}

	for _, file := range files {
		content, readErr := os.ReadFile(file)
		if readErr != nil {
			continue
		}
		if strings.Contains(string(content), "backend \"") {
			rel := relativeOrAbsolute(path, file)
			return newResult(SectionTerraform, "backend configured", StatusPass, fmt.Sprintf("Backend block detected in %s", rel), "Ensure backend credentials and remote state access are configured")
		}
	}

	return newResult(SectionTerraform, "backend configured", StatusWarn, "No backend block detected", "Terraform will use local state unless you configure a backend")
}

func checkTerraformFmt(ctx context.Context, path string) Result {
	if _, err := exec.LookPath("terraform"); err != nil {
		return newResult(SectionTerraform, "terraform fmt", StatusWarn, "Skipped because terraform is not installed")
	}
	if countTerraformFiles(path) == 0 {
		return newResult(SectionTerraform, "terraform fmt", StatusWarn, "Skipped because no Terraform files were found")
	}

	code, output, err := runCommand(ctx, path, "terraform", "fmt", "-check", "-diff", "-no-color")
	if err == nil && code == 0 {
		return newResult(SectionTerraform, "terraform fmt", StatusPass, "Terraform files are formatted")
	}

	if code == 2 || code == 3 {
		details := []string{"Terraform files need formatting", "Run `tfauto fmt --path <dir>` to fix formatting"}
		if trimmed := strings.TrimSpace(output); trimmed != "" {
			details = append(details, trimmed)
		}
		return newResult(SectionTerraform, "terraform fmt", StatusWarn, details...)
	}

	return newResult(SectionTerraform, "terraform fmt", StatusFail, fmt.Sprintf("terraform fmt failed: %v", err))
}

func checkTerraformValidate(ctx context.Context, path string) Result {
	if _, err := exec.LookPath("terraform"); err != nil {
		return newResult(SectionTerraform, "terraform validate", StatusWarn, "Skipped because terraform is not installed")
	}
	if countTerraformFiles(path) == 0 {
		return newResult(SectionTerraform, "terraform validate", StatusWarn, "Skipped because no Terraform files were found")
	}
	if !dirExists(filepath.Join(path, ".terraform")) {
		return newResult(SectionTerraform, "terraform validate", StatusWarn, "Skipped because the project is not initialized", "Run `tfauto validate --path <dir>` or `terraform init` first")
	}

	code, output, err := runCommand(ctx, path, "terraform", "validate", "-no-color")
	if err == nil && code == 0 {
		return newResult(SectionTerraform, "terraform validate", StatusPass, "Terraform configuration is valid")
	}

	details := []string{"Terraform validation failed"}
	if trimmed := strings.TrimSpace(output); trimmed != "" {
		details = append(details, trimmed)
	}
	return newResult(SectionTerraform, "terraform validate", StatusFail, details...)
}

func checkProvidersInstalled(path string) Result {
	providerDir := filepath.Join(path, ".terraform", "providers")
	if !dirExists(filepath.Join(path, ".terraform")) {
		return newResult(SectionTerraform, "providers installed", StatusWarn, "Skipped because the project is not initialized")
	}
	if dirExists(providerDir) {
		return newResult(SectionTerraform, "providers installed", StatusPass, "Provider plugins are present in .terraform/providers")
	}
	return newResult(SectionTerraform, "providers installed", StatusWarn, "Provider plugins are not downloaded yet", "Run `terraform init` to install providers")
}

func checkModulesDownloaded(path string) Result {
	moduleDir := filepath.Join(path, ".terraform", "modules")
	if !dirExists(filepath.Join(path, ".terraform")) {
		return newResult(SectionTerraform, "modules downloaded", StatusWarn, "Skipped because the project is not initialized")
	}
	if dirExists(moduleDir) {
		return newResult(SectionTerraform, "modules downloaded", StatusPass, "Module cache is present in .terraform/modules")
	}
	return newResult(SectionTerraform, "modules downloaded", StatusWarn, "Modules have not been downloaded yet", "Run `terraform init` to download modules")
}

func checkWorkspace(path string) Result {
	if workspace := os.Getenv("TF_WORKSPACE"); workspace != "" {
		return newResult(SectionTerraform, "terraform workspace", StatusPass, fmt.Sprintf("TF_WORKSPACE=%s", workspace))
	}

	workspaceFile := filepath.Join(path, ".terraform", "environment")
	if data, err := os.ReadFile(workspaceFile); err == nil {
		workspace := strings.TrimSpace(string(data))
		if workspace == "" {
			workspace = "default"
		}
		return newResult(SectionTerraform, "terraform workspace", StatusPass, fmt.Sprintf("Current workspace: %s", workspace))
	}

	return newResult(SectionTerraform, "terraform workspace", StatusWarn, "Workspace could not be determined", "Default workspace is likely in use unless TF_WORKSPACE or .terraform/environment says otherwise")
}

func checkVariablePromptRisk(path string) Result {
	requiredVariables, err := requiredVariablesWithoutDefaults(path)
	if err != nil {
		return newResult(SectionTerraform, "variable prompt risk", StatusWarn, fmt.Sprintf("Unable to inspect variables: %v", err))
	}

	tfvarsFiles := matchingFiles(path, func(name string) bool {
		return name == "terraform.tfvars" ||
			name == "terraform.tfvars.json" ||
			strings.HasSuffix(name, ".auto.tfvars") ||
			strings.HasSuffix(name, ".auto.tfvars.json")
	})

	if len(requiredVariables) == 0 {
		return newResult(SectionTerraform, "variable prompt risk", StatusPass, "No required variables without defaults detected")
	}

	details := []string{
		fmt.Sprintf("Required variables without defaults: %s", strings.Join(requiredVariables, ", ")),
	}
	if len(tfvarsFiles) == 0 {
		details = append(details, "No terraform.tfvars or .auto.tfvars files detected")
		return newResult(SectionTerraform, "variable prompt risk", StatusWarn, details...)
	}

	relFiles := make([]string, 0, len(tfvarsFiles))
	for _, file := range tfvarsFiles {
		relFiles = append(relFiles, relativeOrAbsolute(path, file))
	}
	details = append(details, fmt.Sprintf("Variable files detected: %s", strings.Join(relFiles, ", ")))

	return newResult(SectionTerraform, "variable prompt risk", StatusWarn, details...)
}

type awsIdentity struct {
	Account string `json:"Account"`
	Arn     string `json:"Arn"`
}

func getAWSCallerIdentity(ctx context.Context) (awsIdentity, error) {
	cmd := exec.CommandContext(ctx, "aws", "sts", "get-caller-identity", "--output", "json")
	out, err := cmd.Output()
	if err != nil {
		return awsIdentity{}, err
	}

	var payload awsIdentity
	if err := json.Unmarshal(out, &payload); err != nil {
		return awsIdentity{}, err
	}

	return payload, nil
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

func checkAWSCredentials(ctx context.Context, path string) Result {
	awsPath, err := exec.LookPath("aws")
	if err != nil {
		return newResult(SectionAWS, "credentials", StatusWarn, "AWS CLI was not found in PATH", "Install AWS CLI if you want identity verification from tfauto doctor")
	}

	identity, err := getAWSCallerIdentity(ctx)
	if err != nil {
		details := []string{fmt.Sprintf("AWS CLI found at %s", awsPath)}
		details = append(details, awsEnvironmentDetails(path)...)
		details = append(details, "Unable to resolve AWS caller identity")
		return newResult(SectionAWS, "credentials", StatusWarn, details...)
	}

	details := []string{fmt.Sprintf("AWS CLI found at %s", awsPath)}
	details = append(details, awsEnvironmentDetails(path)...)
	if identity.Account != "" {
		details = append(details, fmt.Sprintf("Account: %s", identity.Account))
	}
	return newResult(SectionAWS, "credentials", StatusPass, details...)
}

func checkAWSProfile() Result {
	if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		return newResult(SectionAWS, "profile", StatusPass, fmt.Sprintf("AWS_PROFILE=%s", profile))
	}
	return newResult(SectionAWS, "profile", StatusWarn, "AWS_PROFILE is not set")
}

func checkAWSCallerIdentity(ctx context.Context, path string) Result {
	awsPath, err := exec.LookPath("aws")
	if err != nil {
		return newResult(SectionAWS, "caller identity", StatusWarn, "AWS CLI was not found in PATH")
	}

	identity, err := getAWSCallerIdentity(ctx)
	if err != nil {
		details := []string{fmt.Sprintf("AWS CLI found at %s", awsPath)}
		details = append(details, awsEnvironmentDetails(path)...)
		details = append(details, "Unable to resolve AWS caller identity")
		return newResult(SectionAWS, "caller identity", StatusWarn, details...)
	}

	details := []string{fmt.Sprintf("AWS CLI found at %s", awsPath)}
	details = append(details, awsEnvironmentDetails(path)...)
	if identity.Account != "" {
		details = append(details, fmt.Sprintf("Account: %s", identity.Account))
	}
	if identity.Arn != "" {
		details = append(details, fmt.Sprintf("ARN: %s", identity.Arn))
	}
	return newResult(SectionAWS, "caller identity", StatusPass, details...)
}

func checkAWSRegion(path string) Result {
	if region := os.Getenv("AWS_REGION"); region != "" {
		return newResult(SectionAWS, "region", StatusPass, fmt.Sprintf("Using AWS_REGION=%s", region))
	}
	if region := os.Getenv("AWS_DEFAULT_REGION"); region != "" {
		return newResult(SectionAWS, "region", StatusPass, fmt.Sprintf("Using AWS_DEFAULT_REGION=%s", region))
	}

	files, err := terraformFiles(path)
	if err != nil {
		return newResult(SectionAWS, "region", StatusWarn, fmt.Sprintf("Unable to inspect provider configuration: %v", err))
	}

	for _, file := range files {
		content, readErr := os.ReadFile(file)
		if readErr != nil {
			continue
		}
		if strings.Contains(string(content), "provider \"aws\"") && strings.Contains(string(content), "region") {
			return newResult(SectionAWS, "region", StatusPass, fmt.Sprintf("AWS region appears to be configured in %s", relativeOrAbsolute(path, file)))
		}
	}

	return newResult(SectionAWS, "region", StatusWarn, "No AWS region detected in environment or provider configuration", "Set AWS_REGION or configure region in the provider block")
}

func checkGitBinary() Result {
	path, err := exec.LookPath("git")
	if err != nil {
		return newResult(SectionGit, "git installed", StatusWarn, "git was not found in PATH")
	}
	return newResult(SectionGit, "git installed", StatusPass, fmt.Sprintf("Found at %s", path))
}

func checkGitBranch(ctx context.Context, path string) Result {
	if _, err := exec.LookPath("git"); err != nil {
		return newResult(SectionGit, "current branch", StatusWarn, "Skipped because git is not installed")
	}

	out, err := exec.CommandContext(ctx, "git", "-C", path, "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return newResult(SectionGit, "current branch", StatusWarn, "Unable to determine the current git branch", "Make sure this is a git repository")
	}

	branch := strings.TrimSpace(string(out))
	if branch == "" || branch == "HEAD" {
		return newResult(SectionGit, "current branch", StatusWarn, "Repository is detached or branch could not be determined")
	}

	return newResult(SectionGit, "current branch", StatusPass, branch)
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

func runCommand(ctx context.Context, dir, name string, args ...string) (int, string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err == nil {
		return 0, string(output), nil
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), string(output), err
	}

	return 0, string(output), err
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
