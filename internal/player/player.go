// Package player plays the hymns: a queue of tracks read from files, one
// at a time, faded in and out, with the repeat and the shuffle a queue of
// music wants. Its streams come from the files the app writes of its
// embedded audio; there is nothing to fetch, and nothing to cache.
package player

import (
	"errors"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/audio"

	"codie-bibel/internal/player/media"
)

// RepeatMode is what happens when a track ends: the queue goes on, or the
// one track plays again, or the queue stops.
type RepeatMode int

const (
	RepeatOff RepeatMode = iota
	RepeatAll
	RepeatOne
)

// Which way a load goes past a track that can't play: on to the next, back
// to the one before, or nowhere, staying on it to say why.
const (
	forward = 1
	back    = -1
	stay    = 0
)

// Player is the hymns' player. Everything in it is read and changed on
// the main thread; what plays runs on the audio system's own.
type Player struct {
	post func(func())
	// Open resolves a track to its audio, as a file:// URL of the file
	// the app wrote of its embedded audio.
	Open func(t media.Track) (*media.Stream, error)
	// Changed runs when what plays changed, for the session to keep.
	Changed func()

	queue    []media.Track
	order    []int
	at       int
	repeat   RepeatMode
	shuffle  bool
	volume   float64
	track    *media.Track
	loading  bool
	playing  bool
	err      string
	out      *audio.Player
	stream   stream
	watch    *watched
	info     *media.Stream
	fade     fader
	pending  time.Duration
	seeking  bool
	seekTo   time.Duration
	seekGen  int
	gen      int
	restored bool

	// now says what time it is; tests set their own clock.
	now func() time.Time
}

// New makes the player, and takes the system's media keys for its own.
func New(post func(func())) *Player {
	p := &Player{post: post, volume: 0.8, now: time.Now}
	mygo.NowPlaying.OnCommand(func(cmd mygo.MediaCommand, pos time.Duration) {
		switch cmd {
		case mygo.MediaPlay:
			if !p.Playing() {
				p.TogglePlay()
			}
		case mygo.MediaPause:
			if p.Playing() {
				p.TogglePlay()
			}
		case mygo.MediaNext:
			p.Next()
		case mygo.MediaPrevious:
			p.Prev()
		case mygo.MediaSeek:
			p.Seek(pos)
		}
	})
	return p
}

// announce shows the system what plays.
func (p *Player) announce(t *media.Track) {
	state := mygo.PlaybackPaused
	if p.playing {
		state = mygo.PlaybackPlaying
	}
	mygo.NowPlaying.Set(mygo.NowPlayingInfo{
		Title:    t.Name,
		Artist:   t.ArtistNames(),
		Album:    t.Album.Name,
		Duration: t.Length(),
		Position: p.Position(),
		State:    state,
		Commands: mygo.MediaPlay | mygo.MediaPause | mygo.MediaNext | mygo.MediaPrevious | mygo.MediaSeek,
	})
}

// sync tells the system whether the player plays.
func (p *Player) sync() {
	if p.track == nil {
		mygo.NowPlaying.Clear()
		return
	}
	state := mygo.PlaybackPaused
	if p.playing {
		state = mygo.PlaybackPlaying
	}
	mygo.NowPlaying.SetPlayback(state, p.Position())
	if p.Changed != nil {
		p.Changed()
	}
}

// ended is what a track's own end does: the one again, or the queue's
// next, or nothing.
func (p *Player) ended() {
	if p.repeat == RepeatOne {
		p.out.Seek(0)
		p.out.Play()
		return
	}
	if p.at+1 >= len(p.order) && p.repeat == RepeatOff {
		p.playing = false
		p.sync()
		return
	}
	p.Next()
}

// PlayQueue replaces the queue with tracks and plays the one at start, or
// the first after it that can play.
func (p *Player) PlayQueue(tracks []media.Track, start int) {
	if len(tracks) == 0 {
		return
	}
	start = min(start, len(tracks)-1)
	p.queue = append([]media.Track(nil), tracks...)
	p.order = make([]int, len(p.queue))
	for i := range p.order {
		p.order[i] = i
	}
	if p.shuffle {
		p.shuffleOrder(start)
		start = 0
	}
	p.loadSkipping(max(start, 0), 0, forward, true, 0)
}

// PlayShuffled replaces the queue with the tracks in a random order and
// plays one of them, shuffle on.
func (p *Player) PlayShuffled(tracks []media.Track) {
	p.shuffle = true
	p.PlayQueue(tracks, -1)
}

// shuffleOrder puts the order in a random one with first at its head, or
// with any at its head when first is below 0.
func (p *Player) shuffleOrder(first int) {
	rand.Shuffle(len(p.order), func(i, j int) { p.order[i], p.order[j] = p.order[j], p.order[i] })
	for i, v := range p.order {
		if v == first {
			p.order[0], p.order[i] = p.order[i], p.order[0]
		}
	}
}

// ToggleShuffle turns the shuffle on and off.
func (p *Player) ToggleShuffle() {
	p.shuffle = !p.shuffle
	if len(p.order) == 0 {
		return
	}
	cur := p.order[max(p.at, 0)]
	if p.shuffle {
		p.shuffleOrder(cur)
		p.at = 0
	} else {
		for i := range p.order {
			p.order[i] = i
		}
		p.at = cur
	}
}

// CycleRepeat turns the repeat: off, to the whole queue, to the one
// track, and back to off.
func (p *Player) CycleRepeat() { p.repeat = (p.repeat + 1) % 3 }

// Current returns the track playing or loading, or nil.
func (p *Player) Current() *media.Track { return p.track }

// Queue returns the tracks queued, in the order they were given.
func (p *Player) Queue() []media.Track { return p.queue }

// Playing reports what the listener wants: false while paused.
func (p *Player) Playing() bool { return p.playing }

// Loading reports whether the current track is still on its way: loading,
// or waiting for a seek past what has come in.
func (p *Player) Loading() bool {
	return p.loading || p.watch != nil && p.watch.starved()
}

// Ready reports whether the current track has loaded, to play.
func (p *Player) Ready() bool { return p.out != nil }

// Err returns why the current track did not play, or "".
func (p *Player) Err() string { return p.err }

// Shuffle says whether the shuffle is on.
func (p *Player) Shuffle() bool      { return p.shuffle }
func (p *Player) Repeat() RepeatMode { return p.repeat }

// PlayOrder returns the queue's indices in play order, and where the
// current track is among them, -1 for none.
func (p *Player) PlayOrder() ([]int, int) { return p.order, p.at }

// Next plays the queue's next, round to its head.
func (p *Player) Next() {
	if len(p.order) == 0 {
		return
	}
	at := p.at + 1
	if at >= len(p.order) {
		at = 0
	}
	p.load(at)
}

// Prev plays the track before this one, however far this one has played:
// once a few seconds in, going back to its start only kept the listener
// pressing twice.
func (p *Player) Prev() {
	if len(p.order) == 0 {
		return
	}
	at := p.at - 1
	if at < 0 {
		at = len(p.order) - 1
	}
	p.loadSkipping(at, 0, back, true, 0)
}

// Jump plays the queue's track i, the listener's pick: one that can't play
// stays, saying why.
func (p *Player) Jump(i int) {
	for k, v := range p.order {
		if v == i {
			p.loadSkipping(k, 0, stay, true, 0)
			return
		}
	}
}

// TogglePlay pauses, fading out, or plays, fading in.
func (p *Player) TogglePlay() {
	if p.out == nil {
		switch {
		case p.loading:
			// It plays once loaded, or not.
			p.playing = !p.playing
			p.sync()
		case len(p.order) > 0:
			// Again, as for a track that couldn't play.
			p.loadSkipping(max(p.at, 0), 0, stay, true, p.pending)
		}
		return
	}
	out := p.out
	if p.playing {
		p.playing = false
		p.fade.fade(0, fadeOut, func() {
			p.post(func() {
				if p.out == out && !p.playing {
					out.Pause()
					p.sync()
				}
			})
		})
		p.sync()
		return
	}
	p.playing = true
	if !p.seeking && !out.Playing() {
		if out.Position() >= out.Duration()-time.Second {
			out.Seek(0)
		}
		out.Play()
	}
	p.fade.fade(1, fadeIn, nil)
	p.sync()
}

// seekGap is how long a seek waits for what plays to fade out, which the
// audio player does over 5 ms, as it pauses.
const seekGap = 30 * time.Millisecond

// Seek moves to t. Playing, it pauses and plays again around the jump,
// which the audio player fades over a few milliseconds, rather than
// cutting from one sound to the other with a click.
func (p *Player) Seek(t time.Duration) {
	if p.out == nil {
		p.pending = t
		return
	}
	out := p.out
	// Paused, it jumps at once; but a seek on its way paused it only for
	// the jump, and this one takes its place, playing on from there.
	if !out.Playing() && !p.seeking {
		out.Seek(t)
		p.sync()
		return
	}
	p.seekGen++
	gen := p.seekGen
	p.seeking, p.seekTo = true, t
	out.Pause()
	time.AfterFunc(seekGap, func() {
		p.post(func() {
			if gen != p.seekGen || p.out != out {
				return
			}
			p.seeking = false
			out.Seek(t)
			if p.playing || !p.fade.silent() {
				out.Play()
			}
			p.sync()
		})
	})
	p.sync()
}

// SetVolume sets how loud the hymns play, from 0 to 1.
func (p *Player) SetVolume(v float64) {
	p.volume = v
	// The ear hears loudness in powers.
	p.fade.setGain(v * v)
}

func (p *Player) Position() time.Duration {
	switch {
	case p.out == nil:
		return p.pending
	case p.seeking:
		return p.seekTo
	}
	return p.out.Position()
}

func (p *Player) Duration() time.Duration {
	switch {
	case p.track == nil:
		return 0
	case p.track.Duration > 0 || p.out == nil:
		return p.track.Length()
	}
	return p.out.Duration() // a file that didn't say
}

// stop pauses what plays, and lets it go.
func (p *Player) stop() {
	if out := p.out; out != nil {
		// Paused first, it fades out over a few milliseconds, rather than
		// stopping with a click.
		out.Pause()
		time.AfterFunc(100*time.Millisecond, out.Close)
	}
	p.out, p.stream, p.watch, p.info = nil, nil, nil, nil
	p.seeking = false
	p.fade.attach(nil, 1)
}

// load plays the track at order[at], skipping those that cannot play.
func (p *Player) load(at int) {
	p.loadSkipping(at, 0, forward, true, 0)
}

// loadSkipping loads the track at order[at], to play or not, from from,
// skipping those that cannot play the way dir goes: skipped of them so
// far.
func (p *Player) loadSkipping(at, skipped, dir int, play bool, from time.Duration) {
	t := p.queue[p.order[at]]
	p.stop()
	p.gen++
	gen := p.gen
	p.at = at
	p.track = &t
	p.loading, p.playing, p.err = true, play, ""
	p.pending = from
	p.announce(&t)

	go func() {
		info, s, err := p.open(t)
		p.post(func() {
			if gen != p.gen {
				if s != nil {
					s.Close()
				}
				return
			}
			p.loading = false
			if err != nil {
				p.loadFailed(err, at, skipped, dir)
				return
			}
			p.stream, p.info = s, info
			p.watch = &watched{stream: s, now: p.now, ended: func() { p.post(func() {}) }}
			p.out = audio.NewPlayer(p.watch)
			// From its middle, as a session picked up again, it fades in.
			level := 1.0
			if p.pending > 0 {
				level = 0
			}
			p.fade.attach(p.out, level)
			p.fade.setGain(p.volume * p.volume)
			p.out.OnEnd(func() {
				p.post(func() {
					if gen == p.gen {
						p.ended()
					}
				})
			})
			if p.pending > 0 {
				p.out.Seek(p.pending)
				p.pending = 0
			}
			if p.playing {
				p.out.Play()
				p.fade.fade(1, fadeIn, nil)
			}
			p.sync()
		})
	}()
}

// loadFailed says why the track at order[at] can't play, and skips it the
// way dir goes, unless it stays, or nothing in the queue can play.
func (p *Player) loadFailed(err error, at, skipped, dir int) {
	p.err = errorText(err)
	if dir != stay && skipped+1 < len(p.order) {
		next := (at + dir + len(p.order)) % len(p.order)
		p.loadSkipping(next, skipped+1, dir, p.playing, 0)
		return
	}
	p.playing = false
	p.sync()
}

// errorText says what went wrong with a track, in the reader's terms.
func errorText(err error) string {
	switch {
	case errors.Is(err, media.ErrUnavailable):
		return "This hymn isn't available"
	case errors.Is(err, errNoSource):
		return "Nothing plays this hymn"
	}
	return err.Error()
}

var errNoSource = errors.New("player: no such source")

// open opens track t's audio: the file the app wrote of its embedded
// audio, read whole.
func (p *Player) open(t media.Track) (*media.Stream, stream, error) {
	if p.Open == nil {
		return nil, nil, errNoSource
	}
	info, err := p.Open(t)
	if err != nil {
		return nil, nil, err
	}
	if info == nil {
		return nil, nil, media.ErrUnavailable
	}
	path, isFile := media.PathOf(info.URL)
	if !isFile {
		return nil, nil, media.ErrUnavailable
	}
	dl, err := openFile(path)
	if err != nil {
		return nil, nil, err
	}
	s, err := openStream(dl, t.Duration)
	if err != nil {
		dl.Close()
		return nil, nil, err
	}
	return info, s, nil
}

// openFile opens a song on disk as a download that has all of it in.
func openFile(path string) (*download, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	return opened(f, info.Size()), nil
}

// Session is where playback was left, to pick up again.
type Session struct {
	Queue    []media.Track `json:"queue"`
	Order    []int         `json:"order"`
	At       int           `json:"at"`
	Shuffle  bool          `json:"shuffle"`
	Repeat   RepeatMode    `json:"repeat"`
	Position int64         `json:"position"` // milliseconds
}

// Session returns where playback is, or false with nothing queued.
func (p *Player) Session() (Session, bool) {
	if p.track == nil || p.at < 0 {
		return Session{}, false
	}
	return Session{
		Queue:    p.queue,
		Order:    p.order,
		At:       p.at,
		Shuffle:  p.shuffle,
		Repeat:   p.repeat,
		Position: p.Position().Milliseconds(),
	}, true
}

// Restore picks up a session where it was left, paused, and loads its
// track to play at once when asked.
func (p *Player) Restore(s Session) {
	n := len(s.Queue)
	if n == 0 || len(s.Order) != n || s.At < 0 || s.At >= n {
		return
	}
	seen := make([]bool, n)
	for _, i := range s.Order {
		if i < 0 || i >= n || seen[i] {
			return
		}
		seen[i] = true
	}
	p.queue, p.order = s.Queue, s.Order
	p.shuffle, p.repeat = s.Shuffle, s.Repeat%3
	p.restored = true
	p.loadSkipping(s.At, 0, stay, false, time.Duration(s.Position)*time.Millisecond)
}

// trackLabel is how a track reads in a file name: title, or the artists
// and it, deduplicated.
func trackLabel(t *media.Track) string {
	if t.Name == "" {
		return t.ID
	}
	if artist := t.ArtistNames(); artist != "" && !strings.Contains(t.Name, artist) {
		return artist + " - " + t.Name
	}
	return t.Name
}
