# Health Check Package

The `health` package provides health check and diagnostic utilities for building robust monitoring and repair systems. It includes standardized exit codes, health status reporting, and automated repair operations.

## Features

- **Standardized Exit Codes**: 0 (OK), 1 (Warning), 2 (Error)
- **Health Status Tracking**: Monitor component health (Healthy/Warning/Error)
- **Checker Interface**: Extensible interface for custom health checks
- **Repairer Interface**: Automated repair operations for failed components
- **Diagnostic Reports**: Detailed health check reports with timestamps
- **Status Formatting**: Human-readable output for diagnostics

## Installation

```bash
go get github.com/julianstephens/go-utils/health
```

## Usage

### Basic Health Check

```go
type DatabaseChecker struct{}

func (d *DatabaseChecker) Check() health.Check {
    if err := pingDatabase(); err != nil {
        return health.NewCheckWithError("database", health.StatusError, "Failed to connect", err)
    }
    return health.NewCheck("database", health.StatusHealthy, "Connected")
}

func (d *DatabaseChecker) Name() string { return "database" }

// Run the checker
checkers := []health.Checker{&DatabaseChecker{}}
report := health.RunChecks(checkers...)
fmt.Println(report)
```

### Multiple Checks

```go
checkers := []health.Checker{
    &DatabaseChecker{},
    &CacheChecker{},
    &DiskSpaceChecker{},
}

report := health.RunChecks(checkers...)
fmt.Printf("Status: %v\n", report.ExitCode)
fmt.Printf("Summary: %s\n", report.Summary())
```

### Repair Operations

```go
repairers := map[string]health.Repairer{
    "cache": &CacheRepairer{},
    "connection": &ConnectionRepairer{},
}

// Attempt to repair failed checks
report = health.RepairAll(report, repairers)

// Check if repairs were successful
if report.ExitCode == health.ExitOK {
    fmt.Println("All systems repaired and healthy")
}
```

### Custom Checker Implementation

```go
type ServiceChecker struct {
    name    string
    timeout time.Duration
}

func (s *ServiceChecker) Check() health.Check {
    ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
    defer cancel()
    
    if err := s.checkEndpoint(ctx); err != nil {
        return health.NewCheckWithError(
            s.name,
            health.StatusError,
            fmt.Sprintf("Service unreachable: %v", err),
            err,
        )
    }
    
    return health.NewCheck(s.name, health.StatusHealthy, "Service responding")
}

func (s *ServiceChecker) Name() string { return s.name }
```

### Custom Repairer Implementation

```go
type ServiceRepairer struct {
    name string
}

func (s *ServiceRepairer) Repair() bool {
    if err := restartService(s.name); err != nil {
        return false
    }
    
    // Verify service is healthy after restart
    time.Sleep(time.Second)
    return isServiceHealthy(s.name)
}

func (s *ServiceRepairer) Name() string { return s.name }
```

## Exit Codes

The health package uses standard exit codes:

- `ExitOK` (0): All checks passed, system healthy
- `ExitWarning` (1): One or more warnings detected, operation may be degraded
- `ExitError` (2): One or more critical errors detected

## Status Values

- `StatusHealthy`: Component is operating normally
- `StatusWarning`: Component has degraded performance or minor issues
- `StatusError`: Component has critical failure

## API Reference

### Types

- `Check`: Single health check result with name, status, and optional error
- `Report`: Aggregated report of all health checks with exit code and timestamp
- `Status`: Enum for health status (Healthy/Warning/Error)
- `ExitCode`: Enum for exit codes (0/1/2)

### Interfaces

- `Checker`: Interface for performing health checks
  - `Check() Check`: Execute the health check
  - `Name() string`: Return checker name

- `Repairer`: Interface for repair operations
  - `Repair() bool`: Attempt repair, return success
  - `Name() string`: Return repairer name

### Functions

- `RunChecks(checkers ...Checker) Report`: Execute all checkers and return report
- `RepairAll(report Report, repairers map[string]Repairer) Report`: Attempt repairs on failed checks
- `NewCheck(name, message) Check`: Create a healthy check
- `NewCheckWithError(name, status, message, error) Check`: Create a check with error details

### Report Methods

- `Summary() string`: Get concise summary (e.g., "2 healthy, 1 warning, 0 error")
- `String() string`: Get formatted report with all details

## Integration

Works well with other go-utils packages:

- **logger**: Log health check results with structured fields
- **cliutil**: Display health reports in CLI with colored output
- **tests**: Use test helpers for verifying health check behavior

## Best Practices

- **Consistent naming**: Use descriptive, consistent names for checks (e.g., "database", "cache", "api_server")
- **Timeout context**: Always use context timeouts for external checks (database, HTTP, network)
- **Quick checks first**: Order checkers from fastest to slowest for faster overall reports
- **Clear messages**: Include actionable messages in status descriptions
- **Safe repairs**: Design repairs to be idempotent (safe to run multiple times)
- **Verify after repair**: Confirm the component is actually healthy after repair
- **Log repairs**: Log all repair attempts and results for audit trails

## Security Considerations

- **Error details**: Be cautious about exposing internal error details in reports
- **Credentials**: Never include credentials or secrets in check messages
- **External access**: Validate and sanitize endpoints before checking them
- **Repair permissions**: Ensure repair operations have minimal required permissions

## Performance

- Report generation: O(n) where n = number of checks
- Repair operations: Depends on individual repair implementations
- Concurrent checks: Implement using goroutines if needed
