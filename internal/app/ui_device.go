package app

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/assets"
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
		paired := a.dev != nil && a.dev.Paired()
		if paired {
			// The device calls itself by a serial; the reader doesn't need
			// to read it, so the page says only that one is paired.
			ui.Text(c, T("Paired with your Codie.")).FontSize(15).FontWeight(600)
		} else {
			ui.Row(c).Gap(10).AlignItems(ui.Center).Children(func() {
				ui.Text(c, T("No device is paired. Pair one to hear the Bible read aloud.")).TextColor(a.pal.muted).FontSize(13)
				ui.Box(c).Grow(1)
				if ui.PrimaryButton(c.Key("pair-here"), T("Pair a device")).Clicked() {
					a.wizard.open, a.wizard.step = true, 0
					a.redraw()
				}
			})
			ui.Text(c, T("You can still choose the character and give it a model and a voice now; the pairing sends them over.")).TextColor(a.pal.faint).FontSize(12)
		}
		a.personaCard(c)
		if paired {
			a.actionCard(c)
		}
	})
}

// personaCard is the character the device speaks as, Jesus or the monk,
// with the model and the voice the reader may give it. The device's own
// face and voice answer when none has been given.
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
		for i := range device.Personas {
			p := &device.Personas[i]
			if p.Code != current {
				continue
			}
			// The figure itself, so the reader sees who they are dressing.
			if img := personaPreview(p.Code); img != nil {
				ui.Row(c).Gap(14).AlignItems(ui.Start).Children(func() {
					ui.Image(c, img).Size(148, 148).Fit(ui.Contain).Radius(12)
					ui.Column(c).Gap(4).Children(func() {
						ui.Text(c, p.Summary).TextColor(a.pal.muted).FontSize(12.5).MaxWidth(430)
						a.personaModelVoiceLines(c, p)
					})
				})
			} else {
				ui.Text(c, p.Summary).TextColor(a.pal.muted).FontSize(12.5).MaxWidth(520)
				a.personaModelVoiceLines(c, p)
			}
			// Give it another face, and another voice.
			ui.Row(c).Gap(8).Children(func() {
				if ui.Button(c.Key("model-upload-"+p.Name), T("Upload a model…")).Clicked() {
					a.uploadModel(p)
				}
				voiceKey := "voice-upload-" + p.Name
				if a.dev != nil && a.dev.Paired() {
					if ui.Button(c.Key(voiceKey), T("Upload a voice…")).Clicked() {
						a.uploadVoice(p)
					}
				} else {
					ui.Button(c.Key(voiceKey), T("Upload a voice…")).Disabled(true)
				}
			})
			// Where a model file the reader downloads goes.
			if p.VRM != "" && a.seedVRMPersonal(p) == nil {
				ui.Text(c, Tf("Your own model goes at %s", filepath.Join(userDataDir(""), p.VRM))).TextColor(a.pal.faint).FontSize(11.5)
			}
			if a.dev == nil || !a.dev.Paired() {
				ui.Text(c, T("Cloning a voice needs the paired device; pair one and come back to it.")).TextColor(a.pal.faint).FontSize(11.5)
			}
		}
	})
}

// personaModelVoiceLines is the pair of faint lines under a persona's
// summary: which model file it wears, and which voice speaks.
func (a *app) personaModelVoiceLines(c *ui.Context, p *device.Persona) {
	face := T("Wearing the app's default figure; upload a model to replace it.")
	switch p.Code {
	case "jesus":
		face = T("Wearing the app's figure of Jesus; upload a model to replace it.")
	case "monk":
		face = T("Wearing the app's figure of the monk; upload a model to replace it.")
	}
	if a.seedVRMPersonal(p) != nil {
		face = Tf("Wearing the model file %s.", p.VRM)
	}
	ui.Text(c, face).TextColor(a.pal.faint).FontSize(12)
	voice := T("The device's own voice; upload a sample to clone the reader's.")
	if v := a.settings.CodieVoice[p.Code]; v != "" {
		voice = Tf("Voice %s, cloned on the device.", v)
	}
	ui.Text(c, voice).TextColor(a.pal.faint).FontSize(12)
}

// personaPreview is the figure the persona wears, as a picture: rendered
// once from the very models the app ships, and decoded once.
var personaPreviews = struct {
	once  sync.Once
	jesus ui.ImageSource
	monk  ui.ImageSource
}{}

func personaPreview(code string) ui.ImageSource {
	personaPreviews.once.Do(func() {
		personaPreviews.jesus, _ = ui.DecodeImage(assets.JesusPreview)
		personaPreviews.monk, _ = ui.DecodeImage(assets.MonkPreview)
	})
	if code == "jesus" {
		return personaPreviews.jesus
	}
	return personaPreviews.monk
}

// uploadModel asks for a VRM file, keeps it as the persona's model, and
// dresses the device in it.
func (a *app) uploadModel(p *device.Persona) {
	paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
		Title:       Tf("A model for %s", p.Name),
		ButtonLabel: T("Open"),
		Filters:     []mygo.FileFilter{{Name: "VRM model", Extensions: []string{"vrm"}}},
	})
	if err != nil || len(paths) == 0 {
		return
	}
	b, err := os.ReadFile(paths[0])
	if err != nil {
		a.toast = err.Error()
		a.redraw()
		return
	}
	// Kept where the pairing guide finds it, so the next pairing wears it
	// too.
	name := p.VRM
	if name == "" {
		name = "model.vrm"
	}
	if err := os.WriteFile(userDataDir(name), b, 0o644); err != nil {
		a.toast = err.Error()
		a.redraw()
		return
	}
	with := *p
	a.applyPersona(&with, b)
}

// uploadVoice asks for a sample of a voice, clones it on the device, and
// wears it on the character. The cloning is the device's own work and
// takes its time.
func (a *app) uploadVoice(p *device.Persona) {
	paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
		Title:       Tf("A voice for %s", p.Name),
		ButtonLabel: T("Open"),
		Filters:     []mygo.FileFilter{{Name: "Audio", Extensions: []string{"wav", "mp3", "m4a", "ogg", "flac"}}},
	})
	if err != nil || len(paths) == 0 {
		return
	}
	audio, err := os.ReadFile(paths[0])
	if err != nil {
		a.toast = err.Error()
		a.redraw()
		return
	}
	filename := filepath.Base(paths[0])
	code := p.Code
	a.toast = T("Cloning the voice; this takes a while…")
	a.redraw()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
		defer cancel()
		voiceID, err := a.dev.CloneVoice(ctx, audio, filename, p.Name)
		a.post(func() {
			switch {
			case err == nil && voiceID != "":
				if a.settings.CodieVoice == nil {
					a.settings.CodieVoice = map[string]string{}
				}
				a.settings.CodieVoice[code] = voiceID
				a.settings.save()
				with := *p
				with.VoiceID = voiceID
				a.applyPersona(&with, nil)
			case err != nil:
				a.toast = err.Error()
				a.redraw()
			}
		})
	}()
}

// seedVoiceClone clones a voice on the device for the seed blessing. The
// app's own asks the device; a test puts a fake here.
var seedVoiceClone = func(d *device.Device, audio []byte, name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()
	return d.CloneVoice(ctx, audio, "jesus-voice.wav", name)
}

// seedVoice clones the app's own blessing as the Jesus persona's voice,
// once, at the first pairing: it is the voice the character speaks with
// until the reader clones one of their own. It is the device's work and
// takes its time; until it is done the device's own voice answers.
func (a *app) seedVoice() {
	if a.dev == nil || !a.dev.Paired() || a.settings.CodieSeedVoice {
		return
	}
	if a.settings.CodieVoice["jesus"] != "" {
		a.settings.CodieSeedVoice = true
		a.settings.save()
		return
	}
	a.settings.CodieSeedVoice = true
	a.settings.save()
	d, audio := a.dev, assets.SeedVoice
	go func() {
		voiceID, err := seedVoiceClone(d, audio, "Jesus")
		a.post(func() {
			if err != nil || voiceID == "" {
				return // the device's own voice carries on
			}
			if a.settings.CodieVoice == nil {
				a.settings.CodieVoice = map[string]string{}
			}
			a.settings.CodieVoice["jesus"] = voiceID
			a.settings.save()
			a.syncPersona()
		})
	}()
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
	a.applyPersona(p, nil)
}

// applyPersona sends a persona to the device: its words, its action
// group, its voice, and its face — the model bytes given win over the
// persona's file. It runs as it is asked, and the toast says how it
// went.
func (a *app) applyPersona(p *device.Persona, vrmOverride []byte) {
	if a.dev == nil || !a.dev.Paired() {
		a.toast = T("Pair a device first")
		a.redraw()
		return
	}
	vrm := vrmOverride
	if vrm == nil {
		vrm = a.seedVRM(p)
	}
	with := *p
	if v := a.settings.CodieVoice[p.Code]; v != "" {
		with.VoiceID = v
	}
	persona := with
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

// seedVRMPersonal reads the reader's own model file for a persona, or
// nil when they haven't given one.
func (a *app) seedVRMPersonal(p *device.Persona) []byte {
	if p.VRM == "" {
		return nil
	}
	b, err := os.ReadFile(userDataDir(p.VRM))
	if err != nil {
		return nil
	}
	return b
}

// seedVRM reads the persona's model: the reader's own file where they
// gave one, else the figure the app bundles for them — Jesus his own,
// the monk his vested cleric.
func (a *app) seedVRM(p *device.Persona) []byte {
	if b := a.seedVRMPersonal(p); b != nil {
		return b
	}
	switch p.Code {
	case "jesus":
		return assets.JesusFigure
	case "monk":
		return assets.MonkFigure
	}
	return assets.DefaultFigure
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
	a.applyPersona(&with, nil)
}
