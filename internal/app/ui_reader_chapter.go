package app

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
)

// The chapter: its verses with their numbers in superscript, the words
// of Christ in red, the reader's highlights behind their verses, and the
// verse's menu — listen, highlight, note, copy.

// readerColumn is the width the reading holds, as a book reads: a long
// line tires the eye.
const readerColumn = float32(680)

// chapterView is the chapter open: the breadcrumb, the verses in a
// column a book would hold, and the footer.
func (a *app) chapterView(c *ui.Context) {
	rs := a.reader()
	if rs.book == nil {
		return
	}
	verses := rs.book.Chapters[rs.chapter-1]
	ui.Column(c).Grow(1).Gap(8).Children(func() {
		a.breadcrumb(c)
		ui.Scroll(c.Key("verses")).TrackScroll(&rs.chScroll).Grow(1).Children(func() {
			ui.Column(c).Padding(4, 0).Gap(2).MaxWidth(readerColumn).Children(func() {
				if a.settings.Paragraph {
					a.paragraphView(c, verses)
				} else {
					for v, text := range verses {
						verse := v + 1
						a.verseRow(c, verse, text)
						if rs.menuAt == verse {
							a.verseMenu(c, verse)
						}
						if rs.noteAt == verse {
							a.noteField(c, verse)
						}
					}
				}
			})
		})
		a.chapterFooter(c)
	})
}

// paragraphView reads the chapter as one flowing paragraph, the verse
// numbers superscript inline, as a printed Bible lays it out.
func (a *app) paragraphView(c *ui.Context, verses []string) {
	ui.RichText(c, paragraphRuns(a, verses)...).FontSize(a.readerFontSize()).
		TextColor(a.pal.text).LineHeight(1.7).MaxWidth(readerColumn)
}

// paragraphRuns is the whole chapter as inline runs, superscript numbers
// and the red of Christ's words among them.
func paragraphRuns(a *app, verses []string) []ui.Span {
	var runs []ui.Span
	red := ui.Hex("#c0392b")
	if a.pal.dark {
		red = ui.Hex("#e17055")
	}
	for v, text := range verses {
		verse := v + 1
		verseColor := a.pal.text
		if bible.InRed(a.reader().code, a.reader().chapter, verse) {
			verseColor = red
		}
		runs = append(runs, ui.Span{Text: strconv.Itoa(verse) + " ", Size: a.readerFontSize() - 5, Weight: 600, Color: a.pal.accent})
		runs = append(runs, ui.Span{Text: text + " ", Color: verseColor})
	}
	return runs
}

// breadcrumb says where the reader stands: Read, the book, the chapter —
// each of the first two a way back.
func (a *app) breadcrumb(c *ui.Context) {
	rs := a.reader()
	name, _ := bible.BookName(rs.code)
	ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
		if ui.Button(c.Key("crumb-books"), T("Read")).Clicked() {
			rs.code, rs.chapter, rs.menuAt, rs.noteAt = "", 0, 0, 0
			a.redraw()
		}
		ui.Text(c, "›").TextColor(a.pal.faint).FontSize(12)
		if ui.Button(c.Key("crumb-book"), name.En).Clicked() {
			rs.chapter, rs.menuAt, rs.noteAt = 0, 0, 0
			a.redraw()
		}
		ui.Text(c, "›").TextColor(a.pal.faint).FontSize(12)
		ui.Text(c, strconv.Itoa(rs.chapter)).FontSize(13).FontWeight(600)
		ui.Box(c).Grow(1)
		// Verse by verse, or flowing as a paragraph.
		if a.settings.Paragraph {
			if ui.Button(c.Key("layout"), T("Verse by verse")).Clicked() {
				a.settings.Paragraph = false
				a.settings.save()
				a.redraw()
			}
		} else if ui.Button(c.Key("layout"), T("As a paragraph")).Clicked() {
			a.settings.Paragraph = true
			a.settings.save()
			a.redraw()
		}
	})
}

// verseRow is one verse: its number in superscript, its words — in red
// where they are Christ's, on the highlight's wash where the reader
// marked it — and a tap opening its menu.
func (a *app) verseRow(c *ui.Context, verse int, text string) {
	rs := a.reader()
	row := ui.Row(c.Key("verse-row-"+strconv.Itoa(verse))).AlignItems(ui.Start).Gap(4).
		Padding(2, 4).Radius(6).Label("")
	clicked := row.Clicked()
	bg := ui.Transparent
	if rs.highlights[verse] {
		bg = a.pal.accent.Alpha(0.14)
	} else if row.Hovered() {
		bg = a.pal.hover
	}
	row.Background(bg)
	row.Children(func() {
		ui.Text(c, strconv.Itoa(verse)).FontSize(a.readerFontSize() - 4).
			TextColor(a.pal.accent).FontWeight(600)
		runs := a.verseRuns(verse, text)
		ui.RichText(c, runs...).FontSize(a.readerFontSize()).
			TextColor(a.pal.text).Grow(1).MaxWidth(readerColumn).LineHeight(1.65)
	})
	if clicked {
		rs.menuAt, rs.noteAt = verse, 0
		a.redraw()
	}
}

// verseRuns cuts a verse into the runs the reader shows: the words of
// Christ in red, the rest in the text's color.
func (a *app) verseRuns(verse int, text string) []ui.Span {
	rs := a.reader()
	textColor := a.pal.text
	red := ui.Hex("#c0392b")
	if a.pal.dark {
		red = ui.Hex("#e17055")
	}
	runs := []ui.Span{{Text: text, Color: textColor}}
	if rs.tr != nil && bible.InRed(rs.code, rs.chapter, verse) {
		runs[0].Color = red
	}
	// A note reads as a mark of its own, a small ° in the accent.
	if a.store != nil {
		if _, ok, _ := a.store.Note(rs.tr.ID, rs.code, rs.chapter, verse); ok {
			runs = append(runs, ui.Span{Text: "  °", Color: a.pal.accent})
		}
	}
	return runs
}

// verseMenu is what a tapped verse offers: listen, highlight, note,
// copy.
func (a *app) verseMenu(c *ui.Context, verse int) {
	ui.Row(c.Key("verse-menu")).Padding(4, 20).Gap(8).Children(func() {
		a.listenButton(c, "menu-listen", func() {
			rs := a.reader()
			ref := bible.Ref{Book: rs.code, Chapter: rs.chapter, Verse: verse, EndVerse: verse}
			name, _ := bible.BookName(rs.code)
			a.listenRef(ref, name.En+" "+strconv.Itoa(rs.chapter)+":"+strconv.Itoa(verse))
			rs.menuAt = 0
			a.redraw()
		})
		marked := a.reader().highlights[verse]
		label := T("Highlight")
		var b ui.Element
		if marked {
			b = ui.PrimaryButton(c.Key("menu-highlight"), label)
		} else {
			b = ui.Button(c.Key("menu-highlight"), label)
		}
		if b.Clicked() {
			a.toggleHighlight(verse)
			rs := a.reader()
			rs.menuAt = 0
			a.redraw()
		}
		if ui.Button(c.Key("menu-note"), T("Note")).Clicked() {
			rs := a.reader()
			rs.menuAt, rs.noteAt = 0, verse
			if a.store != nil {
				rs.noteDraft, _, _ = a.store.Note(rs.tr.ID, rs.code, rs.chapter, verse)
			}
			a.redraw()
		}
		if ui.Button(c.Key("menu-copy"), T("Copy")).Clicked() {
			rs := a.reader()
			ref := bible.Ref{Book: rs.code, Chapter: rs.chapter, Verse: verse, EndVerse: verse}
			copyText(`"` + a.verseText(verse) + `" — ` + ref.String())
			rs.menuAt = 0
			a.redraw()
		}
	})
}

// noteField is the field a verse's note goes into, with its save.
func (a *app) noteField(c *ui.Context, verse int) {
	rs := a.reader()
	ui.Column(c.Key("verse-note")).Padding(4, 20).Gap(6).Children(func() {
		ui.TextInput(c.Key("note-input"), &rs.noteDraft).Placeholder(T("A note for this verse")).Grow(1).FontSize(13)
		ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
			if ui.Button(c.Key("note-save"), T("Save")).Clicked() {
				a.saveNote()
			}
			if ui.Button(c.Key("note-drop"), T("Cancel")).Clicked() {
				rs.noteAt, rs.noteDraft = 0, ""
				a.redraw()
			}
		})
	})
}

// chapterFooter walks to the chapters before and after, the reference
// between.
func (a *app) chapterFooter(c *ui.Context) {
	rs := a.reader()
	ref := bible.Ref{Book: rs.code, Chapter: rs.chapter}
	ui.Row(c).Gap(10).AlignItems(ui.Center).Padding(6, 0).Children(func() {
		if ui.Button(c.Key("prev-chapter"), "‹").Clicked() {
			a.prevChapter()
		}
		ui.Text(c, ref.String()).FontSize(13).FontWeight(600)
		if ui.Button(c.Key("next-chapter"), "›").Clicked() {
			a.nextChapter()
		}
		ui.Box(c).Grow(1)
		a.listenButton(c, "listen-chapter", func() { a.listenChapter() })
	})
}

// readerFontSize is how big the chapter's text reads, as the settings
// keep.
func (a *app) readerFontSize() float32 {
	if a.settings.FontSize > 0 {
		return float32(a.settings.FontSize)
	}
	return 17
}
