package gopray

import (
	"fmt"
	"time"

	calc "github.com/mnadev/adhango/pkg/calc"
	data "github.com/mnadev/adhango/pkg/data"
	util "github.com/mnadev/adhango/pkg/util"
)

// App represents the main application with configuration and prayer times
type App struct {
	Config      *Config
	PrayerTimes *calc.PrayerTimes
	Localizer   *Localization
}

// NewApp creates and initializes a new App instance
func NewApp() (*App, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize localizer based on config language
	localizer := GetLocalizer(config.Language)

	app := &App{
		Config:    config,
		Localizer: localizer,
	}

	if err := app.calculatePrayers(); err != nil {
		return nil, fmt.Errorf("failed to calculate prayers: %w", err)
	}

	return app, nil
}

// calculatePrayers computes prayer times based on current configuration
func (a *App) calculatePrayers() error {
	now := time.Now()
	date := data.NewDateComponents(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC))

	// Get calculation method
	calcMethod, ok := calculationMethods[a.Config.Method]
	if !ok {
		return fmt.Errorf("unknown calculation method: %s", a.Config.Method)
	}

	// Get madhab
	madhab, ok := madhabMethods[a.Config.Madhab]
	if !ok {
		return fmt.Errorf("unknown madhab: %s", a.Config.Madhab)
	}

	// Get calculation parameters
	params := calc.GetMethodParameters(calcMethod)
	params.Madhab = madhab

	// Create coordinates
	coords, err := util.NewCoordinates(a.Config.Latitude, a.Config.Longitude)
	if err != nil {
		return fmt.Errorf("invalid coordinates: %w", err)
	}

	// Calculate prayer times
	a.PrayerTimes, err = calc.NewPrayerTimes(coords, date, params)
	if err != nil {
		return fmt.Errorf("failed to calculate prayer times: %w", err)
	}

	// Set timezone
	if err = a.PrayerTimes.SetTimeZone(a.Config.TimeZone); err != nil {
		return fmt.Errorf("invalid timezone %s: %w", a.Config.TimeZone, err)
	}

	return nil
}

// TimeLeftForNextPrayer returns duration until next prayer
func (a *App) TimeLeftForNextPrayer() time.Duration {
	nextPrayer := a.PrayerTimes.NextPrayerNow()
	nextPrayerTime := a.PrayerTimes.TimeForPrayer(nextPrayer)

	// Handle case after Isha prayer (next prayer is tomorrow's Fajr)
	if nextPrayer == 0 {
		nextPrayer = 1 // Fajr
		nextPrayerTime = a.PrayerTimes.TimeForPrayer(nextPrayer).Add(24 * time.Hour)
	}

	timeLeft := time.Until(nextPrayerTime)
	return timeLeft.Round(time.Second)
}

// GetNextPrayerName returns the name of the next prayer
func (a *App) GetNextPrayerName() string {
	nextPrayer := a.PrayerTimes.NextPrayerNow()
	prayerNames := []string{"none", "fajr", "dhuhr", "asr", "maghrib", "isha"}

	if nextPrayer == 0 {
		return a.Localizer.GetPrayerName("fajr") // After Isha, next is tomorrow's Fajr
	}

	if int(nextPrayer) < len(prayerNames) {
		return a.Localizer.GetPrayerName(prayerNames[nextPrayer])
	}

	return "Unknown"
}

// ListAllPrayers displays all prayer times formatted nicely
func (a *App) ListAllPrayers() {
	now := time.Now()
	if a.Config != nil && a.Config.TimeZone != "" {
		if location, err := time.LoadLocation(a.Config.TimeZone); err == nil {
			now = now.In(location)
		}
	}

	dateStr := a.Localizer.FormatGregorianDate(now)
	fmt.Printf("📅 %s %s\n\n", a.Localizer.Get("prayer_times_for"), dateStr)

	fmt.Printf("%s %s:    %s\n", GetPrayerEmoji("fajr"), a.Localizer.GetPrayerName("fajr"), a.PrayerTimes.Fajr.Format("15:04"))
	fmt.Printf("%s %s: %s\n", GetPrayerEmoji("sunrise"), a.Localizer.GetPrayerName("sunrise"), a.PrayerTimes.Sunrise.Format("15:04"))
	fmt.Printf("%s %s:   %s\n", GetPrayerEmoji("dhuhr"), a.Localizer.GetPrayerName("dhuhr"), a.PrayerTimes.Dhuhr.Format("15:04"))
	fmt.Printf("%s %s:     %s\n", GetPrayerEmoji("asr"), a.Localizer.GetPrayerName("asr"), a.PrayerTimes.Asr.Format("15:04"))
	fmt.Printf("%s %s: %s\n", GetPrayerEmoji("maghrib"), a.Localizer.GetPrayerName("maghrib"), a.PrayerTimes.Maghrib.Format("15:04"))
	fmt.Printf("%s %s:    %s\n", GetPrayerEmoji("isha"), a.Localizer.GetPrayerName("isha"), a.PrayerTimes.Isha.Format("15:04"))
}

// HijriDate returns the current Hijri date with proper formatting and month names
func (a *App) HijriDate() (string, error) {
	now := time.Now()

	// Get the local time in the user's configured timezone
	if a.Config != nil && a.Config.TimeZone != "" {
		location, err := time.LoadLocation(a.Config.TimeZone)
		if err == nil {
			now = now.In(location)
		}
	}

	hijriDate := ToHijriDate(now)
	if a.Localizer.Language == "ar" {
		return hijriDate.Format("arabic"), nil
	}
	return hijriDate.Format("simple"), nil
}

// GetDetailedHijriDate returns detailed Hijri date information
func (a *App) GetDetailedHijriDate() (map[string]any, error) {
	now := time.Now()

	// Get the local time in the user's configured timezone
	if a.Config != nil && a.Config.TimeZone != "" {
		location, err := time.LoadLocation(a.Config.TimeZone)
		if err == nil {
			now = now.In(location)
		}
	}

	hijriDate := ToHijriDate(now)

	return map[string]any{
		"day":           hijriDate.Day,
		"month":         hijriDate.Month,
		"year":          hijriDate.Year,
		"weekday":       hijriDate.WeekDay,
		"dayName":       hijriDate.DayName,
		"monthNameEn":   hijriDate.MonthName,
		"monthNameAr":   hijriDate.MonthNameAr,
		"formatted":     hijriDate.Format("simple"),
		"formattedFull": hijriDate.Format("full"),
		"formattedAr":   hijriDate.Format("arabic"),
		"numeric":       hijriDate.Format("numeric"),
	}, nil
} // ShowConfig displays current configuration
func (a *App) ShowConfig() {
	fmt.Printf("📋 %s:\n\n", a.Localizer.Get("current_configuration"))
	fmt.Printf("📍 %s: %.4f°, %.4f°\n", a.Localizer.Get("location"), a.Config.Latitude, a.Config.Longitude)
	fmt.Printf("🕰️  %s: %s\n", a.Localizer.Get("timezone"), a.Config.TimeZone)
	fmt.Printf("📐 %s: %s\n", a.Localizer.Get("method"), a.Config.Method)
	fmt.Printf("⚖️  %s: %s\n", a.Localizer.Get("madhab"), a.Config.Madhab)

	languageDisplay := "English"
	if a.Config.Language == "ar" {
		languageDisplay = "العربية (Arabic)"
	}
	fmt.Printf("🌐 %s: %s\n", a.Localizer.Get("language"), languageDisplay)

	fmt.Printf("\n📁 %s: %s\n", a.Localizer.Get("config_file"), getConfigPath())
}

// RefreshPrayerTimes recalculates prayer times (useful if config changed)
func (a *App) RefreshPrayerTimes() error {
	return a.calculatePrayers()
}
