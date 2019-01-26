package i18n

import "time"


func Ago(when time.Time) string {
	t := time.Since(when)

	years := t / time.Hour * 24 * 365
	days := t / time.Hour * 24
	hours := t / time.Hour
	minutes := t / time.Minute
	seconds := t / time.Second

	if years < 1 && days < 1 && hours < 1 && minutes < 1 {
		return string(seconds) + " Second ago."
	} else {
		if years > 1 {
			return string(years) + " Years ago."
		}
		if days > 1 {
			return string(days) + " Days ago."
		}
		if hours > 1 {
			return string(hours) + " Hours ago."
		}
		if minutes > 1 {
			return string(minutes) + " Minutes ago."
		}
	}
	return "wat"
}
