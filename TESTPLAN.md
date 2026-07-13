# tfauto Test Plan

Living QA checklist for public releases of `tfauto`.

Use this file as the release gate for every tagged build. Mark items complete only after they have been verified on the target release candidate.

## 1. Installation Testing

- [ ] Confirm the release binary is attached to the GitHub Release.
- [ ] Confirm binary names are correct for each OS and architecture.
- [ ] Confirm SHA256 checksum files are attached and match the binaries.
- [ ] Confirm the version string prints correctly from the binary.
- [ ] Confirm the binary runs when invoked from the current directory.
- [ ] Confirm the binary runs when added to `PATH`.
- [ ] Confirm the binary fails clearly when Terraform is not installed.
- [ ] Confirm the binary fails clearly when AWS credentials are not configured.
- [ ] Confirm the binary fails clearly when Git is not installed.
- [ ] Confirm Windows execution works with `.\\tfauto.exe`.
- [ ] Confirm Linux/macOS execution works with `./tfauto`.
- [ ] Confirm executable permissions are correct on Unix-like systems.
- [ ] Confirm the binary does not require source code to run.

## 2. CLI Testing

### Global

- [ ] `tfauto --help`
- [ ] `tfauto --version`
- [ ] `tfauto --json --help`
- [ ] `tfauto --debug --help`
- [ ] `tfauto --timeout 5m --help`
- [ ] `tfauto --chdir <dir> --help`

### Commands

- [ ] `tfauto init`
- [ ] `tfauto plan`
- [ ] `tfauto apply`
- [ ] `tfauto destroy`
- [ ] `tfauto validate`
- [ ] `tfauto fmt`
- [ ] `tfauto doctor`
- [ ] `tfauto config`
- [ ] `tfauto templates`
- [ ] `tfauto template <name>`
- [ ] `tfauto version`

### Expected checks for each command

- [ ] Help text is readable and consistent.
- [ ] Examples are present where useful.
- [ ] Exit code is `0` on success.
- [ ] Exit code is non-zero on failure.
- [ ] Error messages are actionable.
- [ ] Success messages are consistent.
- [ ] JSON output is valid when enabled.

## 3. Argument Testing

- [ ] Missing required flags.
- [ ] Unknown flags.
- [ ] Too many positional arguments.
- [ ] Invalid flag values.
- [ ] Conflicting flag combinations.
- [ ] Empty string values.
- [ ] Relative path handling.
- [ ] Absolute path handling.
- [ ] Paths with spaces.
- [ ] Paths containing special characters.

## 4. Global Flag Testing

- [ ] `--debug`
- [ ] `--timeout`
- [ ] `--json`
- [ ] `--help`
- [ ] `--version`
- [ ] `--chdir`

Verify for each flag:

- [ ] Flag is documented in help.
- [ ] Flag behaves the same across commands.
- [ ] Invalid values fail clearly.
- [ ] JSON mode suppresses noisy human output.

## 5. Terraform Testing

- [ ] Terraform is installed.
- [ ] Terraform is missing.
- [ ] Terraform version is supported.
- [ ] Terraform version is too old.
- [ ] `terraform init` succeeds.
- [ ] `terraform init` fails.
- [ ] `terraform validate` succeeds.
- [ ] `terraform validate` fails.
- [ ] `terraform fmt` succeeds.
- [ ] `terraform fmt` fails.
- [ ] `terraform plan` succeeds.
- [ ] `terraform plan` fails.
- [ ] `terraform apply` succeeds.
- [ ] `terraform apply` fails.
- [ ] `terraform destroy` succeeds.
- [ ] `terraform destroy` fails.
- [ ] Backend configuration is valid.
- [ ] Backend configuration is missing.
- [ ] Providers download successfully.
- [ ] Provider download fails.
- [ ] Modules download successfully.
- [ ] Module download fails.
- [ ] State lock handling is safe.
- [ ] Workspace handling is correct.

## 6. Template Testing

- [ ] `templates` lists all bundled templates.
- [ ] `template <name>` shows metadata and files.
- [ ] Template manifests are readable.
- [ ] Missing `DESCRIPTION.md` is handled safely.
- [ ] Missing manifest is handled safely.
- [ ] Invalid manifest is handled safely.
- [ ] Duplicate template names are rejected or resolved predictably.
- [ ] Empty template directories are handled safely.
- [ ] Nested template directories copy correctly.
- [ ] Embedded templates are available in the binary.
- [ ] New templates remain backward compatible with old ones.

## 7. Configuration Testing

- [ ] Missing config file is handled.
- [ ] Invalid YAML is handled.
- [ ] Invalid JSON is handled.
- [ ] Default config works.
- [ ] Custom config works.
- [ ] Environment variables override config where expected.
- [ ] Missing environment variables fail clearly.

## 8. JSON Output Testing

- [ ] JSON is valid for success cases.
- [ ] JSON is valid for error cases.
- [ ] JSON is machine-readable and stable.
- [ ] JSON includes command name.
- [ ] JSON includes status or exit-code context.
- [ ] JSON output does not mix with human-readable noise.

## 9. Error Handling Testing

- [ ] Permission denied.
- [ ] Disk full.
- [ ] Folder missing.
- [ ] File locked.
- [ ] Network unavailable.
- [ ] Terraform not installed.
- [ ] AWS CLI missing.
- [ ] AWS credentials missing.
- [ ] AWS profile missing.
- [ ] Terraform binary corrupted.
- [ ] Invalid template path.
- [ ] Unsafe overwrite attempt.
- [ ] Destroy confirmation rejected.

## 10. Cross-Platform Testing

- [ ] Windows PowerShell.
- [ ] Windows CMD.
- [ ] Linux Bash.
- [ ] Linux Zsh.
- [ ] macOS Bash.
- [ ] macOS Zsh.
- [ ] Windows path separators work.
- [ ] Unix path separators work.
- [ ] Executable bit behavior is correct on Unix-like systems.
- [ ] Binary naming is correct per OS.

## 11. Performance Testing

- [ ] CLI startup time is acceptable.
- [ ] Help output is fast.
- [ ] `doctor` completes quickly.
- [ ] Template copy speed is acceptable.
- [ ] Large template initialization is acceptable.
- [ ] JSON output generation is fast.

## 12. Security Testing

- [ ] Path traversal is blocked.
- [ ] Unsafe absolute paths are rejected where appropriate.
- [ ] Command injection is not possible through flags.
- [ ] Symlink attacks are handled safely.
- [ ] Unsafe file overwrite is prevented.
- [ ] Dangerous destroy requires confirmation.
- [ ] Template loading does not execute arbitrary code.

## 13. Release Testing

- [ ] Git tag exists for the release.
- [ ] GitHub Actions release workflow ran successfully.
- [ ] Release artifacts were uploaded.
- [ ] Checksums match all artifacts.
- [ ] Version injection works in the binary.
- [ ] Release notes are correct.
- [ ] Downloaded release binary runs successfully.

## 14. User Experience Testing

- [ ] Help text is consistent across commands.
- [ ] Error messages are actionable.
- [ ] Success messages are short and clear.
- [ ] Exit codes match the outcome.
- [ ] Examples are practical and current.
- [ ] JSON mode is predictable.
- [ ] Color output does not break readability.
- [ ] Output feels consistent with professional CLI tools.

## 15. Regression Testing Targets

- [ ] `init`
- [ ] `plan`
- [ ] `apply`
- [ ] `destroy`
- [ ] `validate`
- [ ] `fmt`
- [ ] `doctor`
- [ ] `config`
- [ ] `templates`
- [ ] `template`
- [ ] `version`

## 16. Recommended Automated Coverage

- [ ] Unit tests for `internal/terraform`.
- [ ] Unit tests for `internal/generator`.
- [ ] Unit tests for `internal/doctor`.
- [ ] Command-level tests for Cobra wiring.
- [ ] Template embedding tests.
- [ ] JSON output tests.
- [ ] Release workflow verification.
- [ ] Smoke validation for every bundled template.

## 17. Manual Release Checklist

- [ ] Build release binaries locally.
- [ ] Run `go test ./...`.
- [ ] Run `go build ./...`.
- [ ] Validate every bundled template.
- [ ] Run `doctor` on a clean machine.
- [ ] Test `init`, `validate`, `fmt`, `plan`, `apply`, and `destroy` on a sample project.
- [ ] Confirm release assets are downloadable.
- [ ] Confirm a fresh install works from `PATH`.

## 18. Notes

- Record the release version being tested:
- Record the test environment:
- Record any failures and follow-up fixes:

