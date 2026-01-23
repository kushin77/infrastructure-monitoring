package slo

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// SLI represents a Service Level Indicator
type SLI struct {
	Name    string
	Value   float64
	Samples int
	Window  time.Duration
}

// SLO represents a Service Level Objective
type SLO struct {
	Name        string
	Target      float64       // Target percentage (e.g., 99.95)
	Window      time.Duration // Measurement window (e.g., 30 days)
	Metric      string        // Prometheus metric name
	Filter      string        // PromQL filter expression
	Percentile  int           // For latency SLOs (e.g., 99)
	Threshold   float64       // For latency SLOs (e.g., 100ms)
	Description string
}

// ErrorBudget tracks error budget consumption
type ErrorBudget struct {
	SLO              SLO
	TotalBudget      float64 // Total allowed errors (in minutes)
	BudgetConsumed   float64 // Consumed budget (in minutes)
	BudgetRemaining  float64 // Remaining budget percentage
	BurnRateFast     float64 // Fast burn rate (1-hour window)
	BurnRateSlow     float64 // Slow burn rate (30-day window)
	LastUpdated      time.Time
	ComplianceStatus ComplianceStatus
	mu               sync.RWMutex
}

// ComplianceStatus indicates SLO compliance state
type ComplianceStatus string

const (
	StatusCompliant ComplianceStatus = "compliant"
	StatusAtRisk    ComplianceStatus = "at_risk"
	StatusViolated  ComplianceStatus = "violated"
	StatusUnknown   ComplianceStatus = "unknown"
)

// SLOTracker manages multiple SLOs and error budgets
type SLOTracker struct {
	slos    map[string]*SLO
	budgets map[string]*ErrorBudget
	metrics map[string]MetricsCollector
	mu      sync.RWMutex
}

// MetricsCollector collects metrics for SLO calculation
type MetricsCollector interface {
	// GetErrorRate returns current error rate (0-1)
	GetErrorRate() float64

	// GetLatencyPercentile returns latency at given percentile (ms)
	GetLatencyPercentile(percentile int) float64

	// GetThroughput returns requests per second
	GetThroughput() float64

	// GetCacheHitRate returns cache hit rate (0-1)
	GetCacheHitRate() float64
}

// NewSLOTracker creates a new SLO tracker
func NewSLOTracker() *SLOTracker {
	return &SLOTracker{
		slos:    make(map[string]*SLO),
		budgets: make(map[string]*ErrorBudget),
		metrics: make(map[string]MetricsCollector),
	}
}

// RegisterSLO registers a new SLO
func (st *SLOTracker) RegisterSLO(slo SLO, collector MetricsCollector) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	if slo.Name == "" {
		return fmt.Errorf("SLO name is required")
	}
	if slo.Target <= 0 || slo.Target > 100 {
		return fmt.Errorf("SLO target must be between 0 and 100")
	}
	if slo.Window <= 0 {
		return fmt.Errorf("SLO window must be positive")
	}

	st.slos[slo.Name] = &slo
	st.metrics[slo.Name] = collector

	// Initialize error budget
	budget := st.calculateErrorBudget(slo)
	st.budgets[slo.Name] = budget

	return nil
}

// calculateErrorBudget calculates initial error budget for an SLO
func (st *SLOTracker) calculateErrorBudget(slo SLO) *ErrorBudget {
	// Convert window to minutes
	windowMinutes := slo.Window.Minutes()

	// Total budget = window * (1 - target/100)
	// For 99.95% SLO over 30 days: 30 * 24 * 60 * (1 - 0.9995) = 21.6 minutes
	totalBudget := windowMinutes * (1 - (slo.Target / 100))

	return &ErrorBudget{
		SLO:              slo,
		TotalBudget:      totalBudget,
		LastUpdated:      time.Now(),
		ComplianceStatus: StatusUnknown,
	}
}

// UpdateErrorBudget updates error budget based on current metrics
func (st *SLOTracker) UpdateErrorBudget(sloName string, errorCount, totalCount int64) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	slo, exists := st.slos[sloName]
	if !exists {
		return fmt.Errorf("SLO not found: %s", sloName)
	}

	budget, exists := st.budgets[sloName]
	if !exists {
		return fmt.Errorf("error budget not found: %s", sloName)
	}

	// Calculate current compliance
	if totalCount == 0 {
		budget.ComplianceStatus = StatusUnknown
		budget.LastUpdated = time.Now()
		return nil
	}

	actualCompliance := float64((totalCount-errorCount)*100) / float64(totalCount)

	// Calculate budget consumed
	windowMinutes := slo.Window.Minutes()
	errorRate := float64(errorCount) / float64(totalCount)
	consumedMinutes := errorRate * windowMinutes

	budget.BudgetConsumed = consumedMinutes
	budget.BudgetRemaining = ((budget.TotalBudget - consumedMinutes) / budget.TotalBudget) * 100

	// Update compliance status
	if actualCompliance >= slo.Target {
		budget.ComplianceStatus = StatusCompliant
	} else if budget.BudgetRemaining > 0 {
		budget.ComplianceStatus = StatusAtRisk
	} else {
		budget.ComplianceStatus = StatusViolated
	}

	budget.LastUpdated = time.Now()
	return nil
}

// CalculateBurnRate calculates the burn rate for an error budget
func (st *SLOTracker) CalculateBurnRate(sloName string, window time.Duration, errorCount, totalCount int64) (float64, error) {
	st.mu.RLock()
	slo, exists := st.slos[sloName]
	st.mu.RUnlock()

	if !exists {
		return 0, fmt.Errorf("SLO not found: %s", sloName)
	}

	if totalCount == 0 {
		return 0, nil
	}

	// Burn rate = actual error rate / allowed error rate
	actualErrorRate := float64(errorCount) / float64(totalCount)
	allowedErrorRate := 1 - (slo.Target / 100)

	if allowedErrorRate == 0 {
		return math.Inf(1), nil
	}

	burnRate := actualErrorRate / allowedErrorRate
	return burnRate, nil
}

// GetErrorBudget returns the current error budget for an SLO
func (st *SLOTracker) GetErrorBudget(sloName string) (*ErrorBudget, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	budget, exists := st.budgets[sloName]
	if !exists {
		return nil, fmt.Errorf("error budget not found: %s", sloName)
	}

	return budget, nil
}

// GetAllErrorBudgets returns all error budgets
func (st *SLOTracker) GetAllErrorBudgets() map[string]*ErrorBudget {
	st.mu.RLock()
	defer st.mu.RUnlock()

	result := make(map[string]*ErrorBudget)
	for name, budget := range st.budgets {
		result[name] = budget
	}
	return result
}

// CheckCompliance checks SLO compliance status
func (st *SLOTracker) CheckCompliance(sloName string) (ComplianceStatus, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	budget, exists := st.budgets[sloName]
	if !exists {
		return StatusUnknown, fmt.Errorf("error budget not found: %s", sloName)
	}

	return budget.ComplianceStatus, nil
}

// GetComplianceSummary returns a summary of all SLO compliance
type ComplianceSummary struct {
	TotalSLOs      int
	CompliantSLOs  int
	AtRiskSLOs     int
	ViolatedSLOs   int
	UnknownSLOs    int
	ComplianceRate float64
	ErrorBudgetAvg float64
	WorstBudgetSLO string
	WorstBudgetPct float64
}

func (st *SLOTracker) GetComplianceSummary() ComplianceSummary {
	st.mu.RLock()
	defer st.mu.RUnlock()

	summary := ComplianceSummary{
		TotalSLOs:      len(st.slos),
		WorstBudgetPct: 100,
	}

	totalBudget := 0.0

	for _, budget := range st.budgets {
		totalBudget += budget.BudgetRemaining

		switch budget.ComplianceStatus {
		case StatusCompliant:
			summary.CompliantSLOs++
		case StatusAtRisk:
			summary.AtRiskSLOs++
		case StatusViolated:
			summary.ViolatedSLOs++
		case StatusUnknown:
			summary.UnknownSLOs++
		}

		if budget.BudgetRemaining < summary.WorstBudgetPct {
			summary.WorstBudgetPct = budget.BudgetRemaining
			summary.WorstBudgetSLO = budget.SLO.Name
		}
	}

	if summary.TotalSLOs > 0 {
		summary.ComplianceRate = float64(summary.CompliantSLOs) / float64(summary.TotalSLOs) * 100
		summary.ErrorBudgetAvg = totalBudget / float64(summary.TotalSLOs)
	}

	return summary
}

// StandardSLOs returns a set of recommended SLOs
func StandardSLOs() []SLO {
	return []SLO{
		{
			Name:        "API Availability",
			Target:      99.95,
			Window:      30 * 24 * time.Hour,
			Metric:      "http_requests_total",
			Filter:      "job='api'",
			Description: "API endpoint availability over 30 days",
		},
		{
			Name:        "API Latency P99",
			Target:      99.0,
			Window:      7 * 24 * time.Hour,
			Metric:      "http_request_duration_ms",
			Percentile:  99,
			Threshold:   200,
			Description: "API 99th percentile latency < 200ms over 7 days",
		},
		{
			Name:        "Database Latency P99",
			Target:      99.0,
			Window:      7 * 24 * time.Hour,
			Metric:      "db_query_duration_ms",
			Percentile:  99,
			Threshold:   100,
			Description: "Database query P99 latency < 100ms over 7 days",
		},
		{
			Name:        "Cache Hit Rate",
			Target:      85.0,
			Window:      24 * time.Hour,
			Metric:      "cache_hit_rate",
			Description: "Cache hit rate > 85% over 24 hours",
		},
		{
			Name:        "Error Rate",
			Target:      99.9,
			Window:      24 * time.Hour,
			Metric:      "errors_total",
			Description: "Error rate < 0.1% over 24 hours",
		},
	}
}

// AlertCondition represents an alert condition
type AlertCondition struct {
	Name        string
	Severity    string // critical, warning, info
	Threshold   float64
	Description string
	RunbookURL  string
}

// GetAlertConditions returns alert conditions for error budget
func GetAlertConditions() []AlertCondition {
	return []AlertCondition{
		{
			Name:        "SLOViolation",
			Severity:    "critical",
			Threshold:   0,
			Description: "SLO target not met",
			RunbookURL:  "https://wiki.elevatediq.com/runbooks/slo-violation",
		},
		{
			Name:        "ErrorBudgetBurnRateFast",
			Severity:    "critical",
			Threshold:   5,
			Description: "Error budget burning faster than expected (1h window)",
			RunbookURL:  "https://wiki.elevatediq.com/runbooks/high-burn-rate",
		},
		{
			Name:        "ErrorBudgetBurnRateSlow",
			Severity:    "warning",
			Threshold:   1,
			Description: "Error budget burning at expected rate (30d window)",
			RunbookURL:  "https://wiki.elevatediq.com/runbooks/slow-burn-rate",
		},
		{
			Name:        "ErrorBudgetRunningLow",
			Severity:    "warning",
			Threshold:   10,
			Description: "Error budget remaining < 10%",
			RunbookURL:  "https://wiki.elevatediq.com/runbooks/low-budget",
		},
	}
}
