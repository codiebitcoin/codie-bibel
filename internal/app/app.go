// Package app is Codie Bible's state, actions and views: the Bible, its
// seasons and hymns, and a paired codie device reading it aloud.
package app

import (
	"path/filepath"
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/i18n"
)

// pageKind is one of the app's pages.
type pageKind int

const (
	pageRead     pageKind = iota
	pageDevotions
	pagePray
	pageStories
	pageSeasons
	pageMusic
	pageDevice
	pageSettings
)

// app is the state the window shows. Everything in it is read and changed
// on the main thread: by the views, and by goroutines through post.
type app struct {
	win      *mygo.Window
	pal      palette
	page     pageKind
	settings settings

	navScroll  ui.ScrollState
	pageScroll ui.ScrollState
	scrolled   pageKind // the page pageScroll is of

	toast string

	// postFn, if set, takes what post would run on the main thread, as in
	// tests, which have no window.
	postFn func(func())
}

// post runs fn on the main thread and redraws the window. Without a
// window, as in tests, it drops fn.
func (a *app) post(fn func()) {
	if a.postFn != nil {
		a.postFn(fn)
		return
	}
	if a.win == nil {
		return
	}
	a.win.Update(func() {
		fn()
		a.redraw()
	})
}

func (a *app) redraw() {
	if a.win != nil {
		a.win.Invalidate()
	}
}

func (a *app) showMain() {
	if a.win != nil {
		if a.win.IsMinimized() {
			a.win.Restore()
		}
		a.win.Show()
		a.win.Focus()
		mygo.App.Focus()
	}
}

// userDataDir is where the app keeps what's the reader's.
func userDataDir(sub string) string {
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		dir = "."
	}
	return filepath.Join(dir, sub)
}

// open shows a page.
func (a *app) open(p pageKind) {
	a.page = p
	a.redraw()
}

// Run starts the app: it loads the settings, opens the window, and stays
// until it quits. Opening the app again shows the running one's window.
func Run() error {
	if !mygo.App.RequestSingleInstanceLock() {
		return nil
	}
	a := &app{}
	mygo.App.WhenReady(func() {
		a.settings = loadSettings()
		a.settings.applyAppearance()
		i18n.SetLanguage(i18n.Resolve(a.settings.Language, systemLanguages()))
		a.win = mygo.NewWindow(mygo.WindowOptions{
			Title:                mygo.App.Name(),
			Width:                1120,
			Height:               760,
			MinWidth:             860,
			MinHeight:            560,
			StateKey:             "main",
			TitleBarStyle:        mygo.TitleBarHidden,
			TitleBarHeight:       buttonsHeight,
			TrafficLightPosition: &mygo.Point{X: 20, Y: 20},
			Vibrancy:             mygo.VibrancySidebar,
			BackgroundColor:      "light-dark(#ffffff, #141416)",
			Content:              ui.View(a.view),
		})
		// On the Mac, closing the window hides it, and the Dock's brings
		// it back. Elsewhere there is no tray yet, so closing quits.
		a.win.OnClose(func(e *mygo.CloseEvent) {
			if runtime.GOOS == "darwin" {
				e.PreventDefault()
				a.win.Hide()
			}
		})
	})
	mygo.App.OnWindowAllClosed(func() {})
	mygo.App.OnActivate(func(bool) { a.showMain() })
	mygo.App.OnSecondInstance(func([]string, string) { a.showMain() })
	return mygo.App.Run()
}
