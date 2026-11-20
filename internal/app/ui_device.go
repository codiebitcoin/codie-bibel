package app

import (
	"os"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/device"
)

// The Device page: what the paired codie is, the character it speaks as,
// and the pairing guide that walks a new one in. The guide's four steps
// are the tagline's welcome, the address, the code and the character.

// devicePage is the page: the pairing, the character, and the action
// group. The guide walks over it, from the view's root.
func (a *app) devicePage(c *ui.Context) {
	a.pageHeading(c, T("Device"), func() {
		if ui.Button(c.Key("pair"), T("Pair a device")).Clicked() {
			a.wizard.open = true
			a.wizard.step = 0
			a.redraw()
		}
	})
	top, bottom := a.insets(c)
	ui.Column(c).Grow(1).Padding(top, 24, bottom, 24).Gap(12).Children(func() {
		switch {
		case a.dev == nil || !a.dev.Paired():
			ui.Text(c, T("No device is paired. Pair one to hear the Bible read aloud.")).TextColor(a.pal.muted).FontSize(13)
		default:
			name := a.dev.Name()
			if name == "" {
				name = T("the device")
			}
			ui.Text(c, Tf("Paired with %s.", name)).FontSize(15).FontWeight(600)
			a.personaCard(c)
			a.actionCard(c)
		}
	})
}

// personaCard is the character the device speaks as, Jesus or the monk.
func (a *app) personaCard(c *ui.Context) {
	current := a.settings.CodiePersona
	if current == "" {
		current = device.DefaultPersona().Code
	}
	ui.Column(c).Gap(6).Children(func() {
		ui.Text(c, T("Character")).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
		ui.Row(c).Gap(8).Children(func() {
			for _, p := range device.Personas {
				if p.Code == current {
					if ui.PrimaryButton(c.Key("persona-"+p.Name), p.Name).Clicked() {
						a.choosePersona(p.Code)
					}
					continue
				}
				if ui.Button(c.Key("persona-"+p.Name), p.Name).Clicked() {
					a.choosePersona(p.Code)
				}
			}
		})
		for _, p := range device.Personas {
			if p.Code == current {
				ui.Text(c, p.Summary).TextColor(a.pal.muted).FontSize(12.5).MaxWidth(520)
			}
		}
	})
}

// actionCard is the action group: how the device's own motions read,
// which set by the character's sex, and how strong.
func (a *app) actionCard(c *ui.Context) {
	ui.Column(c).Gap(6).Children(func() {
		ui.Text(c, T("Action group")).FontSize(12).FontWeight(600).TextColor(a.pal.muted)
		ui.Text(c, T("The device plays its own motions as it speaks. How strong they read is yours to set.")).TextColor(a.pal.muted).FontSize(12.5).MaxWidth(520)
		intensity := a.personaIntensity()
		ui.Row(c).Gap(8).Children(func() {
			for _, l := range []struct {
				v    int
				name string
			}{{device.IntensitySubtle, "Subtle"}, {device.IntensityNormal, "Normal"}, {device.IntensityStrong, "Strong"}} {
				if intensity == l.v {
					ui.PrimaryButton(c.Key("intensity-"+l.name), l.name)
					continue
				}
				if ui.Button(c.Key("intensity-"+l.name), l.name).Clicked() {
					a.setPersonaIntensity(l.v)
				}
			}
		})
	})
}

// choosePersona makes the device's character the persona of that code,
// with its model where one has been given, and keeps the choice.
func (a *app) choosePersona(code string) {
	p := device.PersonaByCode(code)
	if p == nil {
		return
	}
	a.settings.CodiePersona = code
	a.settings.save()
	a.applyPersona(p)
}

// applyPersona sends a persona to the device: its words, its action
// group, and its face, where a model file has been given. It runs as it
// is asked, and the toast says how it went.
func (a *app) applyPersona(p *device.Persona) {
	if a.dev == nil || !a.dev.Paired() {
		a.toast = T("Pair a device first")
		a.redraw()
		return
	}
	vrm := a.seedVRM(p)
	persona := *p
	go func() {
		err := appPersonaApply(a.dev, persona, vrm)
		a.post(func() {
			switch {
			case err == nil:
				a.toast = Tf("%s speaks on the device now.", persona.Name)
			default:
				a.toast = err.Error()
			}
			a.redraw()
		})
	}()
}

// seedVRM reads the persona's model file from the reader's data, where
// they may have put one, or reads nothing: the device's own model speaks
// then.
func (a *app) seedVRM(p *device.Persona) []byte {
	if p.VRM == "" {
		return nil
	}
	b, err := os.ReadFile(userDataDir(p.VRM))
	if err != nil {
		return nil
	}
	return b
}

// personaIntensity is the strength of the character's motions, as the
// settings keep it, or the persona's own.
func (a *app) personaIntensity() int {
	if p := device.PersonaByCode(a.settings.CodiePersona); p != nil {
		return p.Group.Intensity
	}
	return device.IntensityStrong
}

// setPersonaIntensity sets the strength of the character's motions, on
// the device and in the settings.
func (a *app) setPersonaIntensity(v int) {
	code := a.settings.CodiePersona
	if code == "" {
		code = device.DefaultPersona().Code
	}
	a.settings.CodiePersona = code
	a.settings.save()
	p := device.PersonaByCode(code)
	if p == nil {
		return
	}
	with := *p
	with.Group.Intensity = v
	a.applyPersona(&with)
}
