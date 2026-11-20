package speech

import (
	"strings"
	"testing"
)

// A passage cuts into turns of a breath: short sentences merge towards
// the limit, and nothing of the passage is lost.
func TestTurns(t *testing.T) {
	p := "In the beginning God created the heaven and the earth. And the earth was without form, and void; and darkness was upon the face of the deep. And God said, Let there be light: and there was light."
	turns := Turns(p)
	if len(turns) != 1 {
		t.Fatalf("the passage cut into %d turns: short ones merge", len(turns))
	}
	if !strings.HasPrefix(turns[0], "In the beginning") || !strings.HasSuffix(turns[0], "there was light.") {
		t.Errorf("the turn reads %q", turns[0])
	}
	// A passage long enough cuts at the ends of its sentences.
	long := strings.Repeat("The mercy of the Lord endures for ever and ever. ", 4)
	turns = Turns(long + long)
	if len(turns) != 2 {
		t.Errorf("the long passage cut into %d turns", len(turns))
	}
	joined := strings.Join(strings.Fields(strings.Join(turns, " ")), " ")
	if joined != strings.Join(strings.Fields(long+long), " ") {
		t.Error("the turns lost some of the passage")
	}
}

// No turn runs past its breath, whatever the sentence did.
func TestTurnsLength(t *testing.T) {
	p := strings.Repeat("And it came to pass, ", 60)
	for _, turn := range Turns(p) {
		if len([]rune(turn)) > MaxLen*2 {
			t.Errorf("a turn runs %d runes", len([]rune(turn)))
		}
	}
}

// Chinese merges and cuts by its own ends of sentences.
func TestTurnsChinese(t *testing.T) {
	p := "神爱世人，甚至将他的独生子赐给他们，叫一切信他的，不至灭亡，反得永生。因为神差他的儿子降世，不是要定世人的罪，乃是要叫世人因他得救。"
	turns := Turns(p)
	if len(turns) != 1 {
		t.Fatalf("the short passage cut into %d turns", len(turns))
	}
	if !strings.HasSuffix(turns[0], "因他得救。") {
		t.Errorf("the turn ends %q", turns[0])
	}
	// Long enough, it cuts after 。 as English cuts after its stops.
	long := strings.Repeat("耶和华是我的牧者，我必不至缺乏。", 24)
	if got := len(Turns(long)); got < 2 {
		t.Errorf("the long Chinese passage cut into %d turns", got)
	}
}

// A plan reads on from where it stood, the turn it stood in again.
func TestPlanResumes(t *testing.T) {
	long := strings.Repeat("The mercy of the Lord endures for ever and ever. ", 4) // ~200 runes
	passage := long + long + long
	if got := len(Turns(passage)); got != 3 {
		t.Fatalf("the passage planned %d turns", got)
	}
	plan := NewPlan(passage, 1)
	turn, ok := plan.Next()
	if !ok || turn != plan.Turns[1] {
		t.Errorf("the plan went on at %q", turn)
	}
	for {
		_, ok := plan.Next()
		if !ok {
			break
		}
	}
	if !plan.Done() {
		t.Error("the plan never finished")
	}
	if _, ok := plan.Next(); ok {
		t.Error("the plan read past its end")
	}
}

// An empty passage plans to nothing.
func TestPlanEmpty(t *testing.T) {
	p := NewPlan("   ", 0)
	if p.Done() != true || len(p.Turns) != 0 {
		t.Error("an empty passage planned turns")
	}
}
