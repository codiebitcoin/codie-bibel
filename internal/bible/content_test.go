package bible

import "testing"

// The words of Christ read red where the speeches run, and the rest of the
// Gospel reads black.
func TestInRed(t *testing.T) {
	red := []struct {
		code   string
		ch, vs int
	}{
		{"MAT", 5, 3},   // the Sermon on the Mount
		{"JHN", 3, 16},  // to Nicodemus
		{"JHN", 17, 5},  // the high priestly prayer
		{"LUK", 15, 21}, // the father's welcome
		{"MRK", 4, 39},  // "Peace, be still"
		{"ACT", 9, 5},   // the voice from heaven
	}
	for _, c := range red {
		if !InRed(c.code, c.ch, c.vs) {
			t.Errorf("%s %d:%d reads black", c.code, c.ch, c.vs)
		}
	}
	black := []struct {
		code   string
		ch, vs int
	}{
		{"JHN", 1, 1},   // the Evangelist's own preface
		{"GEN", 1, 1},   // no red in the Old Testament's opening
		{"MAT", 2, 1},   // the wise men's coming
		{"LUK", 1, 1},   // Luke's dedication
		{"ACT", 28, 30}, // Paul in Rome
		{"ROM", 1, 1},   // no red outside the Gospels
	}
	for _, c := range black {
		if InRed(c.code, c.ch, c.vs) {
			t.Errorf("%s %d:%d reads red", c.code, c.ch, c.vs)
		}
	}
}

// The stories run from creation to the new creation, each with passages to
// read, and every passage names a book the canon has.
func TestStories(t *testing.T) {
	if len(Stories) < 35 {
		t.Fatalf("only %d stories", len(Stories))
	}
	seen := map[string]bool{}
	for _, s := range Stories {
		if s.Code == "" || s.Name == "" || len(s.Refs) == 0 {
			t.Errorf("story %+v is thin", s)
		}
		if seen[s.Code] {
			t.Errorf("two stories share the code %s", s.Code)
		}
		seen[s.Code] = true
		for _, r := range s.Refs {
			if _, ok := Canon[r.Book]; !ok {
				t.Errorf("%s names the unknown book %s", s.Code, r.Book)
			}
		}
	}
	if seen["creation"] && seen["new-creation"] {
		if Stories[0].Code != "creation" || Stories[len(Stories)-1].Code != "new-creation" {
			t.Error("the stories don't run from creation to the new creation")
		}
	}
}
