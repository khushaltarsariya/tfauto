$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$templatesDir = Join-Path $repoRoot "templates"

if (-not (Test-Path $templatesDir)) {
    throw "templates directory not found: $templatesDir"
}

$workDir = Join-Path $repoRoot ".template-smoke"
if (Test-Path $workDir) {
    Remove-Item -Recurse -Force $workDir
}

New-Item -ItemType Directory -Path $workDir | Out-Null

Write-Host "Validating templates from $templatesDir"

Get-ChildItem $templatesDir -Directory | ForEach-Object {
    $templateName = $_.Name
    $target = Join-Path $workDir $templateName

    Write-Host ""
    Write-Host "==> $templateName"

    Copy-Item -Recurse $_.FullName $target

    terraform fmt -check -recursive $target
    terraform ("-chdir=" + $target) init -backend=false -input=false
    terraform ("-chdir=" + $target) validate
}

Write-Host ""
Write-Host "All templates validated successfully."
