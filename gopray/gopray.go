package gopray

import (
	"fmt"
	"time"

	"github.com/hablullah/go-hijri"
	calc "github.com/mnadev/adhango/pkg/calc"
	data "github.com/mnadev/adhango/pkg/data"
	util "github.com/mnadev/adhango/pkg/util"
)

// TODO the help page
// TODO implement left for next prayer
// put it in the gui adkar
// timeZone should be right misspelled?
// how to make it an imported moudle??
// reformat the prayer list

type App struct {
	ConfFile
	PrayerTimes *calc.PrayerTimes
}

func MakeApp() (*App, error) {
	a := &App{}
	a.getConf()
	err := a.calculatePrayers()
	if err != nil {
		return nil, err
	}

	return a, nil
}

// TODO unfinished work, assign params to values. use reflections? if or switch?
func (a *App) calculatePrayers() error {

	date := data.NewDateComponents(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC))

	// TODO: refactor Madhab
	// default???
	madhabs := map[string]calc.AsrJuristicMethod{
		"SHAFI_HANBALI_MALIKI": calc.SHAFI_HANBALI_MALIKI,
		"HANAFI":               calc.HANAFI,
	}
	AsrMadhab := madhabs[a.Madhab]

	// TODO : refactor Method
	// default??
	// print available list and calc default
	methods := map[string]calc.CalculationMethod{
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
	calcMethod := methods[a.Method]

	params := calc.GetMethodParameters(calcMethod)
	params.Madhab = AsrMadhab

	coords, err := util.NewCoordinates(a.Latitude, a.Longitude)
	if err != nil {
		// fmt.Printf("got error %+v", err)
		// return nil
		return err
	}

	a.PrayerTimes, err = calc.NewPrayerTimes(coords, date, params)
	if err != nil {
		// fmt.Printf("got error %+v", err)
		// return nil
		return err
	}

	err = a.PrayerTimes.SetTimeZone(a.TimeZone)
	if err != nil {
		// fmt.Printf("got error %+v", err)
		// return nil
		return err
	}

	return nil
}

func (a *App) TimeLeftForNextPrayer() time.Duration {
	nextPrayer := a.PrayerTimes.NextPrayerNow()
	nextPrayerTime := a.PrayerTimes.TimeForPrayer(nextPrayer)

	// after Isha and before midnight
	if nextPrayer == 0 {
		nextPrayer = 1
		nextPrayerTime = a.PrayerTimes.TimeForPrayer(nextPrayer).Add(time.Hour * 24)
	}

	timeLeft := time.Until(nextPrayerTime)
	return timeLeft.Round(time.Second)

}

func (a *App) ListAllPrayers() {
	// original format
	// Isha: 2015-07-12 21:57:00 -0400 EDT

	fmt.Printf("Fajr: %+v\n", a.PrayerTimes.Fajr.Format("15:04"))
	fmt.Printf("Sunrise: %+v\n", a.PrayerTimes.Sunrise.Format("15:04"))
	fmt.Printf("Dhuhr: %+v\n", a.PrayerTimes.Dhuhr.Format("15:04"))
	fmt.Printf("Asr: %+v\n", a.PrayerTimes.Asr.Format("15:04"))
	fmt.Printf("Maghrib: %+v\n", a.PrayerTimes.Maghrib.Format("15:04"))
	fmt.Printf("Isha: %+v\n", a.PrayerTimes.Isha.Format("15:04"))
}

func (a *App) HijriDate() string {
	hi, _ := hijri.CreateHijriDate(time.Date(2023, 7, 5, 0, 0, 0, 0, time.UTC),
		hijri.Default)
	// TODO ? maybe let the user choose their prefered format?
	hijriDate := fmt.Sprintf("%d-%d-%d", hi.Day, hi.Month, hi.Year)
	return hijriDate
}
