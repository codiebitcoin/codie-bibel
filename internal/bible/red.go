package bible

// The words of Christ in red. The table is the speeches' spans, verse to
// verse, as the Gospels lay them out: the Sermon on the Mount, the
// parables and their explanations, the discourses of John, and the sayings
// that stand alone. A span takes its narration with it where speech and
// narration interleave finely, so the red runs to whole verses; it is a
// reader's aid, not a critical edition, and a refined table can replace it
// without a caller changing.

// redSpan is one span of Christ's speech, ends in.
type redSpan struct {
	Code             string
	StartCh, StartVs int
	EndCh, EndVs     int
}

// redSpans are the spans, by Gospel.
var redSpans = []redSpan{
	// Matthew
	{"MAT", 4, 17, 4, 17},
	{"MAT", 5, 1, 7, 29},
	{"MAT", 9, 2, 9, 6},
	{"MAT", 9, 12, 9, 13},
	{"MAT", 10, 5, 10, 42},
	{"MAT", 11, 4, 11, 30},
	{"MAT", 12, 11, 12, 13},
	{"MAT", 12, 25, 12, 37},
	{"MAT", 12, 39, 12, 45},
	{"MAT", 12, 48, 12, 50},
	{"MAT", 13, 3, 13, 52},
	{"MAT", 15, 3, 15, 20},
	{"MAT", 16, 2, 16, 27},
	{"MAT", 17, 7, 17, 27},
	{"MAT", 18, 3, 18, 35},
	{"MAT", 19, 4, 19, 30},
	{"MAT", 20, 18, 20, 28},
	{"MAT", 21, 2, 21, 44},
	{"MAT", 22, 2, 22, 46},
	{"MAT", 23, 1, 23, 39},
	{"MAT", 24, 2, 24, 51},
	{"MAT", 25, 1, 25, 46},
	{"MAT", 26, 2, 26, 27},
	{"MAT", 26, 31, 26, 75},
	{"MAT", 28, 5, 28, 10},
	{"MAT", 28, 18, 28, 20},
	// Mark
	{"MRK", 1, 15, 1, 17},
	{"MRK", 2, 5, 2, 11},
	{"MRK", 2, 17, 2, 17},
	{"MRK", 2, 19, 2, 22},
	{"MRK", 2, 25, 2, 28},
	{"MRK", 3, 4, 3, 5},
	{"MRK", 3, 23, 3, 35},
	{"MRK", 4, 3, 4, 41},
	{"MRK", 5, 8, 5, 41},
	{"MRK", 6, 4, 6, 38},
	{"MRK", 7, 6, 7, 37},
	{"MRK", 8, 12, 8, 38},
	{"MRK", 9, 12, 9, 50},
	{"MRK", 10, 11, 10, 52},
	{"MRK", 11, 11, 11, 37},
	{"MRK", 12, 6, 12, 44},
	{"MRK", 13, 2, 13, 37},
	{"MRK", 14, 6, 14, 72},
	{"MRK", 16, 6, 16, 18},
	// Luke
	{"LUK", 2, 49, 2, 49},
	{"LUK", 4, 4, 4, 27},
	{"LUK", 5, 8, 5, 36},
	{"LUK", 6, 3, 6, 49},
	{"LUK", 7, 9, 7, 50},
	{"LUK", 8, 21, 8, 56},
	{"LUK", 9, 3, 9, 62},
	{"LUK", 10, 2, 10, 42},
	{"LUK", 11, 4, 11, 54},
	{"LUK", 12, 2, 12, 59},
	{"LUK", 13, 1, 13, 35},
	{"LUK", 14, 3, 14, 35},
	{"LUK", 15, 3, 15, 32},
	{"LUK", 16, 2, 16, 31},
	{"LUK", 17, 3, 17, 37},
	{"LUK", 18, 6, 18, 43},
	{"LUK", 19, 5, 19, 48},
	{"LUK", 20, 3, 20, 47},
	{"LUK", 21, 3, 21, 36},
	{"LUK", 22, 19, 22, 71},
	{"LUK", 23, 28, 23, 46},
	{"LUK", 24, 6, 24, 49},
	// John
	{"JHN", 1, 39, 1, 51},
	{"JHN", 2, 4, 2, 19},
	{"JHN", 3, 3, 3, 21},
	{"JHN", 3, 34, 3, 36},
	{"JHN", 4, 7, 4, 53},
	{"JHN", 5, 6, 5, 47},
	{"JHN", 6, 3, 6, 70},
	{"JHN", 7, 6, 7, 46},
	{"JHN", 8, 2, 8, 58},
	{"JHN", 9, 3, 9, 41},
	{"JHN", 10, 1, 10, 39},
	{"JHN", 11, 4, 11, 43},
	{"JHN", 12, 7, 12, 50},
	{"JHN", 13, 3, 13, 38},
	{"JHN", 14, 1, 14, 31},
	{"JHN", 15, 1, 15, 27},
	{"JHN", 16, 1, 16, 33},
	{"JHN", 17, 1, 17, 26},
	{"JHN", 18, 4, 18, 38},
	{"JHN", 19, 9, 19, 11},
	{"JHN", 20, 15, 20, 29},
	{"JHN", 21, 5, 21, 23},
	// The voice from heaven that Paul heard.
	{"ACT", 9, 4, 9, 6},
}

// red is the spans, by book and chapter, for InRed to walk.
var red = func() map[string][]redSpan {
	m := map[string][]redSpan{}
	for _, s := range redSpans {
		m[s.Code] = append(m[s.Code], s)
	}
	return m
}()

// InRed says whether a verse is among the words of Christ: the reader
// shows those in red, as Bibles in the West have long printed them.
func InRed(code string, chapter, verse int) bool {
	for _, s := range red[code] {
		in := chapter > s.StartCh || chapter == s.StartCh && verse >= s.StartVs
		end := chapter < s.EndCh || chapter == s.EndCh && verse <= s.EndVs
		if in && end {
			return true
		}
	}
	return false
}
