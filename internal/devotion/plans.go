package devotion

import (
	"time"

	"codie-bibel/internal/bible"
)

// Plan is a reading plan: a passage a day for the days it runs. Its days
// are worked out once, from the canon, and read the same every year.
type Plan struct {
	ID   string
	Name string
	Days [][]bible.Ref
}

// plans are the app's reading plans, made once.
var plans = []*Plan{
	throughTheBible(),
	newTestament90(),
	adventureOfAdvent(),
}

// throughTheBible walks the whole Bible in a year: the 1,189 chapters in
// their order, three a day and four where the year's remainder falls.
func throughTheBible() *Plan {
	p := &Plan{ID: "bible-year", Name: "Through the Bible in a Year"}
	refs := allChapters()
	const days = 364
	extra := len(refs) - days*3
	at := 0
	for d := 0; d < days && at < len(refs); d++ {
		n := 3
		if d < extra {
			n = 4
		}
		p.Days = append(p.Days, refs[at:min(at+n, len(refs))])
		at += n
	}
	return p
}

// newTestament90 reads the New Testament in ninety days, three chapters a
// day where they come.
func newTestament90() *Plan {
	p := &Plan{ID: "nt-90", Name: "The New Testament in 90 Days"}
	var refs []bible.Ref
	for _, code := range bible.CanonCodes {
		if bible.Canon[code].OT {
			continue
		}
		for ch := 1; ch <= bible.Canon[code].Chapters; ch++ {
			refs = append(refs, bible.Ref{Book: code, Chapter: ch})
		}
	}
	const days = 90
	per := len(refs) / days
	extra := len(refs) - days*per
	at := 0
	for d := 0; d < days && at < len(refs); d++ {
		n := per
		if d < extra {
			n++
		}
		p.Days = append(p.Days, refs[at:min(at+n, len(refs))])
		at += n
	}
	return p
}

// adventureOfAdvent keeps the twenty-four days of Advent: the promises of
// old to the child in the manger.
func adventureOfAdvent() *Plan {
	p := &Plan{ID: "advent", Name: "Advent: the Road to Bethlehem"}
	one := func(code string, ch, from, to int) []bible.Ref {
		return []bible.Ref{{Book: code, Chapter: ch, Verse: from, EndVerse: to}}
	}
	p.Days = [][]bible.Ref{
		one("GEN", 3, 14, 15), one("GEN", 12, 1, 3), one("GEN", 22, 15, 18),
		one("GEN", 49, 8, 10), one("DEU", 18, 15, 18), one("2SA", 7, 8, 16),
		one("PSA", 2, 1, 12), one("PSA", 22, 1, 18), one("PSA", 72, 1, 11),
		one("ISA", 7, 10, 14), one("ISA", 9, 1, 7), one("ISA", 11, 1, 9),
		one("MIC", 5, 1, 5), one("ISA", 40, 1, 5), one("ISA", 40, 9, 11),
		one("ISA", 52, 13, 15), one("ISA", 53, 1, 12), one("JER", 23, 5, 6),
		one("HAG", 2, 6, 9), one("MAL", 3, 1, 5), one("LUK", 1, 5, 25),
		one("LUK", 1, 26, 56), one("LUK", 2, 1, 20), one("MAT", 2, 1, 12),
	}
	return p
}

// allChapters is every chapter of the canon in its order.
func allChapters() []bible.Ref {
	var refs []bible.Ref
	for _, code := range bible.CanonCodes {
		for ch := 1; ch <= bible.Canon[code].Chapters; ch++ {
			refs = append(refs, bible.Ref{Book: code, Chapter: ch})
		}
	}
	return refs
}

// Plans are the reading plans, in the order the list shows them.
func Plans() []*Plan {
	return plans
}

// ByID is the plan of that ID, or nil for one that isn't there.
func PlanByID(id string) *Plan {
	for _, p := range plans {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// ForDay is the day's passages of the plan, by the day of the year; a
// plan shorter than the year reads it through again.
func (p *Plan) ForDay(day time.Time) []bible.Ref {
	if len(p.Days) == 0 {
		return nil
	}
	return p.Days[(day.YearDay()-1)%len(p.Days)]
}
