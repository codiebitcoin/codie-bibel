package app

import (
	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/i18n"
)

// The pages, until each arrives: Read and its siblings grow over the next
// milestones, each taking the place of its placeholder here.

func (a *app) settingsPage(c *ui.Context) {
	a.pageHeading(c, T("Settings"), nil)
	top, bottom := a.insets(c)
	ui.Column(c).Grow(1).Padding(top, 24, bottom, 24).Gap(16).Children(func() {
		ui.Text(c, T("Language")).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
		ui.Row(c).Wrap().Gap(8).Children(func() {
			for _, l := range i18n.Languages {
				chosen := i18n.Language() == l.Code
				if chosen {
					ui.PrimaryButton(c.Key("lang-"+l.Code), l.Name)
					continue
				}
				if ui.Button(c.Key("lang-"+l.Code), l.Name).Clicked() {
					a.setLanguage(l.Code)
				}
			}
		})
		ui.Text(c, T("The reader")).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, T("Text size")).FontSize(13.5)
			for _, s := range []struct {
				v    float64
				name string
			}{{0, N("Default")}, {18, N("Large")}, {22, N("Larger")}} {
				chosen := a.settings.FontSize == s.v
				if chosen {
					ui.PrimaryButton(c.Key("size-"+s.name), T(s.name))
					continue
				}
				if ui.Button(c.Key("size-"+s.name), T(s.name)).Clicked() {
					a.settings.FontSize = s.v
					a.settings.save()
					a.redraw()
				}
			}
		})
		ui.Row(c).Gap(8).Children(func() {
			label := T("Serif text, as print reads")
			clicked := false
			if a.settings.Serif {
				clicked = ui.PrimaryButton(c.Key("serif"), label).Clicked()
			} else {
				clicked = ui.Button(c.Key("serif"), label).Clicked()
			}
			if clicked {
				a.settings.Serif = !a.settings.Serif
				a.settings.save()
				a.redraw()
			}
		})
	})
}
