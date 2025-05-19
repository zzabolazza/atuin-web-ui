# Build script for atuin-ui
# This script builds the frontend, embeds it in the backend, and compiles for multiple platforms

# Stop on error
$ErrorActionPreference = "Stop"

# Configuration
$frontendDir = "frontend"
$backendDir = "backend"
$distDir = "dist"
$frontendDistDir = "frontend/dist"
$backendEmbedDir = "backend/frontend_dist"
$appName = "atuin-web-ui"

# Ensure dist directory exists
if (-not (Test-Path $distDir)) {
    New-Item -ItemType Directory -Path $distDir | Out-Null
    Write-Host "Created dist directory"
}

# Build frontend
Write-Host "Building frontend..."
Push-Location $frontendDir
try {
    # Install dependencies if needed
    if (-not (Test-Path "node_modules")) {
        Write-Host "Installing frontend dependencies..."
        npm install
    }

    # Build frontend
    npm run build

    if (-not $?) {
        throw "Frontend build failed"
    }

    Write-Host "Frontend built successfully"
} finally {
    Pop-Location
}

# Prepare backend embed directory
Write-Host "Preparing backend embed directory..."
if (Test-Path $backendEmbedDir) {
    Remove-Item -Recurse -Force $backendEmbedDir
}
New-Item -ItemType Directory -Path $backendEmbedDir | Out-Null

# Copy frontend build to backend embed directory
Write-Host "Copying frontend build to backend embed directory..."
Copy-Item -Recurse "$frontendDistDir/*" -Destination $backendEmbedDir

# Ensure the directory structure exists
Write-Host "Ensuring directory structure exists..."
$directories = @("css", "js", "img", "fonts")
foreach ($dir in $directories) {
    $path = Join-Path -Path $backendEmbedDir -ChildPath $dir
    if (-not (Test-Path $path)) {
        New-Item -ItemType Directory -Path $path | Out-Null
        Write-Host "Created directory: $path"
    }
}

# Build backend for different platforms
Write-Host "Building backend for multiple platforms..."

# Define build targets
$targets = @(
    @{GOOS="windows"; GOARCH="amd64"; Suffix=".exe"},
    @{GOOS="darwin"; GOARCH="amd64"; Suffix=""},
    @{GOOS="darwin"; GOARCH="arm64"; Suffix=""},
    @{GOOS="linux"; GOARCH="amd64"; Suffix=""}
)

Write-Host "Building for all platforms"

# Build for each target
foreach ($target in $targets) {
    $goos = $target.GOOS
    $goarch = $target.GOARCH
    $suffix = $target.Suffix
    $outputName = "$appName-$goos-$goarch$suffix"
    $outputPath = "$distDir/$outputName"

    Write-Host "Building for $goos/$goarch..."

    $env:GOOS = $goos
    $env:GOARCH = $goarch

    Push-Location $backendDir
    try {
        go build -o "../$outputPath" -ldflags="-s -w" .

        if (-not $?) {
            throw "Backend build failed for $goos/$goarch"
        }

        Write-Host "Built $outputPath successfully"
    } finally {
        Pop-Location
    }
}

Write-Host "Build completed successfully!"
Write-Host "Binaries are available in the $distDir directory"
