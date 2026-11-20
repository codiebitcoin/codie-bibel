package seasons

import (
	"testing"
	"time"
)

// testLectionary is the table genbible wrote, for the tests to read.
func testLectionary(t *testing.T) *Lectionary {
	t.Helper()
	l, err := NewLectionary()
	if err != nil {
		t.Fatal(err)
	}
	return l
}

// A Sunday's readings come by its place in the year: Advent 1 of the year
// 2026 is Year A, and reads Isaiah's swords into ploughshares.
func TestForAdvent1(t *testing.T) {
	l := testLectionary(t)
	// Advent 1 of the year 2026 falls on November 29, 2026.
	year, code, lessons, ok := l.For(time.Date(2026, time.November, 29, 12, 0, 0, 0, time.UTC))
	if !ok {
		t.Fatal("no readings for the First Sunday of Advent")
	}
	if year != "A" || code != "advent1" {
		t.Errorf("the readings are of year %s, code %s", year, code)
	}
	if lessons[0] != "Isaiah 2:1-5" || lessons[3] != "Matthew 24:36-44" {
		t.Errorf("the lessons read %v", lessons)
	}
}

// A Sunday deep in the after-Pentecost stretch reads by its date range.
func TestForAfterPentecost(t *testing.T) {
	l := testLectionary(t)
	// July 12, 2026 is a Sunday between July 10 and July 16, of the
	// liturgical year 2025, which is Year C.
	year, code, _, ok := l.For(time.Date(2026, time.July, 12, 12, 0, 0, 0, time.UTC))
	if !ok {
		t.Fatal("no readings for the Sunday of July 12")
	}
	if year != "C" {
		t.Errorf("the readings are of year %s", year)
	}
	if code != "proper-7-10" {
		t.Errorf("the code reads %s", code)
	}
}

// Good Friday carries its own readings, though it is no Sunday.
func TestForGoodFriday(t *testing.T) {
	l := testLectionary(t)
	_, code, lessons, ok := l.For(time.Date(2027, time.March, 26, 12, 0, 0, 0, time.UTC))
	if !ok {
		t.Fatal("no readings for Good Friday")
	}
	if code != "good-friday" {
		t.Errorf("the code reads %s", code)
	}
	if len(lessons[0]) == 0 {
		t.Error("Good Friday's first lesson is empty")
	}
}

// A weekday that is no holy day keeps no readings.
func TestForOrdinaryWeekday(t *testing.T) {
	l := testLectionary(t)
	_, _, _, ok := l.For(time.Date(2026, time.November, 17, 12, 0, 0, 0, time.UTC))
	if ok {
		t.Error("a Tuesday in November came up with readings")
	}
}

// The letters turn by the year: 2023 opened Year A, and 2025 reads Year C.
func TestYearLetter(t *testing.T) {
	for _, c := range []struct {
		year   int
		letter string
	}{{2023, "A"}, {2024, "B"}, {2025, "C"}, {2026, "A"}, {2027, "B"}} {
		if got := YearLetter(c.year); got != c.letter {
			t.Errorf("YearLetter(%d) = %s, want %s", c.year, got, c.letter)
		}
	}
}
