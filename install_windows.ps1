# Allow specifying custom destination
$InstallDir = $env:DIR
if (-not $InstallDir) {
    $InstallDir = "$env:USERPROFILE\AppData\Local\Microsoft\WindowsApps"
}

# Detect architecture
$arch = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Get latest release tag from GitHub
$apiUrl = "https://api.github.com/repos/Ealenn/gira/releases/latest"
$tag = (Invoke-RestMethod -Uri $apiUrl).tag_name

# Build filename and URL
$fileName = "gira-windows-$arch.exe"
$downloadUrl = "https://github.com/Ealenn/gira/releases/download/$tag/$fileName"
$destination = Join-Path -Path $InstallDir -ChildPath "gira.exe"

# Download
Write-Host "Downloading $fileName..."
Invoke-WebRequest -Uri $downloadUrl -OutFile $destination

# Make executable (already is, since it's .exe)
Write-Host "✅ gira installed to $destination"
