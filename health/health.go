package health

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ExitCode represents standard exit codes for health check operations
type ExitCode int

const (
	ExitOK      ExitCode = 0 // All systems operational
	ExitWarning ExitCode = 1 // Warning detected, operation may be degraded
	ExitError   ExitCode = 2 // Critical error detected
)

// Status represents the health status of a component
type Status int

const (
	StatusHealthy Status = iota
	StatusWarning
	StatusError
)

func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusWarning:
		return "warning"
	case StatusError:
		return "error"
	default:
		return "unknown"
	}
}

// Check represents a single health check with a name and status
type Check struct {
	Name     string // Name of the check (e.g., "database", "cache")
	Status   Status // Health status
	Message  string // Detailed message about the status
	Error    error  // Error if one occurred
	Repaired bool   // Whether this check was successfully repaired
}

// Report aggregates multiple health checks and provides diagnostic information
type Report struct {
	Timestamp time.Time
	Checks    []Check
	ExitCode  ExitCode
	Message   string
}

// Checker is an interface for performing health checks
type Checker interface {
	// Check performs a health check and returns the Check result
	Check() Check
	// Name returns the name of this checker
	Name() string
}

// Repairer is an interface for repair operations on failed components
type Repairer interface {
	// Repair attempts to fix the issue and returns true if successful
	Repair() bool
	// Name returns the name of this repairer
	Name() string
}

// RunChecks executes all provided checkers and returns a Report
func RunChecks(checkers ...Checker) Report {
	report := Report{
		Timestamp: time.Now(),
		Checks:    make([]Check, 0, len(checkers)),
		ExitCode:  ExitOK,
	}

	for _, checker := range checkers {
		check := checker.Check()
		report.Checks = append(report.Checks, check)

		if check.Status == StatusError && report.ExitCode < ExitError {
			report.ExitCode = ExitError
		} else if check.Status == StatusWarning && report.ExitCode < ExitWarning {
			report.ExitCode = ExitWarning
		}
	}

	report.Message = formatMessage(report)
	return report
}

// RepairAll attempts to repair all failing checks and returns the updated Report
func RepairAll(report Report, repairers map[string]Repairer) Report {
	for i := range report.Checks {
		check := &report.Checks[i]
		if check.Status == StatusError || check.Status == StatusWarning {
			if repairer, ok := repairers[check.Name]; ok {
				if repairer.Repair() {
					check.Repaired = true
					check.Status = StatusHealthy
					check.Message = "Repaired: " + check.Message
				}
			}
		}
	}

	// Recalculate exit code after repairs
	report.ExitCode = ExitOK
	for _, check := range report.Checks {
		if check.Status == StatusError && report.ExitCode < ExitError {
			report.ExitCode = ExitError
		} else if check.Status == StatusWarning && report.ExitCode < ExitWarning {
			report.ExitCode = ExitWarning
		}
	}

	report.Message = formatMessage(report)
	return report
}

// Summary returns a concise summary of the report
func (r Report) Summary() string {
	healthy := 0
	warnings := 0
	errors := 0

	for _, check := range r.Checks {
		switch check.Status {
		case StatusHealthy:
			healthy++
		case StatusWarning:
			warnings++
		case StatusError:
			errors++
		}
	}

	return fmt.Sprintf("%d healthy, %d warning, %d error", healthy, warnings, errors)
}

// String returns a formatted string representation of the report
func (r Report) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Health Check Report (%s)\n", r.Timestamp.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Status: %s\n", statusFromExitCode(r.ExitCode)))
	sb.WriteString(fmt.Sprintf("Summary: %s\n", r.Summary()))
	sb.WriteString("\nChecks:\n")

	// Sort checks by name for consistent output
	checks := make([]Check, len(r.Checks))
	copy(checks, r.Checks)
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].Name < checks[j].Name
	})

	for _, check := range checks {
		status := ""
		switch check.Status {
		case StatusHealthy:
			status = "✓ OK"
		case StatusWarning:
			status = "⚠ WARNING"
		case StatusError:
			status = "✗ ERROR"
		}

		sb.WriteString(fmt.Sprintf("  %s: %s\n", check.Name, status))
		if check.Message != "" {
			sb.WriteString(fmt.Sprintf("      %s\n", check.Message))
		}
		if check.Error != nil {
			sb.WriteString(fmt.Sprintf("      Error: %v\n", check.Error))
		}
		if check.Repaired {
			sb.WriteString("      [REPAIRED]\n")
		}
	}

	return sb.String()
}

// formatMessage creates a summary message based on the report
func formatMessage(r Report) string {
	var statuses []string
	statusMap := make(map[string]int)

	for _, check := range r.Checks {
		statusStr := check.Status.String()
		statusMap[statusStr]++
	}

	// Build message with counts
	if count, ok := statusMap["error"]; ok && count > 0 {
		statuses = append(statuses, fmt.Sprintf("%d error(s)", count))
	}
	if count, ok := statusMap["warning"]; ok && count > 0 {
		statuses = append(statuses, fmt.Sprintf("%d warning(s)", count))
	}
	if count, ok := statusMap["healthy"]; ok && count > 0 {
		statuses = append(statuses, fmt.Sprintf("%d healthy", count))
	}

	if len(statuses) == 0 {
		return "No checks performed"
	}

	return strings.Join(statuses, ", ")
}

// statusFromExitCode converts an exit code to a status string
func statusFromExitCode(code ExitCode) string {
	switch code {
	case ExitOK:
		return "OK"
	case ExitWarning:
		return "WARNING"
	case ExitError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// NewCheck creates a new Check with the given parameters
func NewCheck(name string, status Status, message string) Check {
	return Check{
		Name:    name,
		Status:  status,
		Message: message,
	}
}

// NewCheckWithError creates a new Check with an error
func NewCheckWithError(name string, status Status, message string, err error) Check {
	return Check{
		Name:    name,
		Status:  status,
		Message: message,
		Error:   err,
	}
}
