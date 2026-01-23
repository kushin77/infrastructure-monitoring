package slo

import (
	"math"
	"testing"
	"time"
)

// MockMetricsCollector implements MetricsCollector for testing
type MockMetricsCollector struct {
	errorRate    float64
	latencyP99   float64
	throughput   float64
	cacheHitRate float64
}

func (m *MockMetricsCollector) GetErrorRate() float64 {
	return m.errorRate
}

func (m *MockMetricsCollector) GetLatencyPercentile(percentile int) float64 {
	return m.latencyP99
}

func (m *MockMetricsCollector) GetThroughput() float64 {
	return m.throughput
}

func (m *MockMetricsCollector) GetCacheHitRate() float64 {
	return m.cacheHitRate
}

func TestNewSLOTracker(t *testing.T) {
	tracker := NewSLOTracker()
	if tracker == nil {
		t.Fatal("expected non-nil SLO tracker")
	}
}

func TestRegisterSLO(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}

	collector := &MockMetricsCollector{}
	err := tracker.RegisterSLO(slo, collector)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify SLO was registered
	if err == nil {
		budget, err := tracker.GetErrorBudget("Test SLO")
		if err != nil {
			t.Fatalf("error getting budget: %v", err)
		}
		if budget == nil {
			t.Fatal("expected non-nil error budget")
		}
		if budget.SLO.Name != "Test SLO" {
			t.Errorf("expected SLO name 'Test SLO', got %q", budget.SLO.Name)
		}
	}
}

func TestRegisterSLO_InvalidName(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "",
		Target: 99.9,
		Window: 24 * time.Hour,
	}

	collector := &MockMetricsCollector{}
	err := tracker.RegisterSLO(slo, collector)
	if err == nil {
		t.Fatal("expected error for empty SLO name")
	}
}

func TestRegisterSLO_InvalidTarget(t *testing.T) {
	tests := []struct {
		name   string
		target float64
	}{
		{"negative target", -10},
		{"zero target", 0},
		{"over 100", 101},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracker := NewSLOTracker()
			slo := SLO{
				Name:   "Test SLO",
				Target: tt.target,
				Window: 24 * time.Hour,
			}

			collector := &MockMetricsCollector{}
			err := tracker.RegisterSLO(slo, collector)
			if err == nil {
				t.Fatal("expected error for invalid target")
			}
		})
	}
}

func TestRegisterSLO_InvalidWindow(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: -1 * time.Hour,
	}

	collector := &MockMetricsCollector{}
	err := tracker.RegisterSLO(slo, collector)
	if err == nil {
		t.Fatal("expected error for invalid window")
	}
}

func TestCalculateErrorBudget(t *testing.T) {
	tracker := NewSLOTracker()

	tests := []struct {
		name           string
		target         float64
		window         time.Duration
		expectedBudget float64
	}{
		{
			name:           "99.95% SLO over 30 days",
			target:         99.95,
			window:         30 * 24 * time.Hour,
			expectedBudget: 21.6, // 30*24*60 * (1 - 0.9995)
		},
		{
			name:           "99.9% SLO over 7 days",
			target:         99.9,
			window:         7 * 24 * time.Hour,
			expectedBudget: 10.08, // 7*24*60 * (1 - 0.999)
		},
		{
			name:           "99% SLO over 24 hours",
			target:         99.0,
			window:         24 * time.Hour,
			expectedBudget: 14.4, // 24*60 * (1 - 0.99)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slo := SLO{
				Name:   "Test",
				Target: tt.target,
				Window: tt.window,
			}

			budget := tracker.calculateErrorBudget(slo)

			// Allow small floating point error
			epsilon := 0.01
			if math.Abs(budget.TotalBudget-tt.expectedBudget) > epsilon {
				t.Errorf("expected budget %f, got %f", tt.expectedBudget, budget.TotalBudget)
			}
		})
	}
}

func TestUpdateErrorBudget(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}

	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	// Test with 99.99% success rate (exceeds 99.9% target, compliant)
	err := tracker.UpdateErrorBudget("Test SLO", 1, 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	budget, _ := tracker.GetErrorBudget("Test SLO")
	if budget.ComplianceStatus != StatusCompliant {
		t.Errorf("expected StatusCompliant, got %v", budget.ComplianceStatus)
	}
}

func TestUpdateErrorBudget_ViolationDetection(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}

	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	// Test with 98% success rate (violates 99.9% target)
	err := tracker.UpdateErrorBudget("Test SLO", 200, 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	budget, _ := tracker.GetErrorBudget("Test SLO")
	if budget.ComplianceStatus != StatusViolated {
		t.Errorf("expected StatusViolated, got %v", budget.ComplianceStatus)
	}
}

func TestUpdateErrorBudget_ZeroTotal(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}

	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	err := tracker.UpdateErrorBudget("Test SLO", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	budget, _ := tracker.GetErrorBudget("Test SLO")
	if budget.ComplianceStatus != StatusUnknown {
		t.Errorf("expected StatusUnknown for zero totals, got %v", budget.ComplianceStatus)
	}
}

func TestUpdateErrorBudget_NonexistentSLO(t *testing.T) {
	tracker := NewSLOTracker()
	err := tracker.UpdateErrorBudget("Nonexistent", 10, 100)
	if err == nil {
		t.Fatal("expected error for nonexistent SLO")
	}
}

func TestCalculateBurnRate(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.0,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}

	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	// Test case: 1% error rate with 99% SLO
	// Allowed error rate = 1%, Actual = 1%, Burn rate = 1
	burnRate, err := tracker.CalculateBurnRate("Test SLO", 1*time.Hour, 100, 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBurnRate := 1.0
	epsilon := 0.01
	if math.Abs(burnRate-expectedBurnRate) > epsilon {
		t.Errorf("expected burn rate %f, got %f", expectedBurnRate, burnRate)
	}
}

func TestCalculateBurnRate_HighBurnRate(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.0,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}

	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	// Test case: 2% error rate with 99% SLO
	// Allowed error rate = 1%, Actual = 2%, Burn rate = 2
	burnRate, err := tracker.CalculateBurnRate("Test SLO", 1*time.Hour, 200, 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedBurnRate := 2.0
	epsilon := 0.01
	if math.Abs(burnRate-expectedBurnRate) > epsilon {
		t.Errorf("expected burn rate %f, got %f", expectedBurnRate, burnRate)
	}
}

func TestCalculateBurnRate_ZeroTotal(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.0,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}

	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	burnRate, err := tracker.CalculateBurnRate("Test SLO", 1*time.Hour, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if burnRate != 0 {
		t.Errorf("expected burn rate 0 for zero total, got %f", burnRate)
	}
}

func TestCalculateBurnRate_NonexistentSLO(t *testing.T) {
	tracker := NewSLOTracker()
	_, err := tracker.CalculateBurnRate("Nonexistent", 1*time.Hour, 10, 100)
	if err == nil {
		t.Fatal("expected error for nonexistent SLO")
	}
}

func TestGetErrorBudget(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}

	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	budget, err := tracker.GetErrorBudget("Test SLO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if budget == nil {
		t.Fatal("expected non-nil error budget")
	}
	if budget.SLO.Name != "Test SLO" {
		t.Errorf("expected SLO name 'Test SLO', got %q", budget.SLO.Name)
	}
}

func TestGetErrorBudget_NonexistentSLO(t *testing.T) {
	tracker := NewSLOTracker()
	_, err := tracker.GetErrorBudget("Nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent SLO")
	}
}

func TestGetAllErrorBudgets(t *testing.T) {
	tracker := NewSLOTracker()

	slos := []SLO{
		{Name: "SLO1", Target: 99.9, Window: 24 * time.Hour, Metric: "m1"},
		{Name: "SLO2", Target: 99.0, Window: 7 * 24 * time.Hour, Metric: "m2"},
	}

	collector := &MockMetricsCollector{}
	for _, slo := range slos {
		_ = tracker.RegisterSLO(slo, collector)
	}

	budgets := tracker.GetAllErrorBudgets()
	if len(budgets) != 2 {
		t.Errorf("expected 2 budgets, got %d", len(budgets))
	}
	if _, ok := budgets["SLO1"]; !ok {
		t.Fatal("expected SLO1 in budgets")
	}
	if _, ok := budgets["SLO2"]; !ok {
		t.Fatal("expected SLO2 in budgets")
	}
}

func TestCheckCompliance(t *testing.T) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}

	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	status, err := tracker.CheckCompliance("Test SLO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != StatusUnknown {
		t.Errorf("expected StatusUnknown initially, got %v", status)
	}
}

func TestCheckCompliance_NonexistentSLO(t *testing.T) {
	tracker := NewSLOTracker()
	_, err := tracker.CheckCompliance("Nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent SLO")
	}
}

func TestGetComplianceSummary_Empty(t *testing.T) {
	tracker := NewSLOTracker()
	summary := tracker.GetComplianceSummary()

	if summary.TotalSLOs != 0 {
		t.Errorf("expected 0 SLOs, got %d", summary.TotalSLOs)
	}
	if summary.ComplianceRate != 0 {
		t.Errorf("expected 0 compliance rate, got %f", summary.ComplianceRate)
	}
}

func TestGetComplianceSummary_MultipleSLOs(t *testing.T) {
	tracker := NewSLOTracker()

	slos := []SLO{
		{Name: "SLO1", Target: 99.9, Window: 24 * time.Hour, Metric: "m1"},
		{Name: "SLO2", Target: 99.0, Window: 7 * 24 * time.Hour, Metric: "m2"},
		{Name: "SLO3", Target: 99.95, Window: 30 * 24 * time.Hour, Metric: "m3"},
	}

	collector := &MockMetricsCollector{}
	for _, slo := range slos {
		_ = tracker.RegisterSLO(slo, collector)
	}

	// Update first two SLOs to compliant state (exceed their targets)
	tracker.UpdateErrorBudget("SLO1", 1, 10000) // 99.99% vs 99.9% target = compliant
	tracker.UpdateErrorBudget("SLO2", 5, 10000) // 99.95% vs 99.0% target = compliant

	summary := tracker.GetComplianceSummary()

	if summary.TotalSLOs != 3 {
		t.Errorf("expected 3 total SLOs, got %d", summary.TotalSLOs)
	}
	if summary.CompliantSLOs != 2 {
		t.Errorf("expected 2 compliant SLOs, got %d", summary.CompliantSLOs)
	}
	if summary.UnknownSLOs != 1 {
		t.Errorf("expected 1 unknown SLO, got %d", summary.UnknownSLOs)
	}
}

func TestStandardSLOs(t *testing.T) {
	slos := StandardSLOs()

	if len(slos) == 0 {
		t.Fatal("expected standard SLOs to be non-empty")
	}

	expectedCount := 5
	if len(slos) != expectedCount {
		t.Errorf("expected %d standard SLOs, got %d", expectedCount, len(slos))
	}

	// Verify required SLOs exist
	sloNames := make(map[string]bool)
	for _, slo := range slos {
		sloNames[slo.Name] = true
	}

	requiredSLOs := []string{"API Availability", "API Latency P99", "Database Latency P99", "Cache Hit Rate", "Error Rate"}
	for _, name := range requiredSLOs {
		if !sloNames[name] {
			t.Errorf("expected SLO %q not found", name)
		}
	}
}

func TestGetAlertConditions(t *testing.T) {
	conditions := GetAlertConditions()

	if len(conditions) == 0 {
		t.Fatal("expected alert conditions to be non-empty")
	}

	expectedCount := 4
	if len(conditions) != expectedCount {
		t.Errorf("expected %d alert conditions, got %d", expectedCount, len(conditions))
	}

	// Verify required alert conditions exist
	alertNames := make(map[string]bool)
	for _, alert := range conditions {
		alertNames[alert.Name] = true
	}

	requiredAlerts := []string{"SLOViolation", "ErrorBudgetBurnRateFast", "ErrorBudgetBurnRateSlow", "ErrorBudgetRunningLow"}
	for _, name := range requiredAlerts {
		if !alertNames[name] {
			t.Errorf("expected alert condition %q not found", name)
		}
	}

	// Verify runbook URLs are set
	for _, condition := range conditions {
		if condition.RunbookURL == "" {
			t.Errorf("expected runbook URL for %q", condition.Name)
		}
	}
}

func TestComplianceStatus_String(t *testing.T) {
	tests := []struct {
		status   ComplianceStatus
		expected string
	}{
		{StatusCompliant, "compliant"},
		{StatusAtRisk, "at_risk"},
		{StatusViolated, "violated"},
		{StatusUnknown, "unknown"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.expected {
			t.Errorf("expected status %q, got %q", tt.expected, string(tt.status))
		}
	}
}

// Benchmark tests

func BenchmarkRegisterSLO(b *testing.B) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}
	collector := &MockMetricsCollector{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracker.RegisterSLO(slo, collector)
	}
}

func BenchmarkUpdateErrorBudget(b *testing.B) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}
	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracker.UpdateErrorBudget("Test SLO", 10, 10000)
	}
}

func BenchmarkCalculateBurnRate(b *testing.B) {
	tracker := NewSLOTracker()
	slo := SLO{
		Name:   "Test SLO",
		Target: 99.9,
		Window: 24 * time.Hour,
		Metric: "test_metric",
	}
	collector := &MockMetricsCollector{}
	_ = tracker.RegisterSLO(slo, collector)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tracker.CalculateBurnRate("Test SLO", 1*time.Hour, 100, 10000)
	}
}

func BenchmarkGetComplianceSummary(b *testing.B) {
	tracker := NewSLOTracker()
	slos := StandardSLOs()
	collector := &MockMetricsCollector{}

	for _, slo := range slos {
		_ = tracker.RegisterSLO(slo, collector)
		_ = tracker.UpdateErrorBudget(slo.Name, 5, 10000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tracker.GetComplianceSummary()
	}
}
