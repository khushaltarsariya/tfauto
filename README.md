# tfauto

`tfauto` is a Terraform workflow automation CLI built with Go and Cobra.
It does not replace Terraform. It adds guardrails, reusable templates, and a consistent command flow on top of Terraform CLI.

## Why tfauto

Terraform users often struggle more with workflow mistakes than with HCL itself:
- running commands in the wrong directory
- forgetting `terraform init`
- repeating the same boilerplate project setup
- unsafe `terraform destroy`
- inconsistent structure across projects

`tfauto` standardizes those workflows with a thin CLI layer and reusable templates.

## Current capabilities

- Wrap common Terraform commands: `plan`, `apply`, `destroy`, `validate`, `fmt`
- Scaffold projects from built-in templates with `init`
- List and inspect templates with `templates` and `template`
- Run environment checks with `doctor`
- Ship as a single compiled binary

## Install

### Build from source

Prerequisites:
- Go 1.22+
- Terraform CLI installed and available in `PATH`

```bash
git clone https://github.com/your-username/tfauto.git
cd tfauto
go build -ldflags "-X tfauto/cmd.version=v0.1.0 -X tfauto/cmd.commit=$(git rev-parse --short HEAD) -X tfauto/cmd.date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o tfauto
```

Windows PowerShell:

```powershell
$commit = git rev-parse --short HEAD
$date = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
go build -ldflags "-X tfauto/cmd.version=v0.1.0 -X tfauto/cmd.commit=$commit -X tfauto/cmd.date=$date" -o tfauto.exe
```

## Quick start

```bash
tfauto init --template aws-basic --target ./my-project
tfauto validate --path ./my-project
tfauto fmt --path ./my-project
tfauto plan --path ./my-project
tfauto apply --path ./my-project
```

## Commands

### Project scaffolding

```bash
tfauto templates
tfauto template aws-three-tier-vpc
tfauto init --template aws-three-tier-vpc --target ./network
```

### Terraform workflow

```bash
tfauto validate --path ./network
tfauto fmt --path ./network
tfauto plan --path ./network
tfauto apply --path ./network
tfauto destroy --path ./network
```

### Diagnostics

```bash
tfauto doctor --path ./network
tfauto version
```

## Built-in templates

- `aws-basic`
  Basic EC2 + networking starter
- `aws-vpc`
  Simple VPC starter across two public subnets
- `aws-s3-static-site`
  Public S3 static website starter
- `aws-three-tier-vpc`
  Three-tier VPC with public, app, and DB subnets, NAT, and DB subnet group
- `aws-alb-asg-webapp`
  ALB + Auto Scaling web app stack for an existing VPC

Each template follows the same structure:

```text
templates/<template-name>/
|- DESCRIPTION.md
|- main.tf
|- outputs.tf
|- provider.tf
`- variables.tf
```

## Architecture

```text
User CLI
  -> Cobra command layer (cmd/)
  -> Business logic (internal/)
  -> Terraform execution layer
  -> Terraform CLI
```

Rules used in this project:
- keep `cmd/` thin
- keep business logic in `internal/`
- return errors instead of panicking
- use `context.Context` for command execution

## Release process

This repository includes GitHub Actions for CI and tagged releases.

- `ci.yml` builds the CLI on every push and pull request
- `release.yml` builds Linux, macOS, and Windows binaries when you push a tag like `v0.1.0`

Example release tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

## Launch checklist

Before a public launch:
- set the GitHub repository URL in this README
- verify all templates with real `terraform init` and `terraform validate`
- create the first tagged release such as `v0.1.0`
- add screenshots or terminal demos to the repository page
- publish example usage in the README or docs

## Near-term next steps

- add template validation tests or smoke checks
- add Linux and macOS install instructions with release download examples
- improve template descriptions and examples
- add more production-focused templates such as `aws-rds-postgres` and `aws-ecs-fargate-service`
- improve `doctor` with AWS identity and backend checks

## License

MIT. See [LICENSE](LICENSE).
