package devotion

import "time"

// Verse of the Day: the pool is the verses the church has long kept for a
// day's encouragement, each by its USFM code and verse. The day's is the
// pool's by a turn of the year, the same for everyone the app reads it on,
// whatever the year.
var votdPool = []struct {
	Code   string
	Ch, Vs int
}{
	{"JHN", 3, 16}, {"ROM", 8, 28}, {"PHP", 4, 13}, {"PRO", 3, 5}, {"ISA", 41, 10},
	{"JER", 29, 11}, {"MAT", 11, 28}, {"PSA", 23, 1}, {"1CO", 13, 4}, {"ROM", 8, 38},
	{"EPH", 2, 8}, {"GAL", 5, 22}, {"HEB", 11, 1}, {"JAS", 1, 5}, {"PSA", 46, 1},
	{"ISA", 40, 31}, {"MAT", 6, 33}, {"2CO", 5, 17}, {"PSA", 119, 105}, {"ROM", 12, 2},
	{"JHN", 14, 6}, {"JHN", 15, 5}, {"PSA", 27, 1}, {"COL", 3, 23}, {"1PE", 5, 7},
	{"ZEP", 3, 17}, {"PSA", 34, 8}, {"LAM", 3, 22}, {"MIC", 6, 8}, {"JOS", 1, 9},
	{"PSA", 121, 1}, {"ISA", 26, 3}, {"HEB", 12, 1}, {"1JN", 1, 9}, {"PSA", 51, 10},
	{"MAT", 5, 16}, {"1CO", 16, 13}, {"EPH", 4, 32}, {"PHP", 4, 6}, {"PSA", 139, 14},
	{"2TI", 1, 7}, {"PSA", 19, 1}, {"ROM", 5, 8}, {"PRO", 16, 3}, {"PSA", 37, 4},
	{"1CO", 10, 31}, {"DEU", 6, 5}, {"JHN", 16, 33}, {"PSA", 100, 4}, {"ISA", 53, 5},
	{"PSA", 91, 1}, {"PRO", 22, 6}, {"MAT", 28, 19}, {"ACT", 1, 8}, {"GAL", 2, 20},
	{"1TH", 5, 16}, {"PSA", 63, 1}, {"HEB", 4, 16}, {"JAS", 4, 7}, {"ROM", 1, 16},
	{"PSA", 62, 5}, {"MAT", 11, 29}, {"PSA", 119, 11}, {"2CO", 12, 9}, {"PSA", 55, 22},
	{"1JN", 4, 19}, {"COL", 3, 2}, {"PSA", 118, 24}, {"ISA", 55, 8}, {"1PE", 2, 9},
	{"PRO", 17, 17}, {"MAT", 6, 14}, {"PSA", 84, 11}, {"2CO", 9, 7}, {"PSA", 32, 8},
	{"TIT", 3, 5}, {"REV", 21, 4}, {"PSA", 147, 3}, {"JER", 33, 3}, {"1SA", 16, 7},
	{"PSA", 3, 5}, {"NEH", 8, 10}, {"PSA", 103, 1}, {"MAL", 3, 10}, {"DAN", 12, 3},
	{"JHN", 1, 1}, {"PSA", 8, 3}, {"PRO", 31, 25}, {"JOL", 2, 25}, {"ZEC", 9, 9},
	{"PSA", 42, 1}, {"1CO", 15, 58}, {"ECC", 3, 1}, {"PSA", 5, 3}, {"HOS", 6, 1},
	{"MAT", 7, 7}, {"PSA", 66, 16}, {"LUK", 1, 37}, {"PSA", 71, 3}, {"JOB", 19, 25},
	{"PSA", 90, 12}, {"ISA", 6, 8}, {"PSA", 25, 4}, {"EXO", 14, 14}, {"PSA", 30, 5},
	{"ACT", 16, 31}, {"PSA", 33, 20}, {"PRO", 18, 10}, {"PSA", 73, 26}, {"SNG", 8, 7},
	{"MAT", 5, 9}, {"PSA", 119, 165}, {"ROM", 15, 13}, {"PSA", 116, 1}, {"JUDE", 1, 24},
	{"PSA", 145, 8}, {"2CH", 7, 14}, {"PSA", 126, 5}, {"1JN", 5, 14}, {"PSA", 138, 1},
}

// VerseOfTheDay is the day's verse: the pool's turn by the day of the
// year, stepped by the year so no day reads two years alike.
func VerseOfTheDay(day time.Time) (code string, ch, vs int) {
	y := day.YearDay() - 1 + day.Year()
	return votdPool[y%len(votdPool)].Code, votdPool[y%len(votdPool)].Ch, votdPool[y%len(votdPool)].Vs
}
