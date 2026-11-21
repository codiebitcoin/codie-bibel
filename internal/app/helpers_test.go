package app

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
	"codie-bibel/internal/device"
	"codie-bibel/internal/store"
)

// mainThread takes what has to run on the main thread, which TestMain keeps
// serving: a test runs on another.
var mainThread = make(chan func())

// onMainThread runs fn on the main thread. Without a window nothing else serves
// it, and what MyGo does there, as changing the theme, would wait forever.
func onMainThread(fn func()) {
	done := make(chan struct{})
	mainThread <- func() {
		defer close(done)
		fn()
	}
	<-done
}

// TestMain keeps what the tests write out of the reader's own data, and serves
// the main thread while they run.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "codie-bibel-test")
	if err != nil {
		panic(err)
	}
	mygo.App.SetPath(mygo.PathUserData, dir)
	// The tests never clone a voice on a real device.
	seedVoiceClone = func(*device.Device, []byte, string) (string, error) {
		return "", errors.New("no cloning in tests")
	}
	result := make(chan int)
	go func() { result <- m.Run() }()
	code := 0
	for running := true; running; {
		select {
		case fn := <-mainThread:
			fn()
		case code = <-result:
			running = false
		}
	}
	os.RemoveAll(dir)
	os.Exit(code)
}

// newTestApp is an app for a test to drive: no window, an unpaired
// device, the queue over it, the shipped Bibles, and the reader's store
// in the test's data directory.
func newTestApp() *app {
	a := &app{settings: defaultSettings(), page: pageRead}
	// The posts run where they happen, so a goroutine's answer lands at
	// once; a test that wants them collected sets posts(a).
	a.postFn = func(fn func()) { fn() }
	a.wireDevice()
	a.wireHymns()
	if bibles, err := bible.NewRegistry(os.DirFS(filepath.Join("..", "bible", "data"))); err == nil {
		a.bibles = bibles
	}
	if st, err := store.New(userDataDir("store")); err == nil {
		a.store = st
	}
	return a
}

// posts collects what the app posts to the main thread, for a test to run with
// runPost once what it waits for is on its way.
func posts(a *app) chan func() {
	ch := make(chan func(), 64)
	a.postFn = func(fn func()) { ch <- fn }
	return ch
}

// runPost runs the next thing posted, waiting for it.
func runPost(t *testing.T, ch chan func()) {
	t.Helper()
	select {
	case fn := <-ch:
		fn()
	case <-time.After(2 * time.Second):
		t.Fatal("nothing was posted")
	}
}

// until runs what's posted until cond holds, for what posts more as it lands.
func until(t *testing.T, ch chan func(), cond func() bool) {
	t.Helper()
	for !cond() {
		runPost(t, ch)
	}
}

// settled shows the view of tt at rest, whatever is on its way in.
func settled(tt *ui.Tester) {
	tt.Frame()
	tt.Advance(time.Second)
	tt.Frame()
}
