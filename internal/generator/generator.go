package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyTemplate(name, targetDir string) error {
	src := filepath.Join("templates", name)

	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("Template %s does not exist", name)
	}

	return copyDir(src, targetDir)
}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, os.ModePerm)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)

	if err != nil {
		return err
	}

	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}
