package bible

// Story is one of the Bible's stories: where it reads, and nothing more.
// The title is the app's own; the passage is the Bible's.
type Story struct {
	Code string // stable key, as "creation"
	Name string // as it reads, in English
	Refs []Ref
}

// Stories are the Bible's stories, in the order the Bible tells them:
// creation to the new creation, for reading aloud or for telling.
var Stories = []Story{
	{Code: "creation", Name: "The Creation", Refs: []Ref{{Book: "GEN", Chapter: 1}, {Book: "GEN", Chapter: 2, Verse: 1, EndVerse: 3}}},
	{Code: "eden", Name: "The Garden of Eden", Refs: []Ref{{Book: "GEN", Chapter: 2, Verse: 4, EndVerse: 25}, {Book: "GEN", Chapter: 3}}},
	{Code: "flood", Name: "Noah and the Flood", Refs: []Ref{{Book: "GEN", Chapter: 6, Verse: 5, EndVerse: 22}, {Book: "GEN", Chapter: 7}, {Book: "GEN", Chapter: 8}, {Book: "GEN", Chapter: 9, Verse: 1, EndVerse: 17}}},
	{Code: "babel", Name: "The Tower of Babel", Refs: []Ref{{Book: "GEN", Chapter: 11, Verse: 1, EndVerse: 9}}},
	{Code: "abraham", Name: "The Call of Abraham", Refs: []Ref{{Book: "GEN", Chapter: 12, Verse: 1, EndVerse: 9}}},
	{Code: "sacrifice-isaac", Name: "The Sacrifice of Isaac", Refs: []Ref{{Book: "GEN", Chapter: 22, Verse: 1, EndVerse: 19}}},
	{Code: "joseph", Name: "Joseph and His Brothers", Refs: []Ref{{Book: "GEN", Chapter: 37}, {Book: "GEN", Chapter: 45, Verse: 1, EndVerse: 28}}},
	{Code: "moses-call", Name: "The Burning Bush", Refs: []Ref{{Book: "EXO", Chapter: 3}}},
	{Code: "exodus", Name: "The Passover and the Exodus", Refs: []Ref{{Book: "EXO", Chapter: 12}, {Book: "EXO", Chapter: 14}}},
	{Code: "ten-commandments", Name: "The Ten Commandments", Refs: []Ref{{Book: "EXO", Chapter: 20}}},
	{Code: "jericho", Name: "The Walls of Jericho", Refs: []Ref{{Book: "JOS", Chapter: 6}}},
	{Code: "david-goliath", Name: "David and Goliath", Refs: []Ref{{Book: "1SA", Chapter: 17}}},
	{Code: "david-jonathan", Name: "David and Jonathan", Refs: []Ref{{Book: "1SA", Chapter: 18, Verse: 1, EndVerse: 16}, {Book: "1SA", Chapter: 20}}},
	{Code: "david-psalm", Name: "The Shepherd King", Refs: []Ref{{Book: "PSA", Chapter: 23}}},
	{Code: "solomon-wisdom", Name: "The Wisdom of Solomon", Refs: []Ref{{Book: "1KI", Chapter: 3, Verse: 3, EndVerse: 28}}},
	{Code: "elijah-mount-carmel", Name: "Elijah on Mount Carmel", Refs: []Ref{{Book: "1KI", Chapter: 18}}},
	{Code: "jonah", Name: "Jonah and the Great Fish", Refs: []Ref{{Book: "JON", Chapter: 1}, {Book: "JON", Chapter: 2}, {Book: "JON", Chapter: 3}}},
	{Code: "fiery-furnace", Name: "The Fiery Furnace", Refs: []Ref{{Book: "DAN", Chapter: 3}}},
	{Code: "lions-den", Name: "Daniel in the Lions' Den", Refs: []Ref{{Book: "DAN", Chapter: 6}}},
	{Code: "esther", Name: "Esther the Queen", Refs: []Ref{{Book: "EST", Chapter: 4}, {Book: "EST", Chapter: 5, Verse: 1, EndVerse: 8}, {Book: "EST", Chapter: 7}}},
	{Code: "job", Name: "Job's Trial and Rest", Refs: []Ref{{Book: "JOB", Chapter: 1}, {Book: "JOB", Chapter: 2, Verse: 1, EndVerse: 10}, {Book: "JOB", Chapter: 42, Verse: 7, EndVerse: 17}}},
	{Code: "annunciation", Name: "The Angel's Message", Refs: []Ref{{Book: "LUK", Chapter: 1, Verse: 26, EndVerse: 56}}},
	{Code: "nativity", Name: "The Birth of Jesus", Refs: []Ref{{Book: "LUK", Chapter: 2, Verse: 1, EndVerse: 21}, {Book: "MAT", Chapter: 2, Verse: 1, EndVerse: 12}}},
	{Code: "baptism-jesus", Name: "The Baptism of Jesus", Refs: []Ref{{Book: "MAT", Chapter: 3, Verse: 13, EndVerse: 17}}},
	{Code: "temptation", Name: "The Temptation in the Wilderness", Refs: []Ref{{Book: "MAT", Chapter: 4, Verse: 1, EndVerse: 11}}},
	{Code: "sermon-mount", Name: "The Sermon on the Mount", Refs: []Ref{{Book: "MAT", Chapter: 5}, {Book: "MAT", Chapter: 6}, {Book: "MAT", Chapter: 7}}},
	{Code: "calming-storm", Name: "The Calming of the Storm", Refs: []Ref{{Book: "MRK", Chapter: 4, Verse: 35, EndVerse: 41}}},
	{Code: "feeding-5000", Name: "Five Loaves and Two Fish", Refs: []Ref{{Book: "JHN", Chapter: 6, Verse: 1, EndVerse: 15}}},
	{Code: "good-samaritan", Name: "The Good Samaritan", Refs: []Ref{{Book: "LUK", Chapter: 10, Verse: 25, EndVerse: 37}}},
	{Code: "prodigal-son", Name: "The Prodigal Son", Refs: []Ref{{Book: "LUK", Chapter: 15, Verse: 11, EndVerse: 32}}},
	{Code: "lost-sheep", Name: "The Lost Sheep", Refs: []Ref{{Book: "LUK", Chapter: 15, Verse: 1, EndVerse: 10}}},
	{Code: "sower", Name: "The Sower", Refs: []Ref{{Book: "MRK", Chapter: 4, Verse: 1, EndVerse: 20}}},
	{Code: "lazarus", Name: "The Raising of Lazarus", Refs: []Ref{{Book: "JHN", Chapter: 11}}},
	{Code: "triumphal-entry", Name: "The Road to Jerusalem", Refs: []Ref{{Book: "LUK", Chapter: 19, Verse: 28, EndVerse: 48}}},
	{Code: "last-supper", Name: "The Last Supper", Refs: []Ref{{Book: "JHN", Chapter: 13}}},
	{Code: "crucifixion", Name: "The Crucifixion", Refs: []Ref{{Book: "JHN", Chapter: 19}}},
	{Code: "resurrection", Name: "The Empty Tomb", Refs: []Ref{{Book: "JHN", Chapter: 20}}},
	{Code: "road-to-emmaus", Name: "The Road to Emmaus", Refs: []Ref{{Book: "LUK", Chapter: 24, Verse: 13, EndVerse: 35}}},
	{Code: "pentecost-story", Name: "The Day of Pentecost", Refs: []Ref{{Book: "ACT", Chapter: 2}}},
	{Code: "paul-road", Name: "The Road to Damascus", Refs: []Ref{{Book: "ACT", Chapter: 9, Verse: 1, EndVerse: 22}}},
	{Code: "new-creation", Name: "The New Creation", Refs: []Ref{{Book: "REV", Chapter: 21, Verse: 1, EndVerse: 7}, {Book: "REV", Chapter: 22, Verse: 1, EndVerse: 5}}},
}
