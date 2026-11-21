package app

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
	"codie-bibel/internal/devotion"
)

// The Pray page: the ACTS walk, its steps' scriptures, and the prayers
// the reader keeps. What the confession step holds stays on this
// machine; the device hears only the comfort verses the queue reads.

func (a *app) prayPage(c *ui.Context) {
	a.pageHeading(c, T("Pray"), nil)
	top, bottom := a.insets(c)
	ui.Scroll(c.Key("pray")).Grow(1).Children(func() {
		ui.Column(c).Padding(top, 24, bottom, 24).Gap(16).Children(func() {
			a.prayerWalk(c)
			a.comfortCard(c)
			a.prayerListCard(c)
		})
	})
}

// prayerWalk is the ACTS walk: the step at prayAt, and the steps read
// before it as read.
func (a *app) prayerWalk(c *ui.Context) {
	steps := devotion.PrayerSteps
	at := a.prayStep
	if at >= len(steps) {
		at = len(steps) - 1
	}
	for i, step := range steps {
		chosen := i == at
		row := ui.Column(c.Key("pray-step-"+step.Code)).Padding(12, 14).Gap(8).Radius(12).
			Background(cardBg(a.pal, chosen))
		row.Children(func() {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, strconv.Itoa(i+1)+". "+step.Name).FontSize(14).FontWeight(600)
				ui.Box(c).Grow(1)
				// Any step opens on its name; the walk never traps.
				if !chosen {
					if ui.Button(c.Key("pray-open-"+step.Code), T("Open")).Clicked() {
						a.prayStep = i
						a.redraw()
					}
				}
				if chosen {
					if i > 0 && ui.Button(c.Key("pray-back-"+step.Code), T("Back")).Clicked() {
						a.prayStep--
						a.redraw()
					}
					if !step.Privacy && ui.Button(c.Key("pray-next-"+step.Code), T("Next")).Clicked() {
						a.prayStep++
						a.redraw()
					}
				}
			})
			if chosen {
				ui.Text(c, step.Prompt).FontSize(13).TextColor(a.pal.text)
				for _, r := range step.Refs {
					a.prayerVerseRow(c, r)
				}
				if step.Privacy {
					ui.Text(c, T("What you write here stays on this machine; nothing of it leaves it.")).FontSize(12).TextColor(a.pal.faint)
					ui.TextInput(c.Key("confession-input"), &a.prayText).Placeholder(T("What you confess")).Grow(1).FontSize(13)
					ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
						if ui.Button(c.Key("confession-save"), TC("Save", "prayers")).Clicked() {
							a.prayText = "" // kept nowhere: the field empties, the machine forgets
							a.prayStep++
							a.redraw()
						}
					})
				}
			}
		})
	}
}

// prayerVerseRow is one scripture of a step: its reference, its words,
// and the device reading it.
func (a *app) prayerVerseRow(c *ui.Context, r bible.Ref) {
	tr := a.ensureTranslation()
	name, _ := bible.BookName(r.Book)
	line := r.String()
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, line).FontSize(12.5).FontWeight(600).TextColor(a.pal.accent)
		ui.Box(c).Grow(1)
		if ui.Button(c.Key("pray-listen-"+line), T("Listen")).Clicked() {
			a.listenRef(r, name.En+" "+line)
		}
	})
	if tr != nil {
		endCh := r.EndChapter
		if endCh == 0 {
			endCh = r.Chapter
		}
		endV := r.EndVerse
		if endV == 0 {
			endV = r.Verse
		}
		var text []string
		for v := r.Verse; v <= endV && r.Chapter == endCh; v++ {
			if s, ok := tr.Verse(r.Book, r.Chapter, v); ok {
				text = append(text, s)
			}
		}
		if len(text) > 0 {
			ui.Text(c, joinVerses(text)).FontSize(13).TextColor(a.pal.text).LineHeight(1.5)
		}
	}
}

// comfortCard is what the walk ends with: the scriptures that speak to a
// heavy heart, one of them the device reads aloud.
func (a *app) comfortCard(c *ui.Context) {
	ui.Column(c).Gap(8).Children(func() {
		ui.Text(c, T("Comfort")).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
		for i, r := range devotion.ComfortVerses {
			line := r.String()
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, line).FontSize(12.5)
				ui.Box(c).Grow(1)
				if ui.Button(c.Key("comfort-listen-"+strconv.Itoa(i)), T("Listen")).Clicked() {
					name, _ := bible.BookName(r.Book)
					a.listenRef(r, name.En+" "+line)
				}
			})
		}
	})
}

// prayerListCard is the prayers the reader keeps: added here, marked
// answered, dropped.
func (a *app) prayerListCard(c *ui.Context) {
	ui.Column(c).Gap(8).Children(func() {
		ui.Text(c, T("Prayer List")).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
		ui.Row(c).Gap(8).Children(func() {
			ui.TextInput(c.Key("prayer-input"), &a.prayerDraft).Placeholder(T("A prayer to keep")).Grow(1).FontSize(13)
			if ui.Button(c.Key("prayer-add"), T("Save|prayers")).Clicked() && a.prayerDraft != "" {
				if a.store != nil {
					a.store.AddPrayer(newPrayerID(), a.prayerDraft)
				}
				a.prayerDraft = ""
				a.redraw()
			}
		})
		if a.store == nil {
			return
		}
		prayers, err := a.store.Prayers()
		if err != nil {
			return
		}
		for _, p := range prayers {
			p := p
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				label := p.Text
				if p.Answered {
					label = "✓ " + label
				}
				ui.Text(c, label).FontSize(13)
				ui.Box(c).Grow(1)
				if !p.Answered {
					if ui.Button(c.Key("prayer-answer-"+p.ID), T("Answered")).Clicked() {
						a.store.SetPrayerAnswered(p.ID, true)
						a.redraw()
					}
				}
				if ui.Button(c.Key("prayer-drop-"+p.ID), T("Cancel")).Clicked() {
					a.store.DropPrayer(p.ID)
					a.redraw()
				}
			})
		}
	})
}

// newPrayerID is a prayer's own id.
func newPrayerID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// joinVerses joins verses for a page to read, a space between.
func joinVerses(vs []string) string {
	out := ""
	for i, v := range vs {
		if i > 0 {
			out += " "
		}
		out += v
	}
	return out
}
