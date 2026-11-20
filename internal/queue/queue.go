// Package queue is the device's readings, one after another: a story's
// turns, a psalm's, the prayer's scriptures, the verse of the day. It
// walks the turns of what is asked through the device's voice, and keeps
// where it stands, for a page to show and to go on from.
package queue

import (
	"context"
	"sync"

	"codie-bibel/internal/device"
)

// Item is one thing to read: where it came from, and the turns of it.
type Item struct {
	ID    string   `json:"id"`    // stable, as "story:flood", for the resume
	Title string   `json:"title"` // as the pages read it
	Lines []string `json:"lines"` // the turns, as speech planned them
	At    int      `json:"at"`    // the turn reading, or read next
}

// Done says whether the item has read whole.
func (it *Item) Done() bool { return it.At >= len(it.Lines) }

// Queue walks what is enqueued through the device. The queue's own state
// is read and changed on the main thread; the reading runs on its own,
// and comes back through post.
type Queue struct {
	post func(func())
	dev  device.API

	mu      sync.Mutex
	items   []*Item
	at      int // the item reading, -1 for none
	reading bool
	err     string
	cancel  context.CancelFunc

	// Changed runs whenever what the queue stands at changed.
	Changed func()
}

// New makes the queue, over the device given.
func New(post func(func()), dev device.API) *Queue {
	return &Queue{post: post, dev: dev, at: -1}
}

// Enqueue puts a reading at the queue's end, or in place of its own kind
// already there. Reading begins at once if the queue stood idle, from the
// item's own At, so a story goes on where it stopped.
func (q *Queue) Enqueue(id, title string, turns []string, at int) {
	q.mu.Lock()
	replaced := false
	for _, it := range q.items {
		if it.ID == id {
			it.Title, it.Lines = title, turns
			if at >= 0 && at <= len(turns) {
				it.At = at
			}
			replaced = true
			break
		}
	}
	if !replaced {
		if at < 0 {
			at = 0
		}
		q.items = append(q.items, &Item{ID: id, Title: title, Lines: turns, At: at})
	}
	start := !q.reading
	q.mu.Unlock()
	if start {
		q.kick()
	}
}

// Stop stops the reading and empties the queue.
func (q *Queue) Stop() {
	q.mu.Lock()
	q.items, q.at, q.reading = nil, -1, false
	cancel := q.cancel
	q.cancel = nil
	q.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	q.dev.Stop(context.Background())
	q.changed()
}

// Skip skips the turn reading for the next.
func (q *Queue) Skip() {
	q.mu.Lock()
	cancel := q.cancel
	q.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

// Pause stops the reading where it stands: what was reading keeps its At,
// and stays at the queue's head to go on from.
func (q *Queue) Pause() {
	q.mu.Lock()
	if q.at >= 0 && q.at < len(q.items) {
		it := q.items[q.at]
		q.items = []*Item{it}
		q.at = 0
	}
	q.reading = false
	cancel := q.cancel
	q.cancel = nil
	q.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	q.dev.Stop(context.Background())
	q.changed()
}

// Resume reads on from where the queue paused.
func (q *Queue) Resume() { q.kick() }

// Items is what the queue holds, to show and to keep.
func (q *Queue) Items() []*Item {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]*Item, len(q.items))
	copy(out, q.items)
	return out
}

// Restore puts the queue's readings back, as they were kept; nothing
// reads until something asks.
func (q *Queue) Restore(items []*Item) {
	q.mu.Lock()
	q.items = items
	q.reading = false
	q.at = -1
	q.mu.Unlock()
	q.changed()
}

// Reading says whether the device is reading now.
func (q *Queue) Reading() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.reading
}

// Err is why the reading stopped, if it failed.
func (q *Queue) Err() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.err
}

// changed runs what waits on the queue's movements, on the main thread.
func (q *Queue) changed() {
	if q.Changed != nil {
		q.post(q.Changed)
	}
}

// kick begins the reading, from the queue's head.
func (q *Queue) kick() {
	q.mu.Lock()
	if q.reading || len(q.items) == 0 {
		q.mu.Unlock()
		return
	}
	q.reading = true
	q.err = ""
	q.at = 0
	items := make([]*Item, len(q.items))
	copy(items, q.items)
	ctx, cancel := context.WithCancel(context.Background())
	q.cancel = cancel
	q.mu.Unlock()

	go func() {
		for _, it := range items {
			q.mu.Lock()
			start := it.At
			q.mu.Unlock()
			if !q.readItem(ctx, it, start) {
				return
			}
		}
		// The queue has read whole.
		q.post(func() {
			q.mu.Lock()
			q.items, q.at, q.reading, q.cancel = nil, -1, false, nil
			q.mu.Unlock()
			q.changed()
		})
	}()
	q.changed()
}

// readItem reads an item's turns from start through the device, the
// item's At moving on the main thread as each turn reads. It reads false
// when the reading was stopped or failed, the At standing where it
// stopped.
func (q *Queue) readItem(ctx context.Context, it *Item, start int) bool {
	for i := start; i < len(it.Lines); i++ {
		turn := it.Lines[i]
		err := q.dev.Speak(ctx, turn)
		if ctx.Err() != nil {
			// Stopped or skipped: this turn reads again from here.
			q.post(func() { q.mark(it, i) })
			return false
		}
		if err != nil {
			q.post(func() {
				q.mu.Lock()
				q.reading, q.err, q.at, q.cancel = false, err.Error(), -1, nil
				q.mu.Unlock()
				q.changed()
			})
			return false
		}
		q.post(func() { q.mark(it, i+1) })
	}
	return true
}

// mark moves an item's place in it, on the main thread.
func (q *Queue) mark(it *Item, at int) {
	q.mu.Lock()
	it.At = at
	q.mu.Unlock()
	q.changed()
}
