package seasons

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

//go:embed data/lectionary.json
var lectData []byte

// Lessons are a day's four readings, as the lectionary names them: the
// reference is what reads, and the reader opens it in whatever translation
// is at hand.
type Lessons [4]string

// Lectionary is the Revised Common Lectionary's table: the readings of
// every Sunday and holy day, in the three years' cycle and the fixed holy
// days.
type Lectionary struct {
	years  map[string]map[string]Lessons
	ranges map[string][2]string
}

// NewLectionary reads the table genbible generated.
func NewLectionary() (*Lectionary, error) {
	var raw struct {
		Years  map[string]map[string]Lessons `json:"years"`
		Ranges map[string][2]string          `json:"ranges"`
	}
	if err := json.Unmarshal(lectData, &raw); err != nil {
		return nil, err
	}
	return &Lectionary{years: raw.Years, ranges: raw.Ranges}, nil
}

// YearLetter is the cycle a liturgical year reads: Advent 2023 opened
// Year A, and the letters turn A, B, C by the year.
func YearLetter(year int) string {
	switch ((year % 3) + 3) % 3 {
	case 1:
		return "A"
	case 2:
		return "B"
	default:
		return "C"
	}
}

// For finds the readings of the day: its Sunday's, or the day's own when
// it is a holy day with readings of its own. It returns the cycle the
// readings are of, their code in the table, and the lessons; ok is false
// for a day the lectionary keeps none for, as an ordinary weekday.
func (l *Lectionary) For(day time.Time) (yearLetter, code string, lessons Lessons, ok bool) {
	day = midnight(day)
	// The liturgical year the day sits in, and the letter it reads.
	year := day.Year() - 1
	for !inYear(day, year) {
		year++
	}
	yearLetter = YearLetter(year)
	code = l.codeOf(day, year)
	if lessons, ok = l.years[yearLetter][code]; !ok {
		// Some days have no readings in the letter's year, as the eighth
		// Sunday after the Epiphany has only in the years that have one:
		// the table's other years carry the day.
		for _, other := range []string{"A", "B", "C", "holydays"} {
			if lessons, ok = l.years[other][code]; ok {
				yearLetter = other
				break
			}
		}
	}
	return yearLetter, code, lessons, ok
}

// codeOf works out the day's code in the table, by where the calendar has
// it.
func (l *Lectionary) codeOf(day time.Time, year int) string {
	// The calendar works in UTC days: the reader's local date, read as
	// its own day here.
	y, m, d := day.Date()
	day = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	// A holy day reads its own, whatever day of the week it falls on.
	if code := holydayCode(day); code != "" {
		return code
	}
	easter := Easter(year + 1)
	ash := easter.AddDate(0, 0, -46)
	palm := easter.AddDate(0, 0, -7)
	holySat := easter.AddDate(0, 0, -1)
	maundy := easter.AddDate(0, 0, -3)
	goodFriday := easter.AddDate(0, 0, -2)
	pentecost := easter.AddDate(0, 0, 49)
	ascension := easter.AddDate(0, 0, 39)
	trinity := easter.AddDate(0, 0, 56)
	first := advent1(year)
	lastEpiphany := ash.AddDate(0, 0, -3) // the Sunday before Ash Wednesday
	epiphany := time.Date(year+1, time.January, 6, 0, 0, 0, 0, utc)
	baptism := nextSunday(epiphany)

	// Holy days of the week of days that hang on Easter.
	switch day {
	case ash:
		return "ash-wednesday"
	case maundy:
		return "maundy"
	case goodFriday:
		return "good-friday"
	case holySat:
		return "holy-saturday"
	case ascension:
		return "ascension"
	}
	if day.Weekday() != time.Sunday {
		return "" // an ordinary weekday keeps no readings
	}

	// The Sundays, in the order the year reads them.
	switch {
	case !day.Before(first) && day.Before(time.Date(year, time.December, 25, 0, 0, 0, 0, utc)):
		// Advent's four.
		return "advent" + strconv.Itoa(weeksBetween(first, day)+1)
	case day.Equal(christTheKing(year + 1)):
		return "christ-the-king"
	case day.After(trinity):
		if code, ok := l.properCode(day); ok {
			return code
		}
		return ""
	case day.Equal(trinity):
		return "trinity"
	case day.Equal(pentecost):
		return "pentecost"
	case day.After(easter):
		// The Sundays of Easter, second through seventh.
		return "easter" + strconv.Itoa(weeksBetween(easter, day))
	case day.Equal(palm):
		return "palm"
	case day.After(lastEpiphany) && day.Before(palm):
		// The Sundays in Lent, after the last Sunday after the Epiphany.
		return "lent" + strconv.Itoa(weeksBetween(lastEpiphany, day))
	case day.Equal(lastEpiphany):
		return "last-after-epiphany"
	case day.After(baptism) && day.Before(lastEpiphany):
		return "after-epiphany" + strconv.Itoa(weeksBetween(baptism, day))
	case day.Equal(baptism):
		return "baptism"
	case day.Equal(epiphany):
		return "epiphany"
	case day.After(time.Date(year, time.December, 25, 0, 0, 0, 0, utc)):
		// The Sundays after Christmas, first and second.
		if day.Month() == time.December || day.Day() <= 1 {
			return "after-christmas1"
		}
		return "after-christmas2"
	}
	return ""
}

// properCode finds the code of an after-Pentecost Sunday by its date
// range, as the table keeps them: "Sunday between July 10 and July 16
// inclusive" is proper-7-10, of July 10 to July 16.
func (l *Lectionary) properCode(day time.Time) (string, bool) {
	md := normMD(day.Format("1-02"))
	for code, r := range l.ranges {
		if md >= normMD(r[0]) && md <= normMD(r[1]) {
			return code, true
		}
	}
	return "", false
}

// normMD pads a generated "7-10" to "07-10", what the dates above
// compare by.
func normMD(s string) string {
	month, day, _ := strings.Cut(s, "-")
	if len(month) < 2 {
		month = "0" + month
	}
	if len(day) < 2 {
		day = "0" + day
	}
	return month + "-" + day
}

// holydayCode is the code of a day that is a holy day itself, whatever
// week it sits in, or "" for one that isn't.
func holydayCode(day time.Time) string {
	_, month, d := day.Date()
	switch {
	case month == time.January && d == 1:
		return "holy-name"
	case month == time.January && d == 6:
		return "epiphany"
	case month == time.February && d == 2:
		return "presentation"
	case month == time.March && d == 25:
		return "annunciation"
	case month == time.May && d == 31:
		return "visitation"
	case month == time.September && d == 14:
		return "holy-cross"
	case month == time.November && d == 1:
		return "all-saints"
	case month == time.December && d == 25:
		return "christmas1"
	}
	return ""
}
