package tplfs

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed aws-basic aws-vpc aws-s3-static-site aws-three-tier-vpc aws-alb-asg-webapp aws-rds-postgres aws-ecs-fargate-service aws-lambda-apigateway aws-cloudfront-s3-static-site
var builtins embed.FS

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

func Exists(name string) bool {
	info, err := fs.Stat(builtins, name)
	return err == nil && info.IsDir()
}

func Description(name string) (string, error) {
	data, err := fs.ReadFile(builtins, filepath.ToSlash(filepath.Join(name, "DESCRIPTION.md")))
	if err != nil {
		if os.IsNotExist(err) || strings.Contains(err.Error(), "file does not exist") {
			return "", nil
		}
		return "", fmt.Errorf("read embedded description: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
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
		if rel == "DESCRIPTION.md" {
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
