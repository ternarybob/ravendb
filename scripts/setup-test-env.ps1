# RavenDB Test Environment Setup
# PowerShell script to set up the test environment

param(
    [string]$RavenDBUrl = "http://localhost:8080",
    [string]$DatabaseName = "RavenDBLibTestDB",
    [string]$ConfigFile = "config\test_config.toml"
)

$Green = [System.ConsoleColor]::Green
$Yellow = [System.ConsoleColor]::Yellow
$Red = [System.ConsoleColor]::Red
$Cyan = [System.ConsoleColor]::Cyan

function Write-ColoredOutput {
    param([string]$Message, [System.ConsoleColor]$Color)
    Write-Host $Message -ForegroundColor $Color
}

Write-ColoredOutput "🔧 Setting up RavenDB Test Environment" $Cyan
Write-ColoredOutput "=====================================" $Cyan

# Change to project root
Set-Location $PSScriptRoot\..

# Check if RavenDB server is accessible
Write-ColoredOutput "🌐 Checking RavenDB server connectivity..." $Yellow
try {
    $response = Invoke-WebRequest -Uri $RavenDBUrl -Method GET -TimeoutSec 10
    if ($response.StatusCode -eq 200) {
        Write-ColoredOutput "✅ RavenDB server is accessible at $RavenDBUrl" $Green
    }
} catch {
    Write-ColoredOutput "❌ Cannot connect to RavenDB server at $RavenDBUrl" $Red
    Write-Host "   Please ensure RavenDB is running. You can:"
    Write-Host "   1. Download from: https://ravendb.net/downloads"
    Write-Host "   2. Or use Docker: docker run -p 8080:8080 ravendb/ravendb"
    Write-Host ""
    $continue = Read-Host "Continue anyway? (y/N)"
    if ($continue -ne 'y' -and $continue -ne 'Y') {
        exit 1
    }
}

# Create or update test configuration
Write-ColoredOutput "📝 Creating test configuration file: $ConfigFile" $Yellow

$configContent = @"
[database]
urls = ["$RavenDBUrl"]
database = "$DatabaseName"

[test]
# Test timeout in seconds
timeout = 30
# Clean database before tests
clean_before_tests = true
# Clean database after tests
clean_after_tests = true
"@

$configContent | Out-File -FilePath $ConfigFile -Encoding UTF8
Write-ColoredOutput "✅ Configuration file created: $ConfigFile" $Green

# Display configuration
Write-ColoredOutput "📋 Test Configuration:" $Cyan
Get-Content $ConfigFile | ForEach-Object { Write-Host "  $_" }

# Verify Go environment
Write-ColoredOutput "🔍 Checking Go environment..." $Yellow
$goVersion = go version 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-ColoredOutput "✅ Go found: $goVersion" $Green
} else {
    Write-ColoredOutput "❌ Go not found. Please install Go 1.25 or later" $Red
    exit 1
}

# Check go.mod exists
if (Test-Path "go.mod") {
    Write-ColoredOutput "✅ go.mod found" $Green
} else {
    Write-ColoredOutput "❌ go.mod not found. Run 'go mod init' first" $Red
    exit 1
}

# Download dependencies
Write-ColoredOutput "📦 Downloading Go dependencies..." $Yellow
go mod download
if ($LASTEXITCODE -eq 0) {
    Write-ColoredOutput "✅ Dependencies downloaded" $Green
} else {
    Write-ColoredOutput "❌ Failed to download dependencies" $Red
    exit 1
}

# Verify test files exist
$testFiles = @("database_test.go", "collection_test.go", "test_utils.go")
Write-ColoredOutput "🧪 Checking test files..." $Yellow
foreach ($file in $testFiles) {
    if (Test-Path $file) {
        Write-ColoredOutput "✅ Found: $file" $Green
    } else {
        Write-ColoredOutput "❌ Missing: $file" $Red
    }
}

Write-Host ""
Write-ColoredOutput "🎉 Test environment setup completed!" $Green
Write-Host ""
Write-ColoredOutput "Next steps:" $Cyan
Write-Host "1. Ensure RavenDB server is running at $RavenDBUrl"
Write-Host "2. Run tests: .\scripts\run-tests.ps1"
Write-Host "3. Run tests with coverage: .\scripts\run-tests.ps1 -Coverage"
Write-Host "4. Run specific tests: .\scripts\run-tests.ps1 -Pattern 'TestDatabase.*'"