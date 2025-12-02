# tfauto 🚀
A simple and developer-friendly CLI tool that automates common Terraform workflows.

tfauto helps you:

- Scaffold Terraform project templates
- Run `terraform init`, `plan`, `apply`, `destroy`
- Format & validate Terraform code
- Browse available templates
- View full details of a selected template

Perfect for beginners, DevOps engineers, and anyone who wants a faster Terraform workflow.

---

## ⚡ Features

### ✔ Terraform Commands
- `tfauto init` – Create a new Terraform project from templates  
- `tfauto plan` – Run `terraform plan`  
- `tfauto apply` – Auto-approve applies  
- `tfauto destroy` – Safe destroy with confirmation  
- `tfauto validate` – Run validation  
- `tfauto fmt` – Format or check formatting  

### ✔ Template System
- `tfauto templates` – List all available templates  
- `tfauto template <name>` – View full description of a single template  
- Templates include a `DESCRIPTION.md` so users know what resources will be created

### ✔ Zero dependencies (other than Terraform)
tfauto only requires:
- Go-compiled binary  
- Terraform installed locally
- Configure your aws acoount using the aws-cli

---

## 📦 Installation

### **macOS / Linux / Windows**
Download the binary from the GitHub Releases page:

👉 https://github.com/<your-username>/tfauto/releases

Then:

```bash
chmod +x tfauto
sudo mv tfauto /usr/local/bin/
