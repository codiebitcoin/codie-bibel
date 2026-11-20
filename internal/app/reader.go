package app

import (
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
	"codie-bibel/internal/speech"
)

// readerState is where the Read page stands: the translation, the book
// and chapter it's at, the verses loaded, the highlights of the chapter,
// the verse its menu is open on, and the search.
type readerState struct {
	trID    string
	tr      *bible.Translation
	code    string // the book open, "" for the books grid
	chapter int    // the chapter open, 0 for the chapter list
	book    *bible.Book

	highlights map[int]bool
	menuAt     int // the verse whose menu shows, 0 for none
	noteAt     int // the verse whose note field shows, 0 for none
	noteDraft  string

	trOpen bool // the translation list shows

	searchQ      string
	searched     string // what was searched, as the results show
	results      []bible.Hit
	searchScroll ui.ScrollState
	scroll       ui.ScrollState
	chScroll     ui.ScrollState
}

// reader is the Read page's state, made once.
func (a *app) reader() *readerState {
	if a.rd == nil {
		a.rd = &readerState{}
	}
	return a.rd
}

// ensureTranslation opens the translation the settings name, or the
// first the app ships with, and keeps which it chose.
func (a *app) ensureTranslation() *bible.Translation {
	rs := a.reader()
	if rs.tr != nil && rs.trID == rs.tr.ID {
		return rs.tr
	}
	id := a.settings.Translation
	if id == "" {
		if list := a.bibles.List(); len(list) > 0 {
			id = list[0].ID
		}
	}
	tr, err := a.bibles.Get(id)
	if err != nil {
		// The settings named a translation that isn't there: the first
		// that is.
		if list := a.bibles.List(); len(list) > 0 {
			tr, _ = a.bibles.Get(list[0].ID)
		}
	}
	if tr == nil {
		return nil
	}
	rs.trID, rs.tr = tr.ID, tr
	return tr
}

// setTranslation switches what the pages read, reloading the chapter.
func (a *app) setTranslation(id string) {
	a.settings.Translation = id
	a.settings.save()
	rs := a.reader()
	rs.tr = nil
	rs.trID = ""
	if tr := a.ensureTranslation(); tr != nil && rs.code != "" && rs.chapter > 0 {
		a.openChapter(rs.code, rs.chapter)
	}
	a.redraw()
}

// openBook shows a book's chapters.
func (a *app) openBook(code string) {
	rs := a.reader()
	rs.code, rs.chapter, rs.menuAt, rs.noteAt = code, 0, 0, 0
	rs.chScroll = ui.ScrollState{}
	a.redraw()
}

// openChapter loads a chapter: its verses, its highlights.
func (a *app) openChapter(code string, chapter int) {
	tr := a.ensureTranslation()
	if tr == nil {
		return
	}
	b := tr.Book(code)
	if b == nil || chapter < 1 || chapter > len(b.Chapters) {
		return
	}
	rs := a.reader()
	rs.code, rs.chapter = code, chapter
	rs.book = b
	rs.menuAt, rs.noteAt, rs.noteDraft = 0, 0, ""
	rs.chScroll = ui.ScrollState{}
	if a.store != nil {
		rs.highlights, _ = a.store.Highlights(tr.ID, code, chapter)
	} else {
		rs.highlights = nil
	}
	a.redraw()
}

// nextChapter and prevChapter walk the chapters round the canon.
func (a *app) nextChapter() {
	rs := a.reader()
	if rs.code == "" || rs.chapter == 0 {
		return
	}
	n := len(rs.book.Chapters)
	if rs.chapter < n {
		a.openChapter(rs.code, rs.chapter+1)
		return
	}
	// The first chapter of the next book, round to Genesis.
	i := 0
	for k, c := range bible.CanonCodes {
		if c == rs.code {
			i = k
			break
		}
	}
	next := bible.CanonCodes[(i+1)%len(bible.CanonCodes)]
	a.openChapter(next, 1)
}

func (a *app) prevChapter() {
	rs := a.reader()
	if rs.code == "" || rs.chapter == 0 {
		return
	}
	if rs.chapter > 1 {
		a.openChapter(rs.code, rs.chapter-1)
		return
	}
	i := 0
	for k, c := range bible.CanonCodes {
		if c == rs.code {
			i = k
			break
		}
	}
	prev := bible.CanonCodes[(i-1+len(bible.CanonCodes))%len(bible.CanonCodes)]
	n := bible.Canon[prev].Chapters
	a.openChapter(prev, n)
}

// verseText reads a verse of the open chapter.
func (a *app) verseText(verse int) string {
	rs := a.reader()
	text, _ := rs.tr.Verse(rs.code, rs.chapter, verse)
	return text
}

// toggleHighlight marks a verse, or unmarks it.
func (a *app) toggleHighlight(verse int) {
	if a.store == nil || a.reader().tr == nil {
		return
	}
	rs := a.reader()
	on := !rs.highlights[verse]
	if err := a.store.Highlight(rs.tr.ID, rs.code, rs.chapter, verse, on); err != nil {
		return
	}
	if rs.highlights == nil {
		rs.highlights = map[int]bool{}
	}
	rs.highlights[verse] = on
	a.redraw()
}

// saveNote writes the note draft to its verse, and closes the field.
func (a *app) saveNote() {
	rs := a.reader()
	if rs.noteAt == 0 || a.store == nil {
		rs.noteAt = 0
		return
	}
	a.store.SetNote(rs.tr.ID, rs.code, rs.chapter, rs.noteAt, rs.noteDraft)
	rs.noteAt, rs.noteDraft = 0, ""
	a.redraw()
}

// listenChapter plans the chapter into turns and queues it to the
// device, from where it last stood.
func (a *app) listenChapter() {
	rs := a.reader()
	if rs.tr == nil || rs.code == "" || rs.chapter == 0 {
		return
	}
	ref := bible.Ref{Book: rs.code, Chapter: rs.chapter}
	title := ref.String()
	passage := strings.Join(rs.book.Chapters[rs.chapter-1], " ")
	at := 0
	if a.store != nil {
		at = a.store.Progress(a.listenKey())
	}
	plan := speech.NewPlan(passage, at)
	a.bq.Enqueue("read:"+rs.trID+":"+rs.code+":"+strconv.Itoa(rs.chapter), title, plan.Turns, plan.At)
	a.redraw()
}

// listenRef plans a single reference's passage and queues it.
func (a *app) listenRef(ref bible.Ref, title string) {
	tr := a.ensureTranslation()
	if tr == nil {
		return
	}
	endCh := ref.EndChapter
	if endCh == 0 {
		endCh = ref.Chapter
	}
	var lines []string
	for ch := ref.Chapter; ch <= endCh; ch++ {
		b := tr.Book(ref.Book)
		if b == nil || ch > len(b.Chapters) {
			continue
		}
		verses := b.Chapters[ch-1]
		from, to := 1, len(verses)
		if ch == ref.Chapter && ref.Verse > 0 {
			from = ref.Verse
		}
		if ch == endCh && ref.EndVerse > 0 {
			to = ref.EndVerse
		}
		if from > len(verses) {
			continue
		}
		if to > len(verses) {
			to = len(verses)
		}
		lines = append(lines, strings.Join(verses[from-1:to], " "))
	}
	if len(lines) == 0 {
		return
	}
	passage := strings.Join(lines, " ")
	at := 0
	if a.store != nil {
		at = a.store.Progress("read:" + ref.Book + strconv.Itoa(ref.Chapter))
	}
	plan := speech.NewPlan(passage, at)
	a.bq.Enqueue("read:"+ref.Book+strconv.Itoa(ref.Chapter), title, plan.Turns, plan.At)
	a.redraw()
}

// listenKey is where the open chapter's reading stands in the store.
func (a *app) listenKey() string {
	rs := a.reader()
	return "read:" + rs.trID + ":" + rs.code + ":" + strconv.Itoa(rs.chapter)
}

// queueProgress is what the queue's movements keep: the reading's place.
func (a *app) queueProgress() {
	for _, it := range a.bq.Items() {
		if strings.HasPrefix(it.ID, "read:") && a.store != nil {
			a.store.SetProgress(it.ID[len("read:"):], it.At)
		}
	}
}

// search runs the query over the translation the pages read, into the
// results.
func (a *app) search(query string) {
	tr := a.ensureTranslation()
	if tr == nil {
		return
	}
	rs := a.reader()
	rs.searched = query
	rs.results = bible.Search(tr, query, 200)
	rs.searchScroll = ui.ScrollState{}
	a.redraw()
}

// searchOpen shows the results of the last search, if there was one.
func (a *app) searchOpen() bool {
	rs := a.reader()
	return rs.searched != ""
}

// openSearchResult opens the chapter a result is of.
func (a *app) openSearchResult(h bible.Hit) {
	a.openChapter(h.Ref.Book, h.Ref.Chapter)
}

// bibleRef reads a reference as people do, for the tray's line.
func bibleRef(code string, chapter, verse int) (string, bool) {
	ref := bible.Ref{Book: code, Chapter: chapter, Verse: verse, EndVerse: verse}
	return ref.String(), true
}
