package store

import (
	"testing"
)

// Highlights come and go by the verse, and a chapter's set reads back.
func TestHighlights(t *testing.T) {
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err := s.Highlight("kjv", "JHN", 3, 16, true); err != nil {
		t.Fatal(err)
	}
	if err := s.Highlight("kjv", "JHN", 3, 17, true); err != nil {
		t.Fatal(err)
	}
	got, err := s.Highlights("kjv", "JHN", 3)
	if err != nil {
		t.Fatal(err)
	}
	if !got[16] || !got[17] || got[18] {
		t.Errorf("the highlights read %v", got)
	}
	if err := s.Highlight("kjv", "JHN", 3, 16, false); err != nil {
		t.Fatal(err)
	}
	got, _ = s.Highlights("kjv", "JHN", 3)
	if got[16] {
		t.Error("the highlight stayed after it was taken")
	}
}

// A note is written, read, and dropped by writing an empty one.
func TestNotes(t *testing.T) {
	s, _ := New(t.TempDir())
	defer s.Close()
	if err := s.SetNote("kjv", "GEN", 1, 1, "In the beginning of my notes."); err != nil {
		t.Fatal(err)
	}
	text, ok, _ := s.Note("kjv", "GEN", 1, 1)
	if !ok || text != "In the beginning of my notes." {
		t.Errorf("the note reads %q, %v", text, ok)
	}
	if err := s.SetNote("kjv", "GEN", 1, 1, ""); err != nil {
		t.Fatal(err)
	}
	if _, ok, _ := s.Note("kjv", "GEN", 1, 1); ok {
		t.Error("the note stayed after an empty write")
	}
}

// Prayers keep their order, answer, and go when dropped.
func TestPrayers(t *testing.T) {
	s, _ := New(t.TempDir())
	defer s.Close()
	s.AddPrayer("p1", "for the storm to pass")
	s.AddPrayer("p2", "for a friend's health")
	ps, err := s.Prayers()
	if err != nil || len(ps) != 2 || ps[0].Text != "for a friend's health" {
		t.Fatalf("the prayers read %+v (%v)", ps, err)
	}
	s.SetPrayerAnswered("p1", true)
	ps, _ = s.Prayers()
	if !ps[1].Answered || ps[1].ID != "p1" {
		t.Errorf("the answered prayer reads %+v", ps[1])
	}
	s.DropPrayer("p1")
	ps, _ = s.Prayers()
	if len(ps) != 1 {
		t.Errorf("%d prayers after the drop", len(ps))
	}
}

// Progress keeps where a reading stands.
func TestProgress(t *testing.T) {
	s, _ := New(t.TempDir())
	defer s.Close()
	if s.Progress("story:flood") != 0 {
		t.Error("an unread reading stands somewhere")
	}
	s.SetProgress("story:flood", 7)
	if s.Progress("story:flood") != 7 {
		t.Error("the progress didn't stay")
	}
}

// Settings round-trip through the kv table.
func TestKV(t *testing.T) {
	s, _ := New(t.TempDir())
	defer s.Close()
	if _, ok, _ := s.Get("nope"); ok {
		t.Error("a key that wasn't set read as set")
	}
	s.Put("plan:bible-year", "42")
	v, ok, _ := s.Get("plan:bible-year")
	if !ok || v != "42" {
		t.Errorf("the kv reads %q, %v", v, ok)
	}
}
