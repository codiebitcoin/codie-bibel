package app

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
)

// The search's results: the reference and the verse, the match marked,
// a tap opening the chapter.
func (a *app) searchResults(c *ui.Context) {
	rs := a.reader()
	if len(rs.results) == 0 {
		ui.Text(c, Tf("Nothing for %q.", rs.searched)).TextColor(a.pal.muted).FontSize(13)
		return
	}
	ui.Text(c, Tf("%d verses have it.", len(rs.results))).TextColor(a.pal.faint).FontSize(12)
	ui.Scroll(c.Key("search-results")).TrackScroll(&rs.searchScroll).Grow(1).Children(func() {
		ui.Column(c).Gap(2).Children(func() {
			for _, h := range rs.results {
				h := h
				row := ui.Row(c.Key("hit-"+h.Ref.String())).Padding(6, 8).Gap(8).Radius(8).AlignItems(ui.Start)
				bg := ui.Transparent
				if row.Hovered() {
					bg = a.pal.hover
				}
				row.Background(bg)
				row.Children(func() {
					ui.Column(c).Gap(2).Children(func() {
						ui.Text(c, h.Ref.String()).FontSize(12.5).FontWeight(600).TextColor(a.pal.accent)
						ui.RichText(c, a.markedRun(h)...).FontSize(13).TextColor(a.pal.text)
					})
				})
				if row.Clicked() {
					a.openSearchResult(h)
				}
			}
		})
	})
}

// markedRun marks the match in a hit's text, by the rune offset the
// search found it at.
func (a *app) markedRun(h bible.Hit) []ui.Span {
	runes := []rune(h.Text)
	from := h.Start
	if from < 0 || from > len(runes) {
		return []ui.Span{{Text: h.Text}}
	}
	to := from
	q := []rune(a.reader().searched)
	if len(q) == 0 {
		q = []rune(a.reader().searched)
	}
	for i := from; i < len(runes) && i-from < len(q); i++ {
		to = i + 1
	}
	out := []ui.Span{}
	if from > 0 {
		out = append(out, ui.Span{Text: string(runes[:from])})
	}
	out = append(out, ui.Span{Text: string(runes[from:to]), Background: a.pal.accent.Alpha(0.25)})
	if to < len(runes) {
		out = append(out, ui.Span{Text: string(runes[to:])})
	}
	return out
}

// listenResults queues the results' verses to the device, as one
// reading.
func (a *app) listenResults() {
	rs := a.reader()
	for i, h := range rs.results {
		ref := bible.Ref{Book: h.Ref.Book, Chapter: h.Ref.Chapter, Verse: h.Ref.Verse, EndVerse: h.Ref.Verse}
		a.listenRef(ref, strconv.Itoa(i+1)+". "+h.Ref.String())
	}
}
