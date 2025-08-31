# RavenDB Library Test Runner
# PowerShell script to run all tests for the RavenDB library

param(
    [string]$ConfigFile = "config\test_config.toml",
    [switch]$Coverage = $false,
    [switch]$Verbose = $false,
    [string]$Pattern = "*",
    [int]$Timeout = 60
)

# Colors for output
$Red = [System.ConsoleColor]::Red
$Green = [System.ConsoleColor]::Green
$Yellow = [System.ConsoleColor]::Yellow
$Cyan = [System.ConsoleColor]::Cyan
$White = [System.ConsoleColor]::White

function Write-ColoredOutput {
    param([string]$Message, [System.ConsoleColor]$Color = $White)
    Write-Host $Message -ForegroundColor $Color
}

function Write-Header {
    param([string]$Title)
    Write-Host ""
    Write-ColoredOutput "🧪 $Title" $Cyan
    Write-ColoredOutput ("=" * ($Title.Length + 3)) $Cyan
}

function Write-Step {
    param([string]$Step)
    Write-ColoredOutput "📋 $Step" $Yellow
}

function Write-Success {
    param([string]$Message)
    Write-ColoredOutput "✅ $Message" $Green
}

function Write-Error {
    param([string]$Message)
    Write-ColoredOutput "❌ $Message" $Red
}

# Main script
Write-Header "RavenDB Library Tests"

# Check if test configuration exists
if (!(Test-Path $ConfigFile)) {
    Write-Error "$ConfigFile not found!"
    Write-Host ""
    Write-Host "Please create $ConfigFile with your RavenDB connection settings"
    Write-Host "Example:"
    Write-Host "[database]"
    Write-Host 'urls = ["http://localhost:8080"]'
    Write-Host 'database = "RavenDBLibTestDB"'
    Write-Host ""
    Write-Host "[test]"
    Write-Host "timeout = 30"
    Write-Host "clean_before_tests = true"
    Write-Host "clean_after_tests = true"
    exit 1
}

# Display test configuration
Write-Step "Test Configuration:"
Get-Content $ConfigFile | ForEach-Object { Write-Host "  $_" }
Write-Host ""

# Change to project root
Set-Location $PSScriptRoot\..

# Download dependencies
Write-Step "Downloading dependencies..."
try {
    go mod download
    Write-Success "Dependencies downloaded"
} catch {
    Write-Error "Failed to download dependencies: $_"
    exit 1
}

# Format code
Write-Step "Formatting code..."
go fmt ./...
Write-Success "Code formatted"

# Build the project first
Write-Step "Building project..."
$buildResult = go build ./...
if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed"
    exit 1
}
Write-Success "Build completed"

# Test parameters
$testArgs = @("-v")
if ($Coverage) {
    $testArgs += @("-cover", "-coverprofile=coverage.out")
}
if ($Timeout -gt 0) {
    $testArgs += @("-timeout", "${Timeout}s")
}

# Run specific test patterns
if ($Pattern -eq "*" -or $Pattern -eq "") {
    # Run all tests
    Write-Header "Running All Tests"
    
    # Database tests
    Write-Step "Testing database operations..."
    $dbTestArgs = $testArgs + @("-run", "TestDatabase.*")
    & go test $dbTestArgs
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Database tests passed"
    } else {
        Write-Error "Database tests failed"
    }
    Write-Host ""

    # Collection tests
    Write-Step "Testing collection operations..."
    $collTestArgs = $testArgs + @("-run", "TestCollection.*")
    & go test $collTestArgs
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Collection tests passed"
    } else {
        Write-Error "Collection tests failed"
    }
    Write-Host ""

    # All tests with coverage
    Write-Step "Running complete test suite..."
    & go test $testArgs ./...
    $testExitCode = $LASTEXITCODE
    
} else {
    # Run specific pattern
    Write-Header "Running Tests: $Pattern"
    $patternArgs = $testArgs + @("-run", $Pattern)
    & go test $patternArgs ./...
    $testExitCode = $LASTEXITCODE
}

# Generate coverage report if requested
if ($Coverage -and (Test-Path "coverage.out")) {
    Write-Host ""
    Write-Step "Generating coverage report..."
    
    # Show coverage summary
    $coverageOutput = go tool cover -func=coverage.out
    $totalLine = $coverageOutput | Select-String "total:"
    if ($totalLine) {
        Write-ColoredOutput "📈 Coverage Summary:" $Cyan
        Write-Host "  $($totalLine.Line)"
    }
    
    # Generate HTML coverage report
    go tool cover -html=coverage.out -o coverage.html
    if (Test-Path "coverage.html") {
        Write-Success "HTML coverage report generated: coverage.html"
    }
}

# Final result
Write-Host ""
if ($testExitCode -eq 0) {
    Write-Success "All tests completed successfully!"
} else {
    Write-Error "Some tests failed (exit code: $testExitCode)"
    exit $testExitCode
}

# Show additional information
Write-Host ""
Write-ColoredOutput "💡 Additional Commands:" $Cyan
Write-Host "  Run specific tests:     .\scripts\run-tests.ps1 -Pattern 'TestDatabase.*'"
Write-Host "  Generate coverage:      .\scripts\run-tests.ps1 -Coverage"
Write-Host "  Verbose output:         .\scripts\run-tests.ps1 -Verbose"
Write-Host "  Custom timeout:         .\scripts\run-tests.ps1 -Timeout 120"
Write-Host "  Custom config:          .\scripts\run-tests.ps1 -ConfigFile 'custom_config.toml'"