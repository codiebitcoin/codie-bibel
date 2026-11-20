// Package seasons is the church year: the feasts and the seasons they
// make, when they fall, and what the church reads in them. Easter and
// everything that hangs on it comes from the Gregorian computus; the rest
// are fixed days. The year runs from Advent to Christ the King, so a
// November day sits in the year that names itself by its Christmas.
package seasons

import "time"

// color is a liturgical color, what the season is clothed in.
type color string

const (
	white  color = "white"  // Christmas, Epiphany, Easter: joy
	green  color = "green"  // the growing seasons: Ordinary Time
	purple color = "purple" // Advent, Lent: waiting and penitence
	red    color = "red"    // Holy Week, Pentecost: blood and fire
)

// Feast is one day of the church year.
type Feast struct {
	Name string // as it reads, in English
	Code string // stable key, as "easter"
	Date time.Time
}

// Season is a stretch of the year, from a day to the day before the next
// season's, both ends in the day.
type Season struct {
	Name  string
	Code  string
	Color color
	From  time.Time
	To    time.Time
}

// UTC is where the calendar's days live: the dates the tables make, read
// whatever the wall clock says.
var utc = time.UTC

// Easter is the day of Easter of the year, by the Gregorian computus
// (Meeus, Jones and Butcher).
func Easter(year int) time.Time {
	a := year % 19
	b, c := year/100, year%100
	d, e := b/4, b%4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i, k := c/4, c%4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := (h+l-7*m+114)%31 + 1
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, utc)
}

// advent1 is the year's First Sunday of Advent: the Sunday on or after
// November 27.
func advent1(year int) time.Time {
	return sundayOnOrAfter(time.Date(year, time.November, 27, 0, 0, 0, 0, utc))
}

// christTheKing is the year's last Sunday, a week before its Advent.
func christTheKing(year int) time.Time {
	return advent1(year).AddDate(0, 0, -7)
}

// sundayOnOrAfter is the first Sunday on or after the day.
func sundayOnOrAfter(t time.Time) time.Time {
	for t.Weekday() != time.Sunday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

// nextSunday walks to the next Sunday strictly after t.
func nextSunday(t time.Time) time.Time {
	return sundayOnOrAfter(t.AddDate(0, 0, 1))
}

// Feasts are the feasts of the liturgical year that names itself by the
// Advent it opens with: its Advent and Christmas at the end of the
// calendar year, the rest of them in the months to Easter and past it, in
// the order they come.
func Feasts(year int) []Feast {
	easter := Easter(year + 1)
	jan := func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, utc)
	}
	feasts := []Feast{
		{Code: "advent1", Name: "First Sunday of Advent", Date: advent1(year)},
		{Code: "christmas", Name: "Christmas", Date: jan(year, time.December, 25)},
		{Code: "epiphany", Name: "Epiphany", Date: jan(year+1, time.January, 6)},
		{Code: "ash-wednesday", Name: "Ash Wednesday", Date: easter.AddDate(0, 0, -46)},
		{Code: "palm-sunday", Name: "Palm Sunday", Date: easter.AddDate(0, 0, -7)},
		{Code: "maundy-thursday", Name: "Maundy Thursday", Date: easter.AddDate(0, 0, -3)},
		{Code: "good-friday", Name: "Good Friday", Date: easter.AddDate(0, 0, -2)},
		{Code: "easter", Name: "Easter", Date: easter},
		{Code: "ascension", Name: "Ascension", Date: easter.AddDate(0, 0, 39)},
		{Code: "pentecost", Name: "Pentecost", Date: easter.AddDate(0, 0, 49)},
		{Code: "trinity", Name: "Trinity Sunday", Date: easter.AddDate(0, 0, 56)},
		{Code: "reformation", Name: "Reformation Day", Date: jan(year+1, time.October, 31)},
		{Code: "all-saints", Name: "All Saints' Day", Date: jan(year+1, time.November, 1)},
		{Code: "christ-the-king", Name: "Christ the King", Date: christTheKing(year + 1)},
	}
	sortFeasts(feasts)
	return feasts
}

// sortFeasts puts the feasts in the order their days come.
func sortFeasts(feasts []Feast) {
	for i := 1; i < len(feasts); i++ {
		for j := i; j > 0 && before(feasts[j].Date, feasts[j-1].Date); j-- {
			feasts[j], feasts[j-1] = feasts[j-1], feasts[j]
		}
	}
}

// Seasons are the liturgical year's stretches, from its Advent to the eve
// of the next one's. Lent holds Holy Week; a day of Holy Week reads Holy
// Week, whose red is the color the week wears.
func Seasons(year int) []Season {
	easter := Easter(year + 1)
	ash := easter.AddDate(0, 0, -46)
	palm := easter.AddDate(0, 0, -7)
	holySat := easter.AddDate(0, 0, -1)
	pentecost := easter.AddDate(0, 0, 49)
	trinity := easter.AddDate(0, 0, 56)
	start := advent1(year)
	epiphanyStart := time.Date(year+1, time.January, 6, 0, 0, 0, 0, utc)
	epiphanyEnd := nextSunday(epiphanyStart).AddDate(0, 0, -1)
	return []Season{
		{Name: "Advent", Code: "advent", Color: purple, From: start, To: time.Date(year, time.December, 24, 0, 0, 0, 0, utc)},
		{Name: "Christmas", Code: "christmas", Color: white, From: time.Date(year, time.December, 25, 0, 0, 0, 0, utc), To: time.Date(year+1, time.January, 5, 0, 0, 0, 0, utc)},
		{Name: "Epiphany", Code: "epiphany", Color: white, From: epiphanyStart, To: epiphanyEnd},
		{Name: "Ordinary Time", Code: "ordinary", Color: green, From: nextSunday(epiphanyStart), To: ash.AddDate(0, 0, -1)},
		{Name: "Lent", Code: "lent", Color: purple, From: ash, To: holySat},
		{Name: "Holy Week", Code: "holy-week", Color: red, From: palm, To: holySat},
		{Name: "Easter", Code: "easter", Color: white, From: easter, To: trinity.AddDate(0, 0, -1)},
		{Name: "Pentecost", Code: "pentecost", Color: red, From: pentecost, To: trinity.AddDate(0, 0, -1)},
		{Name: "Ordinary Time", Code: "ordinary", Color: green, From: trinity, To: end(year)},
	}
}

// end is the last day of the liturgical year, the Saturday before the
// next year's Advent.
func end(year int) time.Time {
	return advent1(year+1).AddDate(0, 0, -1)
}

// Today is where a day sits: the season it's in, the feast it is (if it is
// one), and the feast that comes next, with how many days off it is. The
// liturgical year a day belongs to names itself by the Christmas it ends
// with, so a January day sits in the year before its own.
func Today(day time.Time) (season Season, feast *Feast, next Feast, daysTo int) {
	day = midnight(day)
	// The year a day belongs to is never later than its own: a November
	// day sits in the year whose Advent opened weeks before it.
	year := day.Year() - 1
	for !inYear(day, year) {
		year++
	}
	season = seasonOf(day, year)
	for _, f := range Feasts(year) {
		if same(f.Date, day) {
			copy := f
			feast = &copy
		}
	}
	next, daysTo = nextFeast(day, year)
	if daysTo < 0 {
		next, daysTo = nextFeast(day, year+1)
	}
	return season, feast, next, daysTo
}

// inYear says whether the day sits in the liturgical year that Advent
// starts.
func inYear(day time.Time, year int) bool {
	return !before(day, advent1(year)) && before(day, advent1(year+1))
}

// seasonOf picks the season of the day among the year's, whose stretch
// holds it by the making of the year. The seasons read from the end, so
// the week that holds Holy Week inside Lent reads Holy Week, whose red is
// what the week wears.
func seasonOf(day time.Time, year int) Season {
	seasons := Seasons(year)
	for i := len(seasons) - 1; i >= 0; i-- {
		s := seasons[i]
		if !before(day, s.From) && !after(day, s.To) {
			return s
		}
	}
	return Season{Code: "ordinary", Color: green, Name: "Ordinary Time", From: day, To: day}
}

func nextFeast(day time.Time, year int) (Feast, int) {
	for _, f := range Feasts(year) {
		if after(f.Date, day) {
			return f, days(day, f.Date)
		}
	}
	return Feast{}, -1
}

func midnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func same(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

// before says a reads earlier than b, by its day.
func before(a, b time.Time) bool { return days(a, b) > 0 }

// after says a reads later than b, by its day.
func after(a, b time.Time) bool { return days(a, b) < 0 }

// days is how many days b lies ahead of a.
func days(a, b time.Time) int {
	return int(midnight(b).Sub(midnight(a)).Hours() / 24)
}

// weeksBetween is how many whole weeks b lies ahead of a.
func weeksBetween(a, b time.Time) int {
	return days(a, b) / 7
}
