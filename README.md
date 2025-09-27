# 🌙 GoPray CLI

A modern Islamic prayer times CLI tool

## ✨ Features

- **🕐 Accurate Prayer Times** - Five daily prayers + sunrise
- **📅 Hijri Calendar** - Precise Islamic dates with Umm al-Qura calendar
- **🌐 Bilingual Interface** - Full English and Arabic support  
- **🚀 Interactive Setup** - Easy configuration with city presets

## 🚀 Quick Start

```bash
# Download and run interactive setup
./gopray setup

# View today's prayer times
./gopray list

# Check next prayer
./gopray next

# Show current date (Hijri + Gregorian)
./gopray date
```

## 🔧 Installation

```bash
go build -o gopray
./gopray setup
```

### Adding to PATH (Optional)

To use `gopray` from anywhere in your terminal, add the binary to your system PATH:

#### Linux/macOS
```bash
# Move binary to a directory in PATH (requires sudo)
sudo mv gopray /usr/local/bin/

# Or add current directory to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$PATH:$(pwd)"

# Then you can use it globally
gopray list
```

## 📱 Commands

| Command | Description | Example |
|---------|-------------|---------|
| `setup` | Interactive configuration | `./gopray setup` |
| `list` | Show all prayer times | `./gopray list` |
| `next` | Time to next prayer | `./gopray next` |
| `date` | Current Hijri date | `./gopray date -d` |
| `config` | View settings | `./gopray config` |
| `language [en\|ar]` | Change language | `./gopray language ar` |

## 🌐 Languages

Switch between English and Arabic interfaces:

```bash
# العربية
./gopray language ar
./gopray list

# English  
./gopray language en
./gopray list
```

## 📋 Configuration

### Configuration File Location

#### Linux/macOS
```
~/.config/go-pray/conf.yaml
```

#### Windows
```
%APPDATA%\go-pray\conf.yaml
```

```yaml
TimeZone: "Asia/Dubai"
Method: "UMM_AL_QURA"
Madhab: "SHAFI_HANBALI_MALIKI"
Language: "ar"
Latitude: 25.2048
Longitude: 55.2708
```

## 🤝 Contributing

We welcome contributions! Here's how to get started:

### Development Setup

```bash
# Clone the repository
git clone https://github.com/MuhammadSwa/gopray-cli.git
cd gopray-cli

# Install dependencies
go mod tidy

# Build and test
go build -o gopray
./gopray setup
```

### 🏗️ Code Structure

```
gopray/
├── main.go                 # CLI interface (Cobra commands)
├── gopray/
│   ├── gopray.go          # Core app logic & prayer calculations
│   ├── hijri.go           # Umm al-Qura calendar implementation
│   ├── i18n.go            # Internationalization system
│   ├── setup.go           # Interactive configuration setup
│   └── readConfFile.go    # Configuration management
├── go.mod                 # Go modules
└── README.md
```

#### Key Components

- **`main.go`** - Entry point with Cobra CLI framework
- **`gopray.go`** - Prayer time calculations using Adhango library
- **`hijri.go`** - Islamic calendar conversion
- **`i18n.go`** - Bilingual support (English/Arabic) 
- **`setup.go`** - Interactive setup with city presets
- **`readConfFile.go`** - YAML config loading and validation

### Making Changes

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b feature/amazing-feature`
3. **Make** your changes
4. **Test** thoroughly: `go test ./...`
5. **Commit** with clear messages: `git commit -m 'Add amazing feature'`
6. **Push** to your fork: `git push origin feature/amazing-feature`
7. **Submit** a Pull Request

### Contribution Areas

- **🌐 Languages** - Add translations in `i18n.go`
- **📱 Features** - New commands or functionality
- **🐛 Bug Fixes** - Improve reliability

### Adding Translations

To add new language support, extend `gopray/i18n.go`:

```go
// Add your language messages
var [language]Messages = map[string]string{
    "prayer_fajr": "Translation",
    // ... other keys
}
```

## Acknowledgments 
- Built with [Adhango](https://github.com/mnadev/adhango) library for accurate prayer time calculations 
- Uses [saidalisamed hijrical](https://github.com/saidalisamed/utils/blob/master/hijrical/hirjical.go) for Hijri date conversions 


## 📄 License

MIT License - feel free to use and contribute!

---

**Made with ❤️ in Egypt for the Muslim community worldwide** 🕌



