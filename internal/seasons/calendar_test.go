package seasons

import (
	"testing"
	"time"
)

// day is a UTC day of the calendar, for the tests to ask of.
func day(year int, month time.Month, d int) time.Time {
	return time.Date(year, month, d, 12, 0, 0, 0, time.UTC)
}

// Easter comes where the computus puts it: the table is of the Western
// calendar's own years, the extremes included.
func TestEaster(t *testing.T) {
	for _, c := range []struct {
		year  int
		month time.Month
		day   int
	}{
		{2000, time.April, 23},
		{2008, time.March, 23},
		{2016, time.March, 27},
		{2024, time.March, 31},
		{2025, time.April, 20},
		{2026, time.April, 5},
		{2027, time.March, 28},
		{2038, time.April, 25},
	} {
		got := Easter(c.year)
		if got.Month() != c.month || got.Day() != c.day {
			t.Errorf("Easter %d = %s %d, want %s %d", c.year, got.Month(), got.Day(), c.month, c.day)
		}
	}
}

// Lent opens on Ash Wednesday, walks through Palm Sunday and Good Friday,
// and Easter ends the week after: the feasts of a liturgical year read as
// the published calendar has them. The year 2026 runs Advent 2026 to
// Christ the King, November 2027.
func TestFeasts2026(t *testing.T) {
	for _, c := range []struct {
		code  string
		month time.Month
		day   int
		year  int
	}{
		{"advent1", time.November, 29, 2026},
		{"christmas", time.December, 25, 2026},
		{"epiphany", time.January, 6, 2027},
		{"ash-wednesday", time.February, 10, 2027},
		{"palm-sunday", time.March, 21, 2027},
		{"good-friday", time.March, 26, 2027},
		{"easter", time.March, 28, 2027},
		{"pentecost", time.May, 16, 2027},
		{"christ-the-king", time.November, 21, 2027},
	} {
		found := false
		for _, f := range Feasts(2026) {
			if f.Code == c.code {
				found = true
				if f.Date.Month() != c.month || f.Date.Day() != c.day || f.Date.Year() != c.year {
					t.Errorf("%s = %s %d %d, want %s %d %d", c.code, f.Date.Month(), f.Date.Day(), f.Date.Year(), c.month, c.day, c.year)
				}
			}
		}
		if !found {
			t.Errorf("no %s among the feasts of 2026", c.code)
		}
	}
}

// A day in the weeks before Advent sits in Ordinary Time, green, and the
// countdown points at Christ the King.
func TestTodayOrdinary(t *testing.T) {
	season, feast, next, days := Today(day(2026, time.November, 20))
	if season.Code != "ordinary" || season.Color != Green {
		t.Errorf("the season reads %s, %s", season.Code, season.Color)
	}
	if feast != nil {
		t.Errorf("November 20 reads as %s", feast.Name)
	}
	if next.Code != "christ-the-king" || days != 2 {
		t.Errorf("the next feast reads %s in %d days", next.Code, days)
	}
}

// Christmas Day is its own feast, and sits in the Christmas season.
func TestTodayChristmas(t *testing.T) {
	season, feast, _, _ := Today(day(2026, time.December, 25))
	if season.Code != "christmas" {
		t.Errorf("the season reads %s", season.Code)
	}
	if feast == nil || feast.Code != "christmas" {
		t.Errorf("the day reads as %v", feast)
	}
}

// The first Advent Sunday opens the year's purple.
func TestTodayAdvent(t *testing.T) {
	season, _, next, days := Today(day(2026, time.November, 30))
	if season.Code != "advent" {
		t.Errorf("the season reads %s", season.Code)
	}
	if next.Code != "christmas" || days != 25 {
		t.Errorf("the next feast reads %s in %d days", next.Code, days)
	}
}

// Holy Week wears red the week before Easter.
func TestTodayHolyWeek(t *testing.T) {
	season, feast, _, _ := Today(day(2027, time.March, 26))
	if season.Code != "holy-week" || season.Color != Red {
		t.Errorf("the season reads %s, %s", season.Code, season.Color)
	}
	if feast == nil || feast.Code != "good-friday" {
		t.Errorf("the day reads as %v", feast)
	}
}

// The Easter season runs white from Easter Day to the eve of Trinity.
func TestTodayEasterSeason(t *testing.T) {
	season, feast, _, _ := Today(day(2027, time.April, 10))
	if season.Code != "easter" || season.Color != White {
		t.Errorf("the season reads %s, %s", season.Code, season.Color)
	}
	if feast != nil {
		t.Errorf("April 10 reads as %s", feast.Name)
	}
}

// A January day sits in the year before its own: the Epiphany of January
// 2028 is a feast of the year 2027.
func TestTodayJanuary(t *testing.T) {
	season, feast, _, _ := Today(day(2028, time.January, 6))
	if season.Code != "epiphany" {
		t.Errorf("the season reads %s", season.Code)
	}
	if feast == nil || feast.Code != "epiphany" {
		t.Errorf("the day reads as %v", feast)
	}
}
