package terraform

import (
	"fmt"
	"os"
	"os/exec"
)

func runInDir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func Init(path string) error {
	fmt.Println("Terraform Init")
	return runInDir(path, "terraform", "init", "-input=false")
}

func Plan(path string) error {
	fmt.Println("terraform plan")
	return runInDir(path, "terraform", "plan", "-input=false")
}

func Apply(path string) error {
	fmt.Println("terraform apply")
	return runInDir(path, "terraform", "apply", "-auto-approve", "-input=false")
}

func Destroy(path string) error {
	fmt.Println("terraform destroy")
	return runInDir(path, "terraform", "destroy", "-auto-approve", "-input=false")
}
