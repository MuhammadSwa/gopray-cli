package gopray

import (
	"fmt"
	"time"
)

// Localization maps for different languages
type Localization struct {
	Language string
	Messages map[string]string
}

// GetLocalizer returns a localization instance for the given language
func GetLocalizer(language string) *Localization {
	if language == "ar" {
		return &Localization{
			Language: "ar",
			Messages: arabicMessages,
		}
	}
	return &Localization{
		Language: "en",
		Messages: englishMessages,
	}
}

var englishMessages = map[string]string{
	// Prayer names
	"prayer_fajr":    "Fajr",
	"prayer_sunrise": "Sunrise",
	"prayer_dhuhr":   "Dhuhr",
	"prayer_asr":     "Asr",
	"prayer_maghrib": "Maghrib",
	"prayer_isha":    "Isha",

	// UI Messages
	"prayer_times_for":       "Prayer Times for",
	"next_prayer":            "Next prayer",
	"in":                     "in",
	"current_configuration":  "Current Configuration",
	"location":               "Location",
	"timezone":               "Timezone",
	"method":                 "Method",
	"madhab":                 "Madhab",
	"language":               "Language",
	"config_file":            "Config file",
	"todays_date":            "Today's Date",
	"gregorian":              "Gregorian",
	"hijri":                  "Hijri",
	"detailed_date_info":     "Detailed Date Information",
	"arabic_format":          "Arabic Format",
	"numeric_format":         "Numeric Format",
	"components":             "Components",
	"day":                    "Day",
	"month":                  "Month",
	"year":                   "Year",
	"available_calc_methods": "Available Calculation Methods",
	"available_madhabs":      "Available Madhabs (Jurisprudence Schools)",
	"edit_config_message":    "Edit your config file to change these settings",

	// Day names
	"day_monday":    "Monday",
	"day_tuesday":   "Tuesday",
	"day_wednesday": "Wednesday",
	"day_thursday":  "Thursday",
	"day_friday":    "Friday",
	"day_saturday":  "Saturday",
	"day_sunday":    "Sunday",

	// Month names (Gregorian)
	"month_january":   "January",
	"month_february":  "February",
	"month_march":     "March",
	"month_april":     "April",
	"month_may":       "May",
	"month_june":      "June",
	"month_july":      "July",
	"month_august":    "August",
	"month_september": "September",
	"month_october":   "October",
	"month_november":  "November",
	"month_december":  "December",
}

var arabicMessages = map[string]string{
	// Prayer names
	"prayer_fajr":    "الفجر",
	"prayer_sunrise": "الشروق",
	"prayer_dhuhr":   "الظهر",
	"prayer_asr":     "العصر",
	"prayer_maghrib": "المغرب",
	"prayer_isha":    "العشاء",

	// UI Messages
	"prayer_times_for":       "أوقات الصلاة ليوم",
	"next_prayer":            "الصلاة القادمة",
	"in":                     "خلال",
	"current_configuration":  "الإعدادات الحالية",
	"location":               "الموقع",
	"timezone":               "المنطقة الزمنية",
	"method":                 "طريقة الحساب",
	"madhab":                 "المذهب",
	"language":               "اللغة",
	"config_file":            "ملف الإعدادات",
	"todays_date":            "تاريخ اليوم",
	"gregorian":              "الميلادي",
	"hijri":                  "الهجري",
	"detailed_date_info":     "معلومات التاريخ التفصيلية",
	"arabic_format":          "الصيغة العربية",
	"numeric_format":         "الصيغة الرقمية",
	"components":             "المكونات",
	"day":                    "اليوم",
	"month":                  "الشهر",
	"year":                   "السنة",
	"available_calc_methods": "طرق الحساب المتاحة",
	"available_madhabs":      "المذاهب المتاحة (المدارس الفقهية)",
	"edit_config_message":    "قم بتحرير ملف الإعدادات لتغيير هذه الإعدادات",

	// Day names
	"day_monday":    "الإثنين",
	"day_tuesday":   "الثلاثاء",
	"day_wednesday": "الأربعاء",
	"day_thursday":  "الخميس",
	"day_friday":    "الجمعة",
	"day_saturday":  "السبت",
	"day_sunday":    "الأحد",

	// Month names (Gregorian)
	"month_january":   "يناير",
	"month_february":  "فبراير",
	"month_march":     "مارس",
	"month_april":     "أبريل",
	"month_may":       "مايو",
	"month_june":      "يونيو",
	"month_july":      "يوليو",
	"month_august":    "أغسطس",
	"month_september": "سبتمبر",
	"month_october":   "أكتوبر",
	"month_november":  "نوفمبر",
	"month_december":  "ديسمبر",
}

// Get returns the localized message for a given key
func (l *Localization) Get(key string) string {
	if msg, exists := l.Messages[key]; exists {
		return msg
	}
	// Fallback to English if key not found
	if l.Language != "en" {
		if msg, exists := englishMessages[key]; exists {
			return msg
		}
	}
	// Return the key itself if not found in any language
	return key
}

// FormatGregorianDate formats a Gregorian date according to the language
func (l *Localization) FormatGregorianDate(t time.Time) string {
	if l.Language == "ar" {
		// Arabic format: day name, day month year
		dayName := l.GetDayName(t.Weekday())
		monthName := l.GetMonthName(t.Month())
		return fmt.Sprintf("%s، %d %s %d", dayName, t.Day(), monthName, t.Year())
	}
	// English format: day name, month day, year
	dayName := l.GetDayName(t.Weekday())
	monthName := l.GetMonthName(t.Month())
	return fmt.Sprintf("%s, %s %d, %d", dayName, monthName, t.Day(), t.Year())
}

// GetDayName returns the localized day name
func (l *Localization) GetDayName(weekday time.Weekday) string {
	dayKeys := []string{
		"day_sunday", "day_monday", "day_tuesday", "day_wednesday",
		"day_thursday", "day_friday", "day_saturday",
	}
	return l.Get(dayKeys[int(weekday)])
}

// GetMonthName returns the localized month name
func (l *Localization) GetMonthName(month time.Month) string {
	monthKeys := []string{
		"", "month_january", "month_february", "month_march", "month_april",
		"month_may", "month_june", "month_july", "month_august", "month_september",
		"month_october", "month_november", "month_december",
	}
	return l.Get(monthKeys[int(month)])
}

// GetPrayerName returns the localized prayer name
func (l *Localization) GetPrayerName(prayer string) string {
	prayerKey := "prayer_" + prayer
	return l.Get(prayerKey)
}

// IsRTL returns true if the language is right-to-left
func (l *Localization) IsRTL() bool {
	return l.Language == "ar"
}

// GetPrayerEmoji returns appropriate emoji for each prayer
func GetPrayerEmoji(prayer string) string {
	switch prayer {
	case "fajr":
		return "🌅"
	case "sunrise":
		return "☀️"
	case "dhuhr":
		return "🌞"
	case "asr":
		return "🌇"
	case "maghrib":
		return "🌅"
	case "isha":
		return "🌙"
	default:
		return "🕐"
	}
}
