package app

import (
	"strconv"
	"time"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
	"codie-bibel/internal/devotion"
	"codie-bibel/internal/speech"
)

// The Devotions page: the verse of the day, Spurgeon's Morning and
// Evening for the hour, and the reading plans' today.

func (a *app) devotionsPage(c *ui.Context) {
	a.pageHeading(c, T("Devotions"), nil)
	top, bottom := a.insets(c)
	now := time.Now()
	ui.Scroll(c.Key("devotions")).Grow(1).Children(func() {
		ui.Column(c).Padding(top, 24, bottom, 24).Gap(20).Children(func() {
			a.votdCard(c, now)
			a.spurgeonCard(c, now)
			a.plansCard(c, now)
		})
	})
}

// votdCard is the verse of the day: its reference and its words in the
// translation at hand, listen beside it.
func (a *app) votdCard(c *ui.Context, now time.Time) {
	code, ch, vs := devotion.VerseOfTheDay(now)
	tr := a.ensureTranslation()
	ref := bible.Ref{Book: code, Chapter: ch, Verse: vs, EndVerse: vs}
	name, _ := bible.BookName(code)
	a.card(c, T("Verse of the Day"), ref.String(), func() {
		if tr != nil {
			if text, ok := tr.Verse(code, ch, vs); ok {
				ui.Text(c, text).FontSize(15).TextColor(a.pal.text)
			}
		}
	}, func() {
		a.listenRef(ref, T("Verse of the Day")+": "+name.En+" "+strconv.Itoa(ch)+":"+strconv.Itoa(vs))
	})
}

// spurgeonCard is the day's reading of Morning and Evening, the hour's
// half first, the other a click away.
func (a *app) spurgeonCard(c *ui.Context, now time.Time) {
	if a.spurgeon == nil {
		if s, err := devotion.NewSpurgeon(); err == nil {
			a.spurgeon = s
		}
	}
	if a.spurgeon == nil {
		return
	}
	me := a.spurgeon.Read(now)
	half := devotion.Half(now)
	title := T("Morning & Evening") + " — " + i18nHalf(half)
	read := me.AM
	if half == "pm" {
		read = me.PM
	}
	a.card(c, title, read.Verse, func() {
		ui.Text(c, read.Text).FontSize(13.5).TextColor(a.pal.text).LineHeight(1.6)
	}, func() {
		// The whole reading, its text planned into turns.
		plan := speech.NewPlan(read.Verse+". "+read.Text, 0)
		a.bq.Enqueue("spurgeon:"+strconv.Itoa(int(now.Month()))+"-"+strconv.Itoa(now.Day())+"-"+half, title, plan.Turns, 0)
	})
}

// i18nHalf is the half of the day, as the language reads it.
func i18nHalf(half string) string {
	if half == "pm" {
		return T("Evening")
	}
	return T("Morning")
}

// plansCard is the reading plans, each with today's passages.
func (a *app) plansCard(c *ui.Context, now time.Time) {
	ui.Column(c).Gap(8).Children(func() {
		ui.Text(c, T("Reading Plans")).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
		for _, p := range devotion.Plans() {
			ui.Column(c).Padding(12, 14).Gap(6).Radius(12).Background(a.pal.cardShade).Children(func() {
				ui.Text(c, p.Name).FontSize(14).FontWeight(600)
				ui.Column(c).Gap(4).Children(func() {
					for _, r := range p.ForDay(now) {
						name, _ := bible.BookName(r.Book)
						line := name.En + " " + strconv.Itoa(r.Chapter)
						if ui.Button(c.Key("plan-"+p.ID+"-"+line), line).Clicked() {
							a.openChapter(r.Book, r.Chapter)
						}
					}
				})
			})
		}
	})
}

// card is one devotion card: its title, its line under it, its body, and
// what Listen reads.
func (a *app) card(c *ui.Context, title, line string, body func(), listen func()) {
	ui.Column(c).Padding(16, 18).Gap(8).Radius(12).Background(a.pal.cardShade).Children(func() {
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, title).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
			ui.Box(c).Grow(1)
			if ui.Button(c.Key("listen-"+title), T("Listen")).Clicked() {
				listen()
			}
		})
		ui.Text(c, line).FontSize(13).FontWeight(600).TextColor(a.pal.accent)
		body()
	})
}
