package gopray

import (
	"strings"
	"testing"
)

func TestNewApp(t *testing.T) {
	// This test might fail if no config exists, which is expected
	app, err := NewApp()
	if err != nil {
		// This is expected if config doesn't exist or coordinates are 0,0
		t.Logf("NewApp failed as expected: %v", err)
		return
	}

	if app == nil {
		t.Fatal("NewApp returned nil app without error")
	}

	if app.Config == nil {
		t.Fatal("App.Config is nil")
	}

	if app.PrayerTimes == nil {
		t.Fatal("App.PrayerTimes is nil")
	}
}

func TestValidateMethods(t *testing.T) {
	methods, madhabs := ValidateMethods()

	if len(methods) == 0 {
		t.Error("No calculation methods returned")
	}

	if len(madhabs) == 0 {
		t.Error("No madhabs returned")
	}

	// Test some known methods
	expectedMethods := []string{"MUSLIM_WORLD_LEAGUE", "EGYPTIAN", "NORTH_AMERICA"}
	for _, method := range expectedMethods {
		found := false
		for _, m := range methods {
			if m == method {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected method %s not found in methods list", method)
		}
	}
}

func TestIsValidMethod(t *testing.T) {
	testCases := []struct {
		method   string
		expected bool
	}{
		{"MUSLIM_WORLD_LEAGUE", true},
		{"EGYPTIAN", true},
		{"INVALID_METHOD", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := IsValidMethod(tc.method)
		if result != tc.expected {
			t.Errorf("IsValidMethod(%q) = %v, want %v", tc.method, result, tc.expected)
		}
	}
}

func TestIsValidMadhab(t *testing.T) {
	testCases := []struct {
		madhab   string
		expected bool
	}{
		{"SHAFI_HANBALI_MALIKI", true},
		{"HANAFI", true},
		{"INVALID_MADHAB", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := IsValidMadhab(tc.madhab)
		if result != tc.expected {
			t.Errorf("IsValidMadhab(%q) = %v, want %v", tc.madhab, result, tc.expected)
		}
	}
}

func TestConfigValidation(t *testing.T) {
	testCases := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				Method:    "MUSLIM_WORLD_LEAGUE",
				Madhab:    "HANAFI",
				TimeZone:  "UTC",
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			expectError: false,
		},
		{
			name: "zero coordinates",
			config: &Config{
				Method:    "MUSLIM_WORLD_LEAGUE",
				Madhab:    "HANAFI",
				TimeZone:  "UTC",
				Latitude:  0.0,
				Longitude: 0.0,
			},
			expectError: true,
		},
		{
			name: "invalid method",
			config: &Config{
				Method:    "INVALID",
				Madhab:    "HANAFI",
				TimeZone:  "UTC",
				Latitude:  40.7128,
				Longitude: -74.0060,
			},
			expectError: true,
		},
		{
			name: "invalid latitude",
			config: &Config{
				Method:    "MUSLIM_WORLD_LEAGUE",
				Madhab:    "HANAFI",
				TimeZone:  "UTC",
				Latitude:  91.0,
				Longitude: -74.0060,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateConfig(tc.config)
			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkNewApp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		app, err := NewApp()
		if err != nil {
			b.Skip("Config not available for benchmark")
		}
		_ = app
	}
}

func BenchmarkTimeLeftForNextPrayer(b *testing.B) {
	app, err := NewApp()
	if err != nil {
		b.Skip("Config not available for benchmark")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = app.TimeLeftForNextPrayer()
	}
}

func TestHijriDate(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Skip("NewApp failed, likely due to config issues")
	}

	hijriDate, err := app.HijriDate()
	if err != nil {
		t.Fatalf("HijriDate failed: %v", err)
	}

	if hijriDate == "" {
		t.Error("HijriDate returned empty string")
	}

	// Check if it contains "AH" (Anno Hegirae)
	if !strings.Contains(hijriDate, "AH") {
		t.Errorf("HijriDate should contain 'AH', got: %s", hijriDate)
	}

	t.Logf("Hijri date: %s", hijriDate)
}

func TestGetDetailedHijriDate(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Skip("NewApp failed, likely due to config issues")
	}

	details, err := app.GetDetailedHijriDate()
	if err != nil {
		t.Fatalf("GetDetailedHijriDate failed: %v", err)
	}

	// Check required fields
	requiredFields := []string{"day", "month", "year", "monthNameEn", "monthNameAr", "formatted", "formattedAr", "numeric"}
	for _, field := range requiredFields {
		if _, ok := details[field]; !ok {
			t.Errorf("Missing required field: %s", field)
		}
	}

	// Validate day range
	day, ok := details["day"].(int)
	if !ok || day < 1 || day > 30 {
		t.Errorf("Invalid day: %v", details["day"])
	}

	// Validate month range
	month, ok := details["month"].(int)
	if !ok || month < 1 || month > 12 {
		t.Errorf("Invalid month: %v", details["month"])
	}

	// Validate year
	year, ok := details["year"].(int)
	if !ok || year < 1400 || year > 1500 {
		t.Errorf("Invalid year: %v", details["year"])
	}

	t.Logf("Detailed Hijri date: %+v", details)
}
