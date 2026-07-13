package generator

import (
	"fmt"

	tplfs "github.com/khushaltarsariya/tfauto/templates"
)

func CopyTemplate(name, targetDir string) error {
	if !tplfs.Exists(name) {
		return fmt.Errorf("template %s does not exist", name)
	}

	return tplfs.Copy(name, targetDir)
}
