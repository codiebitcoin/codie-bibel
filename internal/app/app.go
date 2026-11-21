// Package app is Codie Bible's state, actions and views: the Bible, its
// seasons and hymns, and a paired codie device reading it aloud.
package app

import (
	"context"
	"path/filepath"
	"runtime"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
	"codie-bibel/internal/device"
	"codie-bibel/internal/devotion"
	"codie-bibel/internal/hymn"
	"codie-bibel/internal/i18n"
	"codie-bibel/internal/player"
	"codie-bibel/internal/player/media"
	"codie-bibel/internal/queue"
	"codie-bibel/internal/store"
)

// pageKind is one of the app's pages.
type pageKind int

const (
	pageRead pageKind = iota
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
	tray     *mygo.Tray // the icon in the menu bar or notification area
	pal      palette
	page     pageKind
	settings settings
	bibles   *bible.Registry

	// The device and what it reads, and the hymns' player.
	dev      *device.Device
	bq       *queue.Queue
	pl       *player.Player
	hymnDir  string
	hymnFile map[string]string

	// The reader's own data, and the Read page's state.
	store        *store.Store
	rd           *readerState
	spurgeon     *devotion.Spurgeon
	spurgeonHalf string
	spurgeonOpen bool

	// The prayer walk: the step it stands at, the confession's field,
	// and the prayer list's field.
	prayStep    int
	prayText    string
	prayerDraft string

	// The story open on the Stories page, by its code.
	storyCode string

	navScroll  ui.ScrollState
	pageScroll ui.ScrollState
	scrolled   pageKind // the page pageScroll is of

	wizard wizardState // the pairing guide, until it's done

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

// syncPersona sends the chosen character back to the device at startup:
// another app may have changed it, and the app's choice is the one that
// stands. Nothing of it shows unless it fails.
func (a *app) syncPersona() {
	if a.dev == nil || !a.dev.Paired() {
		return
	}
	code := a.settings.CodiePersona
	if code == "" {
		code = device.DefaultPersona().Code
	}
	p := device.PersonaByCode(code)
	if p == nil {
		return
	}
	with := *p
	if v := a.settings.CodieVoice[code]; v != "" {
		with.VoiceID = v
	}
	vrm := a.seedVRM(&with)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := a.dev.ApplyPersona(ctx, with, vrm); err != nil {
			a.post(func() { a.toast = err.Error(); a.redraw() })
		}
	}()
}

// open shows a page.
func (a *app) open(p pageKind) {
	a.page = p
	a.redraw()
}

// applyLanguage speaks the settings' language, here and in the system:
// mygo lays its windows from the right for the languages that read so,
// as Arabic does.
func (a *app) applyLanguage() {
	// The system's languages come of mygo's main thread, which a test
	// doesn't serve: without one, the setting or English speaks.
	preferred := []string(nil)
	if a.postFn == nil {
		preferred = systemLanguages()
	}
	lang := i18n.Resolve(a.settings.Language, preferred)
	i18n.SetLanguage(lang)
	if a.win == nil {
		return // tests have no window, and mygo speaks to its main thread
	}
	tags := make([]string, 0, len(i18n.Languages))
	for _, l := range i18n.Languages {
		tags = append(tags, l.Code)
	}
	mygo.App.SetLanguages(tags...)
	mygo.App.SetLocale(lang)
}

// setLanguage keeps the language the reader chose, and speaks it at once.
func (a *app) setLanguage(code string) {
	a.settings.Language = code
	a.settings.save()
	a.applyLanguage()
	a.redraw()
}

// wireDevice makes the device of the saved pairing, and the queue of
// readings it walks.
func (a *app) wireDevice() {
	a.dev = device.New(a.settings.Codie)
	a.bq = queue.New(a.post, a.dev)
	a.bq.Changed = func() {
		a.queueProgress()
		a.redraw()
	}
	a.seedVoice()
}

// wireHymns writes the hymns' audio to the reader's data, once, and makes
// the player that reads it.
func (a *app) wireHymns() {
	a.hymnDir = userDataDir("hymns")
	a.pl = player.New(a.post)
	a.pl.Open = func(t media.Track) (*media.Stream, error) {
		path, ok := a.hymnFile[t.ID]
		if !ok {
			return nil, nil
		}
		ext := filepath.Ext(path)
		if len(ext) > 1 {
			ext = ext[1:]
		}
		return &media.Stream{URL: media.FileURL(path), Type: ext}, nil
	}
	if paths, err := hymn.Ensure(a.hymnDir); err == nil {
		a.hymnFile = paths
	}
	a.pl.SetVolume(a.settings.Volume)
}

// Run starts the app: it loads the settings, opens the window, and stays
// until it quits. Opening the app again shows the running one's window.
func Run(bibles *bible.Registry) error {
	if !mygo.App.RequestSingleInstanceLock() {
		return nil
	}
	a := &app{bibles: bibles}
	mygo.App.WhenReady(func() {
		a.settings = loadSettings()
		a.settings.applyAppearance()
		a.applyLanguage()
		a.wireDevice()
		a.wireHymns()
		if st, err := store.New(userDataDir("store")); err == nil {
			a.store = st
		}
		a.syncPersona()
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
		a.applyTray()
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
