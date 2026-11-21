package app

import (
	"context"
	"time"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/device"
)

// The pairing guide: four steps, from the welcome to the device speaking.
// Its network goes through the pair hooks, which tests put fakes for.

// wizardState is where the guide stands.
type wizardState struct {
	open bool
	step int // 0 welcome, 1 address, 2 code, 3 character

	addr      string
	challenge string
	expires   time.Time
	code      string
	chosen    string // the character's code, as the last step has it

	busy bool
	err  string
}

// pairRequest asks the device at addr to show its pairing code. The app's
// own asks the device; a test puts a fake here.
var appPairRequest = func(d *device.Device, addr string) (string, time.Time, error) {
	return d.RequestChallenge(context.Background(), addr)
}

// pairConfirm completes the pairing with the code the reader typed.
var appPairConfirm = func(d *device.Device, addr, challenge, code string) (*device.Creds, error) {
	return d.ConfirmPairing(context.Background(), addr, challenge, code)
}

// personaApply dresses the character on the device. The app's own sends
// the persona's words, its action group and its model; a test puts a
// fake here.
var appPersonaApply = func(d *device.Device, p device.Persona, vrm []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	return d.ApplyPersona(ctx, p, vrm)
}

// wizardView is the guide, a dialog over the Device page.
func (a *app) wizardView(c *ui.Context) {
	ui.DialogBase(c.Key("pairing-guide"), &a.wizard.open, func(backdrop, panel ui.Element) {
		backdrop.Background(ui.RGBA(0, 0, 0, 0.45))
		panel.Width(560).Padding(24).Radius(16).Gap(14).
			Background(a.pal.bg).Border(1, a.pal.hairline).Shadow(0, 24, 60, -12, ui.RGBA(0, 0, 0, 0.4))
		ui.Column(c).Gap(14).Children(func() {
			switch a.wizard.step {
			case 0:
				a.wizardWelcome(c)
			case 1:
				a.wizardAddress(c)
			case 2:
				a.wizardCode(c)
			case 3:
				a.wizardCharacter(c)
			}
			if a.wizard.err != "" {
				ui.Text(c, a.wizard.err).TextColor(a.pal.accent).FontSize(13)
			}
			a.wizardFooter(c)
		})
	})
}

// wizardWelcome is the first step: what the app is, and what a device
// gives it.
func (a *app) wizardWelcome(c *ui.Context) {
	ui.Text(c, "Make your Codie as Bibel").FontSize(20).FontWeight(700)
	ui.Text(c, T("Pair your codie device, and it reads the Bible aloud: the stories, the psalms, the verses of the day. It speaks as Jesus, or as a monk for the quieter rooms.")).TextColor(a.pal.muted).FontSize(13.5)
	ui.Text(c, T("Turn the device on, and join the same Wi-Fi as this computer.")).TextColor(a.pal.muted).FontSize(13)
}

// wizardAddress is the second step: the device's address on the network.
func (a *app) wizardAddress(c *ui.Context) {
	ui.Text(c, T("The device's address")).FontSize(20).FontWeight(700)
	ui.Text(c, T("Type the address the device's screen shows, as 192.168.1.10:8081.")).TextColor(a.pal.muted).FontSize(13)
	ui.TextInput(c.Key("address"), &a.wizard.addr).Placeholder("192.168.1.10:8081").Label(T("Device address")).Width(280).FontSize(13.5)
	ui.Text(c, T("The port is 8081 unless the device says otherwise.")).TextColor(a.pal.faint).FontSize(12)
}

// wizardCode is the third step: the code the device shows.
func (a *app) wizardCode(c *ui.Context) {
	ui.Text(c, T("The pairing code")).FontSize(20).FontWeight(700)
	ui.Text(c, T("Type the six figures the device's screen shows.")).TextColor(a.pal.muted).FontSize(13)
	ui.TextInput(c.Key("code"), &a.wizard.code).Placeholder("000000").Label(T("Pairing code")).Width(160).FontSize(15)
}

// wizardCharacter is the last step: who speaks.
func (a *app) wizardCharacter(c *ui.Context) {
	ui.Text(c, T("Who speaks?")).FontSize(20).FontWeight(700)
	chosen := a.wizard.chosen
	if chosen == "" {
		chosen = device.DefaultPersona().Code
	}
	for _, p := range device.Personas {
		if p.Code == chosen {
			ui.PrimaryButton(c.Key("choose-"+p.Name), p.Name)
		} else if ui.Button(c.Key("choose-"+p.Name), p.Name).Clicked() {
			a.wizard.chosen = p.Code
			a.redraw()
		}
		ui.Text(c, p.Summary).TextColor(a.pal.muted).FontSize(12.5).MaxWidth(480)
	}
	if seed := device.PersonaByCode(chosen); seed != nil {
		if a.seedVRMPersonal(seed) != nil {
			ui.Text(c, Tf("%s wears your model file %s.", seed.Name, seed.VRM)).TextColor(a.pal.faint).FontSize(12)
		} else {
			ui.Text(c, Tf("%s wears the app's default figure; put your own model at %s to replace it.", seed.Name, seed.VRM)).TextColor(a.pal.faint).FontSize(12)
		}
	}
}

// wizardFooter is the guide's buttons, whatever the step.
func (a *app) wizardFooter(c *ui.Context) {
	ui.Row(c).Gap(8).Justify(ui.End).Children(func() {
		if ui.Button(c.Key("guide-skip"), T("Not now")).Clicked() {
			a.wizard.close()
			a.settings.CodieSetupDone = true
			a.settings.save()
			a.redraw()
		}
		ui.Box(c).Grow(1)
		if a.wizard.step > 0 && !a.wizard.busy {
			if ui.Button(c.Key("guide-back"), T("Back")).Clicked() {
				a.wizard.step--
				a.redraw()
			}
		}
		switch a.wizard.step {
		case 0:
			if ui.PrimaryButton(c.Key("guide-next"), T("Pair a device")).Clicked() {
				a.wizard.step = 1
				a.redraw()
			}
		case 1:
			if a.wizard.busy {
				ui.Text(c, T("Asking the device…")).TextColor(a.pal.muted)
			} else if ui.PrimaryButton(c.Key("guide-next"), T("Next")).Clicked() {
				a.requestChallenge()
			}
		case 2:
			if a.wizard.busy {
				ui.Text(c, T("Pairing…")).TextColor(a.pal.muted)
			} else if ui.PrimaryButton(c.Key("guide-next"), T("Pair")).Clicked() {
				a.confirmPairing()
			}
		case 3:
			if a.wizard.busy {
				ui.Text(c, T("Dressing the character…")).TextColor(a.pal.muted)
			} else if ui.PrimaryButton(c.Key("guide-next"), T("Done")).Clicked() {
				a.finishWizard()
			}
		}
	})
}

// requestChallenge asks the device to show its code.
func (a *app) requestChallenge() {
	addr := trimSpace(a.wizard.addr)
	if addr == "" {
		a.wizard.err = T("Type the device's address first.")
		a.redraw()
		return
	}
	a.wizard.busy, a.wizard.err = true, ""
	a.redraw()
	d := a.dev
	if d == nil {
		d = device.New(nil)
	}
	go func() {
		challenge, expires, err := appPairRequest(d, addr)
		a.post(func() {
			a.wizard.busy = false
			if err != nil {
				a.wizard.err = err.Error()
				a.redraw()
				return
			}
			a.wizard.challenge, a.wizard.expires = challenge, expires
			a.wizard.step = 2
			a.redraw()
		})
	}()
}

// confirmPairing completes the pairing with the code, keeps the device,
// and walks to the character.
func (a *app) confirmPairing() {
	code := trimSpace(a.wizard.code)
	if len(code) != 6 {
		a.wizard.err = T("The code is six figures.")
		a.redraw()
		return
	}
	a.wizard.busy, a.wizard.err = true, ""
	a.redraw()
	d := a.dev
	if d == nil {
		d = device.New(nil)
	}
	go func() {
		creds, err := appPairConfirm(d, trimSpace(a.wizard.addr), a.wizard.challenge, code)
		a.post(func() {
			a.wizard.busy = false
			if err != nil {
				a.wizard.err = err.Error()
				a.redraw()
				return
			}
			a.settings.Codie = creds
			a.settings.save()
			a.wireDevice()
			a.wizard.step = 3
			a.redraw()
		})
	}()
}

// finishWizard dresses the chosen character on the device, and ends the
// guide.
func (a *app) finishWizard() {
	code := a.wizard.chosen
	if code == "" {
		code = device.DefaultPersona().Code
	}
	a.wizard.busy, a.wizard.err = true, ""
	a.redraw()
	a.settings.CodiePersona = code
	a.settings.save()
	p := device.PersonaByCode(code)
	vrm := a.seedVRM(p)
	persona := *p
	go func() {
		err := appPersonaApply(a.dev, persona, vrm)
		a.post(func() {
			a.wizard.busy = false
			a.wizard.close()
			a.settings.CodieSetupDone = true
			a.settings.save()
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

// close shuts the guide.
func (w *wizardState) close() {
	w.open, w.step, w.err, w.busy = false, 0, "", false
	w.addr, w.code, w.challenge, w.chosen = "", "", "", ""
}

// trimSpace is strings.TrimSpace, for the wizard's fields.
func trimSpace(s string) string {
	out := s
	for len(out) > 0 && (out[0] == ' ' || out[0] == '\t') {
		out = out[1:]
	}
	for len(out) > 0 && (out[len(out)-1] == ' ' || out[len(out)-1] == '\t') {
		out = out[:len(out)-1]
	}
	return out
}
