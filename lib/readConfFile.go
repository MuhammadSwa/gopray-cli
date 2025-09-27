package gopray

import (
	"fmt"
	"os"
	"path/filepath"

	calc "github.com/mnadev/adhango/pkg/calc"
	"gopkg.in/yaml.v3"
)

var (
	osConfig, _ = os.UserConfigDir()
	confDir     = filepath.Join(osConfig, "go-pray")
	configPath  = filepath.Join(confDir, "conf.yaml")
)

// calculationMethods maps string names to calculation methods
var calculationMethods = map[string]calc.CalculationMethod{
	"MUSLIM_WORLD_LEAGUE":    calc.MUSLIM_WORLD_LEAGUE,
	"EGYPTIAN":               calc.EGYPTIAN,
	"UMM_AL_QURA":            calc.UMM_AL_QURA,
	"DUBAI":                  calc.DUBAI,
	"MOONSIGHTING_COMMITTEE": calc.MOON_SIGHTING_COMMITTEE,
	"NORTH_AMERICA":          calc.NORTH_AMERICA,
	"KUWAIT":                 calc.KUWAIT,
	"QATAR":                  calc.QATAR,
	"SINGAPORE":              calc.SINGAPORE,
	"OTHER":                  calc.OTHER,
}

// madhabMethods maps string names to juristic methods
var madhabMethods = map[string]calc.AsrJuristicMethod{
	"SHAFI_HANBALI_MALIKI": calc.SHAFI_HANBALI_MALIKI,
	"HANAFI":               calc.HANAFI,
}

// Config represents the application configuration
type Config struct {
	Method    string  `yaml:"Method"`
	Madhab    string  `yaml:"Madhab"`
	TimeZone  string  `yaml:"TimeZone"`
	Latitude  float64 `yaml:"Latitude"`
	Longitude float64 `yaml:"Longitude"`
	Language  string  `yaml:"Language"` // Language preference: "en" or "ar"
}

// ConfFile is deprecated, use Config instead
type ConfFile = Config

// IsValidMethod checks if a calculation method is valid
func IsValidMethod(method string) bool {
	_, ok := calculationMethods[method]
	return ok
}

// IsValidMadhab checks if a madhab is valid
func IsValidMadhab(madhab string) bool {
	_, ok := madhabMethods[madhab]
	return ok
}

// ValidateMethods returns lists of valid calculation methods and madhabs
func ValidateMethods() ([]string, []string) {
	var methods, madhabs []string

	for method := range calculationMethods {
		methods = append(methods, method)
	}

	for madhab := range madhabMethods {
		madhabs = append(madhabs, madhab)
	}

	return methods, madhabs
}

// LoadConfig loads configuration from file, creating defaults if needed
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Check if config exists and create it if needed
	if err := ensureConfigExists(); err != nil {
		return nil, fmt.Errorf("failed to ensure config file exists: %w", err)
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Parse the YAML content
	if err = yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err = validateConfig(config); err != nil {
		// If validation fails due to coordinates not being set, suggest interactive setup
		if config.Latitude == 0.0 && config.Longitude == 0.0 {
			fmt.Println("\n🚀 Welcome to GoPray!")
			fmt.Println("📍 It looks like this is your first time using GoPray.")
			fmt.Println("🔧 Your location coordinates need to be configured for accurate prayer times.")
			fmt.Println("\n💡 You can:")
			fmt.Println("   1. Run 'gopray-cli setup' for an interactive configuration")
			fmt.Printf("   2. Edit the config file manually: %s\n", configPath)
			fmt.Println()
		}
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// ensureConfigExists creates the default config file if it doesn't exist
func ensureConfigExists() error {
	// Check if file exists
	if _, err := os.Stat(configPath); err == nil {
		return nil // File exists
	}

	fmt.Printf("Creating default config file at %s\n", configPath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Default configuration
	defaultConfig := `# GoPray Configuration File
# Edit this file to set your location and preferences

# Your timezone (IANA timezone identifier)
# Examples:  "Africa/Cairo", "Asia/Dubai", "Australia/Sydney"
TimeZone: "Africa/Cairo"

# Prayer calculation method
# Options: MUSLIM_WORLD_LEAGUE, EGYPTIAN, UMM_AL_QURA, DUBAI, 
#          MOONSIGHTING_COMMITTEE, NORTH_AMERICA, KUWAIT, QATAR, SINGAPORE
Method: "EGYPTIAN"

# Juristic method for Asr prayer calculation
# Options: SHAFI_HANBALI_MALIKI, HANAFI
Madhab: "SHAFI_HANBALI_MALIKI"

# Language preference for the interface
# Options: "en" (English), "ar" (Arabic)
Language: "en"

# Your location coordinates (required for accurate prayer times)
# You can find these using Google Maps or similar services
Latitude: 0.0
Longitude: 0.0
`

	// Write the config file
	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Println("✅ Default config file created successfully!")
	fmt.Println("⚠️  Please edit the config file to set your location coordinates.")
	fmt.Printf("📁 Config file location: %s\n", configPath)

	return nil
}

// validateConfig validates the configuration values
func validateConfig(config *Config) error {
	// Check if coordinates are set (not default 0,0)
	if config.Latitude == 0.0 && config.Longitude == 0.0 {
		return fmt.Errorf("location coordinates not set - please update your config file with your actual latitude and longitude")
	}

	// Validate latitude range
	if config.Latitude < -90 || config.Latitude > 90 {
		return fmt.Errorf("invalid latitude %f - must be between -90 and 90", config.Latitude)
	}

	// Validate longitude range
	if config.Longitude < -180 || config.Longitude > 180 {
		return fmt.Errorf("invalid longitude %f - must be between -180 and 180", config.Longitude)
	}

	// Validate method
	if !IsValidMethod(config.Method) {
		validMethods, _ := ValidateMethods()
		return fmt.Errorf("invalid calculation method '%s' - valid methods: %v", config.Method, validMethods)
	}

	// Validate madhab
	if !IsValidMadhab(config.Madhab) {
		_, validMadhabs := ValidateMethods()
		return fmt.Errorf("invalid madhab '%s' - valid madhabs: %v", config.Madhab, validMadhabs)
	}

	// Validate language
	if config.Language == "" {
		config.Language = "en" // Default to English
	}
	if config.Language != "en" && config.Language != "ar" {
		return fmt.Errorf("invalid language '%s' - supported languages: en (English), ar (Arabic)", config.Language)
	}

	return nil
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	return configPath
}

// GetConfigPath is the public version of getConfigPath
func GetConfigPath() string {
	return configPath
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *Config) error {
	// Ensure the config directory exists
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal the config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %w", err)
	}

	// Write the config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
