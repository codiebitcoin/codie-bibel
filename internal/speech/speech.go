// Package speech turns a passage into the small turns the device reads:
// sentences cut to a breath's length, and the position through them, so a
// story or a reading can go on where it stopped.
package speech

import (
	"strings"
	"unicode"
)

// MaxLen is the target length of one turn, in runes: long enough that the
// device's overhead between turns doesn't crowd the reading out, short
// enough for one breath.
const MaxLen = 240

// Turns cuts a passage into turns to read aloud: at the ends of sentences
// first, adjacent short sentences merged towards MaxLen, a sentence too
// long of itself hard-wrapped at spaces and commas. CJK and Latin read
// alike here.
func Turns(passage string) []string {
	passage = strings.TrimSpace(passage)
	if passage == "" {
		return nil
	}
	runes := []rune(passage)

	// The sentences, at their ends' punctuation, the closing quotes and
	// brackets kept with the sentence they close.
	var sent []string
	var cur []rune
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		cur = append(cur, r)
		if !strings.ContainsRune("。！？!?；;", r) {
			continue
		}
		// A period with no space after may be a decimal or an
		// abbreviation: it doesn't end a sentence.
		if r == '.' && i+1 < len(runes) && !unicode.IsSpace(runes[i+1]) {
			continue
		}
		for i+1 < len(runes) && strings.ContainsRune("”’\"')）】」』", runes[i+1]) {
			i++
			cur = append(cur, runes[i])
		}
		sent = append(sent, string(cur))
		cur = nil
	}
	if rest := strings.TrimSpace(string(cur)); rest != "" {
		sent = append(sent, rest)
	}

	// Short sentences merge; a long one wraps.
	var out []string
	var sb strings.Builder
	size := 0
	flush := func() {
		if s := strings.TrimSpace(sb.String()); s != "" {
			out = append(out, s)
		}
		sb.Reset()
		size = 0
	}
	for _, s := range sent {
		n := len([]rune(s))
		if n > MaxLen {
			flush()
			out = append(out, wrap(s, MaxLen)...)
			continue
		}
		if size > 0 && size+n > MaxLen {
			flush()
		}
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(s)
		size += n + 1
	}
	flush()
	return out
}

// wrap cuts a sentence too long of itself at spaces and commas, nearest
// under the limit.
func wrap(s string, limit int) []string {
	runes := []rune(s)
	var out []string
	at := 0
	for at < len(runes) {
		if len(runes)-at <= limit {
			out = append(out, strings.TrimSpace(string(runes[at:])))
			break
		}
		cut := at + limit
		best := -1
		for i := cut; i > at+limit/2; i-- {
			if strings.ContainsRune(" \t，、", runes[i-1]) {
				best = i
				break
			}
		}
		if best < 0 {
			best = cut
		}
		out = append(out, strings.TrimSpace(string(runes[at:best])))
		at = best
	}
	return out
}

// Position is where a reading stands: the turn that plays, of the turns
// the passage has.
type Position struct {
	Turn    int
	OfTurns int
}

// Plan is a passage planned for reading: its turns, and where a reading
// of it has come to.
type Plan struct {
	Text  string
	Turns []string
	At    int
}

// NewPlan plans a passage, from where it last read, so that going on
// reads the turn it stopped in again.
func NewPlan(passage string, at int) *Plan {
	turns := Turns(passage)
	if at < 0 {
		at = 0
	}
	if at > len(turns) {
		at = len(turns)
	}
	return &Plan{Text: passage, Turns: turns, At: at}
}

// Next is the turn to read next, and whether there is one.
func (p *Plan) Next() (string, bool) {
	if p.At >= len(p.Turns) {
		return "", false
	}
	t := p.Turns[p.At]
	p.At++
	return t, true
}

// Done says whether the whole passage has read.
func (p *Plan) Done() bool { return p.At >= len(p.Turns) }

// Position is where the plan stands.
func (p *Plan) Position() Position { return Position{Turn: p.At, OfTurns: len(p.Turns)} }
