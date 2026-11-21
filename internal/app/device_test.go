package app

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"codie-bibel/internal/device"
)

// hasText says the frame reads that text somewhere.
func hasText(tt *ui.Tester, s string) bool {
	for _, t := range tt.Texts() {
		if t == s {
			return true
		}
	}
	return false
}

// waitFor waits for cond, the goroutines' answers coming through the
// app's post, or fails of waiting too long.
func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("what was waited for never came")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// The guide walks the four steps and leaves the device paired: the
// welcome, the address, the code, and the character.
func TestWizardPairs(t *testing.T) {
	a := newTestApp()
	tt := ui.NewTester(a.view, 1120, 760)

	personaApply := appPersonaApply
	appPersonaApply = func(_ *device.Device, p device.Persona, _ []byte) error {
		return nil
	}
	pairRequest, pairConfirm := appPairRequest, appPairConfirm
	appPairRequest = func(_ *device.Device, _ string) (string, time.Time, error) {
		return "challenge-1", time.Now().Add(time.Minute), nil
	}
	appPairConfirm = func(_ *device.Device, addr, challenge, code string) (*device.Creds, error) {
		if addr != "192.168.1.10:8081" || challenge != "challenge-1" || code != "123456" {
			t.Errorf("the pairing went out as %q %q %q", addr, challenge, code)
		}
		return &device.Creds{Addr: addr, LocalSecret: "s", Name: "Codie"}, nil
	}
	defer func() {
		appPairRequest, appPairConfirm, appPersonaApply = pairRequest, pairConfirm, personaApply
	}()

	a.wizard.open = true
	a.wizard.step = 1
	a.wizard.addr = "192.168.1.10:8081"
	settled(tt)
	tt.Click("Next")
	waitFor(t, func() bool { return a.wizard.step == 2 })
	settled(tt)
	// The code goes into the field the reader types into, not the state
	// behind it.
	tt.Click("Pairing code")
	tt.Type("123456")
	waitFor(t, func() bool { return a.wizard.code == "123456" })
	settled(tt)
	tt.Click("Pair")
	waitFor(t, func() bool { return a.wizard.step == 3 })
	if a.settings.Codie == nil {
		t.Error("the credentials weren't kept")
	}
	settled(tt)
	tt.Click("Done")
	waitFor(t, func() bool { return !a.wizard.open })
	if !a.settings.CodieSetupDone {
		t.Error("the guide's done wasn't kept")
	}
}

// Choosing a character keeps it, and the monk's action group is the
// masculine set as the persona has it.
func TestPersonaChoice(t *testing.T) {
	a := newTestApp()
	a.settings.Codie = &device.Creds{Addr: "192.168.1.10:8081", LocalSecret: "s"}
	a.dev = device.New(a.settings.Codie)
	a.choosePersona("monk")
	if a.settings.CodiePersona != "monk" {
		t.Errorf("the character reads %q", a.settings.CodiePersona)
	}
	if p := device.PersonaByCode(a.settings.CodiePersona); p == nil || p.Group.Sex != device.SexMasculine {
		t.Error("the monk's action group isn't the masculine set")
	}
}

// The personas are Jesus first and the monk beside him.
func TestPersonas(t *testing.T) {
	if len(device.Personas) != 2 {
		t.Fatalf("the app keeps %d personas", len(device.Personas))
	}
	if device.Personas[0].Code != "jesus" || device.Personas[1].Code != "monk" {
		t.Errorf("the personas read %q, %q", device.Personas[0].Code, device.Personas[1].Code)
	}
	for _, p := range device.Personas {
		if p.Prompt == "" {
			t.Errorf("%s has no words of its own", p.Code)
		}
	}
}

// The device page says what the pairing is, paired or not.
func TestDevicePageReads(t *testing.T) {
	a := newTestApp()
	a.page = pageDevice
	tt := ui.NewTester(a.view, 1120, 760)
	if !hasText(tt, T("No device is paired. Pair one to hear the Bible read aloud.")) {
		t.Error("the page doesn't say there is no device")
	}
	a.settings.Codie = &device.Creds{Addr: "x", LocalSecret: "s", Name: "Codie"}
	a.wireDevice()
	settled(tt)
	// The device's own name is a serial; the page says only that a Codie
	// is paired, never the serial itself.
	if !hasText(tt, T("Paired with your Codie.")) {
		t.Error("the page doesn't say the pairing")
	}
	if hasText(tt, "Codie-b") {
		t.Error("the page shows the device's serial")
	}
}
