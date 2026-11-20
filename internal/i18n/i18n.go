// Package i18n is the app's languages: those its Bibles read, and so
// those it speaks. The English text is the key: T returns it as it is in
// English, and its translation in any other language. Keys are literal
// strings at the call, for a test to find them all (see catalog_test).
package i18n

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	// System is the setting of following the system's language: see
	// Resolve.
	System = "system"

	// English is what the app is written in, and the key of the catalog.
	English = "en"

	// ChineseSimplified and ChineseTraditional are the two Chinese.
	ChineseSimplified  = "zh-Hans"
	ChineseTraditional = "zh-Hant"

	// The rest, named as the Bibles are.
	German     = "de"
	Spanish    = "es"
	French     = "fr"
	Portuguese = "pt"
	Russian    = "ru"
	Japanese   = "ja"
	Korean     = "ko"
	Arabic     = "ar"
)

// Languages are those the reader can choose, each named in itself.
// English first: the app speaks to the West, and the rest follow in the
// order their Bibles stand in the translation list.
var Languages = []struct{ Code, Name string }{
	{English, "English"},
	{ChineseSimplified, "简体中文"},
	{ChineseTraditional, "繁體中文"},
	{German, "Deutsch"},
	{Spanish, "Español"},
	{French, "Français"},
	{Portuguese, "Português"},
	{Russian, "Русский"},
	{Japanese, "日本語"},
	{Korean, "한국어"},
	{Arabic, "العربية"},
}

// known says whether the app speaks a language.
func known(code string) bool {
	for _, l := range Languages {
		if l.Code == code {
			return true
		}
	}
	return false
}

// Resolve is the language to speak for a setting: the one chosen, or for
// System, the first of those the system prefers that the app has, or
// else English. Chinese reads by its script, as a reader of Traditional
// can't be handed Simplified.
func Resolve(setting string, preferred []string) string {
	if known(setting) {
		return setting
	}
	for _, tag := range preferred {
		tag = strings.ToLower(strings.ReplaceAll(tag, "_", "-"))
		base, region, _ := strings.Cut(tag, "-")
		switch base {
		case "zh":
			switch region {
			case "tw", "hk", "mo", "hant":
				return ChineseTraditional
			}
			return ChineseSimplified
		case "de":
			return German
		case "es":
			return Spanish
		case "fr":
			return French
		case "pt":
			return Portuguese
		case "ru":
			return Russian
		case "ja":
			return Japanese
		case "ko":
			return Korean
		case "ar":
			return Arabic
		case "en":
			return English
		}
	}
	return English
}

var language atomic.Pointer[string]

// SetLanguage switches the language of what T returns from now on. A
// language the app doesn't have is English.
func SetLanguage(code string) {
	if !known(code) {
		code = English
	}
	language.Store(&code)
}

// Language returns the language T translates to.
func Language() string {
	if l := language.Load(); l != nil {
		return *l
	}
	return English
}

// catalog returns the catalog of a language, nil for English, which is
// the keys themselves.
func catalog(code string) map[string]string {
	switch code {
	case ChineseSimplified:
		return zhHans
	case ChineseTraditional:
		return zhHant
	case German:
		return de
	case Spanish:
		return es
	case French:
		return fr
	case Portuguese:
		return pt
	case Russian:
		return ru
	case Japanese:
		return ja
	case Korean:
		return ko
	case Arabic:
		return ar
	}
	return nil
}

// T returns the text in the reader's language.
func T(text string) string {
	if cat := catalog(Language()); cat != nil {
		if s, ok := cat[text]; ok {
			return s
		}
	}
	return text
}

// TC is T of a text that reads differently by what it's about, as Save
// is one word for keeping a note and another for adding a prayer to the
// list. The catalog has it as "text|context", and English reads the
// text.
func TC(text, context string) string {
	if cat := catalog(Language()); cat != nil {
		if s, ok := cat[text+"|"+context]; ok {
			return s
		}
	}
	return text
}

// Tf is T of a format, formatted. Its arguments must stay in the order
// the English has them, or the translation numbers them, as %[2]d.
func Tf(format string, args ...any) string {
	return fmt.Sprintf(T(format), args...)
}

// N marks text for translation where it's kept to show later, through T,
// as in a list: it returns it as it is.
func N(text string) string { return text }

// Ago says how long before now t was.
func Ago(t, now time.Time) string {
	d := now.Sub(t)
	switch {
	case d < time.Minute:
		return T("just now")
	case d < 2*time.Minute:
		return T("a minute ago")
	case d < time.Hour:
		return Tf("%d minutes ago", int(d/time.Minute))
	case d < 2*time.Hour:
		return T("an hour ago")
	case d < 24*time.Hour:
		return Tf("%d hours ago", int(d/time.Hour))
	case d < 48*time.Hour:
		return T("yesterday")
	case d < 30*24*time.Hour:
		return Tf("%d days ago", int(d/(24*time.Hour)))
	}
	if t.Year() == now.Year() {
		return Date(t)
	}
	return t.Format("2006-01-02")
}

// Date says t's month and day, and day of the week in Today. Each
// language reads a date as it does.
func Date(t time.Time) string {
	switch Language() {
	case ChineseSimplified, ChineseTraditional, Japanese, Korean:
		return T(strconv.Itoa(int(t.Month())) + "月" + strconv.Itoa(t.Day()) + "日")
	case Russian:
		return strconv.Itoa(t.Day()) + " " + russianMonths[int(t.Month())]
	case Arabic:
		return strconv.Itoa(t.Day()) + " " + arabicMonths[int(t.Month())]
	}
	return t.Month().String() + " " + strconv.Itoa(t.Day())
}

var russianMonths = [...]string{"", "января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"}

var arabicMonths = [...]string{"", "يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"}

// Today says the day of t, as a page heading does.
func Today(t time.Time) string {
	switch Language() {
	case ChineseSimplified, ChineseTraditional:
		return Date(t) + " " + T(weekdayKey(t))
	case Japanese:
		return Date(t) + " " + japaneseWeekdays[t.Weekday()]
	case Korean:
		return Date(t) + " " + koreanWeekdays[t.Weekday()]
	}
	return localizedWeekday(t) + ", " + Date(t)
}

func weekdayKey(t time.Time) string {
	switch t.Weekday() {
	case time.Sunday:
		return "星期日"
	case time.Monday:
		return "星期一"
	case time.Tuesday:
		return "星期二"
	case time.Wednesday:
		return "星期三"
	case time.Thursday:
		return "星期四"
	case time.Friday:
		return "星期五"
	default:
		return "星期六"
	}
}

var japaneseWeekdays = [...]string{
	time.Sunday: "日", time.Monday: "月", time.Tuesday: "火",
	time.Wednesday: "水", time.Thursday: "木", time.Friday: "金", time.Saturday: "土",
}

var koreanWeekdays = [...]string{
	time.Sunday: "일", time.Monday: "월", time.Tuesday: "화",
	time.Wednesday: "수", time.Thursday: "목", time.Friday: "금", time.Saturday: "토",
}

// localizedWeekday says the weekday in the language, by the language's
// own names, or the system's for one the app has no words of.
func localizedWeekday(t time.Time) string {
	names := weekdayNames[Language()]
	if names == nil {
		return t.Weekday().String()
	}
	return names[t.Weekday()]
}

// weekdayNames are the days of the week, Sunday first, by language.
var weekdayNames = map[string][]string{
	ChineseSimplified:  {"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"},
	ChineseTraditional: {"週日", "週一", "週二", "週三", "週四", "週五", "週六"},
	German:             {"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"},
	Spanish:            {"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"},
	French:             {"dimanche", "lundi", "mardi", "mercredi", "jeudi", "vendredi", "samedi"},
	Portuguese:         {"domingo", "segunda-feira", "terça-feira", "quarta-feira", "quinta-feira", "sexta-feira", "sábado"},
	Russian:            {"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота"},
	Arabic:             {"الأحد", "الإثنين", "الثلاثاء", "الأربعاء", "الخميس", "الجمعة", "السبت"},
}

// Count says a count shortly, as people read it in the language.
func Count(n int) string {
	switch Language() {
	case ChineseTraditional:
		switch {
		case n >= 1e8:
			return strconv.FormatFloat(float64(n)/1e8, 'f', 1, 64) + "億"
		case n >= 1e4:
			return strconv.FormatFloat(float64(n)/1e4, 'f', 1, 64) + "萬"
		}
		return strconv.Itoa(n)
	case Japanese:
		switch {
		case n >= 1e8:
			return strconv.FormatFloat(float64(n)/1e8, 'f', 1, 64) + "億"
		case n >= 1e4:
			return strconv.FormatFloat(float64(n)/1e4, 'f', 1, 64) + "万"
		}
		return strconv.Itoa(n)
	case Korean:
		switch {
		case n >= 1e8:
			return strconv.FormatFloat(float64(n)/1e8, 'f', 1, 64) + "억"
		case n >= 1e4:
			return strconv.FormatFloat(float64(n)/1e4, 'f', 1, 64) + "만"
		}
		return strconv.Itoa(n)
	case ChineseSimplified:
		switch {
		case n >= 1e8:
			return strconv.FormatFloat(float64(n)/1e8, 'f', 1, 64) + "亿"
		case n >= 1e4:
			return strconv.FormatFloat(float64(n)/1e4, 'f', 1, 64) + "万"
		}
		return strconv.Itoa(n)
	}
	switch {
	case n >= 1e9:
		return fmt.Sprintf("%.1fB", float64(n)/1e9)
	case n >= 1e6:
		return fmt.Sprintf("%.1fM", float64(n)/1e6)
	case n >= 1e4:
		return fmt.Sprintf("%.1fK", float64(n)/1e3)
	}
	return strconv.Itoa(n)
}
