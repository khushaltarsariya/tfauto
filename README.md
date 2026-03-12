# tfauto

> **Terraform workflow automation CLI — simplified, safe, and developer-friendly.**

tfauto is a command-line tool built on top of Terraform that standardizes infrastructure workflows, reduces operational mistakes, and provides reusable infrastructure templates. It does **not** replace Terraform — it makes working with Terraform easier and safer.

---

## Why tfauto?

Working with Terraform is powerful, but painful for many developers:

- Forgetting the correct command sequence (`init` → `validate` → `plan` → `apply`)
- Cryptic error messages with no clear fix
- Accidental `terraform destroy` without confirmation
- Repeating the same boilerplate infrastructure code across projects
- Inconsistent project structures across teams

**tfauto solves all of this.**

---

## Features

- ✅ Standardized Terraform workflow commands
- ✅ Reusable AWS infrastructure templates
- ✅ Safe `destroy` with confirmation prompt
- ✅ Environment diagnostics (`doctor` command)
- ✅ Clean, beginner-friendly CLI experience
- ✅ Single binary — no runtime dependencies

---

## Installation

### Build from source

**Prerequisites:** Go 1.21+, Terraform CLI installed

```bash
# Clone the repository
git clone https://github.com/yourusername/tfauto.git
cd tfauto

# Build
go build -ldflags "-X tfauto/cmd.version=v0.1.0" -o tfauto

# Move to PATH (Linux/Mac)
sudo mv tfauto /usr/local/bin/

# Verify installation
tfauto version
```

**Windows:**
```bash
go build -ldflags "-X tfauto/cmd.version=v0.1.0" -o tfauto.exe
```

---

## Quick Start

```bash
# 1. Create a new Terraform project from a template
tfauto init --template aws-basic --target ./my-project

# 2. Validate the configuration
tfauto validate --path ./my-project

# 3. Format your Terraform files
tfauto fmt --path ./my-project

# 4. Preview infrastructure changes
tfauto plan --path ./my-project

# 5. Apply infrastructure
tfauto apply --path ./my-project
```

---

## Commands

### `tfauto init`
Creates a new Terraform project from a template.

```bash
tfauto init --template aws-basic --target ./my-project
```

| Flag | Description |
|------|-------------|
| `--template` | Template name to use (required) |
| `--target` | Directory to create the project in (required) |

---

### `tfauto plan`
Runs `terraform plan` in the specified project directory.

```bash
tfauto plan --path ./my-project
```

---

### `tfauto apply`
Runs `terraform apply` in the specified project directory.

```bash
tfauto apply --path ./my-project
```

---

### `tfauto destroy`
Destroys infrastructure with a required confirmation prompt. You will be asked to confirm before anything is deleted.

```bash
tfauto destroy --path ./my-project
```

---

### `tfauto validate`
Runs `terraform validate` to check your configuration for errors.

```bash
tfauto validate --path ./my-project
```

---

### `tfauto fmt`
Formats all Terraform files in the project directory.

```bash
# Format files
tfauto fmt --path ./my-project

# Check formatting without making changes
tfauto fmt --path ./my-project --check
```

---

### `tfauto templates`
Lists all available built-in templates.

```bash
tfauto templates
```

---

### `tfauto template`
Shows details about a specific template including its description and files.

```bash
tfauto template aws-vpc
```

---

### `tfauto doctor`
Runs environment diagnostics to check everything is set up correctly.

```bash
tfauto doctor
```

Checks:
- Terraform is installed
- Terraform version compatibility
- AWS credentials are configured
- Terraform files exist in the target path

---

### `tfauto version`
Displays the current tfauto version.

```bash
tfauto version
```

---

## Templates

tfauto comes with the following built-in AWS templates:

### `aws-basic`
Simple AWS EC2 infrastructure. Great starting point for learning Terraform.

```bash
tfauto init --template aws-basic --target ./my-ec2-project
```

**Creates:**
- EC2 instance
- Security group
- Key pair configuration

---

### `aws-vpc`
Production-ready VPC setup.

```bash
tfauto init --template aws-vpc --target ./my-vpc-project
```

**Creates:**
- VPC
- Public and private subnets
- Internet gateway
- Route tables

---

### `aws-s3-static-site`
S3-based static website hosting.

```bash
tfauto init --template aws-s3-static-site --target ./my-site
```

**Creates:**
- S3 bucket
- Public access configuration
- Static website hosting settings

---

Each template includes:

```
template-name/
├── main.tf
├── variables.tf
├── outputs.tf
├── provider.tf
└── DESCRIPTION.md
```

---

## Prerequisites

| Dependency | Required | Notes |
|------------|----------|-------|
| Go 1.21+ | Build only | For building from source |
| Terraform CLI | ✅ Required | [Install Terraform](https://developer.hashicorp.com/terraform/install) |
| AWS CLI + credentials | For AWS templates | [Configure AWS](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-quickstart.html) |

Run `tfauto doctor` to verify your environment is ready.

---

## Project Structure

```
tfauto/
│
├── main.go
├── cmd/              # CLI command definitions
│   ├── root.go
│   ├── init.go
│   ├── plan.go
│   ├── apply.go
│   ├── destroy.go
│   ├── validate.go
│   ├── fmt.go
│   ├── templates.go
│   ├── template.go
│   ├── doctor.go
│   └── version.go
│
├── internal/
│   ├── terraform/    # Terraform execution logic
│   └── generator/    # Template copying and scaffolding
│
├── templates/        # Built-in infrastructure templates
│   ├── aws-basic/
│   ├── aws-vpc/
│   └── aws-s3-static-site/
│
├── go.mod
└── README.md
```

---

## Roadmap

- [x] Phase 1 — Core CLI workflow (current)
- [ ] Phase 2 — Intelligent error translation and state safety
- [ ] Phase 3 — Expanded multi-cloud templates (Azure, GCP)
- [ ] Phase 4 — Team workflow enforcement (`.tfauto.yaml`, hooks, audit logging)
- [ ] Phase 5 — AI-assisted plan explanation and security scanning

---

## Contributing

Contributions are welcome! Please open an issue first to discuss what you'd like to change.

```bash
# Fork and clone the repo
git clone https://github.com/yourusername/tfauto.git

# Create a feature branch
git checkout -b feature/my-feature

# Make your changes, then run tests
go test ./...

# Submit a pull request
```

Please keep these principles in mind:
- No business logic in `cmd/` — only CLI interaction
- All core logic lives in `internal/`
- Always use `context.Context` for command execution
- Return errors instead of panicking

---

## License

MIT License — see [LICENSE](./LICENSE) for details.

---

## Acknowledgements

Built with:
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Terraform](https://www.terraform.io/) — Infrastructure engine

---

> **tfauto** is not affiliated with HashiCorp or Terraform. Terraform is a trademark of HashiCorp.
