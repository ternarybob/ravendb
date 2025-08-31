# RavenDB Test Environment Manager
# PowerShell script to manage different test configurations and environments

param(
    [Parameter(Position=0)]
    [ValidateSet("list", "create", "test", "setup")]
    [string]$Action = "list",
    
    [Parameter(Position=1)]
    [string]$Environment = "test",
    
    [string]$RavenDBUrl = "http://localhost:8080",
    [string]$DatabaseName,
    [switch]$Coverage = $false,
    [string]$Pattern = "*"
)

$Green = [System.ConsoleColor]::Green
$Yellow = [System.ConsoleColor]::Yellow
$Red = [System.ConsoleColor]::Red
$Cyan = [System.ConsoleColor]::Cyan
$White = [System.ConsoleColor]::White

function Write-ColoredOutput {
    param([string]$Message, [System.ConsoleColor]$Color = $White)
    Write-Host $Message -ForegroundColor $Color
}

function Write-Header {
    param([string]$Title)
    Write-Host ""
    Write-ColoredOutput "🔧 $Title" $Cyan
    Write-ColoredOutput ("=" * ($Title.Length + 3)) $Cyan
}

# Change to project root
Set-Location $PSScriptRoot\..

# Available configurations
$availableConfigs = @{
    "test" = @{
        "file" = "config\test_config.toml"
        "description" = "Standard test configuration"
        "database" = "RavenDBLibTestDB"
    }
    "local" = @{
        "file" = "config\local_config.toml"
        "description" = "Local development environment"
        "database" = "RavenDBLib_Local"
    }
    "docker" = @{
        "file" = "config\docker_config.toml"
        "description" = "Docker containerized environment"
        "database" = "RavenDBLib_Docker"
    }
    "ci" = @{
        "file" = "config\ci_config.toml"
        "description" = "Continuous Integration environment"
        "database" = "RavenDBLib_CI"
    }
}

switch ($Action) {
    "list" {
        Write-Header "Available Test Environments"
        foreach ($env in $availableConfigs.Keys) {
            $config = $availableConfigs[$env]
            $exists = Test-Path $config.file
            $status = if ($exists) { "✅" } else { "❌" }
            
            Write-ColoredOutput "$status $env" $(if ($exists) { $Green } else { $Red })
            Write-Host "   📄 File: $($config.file)"
            Write-Host "   📝 Description: $($config.description)"
            Write-Host "   🗄️  Database: $($config.database)"
            Write-Host ""
        }
        
        Write-ColoredOutput "Usage Examples:" $Cyan
        Write-Host "  List environments:     .\scripts\test-env-manager.ps1 list"
        Write-Host "  Create environment:    .\scripts\test-env-manager.ps1 create local"
        Write-Host "  Setup environment:     .\scripts\test-env-manager.ps1 setup test"
        Write-Host "  Run tests:             .\scripts\test-env-manager.ps1 test docker"
    }
    
    "create" {
        if (-not $availableConfigs.ContainsKey($Environment)) {
            Write-ColoredOutput "❌ Unknown environment: $Environment" $Red
            Write-Host "Available environments: $($availableConfigs.Keys -join ', ')"
            exit 1
        }
        
        $config = $availableConfigs[$Environment]
        Write-Header "Creating Configuration: $Environment"
        
        # Set database name if not provided
        if (-not $DatabaseName) {
            $DatabaseName = $config.database
        }
        
        # Create config content
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

        # Ensure config directory exists
        if (-not (Test-Path "config")) {
            New-Item -ItemType Directory -Path "config" | Out-Null
        }
        
        # Write configuration file
        $configContent | Out-File -FilePath $config.file -Encoding UTF8
        Write-ColoredOutput "✅ Created: $($config.file)" $Green
        
        # Display configuration
        Write-ColoredOutput "📋 Configuration:" $Yellow
        Get-Content $config.file | ForEach-Object { Write-Host "  $_" }
    }
    
    "setup" {
        if (-not $availableConfigs.ContainsKey($Environment)) {
            Write-ColoredOutput "❌ Unknown environment: $Environment" $Red
            exit 1
        }
        
        $config = $availableConfigs[$Environment]
        Write-Header "Setting up Environment: $Environment"
        
        # Create config if it doesn't exist
        if (-not (Test-Path $config.file)) {
            Write-ColoredOutput "📝 Configuration file not found, creating..." $Yellow
            & $PSCommandPath create $Environment -RavenDBUrl $RavenDBUrl -DatabaseName $config.database
        }
        
        # Run setup script
        & ".\scripts\setup-test-env.ps1" -ConfigFile $config.file
    }
    
    "test" {
        if (-not $availableConfigs.ContainsKey($Environment)) {
            Write-ColoredOutput "❌ Unknown environment: $Environment" $Red
            exit 1
        }
        
        $config = $availableConfigs[$Environment]
        Write-Header "Running Tests: $Environment"
        
        # Check if config exists
        if (-not (Test-Path $config.file)) {
            Write-ColoredOutput "❌ Configuration file not found: $($config.file)" $Red
            Write-Host "Run: .\scripts\test-env-manager.ps1 create $Environment"
            exit 1
        }
        
        # Run tests with the specified configuration
        $testArgs = @("-ConfigFile", $config.file, "-Pattern", $Pattern)
        if ($Coverage) {
            $testArgs += "-Coverage"
        }
        
        & ".\scripts\run-tests.ps1" @testArgs
    }
    
    default {
        Write-ColoredOutput "❌ Unknown action: $Action" $Red
        Write-Host "Available actions: list, create, setup, test"
        exit 1
    }
}