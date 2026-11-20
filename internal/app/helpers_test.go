package app

import (
	"os"
	"testing"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
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

// newTestApp is an app for a test to drive: no window, no device.
func newTestApp() *app {
	return &app{settings: defaultSettings(), page: pageRead}
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
