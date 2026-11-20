package devotion

import "codie-bibel/internal/bible"

// The guided prayer walks Adoration, Confession, Thanksgiving and
// Supplication, each step with its scriptures to read and to pray by.
// The device reads the scriptures aloud; what is confessed stays on this
// machine.

// PrayerStep is one of the prayer's steps.
type PrayerStep struct {
	Code    string // stable key, as "confession"
	Name    string // as it reads, in English
	Prompt  string // what the step asks of the one praying
	Refs    []bible.Ref
	Privacy bool // what is written here stays on this machine
}

// PrayerSteps are the prayer's steps, in the order the prayer walks.
var PrayerSteps = []PrayerStep{
	{Code: "adoration", Name: "Adoration", Prompt: "Praise God for who he is.",
		Refs: []bible.Ref{{Book: "PSA", Chapter: 95, Verse: 1, EndVerse: 7}, {Book: "PSA", Chapter: 145, Verse: 1, EndVerse: 9}, {Book: "REV", Chapter: 4, Verse: 8, EndVerse: 11}}},
	{Code: "confession", Name: "Confession", Prompt: "Name what you have done, and what you have left undone.", Privacy: true,
		Refs: []bible.Ref{{Book: "PSA", Chapter: 51, Verse: 1, EndVerse: 10}, {Book: "1JN", Chapter: 1, Verse: 8, EndVerse: 10}, {Book: "PSA", Chapter: 32, Verse: 1, EndVerse: 5}, {Book: "PRO", Chapter: 28, Verse: 13}}},
	{Code: "thanksgiving", Name: "Thanksgiving", Prompt: "Thank God for what he has given, and what he has carried you through.",
		Refs: []bible.Ref{{Book: "1TH", Chapter: 5, Verse: 16, EndVerse: 18}, {Book: "PSA", Chapter: 103, Verse: 1, EndVerse: 5}, {Book: "PHP", Chapter: 4, Verse: 6}}},
	{Code: "supplication", Name: "Supplication", Prompt: "Ask, for yourself and for others.",
		Refs: []bible.Ref{{Book: "MAT", Chapter: 6, Verse: 9, EndVerse: 13}, {Book: "JAS", Chapter: 1, Verse: 5}, {Book: "PHP", Chapter: 4, Verse: 19}, {Book: "1TI", Chapter: 2, Verse: 1, EndVerse: 2}}},
}

// ComfortVerses are the scriptures that speak to a heavy heart, that the
// device may read one aloud. They answer to what the prayer carried, in
// the order they are likely to fit: none is a judgment, and none is read
// back to the one who wrote it as an answer from God.
var ComfortVerses = []bible.Ref{
	{Book: "MAT", Chapter: 11, Verse: 28, EndVerse: 30},
	{Book: "PSA", Chapter: 34, Verse: 18},
	{Book: "ISA", Chapter: 41, Verse: 10},
	{Book: "PSA", Chapter: 147, Verse: 3},
	{Book: "2CO", Chapter: 1, Verse: 3, EndVerse: 4},
	{Book: "PSA", Chapter: 51, Verse: 17},
	{Book: "ROM", Chapter: 8, Verse: 1},
	{Book: "PSA", Chapter: 103, Verse: 11, EndVerse: 14},
}
