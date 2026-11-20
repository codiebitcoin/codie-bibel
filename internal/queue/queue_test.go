package queue

import (
	"context"
	"errors"
	"testing"
	"time"

	"codie-bibel/internal/device"
)

// fakeDevice is a device that records what it was asked to read, and
// speaks nothing: each Speak waits for its own length, so a stop lands
// mid-reading. failAt is the turn that fails, or -1 for none.
type fakeDevice struct {
	spoken []string
	stops  int
	each   time.Duration
	failAt int
}

func (f *fakeDevice) Speak(ctx context.Context, text string) error {
	if f.failAt == len(f.spoken) {
		return errors.New("the device went quiet")
	}
	f.spoken = append(f.spoken, text)
	if f.each > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(f.each):
		}
	}
	return nil
}

func (f *fakeDevice) Chat(ctx context.Context, text string) error {
	f.spoken = append(f.spoken, "chat:"+text)
	return nil
}

func (f *fakeDevice) Stop(ctx context.Context) error { f.stops++; return nil }

func (f *fakeDevice) Status(ctx context.Context) error { return nil }

// testQueue is a queue whose posted work runs where it happens, so the
// tests read the item's place as the turns read. The app runs it on the
// main thread.
func testQueue(dev device.API) *Queue {
	return New(func(fn func()) { fn() }, dev)
}

// wait until cond holds, or the test fails of waiting too long.
func wait(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("the queue never came to rest")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// A turn after turn reads through the device, and the queue empties when
// all have read.
func TestQueueReadsThrough(t *testing.T) {
	dev := &fakeDevice{failAt: -1}
	q := testQueue(dev)
	q.Enqueue("story:flood", "The Flood", []string{"turn one.", "turn two.", "turn three."}, -1)
	wait(t, func() bool { return len(dev.spoken) == 3 && !q.Reading() })
	if dev.spoken[0] != "turn one." || dev.spoken[2] != "turn three." {
		t.Errorf("the turns read out of order: %v", dev.spoken)
	}
	if items := q.Items(); len(items) != 0 {
		t.Errorf("the queue kept %d items after reading whole", len(items))
	}
}

// A pause leaves the item where it stood, and the device asked to hush.
func TestQueuePauseKeepsPlace(t *testing.T) {
	dev := &fakeDevice{failAt: -1, each: 50 * time.Millisecond}
	q := testQueue(dev)
	q.Enqueue("devotions:me", "Morning", []string{"one.", "two.", "three.", "four."}, -1)
	time.Sleep(130 * time.Millisecond) // two turns and into the third
	q.Pause()
	if dev.stops == 0 {
		t.Error("the device wasn't asked to stop")
	}
	items := q.Items()
	if len(items) != 1 || items[0].ID != "devotions:me" || items[0].At < 1 {
		t.Fatalf("the item stands at %+v", items)
	}
	if q.Reading() {
		t.Error("the queue reads after its pause")
	}
}

// A stop empties the queue, and hushes the device.
func TestQueueStopEmpties(t *testing.T) {
	dev := &fakeDevice{failAt: -1, each: 50 * time.Millisecond}
	q := testQueue(dev)
	q.Enqueue("devotions:me", "Morning", []string{"one.", "two.", "three."}, -1)
	time.Sleep(60 * time.Millisecond)
	q.Stop()
	if len(q.Items()) != 0 {
		t.Errorf("the queue kept %d items", len(q.Items()))
	}
}

// A paused reading goes on where it stopped, reading the turn it stood in
// again.
func TestQueueResumesWhereItPaused(t *testing.T) {
	dev := &fakeDevice{failAt: -1, each: 40 * time.Millisecond}
	q := testQueue(dev)
	q.Enqueue("devotions:me", "Morning", []string{"one.", "two.", "three."}, -1)
	time.Sleep(110 * time.Millisecond)
	q.Pause()
	already := len(dev.spoken)
	if q.Reading() {
		t.Error("the queue reads after its pause")
	}
	q.Enqueue("devotions:me", "Morning", []string{"one.", "two.", "three."}, -1)
	wait(t, func() bool { return !q.Reading() })
	// The turn it stood in reads again: the count grows past what had
	// read, and the whole of it ends read.
	if len(dev.spoken) <= already {
		t.Errorf("resuming read nothing more of %d", already)
	}
	if len(q.Items()) != 0 {
		t.Errorf("the queue kept %d items after reading whole", len(q.Items()))
	}
}

// A device that fails stops the reading, and says why.
func TestQueueFailureStops(t *testing.T) {
	dev := &fakeDevice{failAt: 1}
	q := testQueue(dev)
	q.Enqueue("psalm", "Psalm 23", []string{"one.", "two."}, -1)
	wait(t, func() bool { return !q.Reading() })
	if q.Err() == "" {
		t.Error("the failure reads as nothing")
	}
}

// The device API the queue walks is the one the device answers.
var _ device.API = (*fakeDevice)(nil)
