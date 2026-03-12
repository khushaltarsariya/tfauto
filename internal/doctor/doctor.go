package doctor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
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

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
