package bible

// canon is one book of the canon: its USFM code, which is the app's stable
// key across languages, and how many chapters it has. The chapter counts are
// of the standard Protestant text; a translation that disagrees with them
// has lost or found something, and genbible refuses it.
type canon struct {
	Chapters int
	OT       bool
}

// Canon is the 66 books, by USFM code.
var Canon = map[string]canon{
	"GEN": {50, true}, "EXO": {40, true}, "LEV": {27, true}, "NUM": {36, true}, "DEU": {34, true},
	"JOS": {24, true}, "JDG": {21, true}, "RUT": {4, true}, "1SA": {31, true}, "2SA": {24, true},
	"1KI": {22, true}, "2KI": {25, true}, "1CH": {29, true}, "2CH": {36, true}, "EZR": {10, true},
	"NEH": {13, true}, "EST": {10, true}, "JOB": {42, true}, "PSA": {150, true}, "PRO": {31, true},
	"ECC": {12, true}, "SNG": {8, true}, "ISA": {66, true}, "JER": {52, true}, "LAM": {5, true},
	"EZK": {48, true}, "DAN": {12, true}, "HOS": {14, true}, "JOL": {3, true}, "AMO": {9, true},
	"OBA": {1, true}, "JON": {4, true}, "MIC": {7, true}, "NAM": {3, true}, "HAB": {3, true},
	"ZEP": {3, true}, "HAG": {2, true}, "ZEC": {14, true}, "MAL": {4, true},
	"MAT": {28, false}, "MRK": {16, false}, "LUK": {24, false}, "JHN": {21, false}, "ACT": {28, false},
	"ROM": {16, false}, "1CO": {16, false}, "2CO": {13, false}, "GAL": {6, false}, "EPH": {6, false},
	"PHP": {4, false}, "COL": {4, false}, "1TH": {5, false}, "2TH": {3, false}, "1TI": {6, false},
	"2TI": {4, false}, "TIT": {3, false}, "PHM": {1, false}, "HEB": {13, false}, "JAS": {5, false},
	"1PE": {5, false}, "2PE": {3, false}, "1JN": {5, false}, "2JN": {1, false}, "3JN": {1, false},
	"JUD": {1, false}, "REV": {22, false},
}

// CanonCodes are the 66 books in their order, OT then NT.
var CanonCodes = []string{
	"GEN", "EXO", "LEV", "NUM", "DEU", "JOS", "JDG", "RUT", "1SA", "2SA",
	"1KI", "2KI", "1CH", "2CH", "EZR", "NEH", "EST", "JOB", "PSA", "PRO",
	"ECC", "SNG", "ISA", "JER", "LAM", "EZK", "DAN", "HOS", "JOL", "AMO",
	"OBA", "JON", "MIC", "NAM", "HAB", "ZEP", "HAG", "ZEC", "MAL",
	"MAT", "MRK", "LUK", "JHN", "ACT", "ROM", "1CO", "2CO", "GAL", "EPH",
	"PHP", "COL", "1TH", "2TH", "1TI", "2TI", "TIT", "PHM", "HEB", "JAS",
	"1PE", "2PE", "1JN", "2JN", "3JN", "JUD", "REV",
}

// bookName is one book's name in the languages the app speaks, with the short
// forms people write.
type bookName struct {
	Code string
	En   string
	EnAb string
	Zh   string
	ZhAb string
}

// BookNames are the 66 books, in canonical order. The abbreviations are what
// a reference may be written with beyond the full name, as "Ps" is Psalms.
var BookNames = []bookName{
	{"GEN", "Genesis", "gen", "创世记", "创"},
	{"EXO", "Exodus", "exod", "出埃及记", "出"},
	{"LEV", "Leviticus", "lev", "利未记", "利"},
	{"NUM", "Numbers", "num", "民数记", "民"},
	{"DEU", "Deuteronomy", "deut", "申命记", "申"},
	{"JOS", "Joshua", "josh", "约书亚记", "书"},
	{"JDG", "Judges", "judg", "士师记", "士"},
	{"RUT", "Ruth", "ruth", "路得记", "得"},
	{"1SA", "1 Samuel", "1sam", "撒母耳记上", "撒上"},
	{"2SA", "2 Samuel", "2sam", "撒母耳记下", "撒下"},
	{"1KI", "1 Kings", "1kgs", "列王纪上", "王上"},
	{"2KI", "2 Kings", "2kgs", "列王纪下", "王下"},
	{"1CH", "1 Chronicles", "1chr", "历代志上", "代上"},
	{"2CH", "2 Chronicles", "2chr", "历代志下", "代下"},
	{"EZR", "Ezra", "ezra", "以斯拉记", "拉"},
	{"NEH", "Nehemiah", "neh", "尼希米记", "尼"},
	{"EST", "Esther", "esth", "以斯帖记", "斯"},
	{"JOB", "Job", "job", "约伯记", "伯"},
	{"PSA", "Psalms", "ps", "诗篇", "诗"},
	{"PRO", "Proverbs", "prov", "箴言", "箴"},
	{"ECC", "Ecclesiastes", "eccl", "传道书", "传"},
	{"SNG", "Song of Songs", "song", "雅歌", "歌"},
	{"ISA", "Isaiah", "isa", "以赛亚书", "赛"},
	{"JER", "Jeremiah", "jer", "耶利米书", "耶"},
	{"LAM", "Lamentations", "lam", "耶利米哀歌", "哀"},
	{"EZK", "Ezekiel", "ezek", "以西结书", "结"},
	{"DAN", "Daniel", "dan", "但以理书", "但"},
	{"HOS", "Hosea", "hos", "何西阿书", "何"},
	{"JOL", "Joel", "joel", "约珥书", "珥"},
	{"AMO", "Amos", "amos", "阿摩司书", "摩"},
	{"OBA", "Obadiah", "obad", "俄巴底亚书", "俄"},
	{"JON", "Jonah", "jonah", "约拿书", "拿"},
	{"MIC", "Micah", "mic", "弥迦书", "弥"},
	{"NAM", "Nahum", "nah", "那鸿书", "鸿"},
	{"HAB", "Habakkuk", "hab", "哈巴谷书", "哈"},
	{"ZEP", "Zephaniah", "zeph", "西番雅书", "番"},
	{"HAG", "Haggai", "hag", "哈该书", "该"},
	{"ZEC", "Zechariah", "zech", "撒迦利亚书", "亚"},
	{"MAL", "Malachi", "mal", "玛拉基书", "玛"},
	{"MAT", "Matthew", "matt", "马太福音", "太"},
	{"MRK", "Mark", "mark", "马可福音", "可"},
	{"LUK", "Luke", "luke", "路加福音", "路"},
	{"JHN", "John", "john", "约翰福音", "约"},
	{"ACT", "Acts", "acts", "使徒行传", "徒"},
	{"ROM", "Romans", "rom", "罗马书", "罗"},
	{"1CO", "1 Corinthians", "1cor", "哥林多前书", "林前"},
	{"2CO", "2 Corinthians", "2cor", "哥林多后书", "林后"},
	{"GAL", "Galatians", "gal", "加拉太书", "加"},
	{"EPH", "Ephesians", "eph", "以弗所书", "弗"},
	{"PHP", "Philippians", "phil", "腓立比书", "腓"},
	{"COL", "Colossians", "col", "歌罗西书", "西"},
	{"1TH", "1 Thessalonians", "1thess", "帖撒罗尼迦前书", "帖前"},
	{"2TH", "2 Thessalonians", "2thess", "帖撒罗尼迦后书", "帖后"},
	{"1TI", "1 Timothy", "1tim", "提摩太前书", "提前"},
	{"2TI", "2 Timothy", "2tim", "提摩太后书", "提后"},
	{"TIT", "Titus", "titus", "提多书", "多"},
	{"PHM", "Philemon", "phlm", "腓利门书", "门"},
	{"HEB", "Hebrews", "heb", "希伯来书", "来"},
	{"JAS", "James", "jas", "雅各书", "雅"},
	{"1PE", "1 Peter", "1pet", "彼得前书", "彼前"},
	{"2PE", "2 Peter", "2pet", "彼得后书", "彼后"},
	{"1JN", "1 John", "1john", "约翰一书", "约一"},
	{"2JN", "2 John", "2john", "约翰二书", "约二"},
	{"3JN", "3 John", "3john", "约翰三书", "约三"},
	{"JUD", "Jude", "jude", "犹大书", "犹"},
	{"REV", "Revelation", "rev", "启示录", "启"},
}

// BookName is the book of a code, in canonical order.
func BookName(code string) (bookName, bool) {
	for _, b := range BookNames {
		if b.Code == code {
			return b, true
		}
	}
	return bookName{}, false
}

// enAliases are the English spellings a reference may use beyond the names
// and abbreviations above.
var enAliases = map[string]string{
	"psalm": "PSA", "psalms": "PSA", "psa": "PSA", "pslm": "PSA", "psm": "PSA",
	"song of solomon": "SNG", "canticles": "SNG", "canticle": "SNG",
	"genesis": "GEN", "exodus": "EXO", "revelation": "REV",
	"jn": "JHN", "jno": "JHN", "matt": "MAT", "mt": "MAT", "mk": "MRK",
	"1co": "1CO", "2co": "2CO", "1cor": "1CO", "2cor": "2CO",
	"lk": "LUK", "lu": "LUK", "james": "JAS", "jude": "JUD",
	"eccl": "ECC", "ecclesiastes": "ECC", "qoheleth": "ECC",
	"1thes": "1TH", "2thes": "2TH", "1th": "1TH", "2th": "2TH",
	"philem": "PHM", "phil": "PHP", "php": "PHP",
	"1 jn": "1JN", "2 jn": "2JN", "3 jn": "3JN",
}

// BookByEnName finds a book by its English name or abbreviation, matched
// without case. It walks BookNames twice: full names, then abbreviations.
func BookByEnName(name string) (string, bool) {
	for _, b := range BookNames {
		if equalFold(b.En, name) {
			return b.Code, true
		}
	}
	for _, b := range BookNames {
		if equalFold(b.EnAb, name) {
			return b.Code, true
		}
	}
	if code, ok := enAliases[lower(name)]; ok {
		return code, true
	}
	return "", false
}

// lower is strings.ToLower, here so the tables read on their own.
func lower(s string) string {
	out := []byte(s)
	for i, c := range out {
		if 'A' <= c && c <= 'Z' {
			out[i] = c + 'a' - 'A'
		}
	}
	return string(out)
}

// BookByZhName finds a book by its Chinese name or abbreviation.
func BookByZhName(name string) (string, bool) {
	for _, b := range BookNames {
		if b.Zh == name {
			return b.Code, true
		}
	}
	for _, b := range BookNames {
		if b.ZhAb == name {
			return b.Code, true
		}
	}
	return "", false
}

// equalFold is strings.EqualFold, here so the tables read on their own.
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
