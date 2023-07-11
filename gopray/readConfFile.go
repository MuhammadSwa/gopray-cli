package gopray

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	OsConfig, _ = os.UserConfigDir()
	confDir     = filepath.Join(OsConfig, "go-pray")
	yamlPath    = filepath.Join(confDir, "conf-new.yaml")
)

type ConfFile struct {
	Method    string  `yaml:"Method"`
	Madhab    string  `yaml:"Madhab"`
	TimeZone  string  `yaml:"TimeZone"`
	Latitude  float64 `yaml:"Latitude"`
	Longitude float64 `yaml:"Longitude"`
}

func (c *ConfFile) getConf() {
	c.checkIfConfigExist()
	file, err := os.ReadFile(yamlPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

// TODO check if the yaml config file exists in it's default place ???? where to read the def config file??
func (c *ConfFile) checkIfConfigExist() {
	_, err := os.Stat(yamlPath)
	if err != nil {
		fmt.Printf("Creating the default config file in %s\n", yamlPath)

		var defaultConf = `# city (string) - A city name. Example: London
city: "Cairo"

# country (string) - A country name or 2 character alpha ISO 3166 code. Examples: GB or United Kindom
country: "Egypt"

# the full path to where to store the json file.
# If left empty it will be stored in same dir with conf file
dataPath: "" 

method: 5
###################################
# method (number) -
	# A prayer times calculation method. Methods identify various schools of thought about how to compute the timings.
	# If not specified, it defaults to the closest authority based on the location or co-ordinates specified in the API call.
	# 1 - University of Islamic Sciences, Karachi
	# 2 - Islamic Society of North America
	# 3 - Muslim World League
	# 4 - Umm Al-Qura University, Makkah
	# 5 - Egyptian General Authority of Survey
	# 7 - Institute of Geophysics, University of Tehran
	# 8 - Gulf Region
	# 9 - Kuwait
	# 10 - Qatar
	# 11 - Majlis Ugama Islam Singapura, Singapore
	# 12 - Union Organization islamic de France
	# 13 - Diyanet İşleri Başkanlığı, Turkey
	# 14 - Spiritual Administration of Muslims of Russia
	# 15 - Moonsighting Committee Worldwide (also requires shafaq paramteer) 
`
		file, _ := os.OpenFile(yamlPath, os.O_CREATE|os.O_WRONLY, 0644)
		file.Write([]byte(defaultConf))
		defer file.Close()
	}
}
