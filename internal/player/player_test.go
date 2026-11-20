package player

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/audio"

	"codie-bibel/internal/player/media"
)

// silence is a long song of nothing.
type silence struct{}

func (silence) SampleRate() int               { return 48000 }
func (silence) Channels() int                 { return 2 }
func (silence) Len() int64                    { return 48000 * 600 }
func (silence) SetPosition(int64) error       { return nil }
func (silence) Read(b []float32) (int, error) { clear(b); return len(b), nil }

// posted makes a player of two silent hymns, playing, whose posted work a
// test runs itself.
func posted(t *testing.T) (chan func(), *Player) {
	t.Helper()
	ch := make(chan func(), 16)
	p := New(func(fn func()) { ch <- fn })
	p.now = time.Now
	path, err := filepath.Abs(filepath.Join("testdata", "seek.mp3"))
	if err != nil {
		t.Fatal(err)
	}
	p.Open = func(t media.Track) (*media.Stream, error) {
		return &media.Stream{URL: media.FileURL(path), Type: "mp3"}, nil
	}
	tracks := []media.Track{
		{ID: "1", Name: "one", Duration: 6000},
		{ID: "2", Name: "two", Duration: 6000},
		{ID: "3", Name: "three", Duration: 6000},
	}
	p.PlayQueue(tracks, 0)
	p.out = audio.NewPlayer(silence{})
	out := p.out
	out.SetVolume(0)
	out.Play()
	p.fade.attach(out, 1)
	p.playing = true
	p.loading = false
	return ch, p
}

// next runs what the player posted, waiting for it.
func next(t *testing.T, ch chan func()) {
	t.Helper()
	select {
	case fn := <-ch:
		fn()
	case <-time.After(2 * time.Second):
		t.Fatal("nothing posted")
	}
}

// Playing from an empty queue plays nothing.
func TestPlayEmptyQueue(t *testing.T) {
	p := &Player{at: -1}
	p.PlayQueue(nil, 0)
	if p.Current() != nil {
		t.Error("playing from an empty queue")
	}
}

// A track's end brings the next one up, and the queue's end with the
// repeat off stops.
func TestEndGoesToNext(t *testing.T) {
	ch, p := posted(t)
	p.repeat = RepeatAll
	p.ended()
	if p.track.Name != "two" {
		t.Errorf("the queue went to %q", p.track.Name)
	}
	p.at = len(p.order) - 1
	p.ended()
	if p.track.Name != "one" {
		t.Errorf("the queue went round to %q", p.track.Name)
	}
	// Off, the queue's end stops the player.
	p.repeat = RepeatOff
	p.at = len(p.order) - 1
	p.ended()
	if p.playing {
		t.Error("the queue played on past its end")
	}
	next(t, ch)
}

// The one track repeats at its end: the same track, playing.
func TestRepeatOne(t *testing.T) {
	_, p := posted(t)
	p.repeat = RepeatOne
	p.ended()
	if p.track.Name != "one" || !p.playing {
		t.Errorf("the one track didn't play again: %v, %v", p.track.Name, p.playing)
	}
}

// A seek before a track has loaded stays pending, and reads as the
// position until it plays.
func TestSeekPending(t *testing.T) {
	_, p := posted(t)
	p.out = nil
	p.Seek(3 * time.Second)
	if p.pending != 3*time.Second {
		t.Errorf("the seek didn't stay pending: %v", p.pending)
	}
	if p.Position() != 3*time.Second {
		t.Errorf("the position reads %v", p.Position())
	}
}

// Next and Prev walk the queue round.
func TestNextPrev(t *testing.T) {
	_, p := posted(t)
	p.Next()
	if p.track.Name != "two" {
		t.Errorf("Next went to %q", p.track.Name)
	}
	p.Prev()
	if p.track.Name != "one" {
		t.Errorf("Prev went to %q", p.track.Name)
	}
	p.at = 0
	p.Prev()
	if p.track.Name != "three" {
		t.Errorf("Prev round the queue went to %q", p.track.Name)
	}
}

// A file URL reads back as the path it came from, and a URL that is no
// file doesn't.
func TestPathOf(t *testing.T) {
	path, err := filepath.Abs(filepath.Join("a", "b.mp3"))
	if err != nil {
		t.Fatal(err)
	}
	u := media.FileURL(path)
	got, ok := media.PathOf(u)
	if !ok || got != path {
		t.Errorf("PathOf(%q) = %q, %v", u, got, ok)
	}
	if _, ok := media.PathOf("https://example.com/a.mp3"); ok {
		t.Error("a URL of the web read as a file")
	}
}

// The stream of a real MP3 opens, reads its length, and seeks clean.
func TestOpenSeekMP3(t *testing.T) {
	path := filepath.Join("testdata", "seek.mp3")
	f, err := os.Open(path)
	if err != nil {
		t.Skip("no testdata here")
	}
	defer f.Close()
	fi, _ := f.Stat()
	dl := opened(f, fi.Size())
	s, err := openStream(dl, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if s.Len() <= 0 {
		t.Error("the stream reads no length")
	}
	mid := s.Len() / 2
	if err := s.SetPosition(mid); err != nil {
		t.Fatal(err)
	}
	buf := make([]float32, 2048)
	if _, err := s.Read(buf); err != nil {
		t.Errorf("reading from the middle: %v", err)
	}
}
