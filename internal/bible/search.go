package bible

// Hit is one place a search found.
type Hit struct {
	Ref   Ref
	Text  string
	Start int // where the query begins in Text, in runes
}

// Search reads a translation through for the query, verse by verse, in the
// order the Bible has its books. The text is small enough to read whole:
// a translation is a few million runes, and a scan of it is well under a
// tenth of a second. Cased scripts match without case; scripts without
// case, as Chinese, match as they are.
func Search(t *Translation, query string, limit int) []Hit {
	query = fold(query)
	if query == "" {
		return nil
	}
	var hits []Hit
	for _, b := range t.Books {
		name, _ := BookName(b.Code)
		_ = name
		for ch, verses := range b.Chapters {
			for vs, text := range verses {
				i := indexOfFold(text, query)
				if i < 0 {
					continue
				}
				hits = append(hits, Hit{
					Ref:   Ref{Book: b.Code, Chapter: ch + 1, Verse: vs + 1, EndVerse: vs + 1},
					Text:  text,
					Start: i,
				})
				if len(hits) >= limit {
					return hits
				}
			}
		}
	}
	return hits
}

// CountVerses is how many verses a translation has, for the app to show
// what a search walked.
func (t *Translation) CountVerses() int {
	n := 0
	for _, b := range t.Books {
		for _, ch := range b.Chapters {
			n += len(ch)
		}
	}
	return n
}
