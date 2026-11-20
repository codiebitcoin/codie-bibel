package seasons

import "codie-bibel/internal/bible"

// feastVerses are the scriptures a feast reads, what its day is about.
// They are the app's own picking, the classic passages of each feast.
var feastVersesTable = map[string][]bible.Ref{
	"advent1":         {{Book: "ISA", Chapter: 2, Verse: 1, EndVerse: 5}},
	"christmas":       {{Book: "LUK", Chapter: 2, Verse: 1, EndVerse: 20}},
	"epiphany":        {{Book: "MAT", Chapter: 2, Verse: 1, EndVerse: 12}},
	"ash-wednesday":   {{Book: "JOEL", Chapter: 2, Verse: 12, EndVerse: 19}},
	"palm-sunday":     {{Book: "MAT", Chapter: 21, Verse: 1, EndVerse: 11}},
	"maundy-thursday": {{Book: "JHN", Chapter: 13, Verse: 1, EndVerse: 20}},
	"good-friday":     {{Book: "JHN", Chapter: 19, Verse: 16, EndVerse: 30}},
	"easter":          {{Book: "MAT", Chapter: 28, Verse: 1, EndVerse: 10}},
	"ascension":       {{Book: "LUK", Chapter: 24, Verse: 50, EndVerse: 53}},
	"pentecost":       {{Book: "ACT", Chapter: 2, Verse: 1, EndVerse: 21}},
	"reformation":     {{Book: "JHN", Chapter: 8, Verse: 31, EndVerse: 36}},
	"all-saints":      {{Book: "REV", Chapter: 7, Verse: 9, EndVerse: 17}},
	"christ-the-king": {{Book: "JHN", Chapter: 18, Verse: 33, EndVerse: 37}},
}

// FeastVerses are the scriptures of a feast by its code, or none.
func FeastVerses(code string) []bible.Ref {
	return feastVersesTable[code]
}
