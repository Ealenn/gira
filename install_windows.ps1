# Gira Installer for Windows (PowerShell)

# Set strict mode for better error handling
Set-StrictMode -Version Latest

Write-Host ""
Write-Host "=== Gira Installer ===" -ForegroundColor Cyan
Write-Host ""

# Allow specifying a custom destination directory (default: WindowsApps)
$InstallDir = $env:DIR
if (-not $InstallDir) {
    $InstallDir = Join-Path $env:USERPROFILE "AppData\Local\Microsoft\WindowsApps"
}

# Ensure the install directory exists
if (-not (Test-Path $InstallDir)) {
    Write-Host "Creating install directory: $InstallDir"
    New-Item -Path $InstallDir -ItemType Directory -Force | Out-Null
}

# Detect architecture
$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Get latest release tag from GitHub API
$apiUrl = "https://api.github.com/repos/Ealenn/gira/releases/latest"
try {
    $tag = (Invoke-RestMethod -Uri $apiUrl -UseBasicParsing).tag_name
} catch {
    Write-Error "Failed to fetch the latest Gira release from GitHub."
    exit 1
}

# Build the download URL and destination path
$fileName = "gira-windows-$arch.exe"
$downloadUrl = "https://github.com/Ealenn/gira/releases/download/$tag/$fileName"
$destination = Join-Path -Path $InstallDir -ChildPath "gira.exe"

Write-Host ""
Write-Host "Downloading $fileName from $downloadUrl..." -ForegroundColor Yellow

try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $destination -UseBasicParsing
} catch {
    Write-Error "Failed to download the Gira binary. Please verify that the release exists for your platform."
    Write-Host "More information: https://github.com/Ealenn/gira" -ForegroundColor Cyan
    exit 1
}

Write-Host ""
Write-Host "$([char]0x2713) Gira $tag installed to $destination" -ForegroundColor Green
Write-Host ""

# Suggest autocompletion for PowerShell
Write-Host "Gira supports autocompletion for PowerShell." -ForegroundColor Cyan
Write-Host "To enable autocompletion, add the following to your PowerShell profile:" -ForegroundColor Yellow
Write-Host ""
Write-Host '    Invoke-Expression -Command (gira completion powershell | Out-String)' -ForegroundColor White
Write-Host ""
Write-Host "You can edit your profile with:" -ForegroundColor Gray
Write-Host "    notepad $PROFILE" -ForegroundColor White
Write-Host ""

# Suggest adding to PATH if not already present
$envPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($envPath -notlike "*$InstallDir*") {
    Write-Host "To use 'gira' everywhere, add the following directory to your PATH environment variable:" -ForegroundColor Yellow
    Write-Host "    $InstallDir" -ForegroundColor Cyan
    Write-Host ""
}
