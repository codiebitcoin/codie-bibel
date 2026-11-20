package app

import (
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

const sidebarWidth = 232

func (a *app) view(c *ui.Context) {
	t := *c.Theme()
	a.pal = paletteFor(&t)
	p := a.pal
	useAccent(&t, p)
	t.Selection = p.accent.Alpha(0.3)
	c.SetTheme(&t)

	if c.Vibrancy() {
		c.Root().Background(ui.Transparent)
	}
	root := ui.Box(c).Fill()
	root.Children(func() {
		page := ui.Row(c).Fill().AlignItems(ui.Stretch)
		page.Children(func() {
			side := ui.Column(c).Width(sidebarWidth)
			if !c.Vibrancy() {
				side.Background(p.sidebar)
			}
			side.Children(func() { a.sidebar(c) })
			ui.Box(c.Key("main")).Grow(1).Background(p.bg).Children(func() {
				a.content(c)
			})
		})
	})
	if a.toast != "" {
		c.Toast(a.toast)
		a.toast = ""
	}
}

// navItem is one entry of the sidebar.
type navItem struct {
	label string
	page  pageKind
}

// sidebarItems are the app's pages, in the order the reader walks them.
func sidebarItems() []navItem {
	return []navItem{
		{T("Read"), pageRead},
		{T("Devotions"), pageDevotions},
		{T("Pray"), pagePray},
		{T("Stories"), pageStories},
		{T("Seasons"), pageSeasons},
		{T("Music"), pageMusic},
		{T("Device"), pageDevice},
	}
}

func (a *app) sidebar(c *ui.Context) {
	// Under the window buttons where they sit over the sidebar, as the
	// traffic lights do; elsewhere, as close to the top as to the sides.
	top := float32(14)
	if c.TitleBar().Left > 0 {
		top = topBar(c) + 4
	}
	// Where nothing in it takes the press, the sidebar moves the window, as
	// the title bar it reaches up into would.
	ui.Column(c).Fill().Padding(top, 14, 14, 14).Gap(4).DragWindow().Children(func() {
		ui.Row(c).Height(28).PaddingX(7).Gap(6).AlignItems(ui.Center).Children(func() {
			ui.Text(c, mygo.App.Name()).FontSize(12.5).FontWeight(600).TextColor(a.pal.muted)
		})

		ui.Scroll(c.Key("nav")).TrackScroll(&a.navScroll).Grow(1).Margin(0, -6).Padding(0, 6).Gap(4).Children(func() {
			for _, n := range sidebarItems() {
				a.navRow(c, n.label, n.page)
			}
		})

		a.navRow(c, T("Settings"), pageSettings)
	})
}

// navRow is one row of the sidebar, which shows its page.
func (a *app) navRow(c *ui.Context, label string, page pageKind) {
	p := a.pal
	selected := a.page == page
	b := ui.ButtonBase(c.Key(label)).Height(34).PaddingX(10).Gap(10).Radius(8).AlignItems(ui.Center).Justify(ui.Start).FocusRing(false).Label(label)
	bg := ui.Transparent
	switch {
	case selected:
		bg = p.press
	case b.Hovered():
		bg = p.hover
	}
	b.Background(b.AnimateColor("bg", bg, ui.Animation{Duration: 120 * time.Millisecond}))
	color := p.text
	if selected {
		color = p.accent
	}
	b.Children(func() {
		ui.Text(c, label).FontSize(13.5).FontWeight(500).TextColor(color)
	})
	if b.Clicked() {
		a.open(page)
	}
}

func (a *app) content(c *ui.Context) {
	// A new page starts at its top.
	if a.scrolled != a.page {
		a.pageScroll, a.scrolled = ui.ScrollState{}, a.page
	}
	switch a.page {
	case pageRead:
		a.readPage(c)
	case pageDevotions:
		a.devotionsPage(c)
	case pagePray:
		a.prayPage(c)
	case pageStories:
		a.storiesPage(c)
	case pageSeasons:
		a.seasonsPage(c)
	case pageMusic:
		a.musicPage(c)
	case pageDevice:
		a.devicePage(c)
	case pageSettings:
		a.settingsPage(c)
	}
}

// coming is what a page says until its content arrives: one line of what it
// will hold.
func (a *app) coming(c *ui.Context, what string) {
	top, bottom := a.insets(c)
	ui.Column(c).Grow(1).Padding(top, 24, bottom, 24).Gap(12).Children(func() {
		ui.Text(c, what).TextColor(a.pal.faint).FontSize(13)
	})
}
