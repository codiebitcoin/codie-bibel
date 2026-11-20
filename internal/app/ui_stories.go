package app

import (
	"strconv"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/bible"
	"codie-bibel/internal/speech"
)

// The Stories page: the Bible's stories in the order the Bible tells
// them, one open at a time, the device telling it with its turns and
// going on from where it stopped.

func (a *app) storiesPage(c *ui.Context) {
	a.pageHeading(c, T("Stories"), nil)
	top, bottom := a.insets(c)
	if a.storyCode == "" {
		ui.Scroll(c.Key("stories")).Grow(1).Children(func() {
			ui.Column(c).Padding(top, 24, bottom, 24).Gap(6).Children(func() {
				for i, s := range bible.Stories {
					open := a.storyCode == s.Code
					var b ui.Element
					if open {
						b = ui.PrimaryButton(c.Key("story-"+s.Code), s.Name)
					} else {
						b = ui.Button(c.Key("story-"+s.Code), s.Name)
					}
					if b.Clicked() {
						a.storyCode = s.Code
						a.redraw()
					}
					_ = i
					// The reference line under the name.
					ui.Text(c, storyRefsLine(s)).TextColor(a.pal.faint).FontSize(12)
				}
			})
		})
		return
	}
	a.storyView(c, top, bottom)
}

// storyView is the story open: its passages as paragraphs, Listen
// queuing them with the place it stood, and the back to the list.
func (a *app) storyView(c *ui.Context, top, bottom float32) {
	var story *bible.Story
	for i := range bible.Stories {
		if bible.Stories[i].Code == a.storyCode {
			story = &bible.Stories[i]
		}
	}
	if story == nil {
		a.storyCode = ""
		a.redraw()
		return
	}
	key := "story:" + story.Code
	ui.Column(c).Grow(1).Padding(top, 24, bottom, 24).Gap(12).Children(func() {
		ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
			if ui.Button(c.Key("stories-back"), T("Stories")).Clicked() {
				a.storyCode = ""
				a.redraw()
			}
			ui.Text(c, story.Name).FontSize(17).FontWeight(700)
			ui.Box(c).Grow(1)
			if a.bq.Reading() {
				ui.Text(c, T("Reading…")).TextColor(a.pal.faint).FontSize(12)
				if ui.Button(c.Key("story-pause"), T("Pause")).Clicked() {
					a.bq.Pause()
				}
			} else if len(a.bq.Items()) > 0 {
				if ui.Button(c.Key("story-resume"), T("Play")).Clicked() {
					a.bq.Resume()
				}
			}
		})
		ui.Scroll(c.Key("story-body")).Grow(1).Children(func() {
			ui.Column(c).Gap(14).Children(func() {
				for _, r := range story.Refs {
					name, _ := bible.BookName(r.Book)
					endCh := r.EndChapter
					if endCh == 0 {
						endCh = r.Chapter
					}
					for ch := r.Chapter; ch <= endCh; ch++ {
						tr := a.ensureTranslation()
						if tr == nil {
							continue
						}
						b := tr.Book(r.Book)
						if b == nil || ch > len(b.Chapters) {
							continue
						}
						heading := name.En + " " + strconv.Itoa(ch)
						if r.EndChapter == 0 && r.Verse > 0 {
							heading += ":" + strconv.Itoa(r.Verse)
						}
						ui.Text(c, heading).FontSize(12).FontWeight(600).TextColor(a.pal.accent)
						ui.Text(c, joinVerses(b.Chapters[ch-1])).FontSize(a.readerFontSize()).
							TextColor(a.pal.text).LineHeight(1.6)
					}
				}
			})
		})
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			progress := a.store.Progress(key)
			if ui.Button(c.Key("story-listen"), T("Listen")).Clicked() {
				var lines []string
				tr := a.ensureTranslation()
				for _, r := range story.Refs {
					if tr == nil {
						break
					}
					endCh := r.EndChapter
					if endCh == 0 {
						endCh = r.Chapter
					}
					for ch := r.Chapter; ch <= endCh; ch++ {
						b := tr.Book(r.Book)
						if b == nil || ch > len(b.Chapters) {
							continue
						}
						lines = append(lines, joinVerses(b.Chapters[ch-1]))
					}
				}
				passage := joinVerses(lines)
				plan := speech.NewPlan(passage, progress)
				a.bq.Enqueue(key, story.Name, plan.Turns, plan.At)
			}
			if progress > 0 {
				ui.Text(c, T("You stood at turn")+" "+strconv.Itoa(progress)).TextColor(a.pal.faint).FontSize(12)
				if ui.Button(c.Key("story-reset"), T("Back")).Clicked() {
					a.store.SetProgress(key, 0)
					a.redraw()
				}
			}
		})
	})
}

// storyRefsLine is the story's references, as its list reads them.
func storyRefsLine(s bible.Story) string {
	out := ""
	for i, r := range s.Refs {
		if i > 0 {
			out += "; "
		}
		out += r.String()
	}
	return out
}
