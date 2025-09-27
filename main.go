package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MuhammadSwa/gopray-cli/lib"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	app     *gopray.App
)

func init() {
	var err error
	app, err = gopray.NewApp()
	if err != nil {
		log.Printf("Warning: Failed to initialize app: %v", err)
		log.Println("You may need to configure your location settings.")
	}
}

var rootCmd = &cobra.Command{
	Use:     "gopray",
	Short:   "Islamic prayer times CLI tool",
	Long:    `GoPray is a command-line tool for displaying Islamic prayer times and Hijri dates.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display all prayer times for today",
	Long:  `Shows all five daily prayer times (Fajr, Dhuhr, Asr, Maghrib, Isha) along with sunrise time.`,
	Run: func(cmd *cobra.Command, args []string) {
		if app == nil {
			fmt.Fprintln(os.Stderr, "Error: App not initialized. Please check your configuration.")
			os.Exit(1)
		}
		app.ListAllPrayers()
	},
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Show time remaining until next prayer",
	Long:  `Displays the time remaining until the next prayer time.`,
	Run: func(cmd *cobra.Command, args []string) {
		if app == nil {
			fmt.Fprintln(os.Stderr, "Error: App not initialized. Please check your configuration.")
			os.Exit(1)
		}
		duration := app.TimeLeftForNextPrayer()
		prayerName := app.GetNextPrayerName()

		if app.Config.Language == "ar" {
			fmt.Printf("⏰ %s: %s %s %s\n", app.Localizer.Get("next_prayer"), prayerName, app.Localizer.Get("in"), formatDuration(duration))
		} else {
			fmt.Printf("⏰ %s: %s %s %s\n", app.Localizer.Get("next_prayer"), prayerName, app.Localizer.Get("in"), formatDuration(duration))
		}
	},
}

var dateCmd = &cobra.Command{
	Use:   "date",
	Short: "Display current Hijri date",
	Long:  `Shows the current date in the Islamic (Hijri) calendar alongside the Gregorian date.`,
	Run: func(cmd *cobra.Command, args []string) {
		if app == nil {
			fmt.Fprintln(os.Stderr, "Error: App not initialized. Please check your configuration.")
			os.Exit(1)
		}

		// Get current time in the user's timezone
		now := time.Now()
		if app.Config != nil && app.Config.TimeZone != "" {
			if location, err := time.LoadLocation(app.Config.TimeZone); err == nil {
				now = now.In(location)
			}
		}

		// Check for detailed flag
		detailed, _ := cmd.Flags().GetBool("detailed")

		if detailed {
			hijriDetails, err := app.GetDetailedHijriDate()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting Hijri date: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("📅 %s:\n\n", app.Localizer.Get("detailed_date_info"))
			fmt.Printf("🗓️  %s: %s\n", app.Localizer.Get("gregorian"), app.Localizer.FormatGregorianDate(now))

			if app.Config.Language == "ar" {
				fmt.Printf("🌙 %s: %s\n", app.Localizer.Get("hijri"), hijriDetails["formattedAr"])
				fmt.Printf("📊 %s: %s\n", app.Localizer.Get("numeric_format"), hijriDetails["numeric"])
			} else {
				fmt.Printf("🌙 %s: %s\n", app.Localizer.Get("hijri"), hijriDetails["formatted"])
				fmt.Printf("🌙 %s: %s\n", app.Localizer.Get("arabic_format"), hijriDetails["formattedAr"])
				fmt.Printf("📊 %s: %s\n", app.Localizer.Get("numeric_format"), hijriDetails["numeric"])
			}

			fmt.Printf("📝 %s: %s %v, %s %v (%s), %s %v\n",
				app.Localizer.Get("components"),
				app.Localizer.Get("day"), hijriDetails["day"],
				app.Localizer.Get("month"), hijriDetails["month"],
				hijriDetails["monthNameEn"],
				app.Localizer.Get("year"), hijriDetails["year"])
		} else {
			hijriDate, err := app.HijriDate()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting Hijri date: %v\n", err)
				os.Exit(1)
			}

			// Display both dates for context
			fmt.Printf("📅 %s:\n", app.Localizer.Get("todays_date"))
			fmt.Printf("🗓️  %s: %s\n", app.Localizer.Get("gregorian"), app.Localizer.FormatGregorianDate(now))
			fmt.Printf("🌙 %s: %s\n", app.Localizer.Get("hijri"), hijriDate)
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show configuration information",
	Long:  `Displays current configuration settings and file location.`,
	Run: func(cmd *cobra.Command, args []string) {
		if app == nil {
			fmt.Fprintln(os.Stderr, "Error: App not initialized.")
			os.Exit(1)
		}
		app.ShowConfig()
	},
}

var methodsCmd = &cobra.Command{
	Use:   "methods",
	Short: "List available calculation methods and madhabs",
	Long:  `Shows all supported prayer calculation methods and jurisprudence schools (madhabs).`,
	Run: func(cmd *cobra.Command, args []string) {
		methods, madhabs := gopray.ValidateMethods()

		// Use localizer if app is initialized
		var localizer *gopray.Localization
		if app != nil {
			localizer = app.Localizer
		} else {
			localizer = gopray.GetLocalizer("en") // Default to English
		}

		fmt.Printf("📐 %s:\n", localizer.Get("available_calc_methods"))
		for i, method := range methods {
			fmt.Printf("  %d. %s\n", i+1, method)
		}

		fmt.Printf("\n⚖️  %s:\n", localizer.Get("available_madhabs"))
		for i, madhab := range madhabs {
			fmt.Printf("  %d. %s\n", i+1, madhab)
		}

		fmt.Printf("\n📁 %s: %s\n", localizer.Get("config_file"), gopray.GetConfigPath())
		fmt.Printf("\n%s.\n", localizer.Get("edit_config_message"))
	},
}

var languageCmd = &cobra.Command{
	Use:   "language [en|ar]",
	Short: "Change interface language / تغيير لغة الواجهة",
	Long:  `Change the interface language between English (en) and Arabic (ar) / تغيير لغة الواجهة بين الإنجليزية (en) والعربية (ar)`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lang := args[0]
		if lang != "en" && lang != "ar" {
			fmt.Fprintln(os.Stderr, "Error: Language must be 'en' (English) or 'ar' (Arabic)")
			fmt.Fprintln(os.Stderr, "خطأ: يجب أن تكون اللغة 'en' (إنجليزية) أو 'ar' (عربية)")
			os.Exit(1)
		}

		// Load current config
		config, err := gopray.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		// Update language
		config.Language = lang

		// Save updated config
		err = gopray.SaveConfig(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		if lang == "ar" {
			fmt.Println("✅ تم تغيير اللغة إلى العربية بنجاح!")
			fmt.Println("🔄 يرجى تشغيل الأمر مرة أخرى لرؤية الواجهة باللغة العربية.")
		} else {
			fmt.Println("✅ Language changed to English successfully!")
			fmt.Println("🔄 Please run the command again to see the interface in English.")
		}
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive configuration setup / إعداد تفاعلي للإعدادات",
	Long:  `Run an interactive setup to configure your location, timezone, calculation method, and other preferences / تشغيل إعداد تفاعلي لتكوين موقعك والمنطقة الزمنية وطريقة الحساب والتفضيلات الأخرى`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gopray.InteractiveSetup(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Add flags to commands
	dateCmd.Flags().BoolP("detailed", "d", false, "Show detailed date information including Arabic names")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(nextCmd)
	rootCmd.AddCommand(dateCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(methodsCmd)
	rootCmd.AddCommand(languageCmd)
	rootCmd.AddCommand(setupCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
