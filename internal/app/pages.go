package app

import (
	"github.com/egoist/mygo/ui"
)

// The pages, until each arrives: Read and its siblings grow over the next
// milestones, each taking the place of its placeholder here.

func (a *app) readPage(c *ui.Context) {
	a.pageHeading(c, T("Read"), nil)
	a.coming(c, T("The Bible, book by book, chapter by chapter."))
}

func (a *app) devotionsPage(c *ui.Context) {
	a.pageHeading(c, T("Devotions"), nil)
	a.coming(c, T("The verse of the day, Morning and Evening, and reading plans."))
}

func (a *app) prayPage(c *ui.Context) {
	a.pageHeading(c, T("Pray"), nil)
	a.coming(c, T("A guided prayer, and the prayers you keep."))
}

func (a *app) storiesPage(c *ui.Context) {
	a.pageHeading(c, T("Stories"), nil)
	a.coming(c, T("The Bible's stories, told aloud."))
}

func (a *app) seasonsPage(c *ui.Context) {
	a.pageHeading(c, T("Seasons"), nil)
	a.coming(c, T("The church year, its feasts, and this week's readings."))
}

func (a *app) musicPage(c *ui.Context) {
	a.pageHeading(c, T("Music"), nil)
	a.coming(c, T("Hymns to play quietly in the background."))
}

func (a *app) settingsPage(c *ui.Context) {
	a.pageHeading(c, T("Settings"), nil)
	a.coming(c, T("Language, appearance, and the reader."))
}
