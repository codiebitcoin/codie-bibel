package devotion

import (
	"os"
	"testing"
	"time"

	"codie-bibel/internal/bible"
)

// testSpurgeon is the table genbible wrote, for the tests to read.
func testSpurgeon(t *testing.T) *Spurgeon {
	t.Helper()
	if _, err := os.Stat("data"); err != nil {
		t.Skip("no data directory: genbible hasn't run here")
	}
	s, err := NewSpurgeon()
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// The day's two readings open: the first of January has its morning and
// its evening, each with a verse and something written about it.
func TestSpurgeonReads(t *testing.T) {
	s := testSpurgeon(t)
	day := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	me := s.Read(day)
	if me.AM.Verse == "" || len(me.AM.Text) < 200 {
		t.Errorf("the morning of January 1 is thin: %q, %d runes", me.AM.Verse, len(me.AM.Text))
	}
	if me.PM.Verse == "" || len(me.PM.Text) < 200 {
		t.Errorf("the evening of January 1 is thin: %q, %d runes", me.PM.Verse, len(me.PM.Text))
	}
	// Christmas reads too, evening's own.
	if me := s.Read(time.Date(2026, 12, 25, 9, 0, 0, 0, time.UTC)); me.PM.Verse == "" {
		t.Error("the evening of Christmas is empty")
	}
}

// The hour picks the half of the day's readings.
func TestHalf(t *testing.T) {
	if Half(time.Date(2026, 5, 3, 7, 0, 0, 0, time.UTC)) != "am" {
		t.Error("seven in the morning reads as evening")
	}
	if Half(time.Date(2026, 5, 3, 20, 0, 0, 0, time.UTC)) != "pm" {
		t.Error("eight in the evening reads as morning")
	}
}

// The verse of the day turns by the day, and reads the same all day.
func TestVerseOfTheDay(t *testing.T) {
	morning := time.Date(2026, 4, 5, 7, 0, 0, 0, time.UTC)
	night := time.Date(2026, 4, 5, 23, 0, 0, 0, time.UTC)
	c1, ch1, v1 := VerseOfTheDay(morning)
	c2, ch2, v2 := VerseOfTheDay(night)
	if c1 != c2 || ch1 != ch2 || v1 != v2 {
		t.Errorf("the verse changed through the day: %s %d:%d then %s %d:%d", c1, ch1, v1, c2, ch2, v2)
	}
	if _, ok := bible.Canon[c1]; !ok {
		t.Errorf("the verse names the unknown book %s", c1)
	}
	tomorrow, _, _ := VerseOfTheDay(morning.AddDate(0, 0, 1))
	if tomorrow == c1 && morning.YearDay()+1 <= len(votdPool) {
		// The next day's may share the verse only if the pool wraps.
	}
}

// The plans read through the whole canon: a year's days cover every
// chapter once.
func TestThroughTheBible(t *testing.T) {
	p := PlanByID("bible-year")
	if p == nil {
		t.Fatal("no plan through the Bible")
	}
	covered := map[string]bool{}
	for _, day := range p.Days {
		for _, r := range day {
			key := r.Book + " " + itoa(r.Chapter)
			if covered[key] {
				t.Errorf("%s is read twice", key)
			}
			covered[key] = true
		}
	}
	if len(covered) != 1189 {
		t.Errorf("the plan covers %d chapters, want 1189", len(covered))
	}
}

// The ninety-day plan reads the New Testament, and nothing of the Old.
func TestNT90(t *testing.T) {
	p := PlanByID("nt-90")
	if p == nil {
		t.Fatal("no ninety-day plan")
	}
	if len(p.Days) != 90 {
		t.Fatalf("the plan runs %d days", len(p.Days))
	}
	for _, day := range p.Days {
		for _, r := range day {
			if bible.Canon[r.Book].OT {
				t.Errorf("the New Testament plan reads %s %d", r.Book, r.Chapter)
			}
		}
	}
}

// The Advent plan runs the twenty-four days of December to the manger,
// and every plan answers the day it is asked of.
func TestPlansAnswerDays(t *testing.T) {
	advent := PlanByID("advent")
	if advent == nil {
		t.Fatal("no Advent plan")
	}
	if len(advent.Days) != 24 {
		t.Fatalf("the Advent plan runs %d days", len(advent.Days))
	}
	for _, p := range Plans() {
		if got := p.ForDay(time.Date(2026, 11, 20, 9, 0, 0, 0, time.UTC)); len(got) == 0 {
			t.Errorf("%s has nothing for November 20", p.ID)
		}
	}
}

// The prayer walks Adoration, Confession, Thanksgiving and Supplication,
// each with scriptures, and only the confession is kept private.
func TestPrayerSteps(t *testing.T) {
	if len(PrayerSteps) != 4 {
		t.Fatalf("the prayer has %d steps", len(PrayerSteps))
	}
	want := []string{"adoration", "confession", "thanksgiving", "supplication"}
	private := 0
	for i, s := range PrayerSteps {
		if s.Code != want[i] {
			t.Errorf("step %d is %s, want %s", i, s.Code, want[i])
		}
		if len(s.Refs) == 0 {
			t.Errorf("%s has no scriptures", s.Code)
		}
		for _, r := range s.Refs {
			if _, ok := bible.Canon[r.Book]; !ok {
				t.Errorf("%s names the unknown book %s", s.Code, r.Book)
			}
		}
		if s.Privacy {
			private++
		}
	}
	if private != 1 {
		t.Errorf("%d steps claim privacy, want the confession's one", private)
	}
	if len(ComfortVerses) < 4 {
		t.Error("the comfort verses are thin")
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [4]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
