package bible

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Ref is a place in the Bible: from chapter Verse of Book to chapter
// EndChapter's EndVerse, its ends inclusive. One verse has EndChapter 0,
// which reads as Chapter; one chapter has both Verse and EndVerse 0.
type Ref struct {
	Book       string // USFM code, as GEN
	Chapter    int
	Verse      int
	EndChapter int
	EndVerse   int
}

// String writes the reference as people do, as "John 3:16" or "Psalm 23:1-6".
// It reads the book's name in English; the views write their own language
// from the code and numbers.
func (r Ref) String() string {
	name, _ := BookName(r.Book)
	if name.En == "" {
		name.En = r.Book
	}
	// A reference of one verse ends at the chapter it starts in.
	end := r.EndChapter
	if end == 0 {
		end = r.Chapter
	}
	if r.Verse == 0 {
		if end != r.Chapter {
			return fmt.Sprintf("%s %d-%d", name.En, r.Chapter, end)
		}
		return fmt.Sprintf("%s %d", name.En, r.Chapter)
	}
	if end != r.Chapter {
		return fmt.Sprintf("%s %d:%d-%d:%d", name.En, r.Chapter, r.Verse, end, r.EndVerse)
	}
	if r.EndVerse == r.Verse {
		return fmt.Sprintf("%s %d:%d", name.En, r.Chapter, r.Verse)
	}
	return fmt.Sprintf("%s %d:%d-%d", name.En, r.Chapter, r.Verse, r.EndVerse)
}

// ParseRef reads a reference, as "John 3:16", "Ps 23:1-6", "Psalm 23",
// "Matthew 5-7", "约 3:16" or "创世记 1", in English or Simplified Chinese.
// It finds the book by name and abbreviation; the rest is numbers.
func ParseRef(s string) (Ref, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Ref{}, false
	}
	// The book is whatever leads before the first number; its name may have
	// its own digits, as 1 Corinthians does, so the split walks from the
	// first digit that follows a space or has a colon or dash after it.
	book, rest, ok := splitBook(s)
	if !ok {
		return Ref{}, false
	}
	code, ok := matchBook(book)
	if !ok {
		return Ref{}, false
	}
	if strings.TrimSpace(rest) == "" {
		return Ref{Book: code, Chapter: 1}, false // a book alone isn't a place
	}
	return parseRest(code, strings.TrimSpace(rest))
}

// splitBook splits a reference into the book's name and what follows: the
// lead up to the first chapter number. A digit at the very start is the
// numeral of the book's own name, as the 1 of "1 Corinthians" or "1co 3:16";
// what follows it, spaces in it too, is the book's name. Elsewhere the book
// ends where a number begins, as the 3 of "约3:16" does.
func splitBook(s string) (book, rest string, ok bool) {
	s = strings.TrimSpace(s)
	i := 0
	if i < len(s) && isDigit(s[i]) {
		// The book's own numeral: "1co" has its letter straight after, and
		// "1 Corinthians" its name before the next number comes.
		if len(s) > i+1 && isLetter(s[i+1]) {
			i++
		} else if len(s) > i+1 && s[i+1] == ' ' {
			letter := strings.IndexFunc(s[i+1:], func(r rune) bool { return unicode.IsLetter(r) })
			digit := strings.IndexFunc(s[i+1:], func(r rune) bool { return unicode.IsDigit(r) })
			if letter >= 0 && (digit < 0 || letter < digit) {
				i++
			} else {
				return "", "", false
			}
		} else {
			return "", "", false
		}
	}
	for i < len(s) {
		r, size := utf8.DecodeRuneInString(s[i:])
		if unicode.IsDigit(r) {
			break
		}
		i += size
	}
	book = strings.TrimSpace(s[:i])
	if book == "" {
		return "", "", false
	}
	return book, strings.TrimSpace(s[i:]), true
}

func isDigit(b byte) bool { return '0' <= b && b <= '9' }

// parseRest reads what follows the book: "3:16", "3:16-20", "3", "5-7",
// "3:4-5:12". After a dash, what carries its own chapter reads as chapter
// and verse; else, if the left had a verse, it reads as verses of the same
// chapter, and else as a chapter.
func parseRest(code, rest string) (Ref, bool) {
	from, to, hasTo := strings.Cut(rest, "-")
	ch1, v1, ok := parseChVerse(strings.TrimSpace(from))
	if !ok {
		return Ref{}, false
	}
	r := Ref{Book: code, Chapter: ch1, Verse: v1}
	to = strings.TrimSpace(to)
	if !hasTo {
		if v1 == 0 {
			return r, true // a whole chapter
		}
		r.EndVerse = v1 // one verse reads to itself
		return r, true
	}
	var ch2, v2 int
	if _, _, hasColon := strings.Cut(to, ":"); hasColon {
		ch2, v2, ok = parseChVerse(to)
		if !ok {
			return Ref{}, false
		}
	} else if n, err := strconv.Atoi(to); err == nil && n > 0 {
		if v1 != 0 {
			ch2, v2 = ch1, n
		} else {
			ch2, v2 = n, 0
		}
	} else {
		return Ref{}, false
	}
	r.EndChapter, r.EndVerse = ch2, v2
	if r.EndChapter < r.Chapter || r.EndChapter == r.Chapter && r.EndVerse != 0 && r.EndVerse < r.Verse {
		return Ref{}, false
	}
	return r, true
}

// parseChVerse reads "3" or "3:16" as a chapter and a verse, the verse 0 if
// it isn't there.
func parseChVerse(s string) (ch, vs int, ok bool) {
	if s == "" {
		return 0, 0, false
	}
	if a, b, has := strings.Cut(s, ":"); has {
		ch, err := strconv.Atoi(strings.TrimSpace(a))
		if err != nil {
			return 0, 0, false
		}
		vs, err := strconv.Atoi(strings.TrimSpace(b))
		if err != nil {
			return 0, 0, false
		}
		return ch, vs, ch > 0 && vs > 0
	}
	ch, err := strconv.Atoi(s)
	if err != nil {
		return 0, 0, false
	}
	return ch, 0, ch > 0
}

// matchBook finds a book by name, in English or Simplified Chinese, after
// taking the punctuation and extra spaces out.
func matchBook(name string) (string, bool) {
	name = strings.Join(strings.Fields(name), " ")
	name = strings.Trim(name, ".,;:")
	if name == "" {
		return "", false
	}
	if code, ok := BookByEnName(name); ok {
		return code, true
	}
	if code, ok := BookByZhName(name); ok {
		return code, true
	}
	// Chinese references write the number of a book into the name, as 约一
	// stands for 约翰一书; the tables have those already. A trailing 书 the
	// abbreviations drop is no book less.
	if strings.HasSuffix(name, "书") {
		if code, ok := BookByZhName(name); ok {
			return code, true
		}
	}
	return "", false
}

func indexDigit(s string) int {
	for i, r := range s {
		if unicode.IsDigit(r) {
			return i
		}
	}
	return -1
}

func isLetter(b byte) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z'
}
