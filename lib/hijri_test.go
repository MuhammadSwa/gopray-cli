package gopray

import (
	"testing"
	"time"
)

func TestToHijriDate(t *testing.T) {
	// Test with a known date
	// September 27, 2025 should be 5 Rabi' al-Thani 1447
	testDate := time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC)
	hijriDate := ToHijriDate(testDate)

	if hijriDate.Day != 5 {
		t.Errorf("Expected day 5, got %d", hijriDate.Day)
	}

	if hijriDate.Month != 4 {
		t.Errorf("Expected month 4 (Rabi' al-Thani), got %d", hijriDate.Month)
	}

	if hijriDate.Year != 1447 {
		t.Errorf("Expected year 1447, got %d", hijriDate.Year)
	}

	if hijriDate.MonthName != "Rabi' al-Thani" {
		t.Errorf("Expected month name 'Rabi' al-Thani', got '%s'", hijriDate.MonthName)
	}
}

func TestHijriDateFormat(t *testing.T) {
	testDate := time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC)
	hijriDate := ToHijriDate(testDate)

	testCases := []struct {
		format   string
		expected string
	}{
		{"simple", "5 Rabi' al-Thani 1447 AH"},
		{"numeric", "05/04/1447"},
		{"arabic", "5 ربيع الثاني 1447 هـ"},
	}

	for _, tc := range testCases {
		result := hijriDate.Format(tc.format)
		if result != tc.expected {
			t.Errorf("Format %s: expected '%s', got '%s'", tc.format, tc.expected, result)
		}
	}
}

func TestSimpleHijriDate(t *testing.T) {
	testDate := time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC)
	result := SimpleHijriDate(testDate)

	// Should include weekday in the format
	if len(result) == 0 {
		t.Error("SimpleHijriDate returned empty string")
	}

	// Should contain the year 1447
	if !contains(result, "1447") {
		t.Errorf("Expected result to contain '1447', got '%s'", result)
	}
}

func TestRawHijriDate(t *testing.T) {
	testDate := time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC)
	raw := RawHijriDate(testDate)

	if len(raw) != 4 {
		t.Errorf("Expected 4 elements in raw date, got %d", len(raw))
	}

	// raw[1] should be day (5)
	if raw[1] != 5 {
		t.Errorf("Expected day 5, got %d", raw[1])
	}

	// raw[2] should be month (4)
	if raw[2] != 4 {
		t.Errorf("Expected month 4, got %d", raw[2])
	}

	// raw[3] should be year (1447)
	if raw[3] != 1447 {
		t.Errorf("Expected year 1447, got %d", raw[3])
	}
}

func TestHijriDateBoundaries(t *testing.T) {
	// Test with dates near the boundaries of the Umm al-Qura calendar
	testCases := []struct {
		name string
		date time.Time
	}{
		{"Recent date", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"Current era", time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC)},
		{"Future date", time.Date(2030, 12, 31, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hijriDate := ToHijriDate(tc.date)

			// Basic sanity checks
			if hijriDate.Day < 1 || hijriDate.Day > 30 {
				t.Errorf("%s: Invalid day %d", tc.name, hijriDate.Day)
			}

			if hijriDate.Month < 1 || hijriDate.Month > 12 {
				t.Errorf("%s: Invalid month %d", tc.name, hijriDate.Month)
			}

			if hijriDate.Year < 1400 || hijriDate.Year > 1500 {
				t.Errorf("%s: Year %d seems out of expected range", tc.name, hijriDate.Year)
			}

			// Check that month name is set
			if hijriDate.MonthName == "" {
				t.Errorf("%s: Month name is empty", tc.name)
			}
		})
	}
}

func TestAppHijriDate(t *testing.T) {
	// Create a test app (this might fail if config doesn't exist)
	app, err := NewApp()
	if err != nil {
		t.Skip("Skipping test - app initialization failed (likely no config)")
	}

	hijriDate, err := app.HijriDate()
	if err != nil {
		t.Errorf("App.HijriDate() failed: %v", err)
	}

	if hijriDate == "" {
		t.Error("App.HijriDate() returned empty string")
	}

	// Should contain "AH"
	if !contains(hijriDate, "AH") {
		t.Errorf("Expected result to contain 'AH', got '%s'", hijriDate)
	}
}

func TestAppGetDetailedHijriDate(t *testing.T) {
	// Create a test app (this might fail if config doesn't exist)
	app, err := NewApp()
	if err != nil {
		t.Skip("Skipping test - app initialization failed (likely no config)")
	}

	details, err := app.GetDetailedHijriDate()
	if err != nil {
		t.Errorf("App.GetDetailedHijriDate() failed: %v", err)
	}

	if details == nil {
		t.Error("App.GetDetailedHijriDate() returned nil")
	}

	// Check that required keys exist
	requiredKeys := []string{"day", "month", "year", "monthNameEn", "monthNameAr", "formatted"}
	for _, key := range requiredKeys {
		if _, exists := details[key]; !exists {
			t.Errorf("Expected key '%s' not found in details", key)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i < len(s)-len(substr)+1; i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}

// Benchmark tests
func BenchmarkToHijriDate(b *testing.B) {
	testDate := time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ToHijriDate(testDate)
	}
}

func BenchmarkHijriDateFormat(b *testing.B) {
	testDate := time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC)
	hijriDate := ToHijriDate(testDate)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hijriDate.Format("simple")
	}
}
