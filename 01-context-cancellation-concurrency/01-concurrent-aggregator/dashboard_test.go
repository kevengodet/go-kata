package main

import (
	"context"
	"testing"
	"time"
)

func TestDashboard(t *testing.T) {
	tests := []struct {
		name         string
		profileDelay time.Duration
		profileError string
		orderDelay   time.Duration
		orderError   string
		timeout      time.Duration
		expectError  bool
	}{
		// Timeout cases
		{"ProfileTimeout", 100 * time.Millisecond, "", 0, "", 50 * time.Millisecond, true},
		{"OrderTimeout", 0, "", 100 * time.Millisecond, "", 50 * time.Millisecond, true},
		{"BothTimeout", 100 * time.Millisecond, "", 100 * time.Millisecond, "", 50 * time.Millisecond, true},
		{"ProfileOk", 50 * time.Millisecond, "", 0, "", 100 * time.Millisecond, false},
		{"OrderOk", 0, "", 50 * time.Millisecond, "", 100 * time.Millisecond, false},
		{"BothOk", 0, "", 0, "", 100 * time.Millisecond, false},

		// Error cases
		{"ProfileError", 0, "profile failed", 200 * time.Millisecond, "", 100 * time.Millisecond, true},
		{"OrderError", 200 * time.Millisecond, "", 0, "order failed", 100 * time.Millisecond, true},
		{"BothError", 200 * time.Millisecond, "profile failed", 200 * time.Millisecond, "order failed", 100 * time.Millisecond, true},
		{"Success", 0, "", 0, "", 100 * time.Millisecond, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dashboard := New(
				mockProfile(
					MockTimeout(tt.profileDelay),
					MockError(tt.profileError),
				),
				mockOrder(
					MockTimeout(tt.orderDelay),
					MockError(tt.orderError),
				),
				WithTimeout(tt.timeout),
			)

			startTime := time.Now()
			s, err := dashboard.Aggregate(context.Background(), 1)
			elapsed := time.Since(startTime)

			// Add a small buffer for timing fluctuations
			if elapsed > tt.timeout+50*time.Millisecond {
				t.Errorf("Expected elapsed time to be less than %v, got %v", tt.timeout+50*time.Millisecond, elapsed)
			}

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil (response: %q)", s)
				}
				if s != "" {
					t.Errorf("Expected empty response on error, got %q", s)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if s == "" {
					t.Errorf("Expected non-empty response, got %q", s)
				}
			}
		})
	}
}
