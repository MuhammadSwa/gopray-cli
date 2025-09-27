package gopray

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// CityPreset represents a preset configuration for a city
type CityPreset struct {
	Name      string
	NameAr    string
	Latitude  float64
	Longitude float64
	Timezone  string
	Method    string
}

// Popular city presets
var cityPresets = []CityPreset{
	{"Cairo, Egypt", "القاهرة، مصر", 30.0444, 31.2357, "Africa/Cairo", "EGYPTIAN"},
	{"Dubai, UAE", "دبي، الإمارات", 25.2048, 55.2708, "Asia/Dubai", "UMM_AL_QURA"},
	{"Riyadh, Saudi Arabia", "الرياض، السعودية", 24.7136, 46.6753, "Asia/Riyadh", "UMM_AL_QURA"},
	{"Istanbul, Turkey", "إسطنبول، تركيا", 41.0082, 28.9784, "Europe/Istanbul", "MUSLIM_WORLD_LEAGUE"},
	{"London, UK", "لندن، المملكة المتحدة", 51.5074, -0.1278, "Europe/London", "MUSLIM_WORLD_LEAGUE"},
	{"New York, USA", "نيويورك، الولايات المتحدة", 40.7128, -74.0060, "America/New_York", "NORTH_AMERICA"},
	{"Kuala Lumpur, Malaysia", "كوالالمبور، ماليزيا", 3.1390, 101.6869, "Asia/Kuala_Lumpur", "SINGAPORE"},
	{"Jakarta, Indonesia", "جاكرتا، إندونيسيا", -6.2088, 106.8456, "Asia/Jakarta", "SINGAPORE"},
}

// InteractiveSetup runs an interactive configuration setup
func InteractiveSetup() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🌙 Welcome to GoPray Interactive Setup!")
	fmt.Println("📋 This will help you configure your prayer times settings.")

	// Language selection first
	language := promptLanguage(reader)
	localizer := GetLocalizer(language)

	// Create new config with defaults
	config := &Config{
		Language: language,
		Method:   "EGYPTIAN",
		Madhab:   "SHAFI_HANBALI_MALIKI",
		TimeZone: "Africa/Cairo",
	}

	fmt.Printf("\n🌐 %s: %s\n", localizer.Get("language"), getLanguageDisplay(language))

	// Location setup
	if err := promptLocation(reader, config, localizer); err != nil {
		return err
	}

	// Timezone setup
	if err := promptTimezone(reader, config, localizer); err != nil {
		return err
	}

	// Method setup
	if err := promptMethod(reader, config, localizer); err != nil {
		return err
	}

	// Madhab setup
	if err := promptMadhab(reader, config, localizer); err != nil {
		return err
	}

	// Show summary and confirm
	fmt.Printf("\n📋 %s:\n", localizer.Get("current_configuration"))
	fmt.Printf("📍 %s: %.4f°, %.4f°\n", localizer.Get("location"), config.Latitude, config.Longitude)
	fmt.Printf("🕰️  %s: %s\n", localizer.Get("timezone"), config.TimeZone)
	fmt.Printf("📐 %s: %s\n", localizer.Get("method"), config.Method)
	fmt.Printf("⚖️  %s: %s\n", localizer.Get("madhab"), config.Madhab)
	fmt.Printf("🌐 %s: %s\n", localizer.Get("language"), getLanguageDisplay(config.Language))

	// Confirm save
	if language == "ar" {
		fmt.Print("\n💾 هل تريد حفظ هذه الإعدادات؟ [y/N]: ")
	} else {
		fmt.Print("\n💾 Save these settings? [y/N]: ")
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" || input == "yes" || input == "نعم" {
		if err := SaveConfig(config); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		if language == "ar" {
			fmt.Println("✅ تم حفظ الإعدادات بنجاح!")
			fmt.Println("🚀 يمكنك الآن استخدام gopray-cli لعرض أوقات الصلاة.")
		} else {
			fmt.Println("✅ Configuration saved successfully!")
			fmt.Println("🚀 You can now use gopray-cli to view prayer times.")
		}
	} else {
		if language == "ar" {
			fmt.Println("❌ تم إلغاء حفظ الإعدادات.")
		} else {
			fmt.Println("❌ Configuration not saved.")
		}
	}

	return nil
}

func promptLanguage(reader *bufio.Reader) string {
	fmt.Println("🌐 Choose your language / اختر لغتك:")
	fmt.Println("1. English")
	fmt.Println("2. العربية (Arabic)")
	fmt.Print("Enter choice [1-2]: ")

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1", "en", "english":
			return "en"
		case "2", "ar", "arabic", "عربي", "عربية":
			return "ar"
		default:
			fmt.Print("❌ Invalid choice. Please enter 1 or 2: ")
		}
	}
}

func promptLocation(reader *bufio.Reader, config *Config, localizer *Localization) error {
	if localizer.Language == "ar" {
		fmt.Printf("\n📍 إعداد الموقع:\n")
		fmt.Println("يمكنك اختيار مدينة من القائمة أو إدخال الإحداثيات يدوياً.")
		fmt.Println("\n🏙️ المدن الشائعة:")
	} else {
		fmt.Printf("\n📍 Location Setup:\n")
		fmt.Println("You can choose a city from the list or enter coordinates manually.")
		fmt.Println("\n🏙️ Popular cities:")
	}

	// Display city presets
	for i, city := range cityPresets {
		if localizer.Language == "ar" {
			fmt.Printf("  %d. %s\n", i+1, city.NameAr)
		} else {
			fmt.Printf("  %d. %s\n", i+1, city.Name)
		}
	}

	if localizer.Language == "ar" {
		fmt.Printf("  %d. إدخال الإحداثيات يدوياً\n", len(cityPresets)+1)
		fmt.Printf("\nاختر مدينة أو إدخال يدوي [1-%d]: ", len(cityPresets)+1)
	} else {
		fmt.Printf("  %d. Enter coordinates manually\n", len(cityPresets)+1)
		fmt.Printf("\nChoose a city or manual entry [1-%d]: ", len(cityPresets)+1)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if choice, err := strconv.Atoi(input); err == nil && choice >= 1 && choice <= len(cityPresets) {
		// City preset selected
		preset := cityPresets[choice-1]
		config.Latitude = preset.Latitude
		config.Longitude = preset.Longitude
		config.TimeZone = preset.Timezone
		config.Method = preset.Method

		if localizer.Language == "ar" {
			fmt.Printf("✅ تم اختيار: %s\n", preset.NameAr)
		} else {
			fmt.Printf("✅ Selected: %s\n", preset.Name)
		}
		return nil
	} else if choice, err := strconv.Atoi(input); err == nil && choice == len(cityPresets)+1 {
		// Manual coordinates entry
		return promptManualCoordinates(reader, config, localizer)
	} else {
		if localizer.Language == "ar" {
			fmt.Println("❌ اختيار غير صحيح، سيتم الانتقال إلى الإدخال اليدوي")
		} else {
			fmt.Println("❌ Invalid choice, proceeding with manual entry")
		}
		return promptManualCoordinates(reader, config, localizer)
	}
}

func promptManualCoordinates(reader *bufio.Reader, config *Config, localizer *Localization) error {
	if localizer.Language == "ar" {
		fmt.Printf("\n📍 إدخال الإحداثيات يدوياً:\n")
		fmt.Println("يمكنك العثور على إحداثيات موقعك باستخدام خرائط جوجل.")
	} else {
		fmt.Printf("\n📍 Manual Coordinates Entry:\n")
		fmt.Println("You can find your location coordinates using Google Maps.")
	}

	// Latitude
	for {
		if localizer.Language == "ar" {
			fmt.Print("خط العرض (Latitude) [-90 to 90]: ")
		} else {
			fmt.Print("Latitude [-90 to 90]: ")
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if lat, err := strconv.ParseFloat(input, 64); err == nil {
			if lat >= -90 && lat <= 90 {
				config.Latitude = lat
				break
			}
		}

		if localizer.Language == "ar" {
			fmt.Println("❌ يرجى إدخال رقم صحيح بين -90 و 90")
		} else {
			fmt.Println("❌ Please enter a valid number between -90 and 90")
		}
	}

	// Longitude
	for {
		if localizer.Language == "ar" {
			fmt.Print("خط الطول (Longitude) [-180 to 180]: ")
		} else {
			fmt.Print("Longitude [-180 to 180]: ")
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if lng, err := strconv.ParseFloat(input, 64); err == nil {
			if lng >= -180 && lng <= 180 {
				config.Longitude = lng
				break
			}
		}

		if localizer.Language == "ar" {
			fmt.Println("❌ يرجى إدخال رقم صحيح بين -180 و 180")
		} else {
			fmt.Println("❌ Please enter a valid number between -180 and 180")
		}
	}

	return nil
}

func promptTimezone(reader *bufio.Reader, config *Config, localizer *Localization) error {
	if localizer.Language == "ar" {
		fmt.Printf("\n🕰️ إعداد المنطقة الزمنية:\n")
		fmt.Println("أمثلة: Africa/Cairo, Asia/Dubai, Asia/Riyadh, Europe/London")
		fmt.Printf("المنطقة الزمنية الحالية: %s\n", config.TimeZone)
		fmt.Print("أدخل المنطقة الزمنية الجديدة (اتركها فارغة للاحتفاظ بالحالية): ")
	} else {
		fmt.Printf("\n🕰️ Timezone Setup:\n")
		fmt.Println("Examples: Africa/Cairo, Asia/Dubai, Asia/Riyadh, Europe/London")
		fmt.Printf("Current timezone: %s\n", config.TimeZone)
		fmt.Print("Enter new timezone (leave empty to keep current): ")
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		// Validate timezone
		if _, err := time.LoadLocation(input); err == nil {
			config.TimeZone = input
			if localizer.Language == "ar" {
				fmt.Printf("✅ تم تحديث المنطقة الزمنية إلى: %s\n", input)
			} else {
				fmt.Printf("✅ Timezone updated to: %s\n", input)
			}
		} else {
			if localizer.Language == "ar" {
				fmt.Printf("⚠️ منطقة زمنية غير صحيحة، سيتم الاحتفاظ بـ: %s\n", config.TimeZone)
			} else {
				fmt.Printf("⚠️ Invalid timezone, keeping: %s\n", config.TimeZone)
			}
		}
	}

	return nil
}

func promptMethod(reader *bufio.Reader, config *Config, localizer *Localization) error {
	methods, _ := ValidateMethods()

	if localizer.Language == "ar" {
		fmt.Printf("\n📐 طرق حساب أوقات الصلاة:\n")
	} else {
		fmt.Printf("\n📐 Prayer Calculation Methods:\n")
	}

	for i, method := range methods {
		fmt.Printf("%d. %s\n", i+1, method)
	}

	if localizer.Language == "ar" {
		fmt.Printf("الطريقة الحالية: %s\n", config.Method)
		fmt.Printf("اختر طريقة الحساب [1-%d] (اتركها فارغة للاحتفاظ بالحالية): ", len(methods))
	} else {
		fmt.Printf("Current method: %s\n", config.Method)
		fmt.Printf("Choose calculation method [1-%d] (leave empty to keep current): ", len(methods))
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		if choice, err := strconv.Atoi(input); err == nil && choice >= 1 && choice <= len(methods) {
			config.Method = methods[choice-1]
		} else {
			if localizer.Language == "ar" {
				fmt.Printf("⚠️ اختيار غير صحيح، سيتم الاحتفاظ بـ: %s\n", config.Method)
			} else {
				fmt.Printf("⚠️ Invalid choice, keeping: %s\n", config.Method)
			}
		}
	}

	return nil
}

func promptMadhab(reader *bufio.Reader, config *Config, localizer *Localization) error {
	_, madhabs := ValidateMethods()

	if localizer.Language == "ar" {
		fmt.Printf("\n⚖️ المذاهب الفقهية:\n")
	} else {
		fmt.Printf("\n⚖️ Jurisprudence Schools (Madhabs):\n")
	}

	for i, madhab := range madhabs {
		fmt.Printf("%d. %s\n", i+1, madhab)
	}

	if localizer.Language == "ar" {
		fmt.Printf("المذهب الحالي: %s\n", config.Madhab)
		fmt.Printf("اختر المذهب [1-%d] (اتركها فارغة للاحتفاظ بالحالية): ", len(madhabs))
	} else {
		fmt.Printf("Current madhab: %s\n", config.Madhab)
		fmt.Printf("Choose madhab [1-%d] (leave empty to keep current): ", len(madhabs))
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "" {
		if choice, err := strconv.Atoi(input); err == nil && choice >= 1 && choice <= len(madhabs) {
			config.Madhab = madhabs[choice-1]
		} else {
			if localizer.Language == "ar" {
				fmt.Printf("⚠️ اختيار غير صحيح، سيتم الاحتفاظ بـ: %s\n", config.Madhab)
			} else {
				fmt.Printf("⚠️ Invalid choice, keeping: %s\n", config.Madhab)
			}
		}
	}

	return nil
}

func getLanguageDisplay(language string) string {
	if language == "ar" {
		return "العربية (Arabic)"
	}
	return "English"
}
