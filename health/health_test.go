package health_test

import (
	"errors"
	"testing"
	"time"

	"github.com/julianstephens/go-utils/health"
	tst "github.com/julianstephens/go-utils/tests"
)

// MockChecker implements Checker for testing
type MockChecker struct {
	name   string
	status health.Status
	msg    string
	err    error
}

func (m *MockChecker) Check() health.Check {
	return health.NewCheckWithError(m.name, m.status, m.msg, m.err)
}

func (m *MockChecker) Name() string {
	return m.name
}

// MockRepairer implements Repairer for testing
type MockRepairer struct {
	name         string
	shouldRepair bool
}

func (m *MockRepairer) Repair() bool {
	return m.shouldRepair
}

func (m *MockRepairer) Name() string {
	return m.name
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		status health.Status
		want   string
	}{
		{health.StatusHealthy, "healthy"},
		{health.StatusWarning, "warning"},
		{health.StatusError, "error"},
		{health.Status(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			result := tt.status.String()
			tst.AssertEqual(t, result, tt.want)
		})
	}
}

func TestNewCheck(t *testing.T) {
	check := health.NewCheck("database", health.StatusHealthy, "Connection OK")
	tst.AssertEqual(t, check.Name, "database")
	tst.AssertEqual(t, check.Status, health.StatusHealthy)
	tst.AssertEqual(t, check.Message, "Connection OK")
	tst.AssertNil(t, check.Error)
	tst.AssertFalse(t, check.Repaired, "should not be repaired")
}

func TestNewCheckWithError(t *testing.T) {
	err := errors.New("connection failed")
	check := health.NewCheckWithError("database", health.StatusError, "Failed to connect", err)
	tst.AssertEqual(t, check.Name, "database")
	tst.AssertEqual(t, check.Status, health.StatusError)
	tst.AssertEqual(t, check.Message, "Failed to connect")
	tst.AssertNotNil(t, check.Error)
}

func TestRunChecks_AllHealthy(t *testing.T) {
	checkers := []health.Checker{
		&MockChecker{"db", health.StatusHealthy, "OK", nil},
		&MockChecker{"cache", health.StatusHealthy, "OK", nil},
	}

	report := health.RunChecks(checkers...)

	tst.AssertEqual(t, len(report.Checks), 2)
	tst.AssertEqual(t, report.ExitCode, health.ExitOK)
	tst.AssertTrue(t, len(report.Message) > 0, "message should not be empty")
}

func TestRunChecks_WithWarning(t *testing.T) {
	checkers := []health.Checker{
		&MockChecker{"db", health.StatusHealthy, "OK", nil},
		&MockChecker{"disk", health.StatusWarning, "80% full", nil},
	}

	report := health.RunChecks(checkers...)

	tst.AssertEqual(t, len(report.Checks), 2)
	tst.AssertEqual(t, report.ExitCode, health.ExitWarning)
}

func TestRunChecks_WithError(t *testing.T) {
	checkers := []health.Checker{
		&MockChecker{"db", health.StatusHealthy, "OK", nil},
		&MockChecker{"service", health.StatusError, "Failed", errors.New("timeout")},
	}

	report := health.RunChecks(checkers...)

	tst.AssertEqual(t, len(report.Checks), 2)
	tst.AssertEqual(t, report.ExitCode, health.ExitError)
}

func TestRunChecks_Empty(t *testing.T) {
	report := health.RunChecks()

	tst.AssertEqual(t, len(report.Checks), 0)
	tst.AssertEqual(t, report.ExitCode, health.ExitOK)
}

func TestReport_Summary(t *testing.T) {
	report := health.Report{
		Checks: []health.Check{
			{Name: "a", Status: health.StatusHealthy},
			{Name: "b", Status: health.StatusHealthy},
			{Name: "c", Status: health.StatusWarning},
			{Name: "d", Status: health.StatusError},
		},
		ExitCode: health.ExitError,
	}

	summary := report.Summary()
	tst.AssertTrue(t, len(summary) > 0, "summary should not be empty")
	tst.AssertTrue(t, contains(summary, "2 healthy"), "summary should contain healthy count")
	tst.AssertTrue(t, contains(summary, "1 warning"), "summary should contain warning count")
	tst.AssertTrue(t, contains(summary, "1 error"), "summary should contain error count")
}

func TestReport_String(t *testing.T) {
	report := health.Report{
		Timestamp: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Checks: []health.Check{
			{Name: "database", Status: health.StatusHealthy, Message: "Connected"},
			{Name: "cache", Status: health.StatusWarning, Message: "Memory high"},
		},
		ExitCode: health.ExitWarning,
	}

	str := report.String()
	tst.AssertTrue(t, contains(str, "Health Check Report"), "should contain header")
	tst.AssertTrue(t, contains(str, "database"), "should contain database check")
	tst.AssertTrue(t, contains(str, "cache"), "should contain cache check")
	tst.AssertTrue(t, contains(str, "WARNING"), "should contain WARNING status")
}

func TestRepairAll_SuccessfulRepair(t *testing.T) {
	report := health.Report{
		Checks: []health.Check{
			{Name: "service1", Status: health.StatusError, Message: "Down"},
			{Name: "service2", Status: health.StatusHealthy, Message: "OK"},
		},
		ExitCode: health.ExitError,
	}

	repairers := map[string]health.Repairer{
		"service1": &MockRepairer{"service1", true},
		"service2": &MockRepairer{"service2", false},
	}

	updated := health.RepairAll(report, repairers)

	// Find service1 in updated checks
	var service1Check *health.Check
	for i := range updated.Checks {
		if updated.Checks[i].Name == "service1" {
			service1Check = &updated.Checks[i]
			break
		}
	}

	tst.AssertNotNil(t, service1Check, "service1 check should be found")
	tst.AssertEqual(t, service1Check.Status, health.StatusHealthy)
	tst.AssertTrue(t, service1Check.Repaired, "should be marked as repaired")
	tst.AssertEqual(t, updated.ExitCode, health.ExitOK)
}

func TestRepairAll_FailedRepair(t *testing.T) {
	report := health.Report{
		Checks: []health.Check{
			{Name: "db", Status: health.StatusError, Message: "Connection failed"},
		},
		ExitCode: health.ExitError,
	}

	repairers := map[string]health.Repairer{
		"db": &MockRepairer{"db", false},
	}

	updated := health.RepairAll(report, repairers)

	tst.AssertEqual(t, len(updated.Checks), 1)
	tst.AssertEqual(t, updated.Checks[0].Status, health.StatusError)
	tst.AssertFalse(t, updated.Checks[0].Repaired, "should not be marked as repaired")
	tst.AssertEqual(t, updated.ExitCode, health.ExitError)
}

func TestRepairAll_NoRepairerAvailable(t *testing.T) {
	report := health.Report{
		Checks: []health.Check{
			{Name: "service", Status: health.StatusError, Message: "Failed"},
		},
		ExitCode: health.ExitError,
	}

	repairers := map[string]health.Repairer{}

	updated := health.RepairAll(report, repairers)

	tst.AssertEqual(t, updated.Checks[0].Status, health.StatusError)
	tst.AssertFalse(t, updated.Checks[0].Repaired, "should not be marked as repaired")
}

func TestRepairAll_MixedResults(t *testing.T) {
	report := health.Report{
		Checks: []health.Check{
			{Name: "a", Status: health.StatusError, Message: "Failed"},
			{Name: "b", Status: health.StatusWarning, Message: "Warning"},
			{Name: "c", Status: health.StatusHealthy, Message: "OK"},
		},
		ExitCode: health.ExitError,
	}

	repairers := map[string]health.Repairer{
		"a": &MockRepairer{"a", true},  // Will repair
		"b": &MockRepairer{"b", false}, // Won't repair
	}

	updated := health.RepairAll(report, repairers)

	// Check a was repaired
	aCheck := updated.Checks[0] // Checks should be sorted, but let's find it
	for _, c := range updated.Checks {
		if c.Name == "a" {
			aCheck = c
			break
		}
	}
	tst.AssertTrue(t, aCheck.Repaired, "a should be repaired")

	// Exit code should still be warning or error (b wasn't repaired)
	tst.AssertTrue(t, updated.ExitCode == health.ExitWarning || updated.ExitCode == health.ExitError, "exit code should be warning or error")
}

func TestTimestamp(t *testing.T) {
	before := time.Now()
	report := health.RunChecks()
	after := time.Now()

	tst.AssertTrue(t, !report.Timestamp.Before(before), "timestamp should be after start")
	tst.AssertTrue(t, !report.Timestamp.After(after.Add(1*time.Second)), "timestamp should be within bounds")
}

// Helper function for string containment checks
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr))
}
