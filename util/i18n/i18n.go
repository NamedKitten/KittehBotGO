package i18n

//import "fmt"
import "time"
import "github.com/nicksnyder/go-i18n/i18n"
import "github.com/NamedKitten/KittehBotGo/util/static"

func init() {
	files, _ := static.WalkDirs("translations/", false)
	for _, filename := range files {
		file, err := static.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		i18n.ParseTranslationFileBytes(filename, file)

	}
	//fmt.Println("Initialising map.")
	//Translations = make(map[string]map[string]string)
}

func Ago(when time.Time, T i18n.TranslateFunc) string {
	t := time.Since(when)

	years := t / time.Hour * 24 * 365
	days := t / time.Hour * 24
	hours := t / time.Hour
	minutes := t / time.Minute
	seconds := t / time.Second

	if years < 1 && days < 1 && hours < 1 && minutes < 1 {
		return T("humanize_second_ago", seconds)
	} else {
		if years > 1 {
			return T("humanize_year_ago", years)
		}
		if days > 1 {
			return T("humanize_day_ago", days)
		}
		if hours > 1 {
			return T("humanize_hour_ago", hours)
		}
		if minutes > 1 {
			return T("humanize_minute_ago", minutes)
		}
	}
	return "wat"
}
