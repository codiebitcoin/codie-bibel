package app

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
)

// The Read page: the books grid, then a book's chapters, then the
// chapter itself, with the translation picker and the search above them.
// The chapter and the results are their own files.

func (a *app) readPage(c *ui.Context) {
	a.pageHeading(c, T("Read"), func() { a.translationPicker(c) })
	top, bottom := a.insets(c)
	body := ui.Column(c).Grow(1).Padding(top, 24, bottom, 24).Gap(12)
	body.Children(func() {
		a.continueCard(c)
		a.searchField(c)
		switch {
		case a.searchOpen():
			a.searchResults(c)
		case a.reader().chapter > 0:
			a.chapterView(c)
		case a.reader().code != "":
			a.chapterList(c)
		default:
			a.booksGrid(c)
		}
	})
}

// continueCard is the way back to where the reader stood: the last
// chapter read, one tap to go on.
func (a *app) continueCard(c *ui.Context) {
	if a.store == nil || a.reader().code != "" {
		return
	}
	book, _, _ := a.store.Get("last:book")
	chStr, _, _ := a.store.Get("last:chapter")
	ch := 0
	if chStr != "" {
		ch, _ = strconv.Atoi(chStr)
	}
	if book == "" || ch < 1 {
		return
	}
	name, _ := bible.BookName(book)
	ref := bible.Ref{Book: book, Chapter: ch}
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		if ui.PrimaryButton(c.Key("continue"), T("Continue reading")+": "+ref.String()).Clicked() {
			a.openChapter(book, ch)
		}
	})
	_ = name
}

// translationPicker is the button over the pages that drops the list of
// translations, the one reading checked.
func (a *app) translationPicker(c *ui.Context) {
	tr := a.ensureTranslation()
	name := ""
	if tr != nil {
		name = tr.Short
	}
	chosen := ""
	if tr != nil {
		chosen = tr.ID
	}
	sel := ui.SelectBase(c.Key("translation"), &chosen)
	b := sel.Trigger.Height(30).PaddingX(10).Gap(6).Radius(8).AlignItems(ui.Center).
		Background(a.pal.cardShade).Label(T("Translation") + ": " + name)
	b.Children(func() {
		ui.Text(c, name).FontSize(13).FontWeight(600)
		ui.Text(c, "▾").FontSize(10).TextColor(a.pal.faint)
	})
	sel.Popup(func(panel ui.Element) {
		panel.Margin(6, 0, 0, 0).Padding(5).Gap(1).MinWidth(220).Radius(12).
			Background(a.pal.bg).Border(1, a.pal.hairline).Shadow(0, 12, 32, -6, ui.RGBA(0, 0, 0, 0.25))
		lastLang := ""
		for _, info := range a.bibles.List() {
			if info.Language != lastLang {
				lastLang = info.Language
				ui.Text(c, info.Language).FontSize(10.5).FontWeight(600).TextColor(a.pal.faint).Padding(6, 8, 2, 8)
			}
			item := sel.Item(info.ID).Height(32).PaddingX(10).Gap(8).Radius(7).AlignItems(ui.Center).Label(info.Name + " (" + info.Short + ")")
			if item.Highlighted() {
				item.Background(a.pal.press)
			}
			item.Children(func() {
				ui.Text(c, info.Name).FontSize(13)
				ui.Box(c).Grow(1)
				ui.Text(c, info.Short).FontSize(11.5).TextColor(a.pal.faint)
			})
		}
	})
	if sel.Trigger.Changed() {
		a.setTranslation(chosen)
	}
}

// searchField is the field the query goes into, with the search beside
// it; typing anew and submitting runs it again.
func (a *app) searchField(c *ui.Context) {
	rs := a.reader()
	ui.Row(c).Gap(8).Children(func() {
		in := ui.TextInput(c.Key("search-field"), &rs.searchQ).Placeholder(T("Search")).Width(280).FontSize(13.5)
		if in.Submitted() {
			a.search(rs.searchQ)
		}
		if ui.Button(c.Key("search-run"), T("Search")).Clicked() {
			a.search(rs.searchQ)
		}
		if a.searchOpen() {
			if ui.Button(c.Key("search-close"), T("Back")).Clicked() {
				rs.searched, rs.results = "", nil
				a.redraw()
			}
		}
	})
}

// booksGrid is the canon: the Old Testament's books, then the New's,
// each its button.
func (a *app) booksGrid(c *ui.Context) {
	rs := a.reader()
	tr := a.ensureTranslation()
	if tr == nil {
		a.centered(c, T("No Bible is in the app."), nil)
		return
	}
	ui.Scroll(c.Key("books")).TrackScroll(&rs.scroll).Grow(1).Children(func() {
		for _, group := range []struct {
			name string
			ot   bool
		}{{T("Old Testament"), true}, {T("New Testament"), false}} {
			ui.Column(c).Gap(6).Padding(8, 0).Children(func() {
				ui.Text(c, group.name).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
				ui.Row(c).Wrap().Gap(6).Children(func() {
					for _, code := range bibleCodesOT(bible.CanonCodes, group.ot) {
						name, _ := bible.BookName(code)
						if ui.Button(c.Key("book-"+code), name.En).Clicked() {
							a.openBook(code)
						}
					}
				})
			})
		}
	})
}

// bibleCodesOT takes the codes of one testament, in order.
func bibleCodesOT(codes []string, ot bool) []string {
	var out []string
	for _, c := range codes {
		if bible.Canon[c].OT == ot {
			out = append(out, c)
		}
	}
	return out
}

// chapterList is a book's chapters, eight a row.
func (a *app) chapterList(c *ui.Context) {
	rs := a.reader()
	tr := a.ensureTranslation()
	b := tr.Book(rs.code)
	if b == nil {
		return
	}
	name, _ := bible.BookName(rs.code)
	ui.Scroll(c.Key("chapters")).TrackScroll(&rs.scroll).Grow(1).Children(func() {
		ui.Column(c).Gap(8).Children(func() {
			ui.Text(c, name.En).FontSize(17).FontWeight(700)
			for start := 1; start <= len(b.Chapters); start += 8 {
				end := start + 8
				if end > len(b.Chapters)+1 {
					end = len(b.Chapters) + 1
				}
				ui.Row(c).Wrap().Gap(6).Children(func() {
					for ch := start; ch < end; ch++ {
						label := strconv.Itoa(ch)
						if ui.Button(c.Key("chapter-"+strconv.Itoa(ch)), label).Clicked() {
							a.openChapter(rs.code, ch)
						}
					}
				})
			}
		})
	})
}
