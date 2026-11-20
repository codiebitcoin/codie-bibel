package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	clt01 "workshop_inner_sdk"
)

// The app's slot on the device: the first of the five the device keeps.
const Slot = 1

// Creds is what a paired device left with the app, to keep in the
// settings and to speak for again.
type Creds struct {
	Addr        string `json:"addr"`
	AuthScope   string `json:"authScope"`
	LocalSecret string `json:"localSecret"`
	LocalSessID string `json:"localSessionID"`
	UserID      string `json:"userID"`
	DeviceID    string `json:"deviceID"`
	Name        string `json:"name"`
}

// Paired says whether the credentials are enough to speak for.
func (c *Creds) Paired() bool { return c != nil && c.Addr != "" && c.LocalSecret != "" }

// Device is the paired codie: what it is asked, in the app's terms. The
// calls go to the workshop SDK's client; a test puts a fake in their
// place, by the API.
type Device struct {
	creds  *Creds
	client *clt01.Client

	utterance string // the reading aloud on its way, to cancel
}

// API is what the app asks of a device's voice: what the queue and the
// pages call, and a test fakes.
type API interface {
	// Speak reads a turn aloud, and returns once it has read.
	Speak(ctx context.Context, text string) error
	// Chat says the text to the character, who answers aloud of itself.
	Chat(ctx context.Context, text string) error
	// Stop stops what is being read.
	Stop(ctx context.Context) error
	// Status says whether the device is reachable.
	Status(ctx context.Context) error
}

var _ API = (*Device)(nil)

// New makes the device of a saved pairing, or of none yet.
func New(creds *Creds) *Device {
	d := &Device{creds: creds}
	d.dial()
	return d
}

// dial makes the client the credentials speak through.
func (d *Device) dial() {
	if !d.creds.Paired() {
		d.client = nil
		return
	}
	c, err := clt01.New(d.creds.Addr)
	if err != nil {
		d.client = nil
		return
	}
	c.Apply(&clt01.Credentials{
		AuthScope:      d.creds.AuthScope,
		LocalSecret:    d.creds.LocalSecret,
		LocalSessionID: d.creds.LocalSessID,
		DeviceID:       d.creds.DeviceID,
		DeviceAddr:     d.creds.Addr,
		UserID:         d.creds.UserID,
	})
	d.client = c
}

// Paired says whether a device is paired.
func (d *Device) Paired() bool { return d.creds.Paired() && d.client != nil }

// Name is what the device calls itself, as the pairing found it.
func (d *Device) Name() string {
	if d.creds == nil {
		return ""
	}
	return d.creds.Name
}

// Pair walks the pairing: the device was asked to show a code, and code
// is what its reader typed. The credentials it hands out are permanent,
// and kept by the caller in the settings.
func (d *Device) Pair(ctx context.Context, addr string, code string) (*Creds, error) {
	c, err := clt01.New(addr)
	if err != nil {
		return nil, err
	}
	challenge, _, err := c.RequestPairing()
	if err != nil {
		return nil, err
	}
	creds, err := c.ConfirmPairing(challenge, code, clt01.AuthScopePermanent)
	if err != nil {
		return nil, err
	}
	d.creds = &Creds{
		Addr:        addr,
		AuthScope:   creds.AuthScope,
		LocalSecret: creds.LocalSecret,
		LocalSessID: creds.LocalSessionID,
		UserID:      creds.UserID,
		DeviceID:    creds.DeviceID,
		Name:        creds.Name,
	}
	dial(d)
	return d.creds, nil
}

// RequestChallenge asks the device to show a pairing code, for a wizard
// that asks the reader to type it before going on.
func (d *Device) RequestChallenge(ctx context.Context, addr string) (challengeID string, expires time.Time, err error) {
	c, err := clt01.New(addr)
	if err != nil {
		return "", time.Time{}, err
	}
	id, expiresMs, err := c.RequestPairing()
	if err != nil {
		return "", time.Time{}, err
	}
	return id, time.UnixMilli(expiresMs), nil
}

// ConfirmPairing completes a pairing begun with RequestChallenge.
func (d *Device) ConfirmPairing(ctx context.Context, addr, challengeID, code string) (*Creds, error) {
	c, err := clt01.New(addr)
	if err != nil {
		return nil, err
	}
	creds, err := c.ConfirmPairing(challengeID, code, clt01.AuthScopePermanent)
	if err != nil {
		return nil, err
	}
	d.creds = &Creds{
		Addr:        addr,
		AuthScope:   creds.AuthScope,
		LocalSecret: creds.LocalSecret,
		LocalSessID: creds.LocalSessionID,
		UserID:      creds.UserID,
		DeviceID:    creds.DeviceID,
		Name:        creds.Name,
	}
	dial(d)
	return d.creds, nil
}

// dial is Device.dial for a device another of its calls made.
func dial(d *Device) { d.dial() }

// cc is the call's context on the device: the app's slot, by the
// character's name.
func (d *Device) cc(name string) *clt01.CallContext {
	return &clt01.CallContext{MachineID: d.creds.DeviceID, Slot: Slot, CharacterName: name}
}

// Speak reads a turn aloud, and returns once it has read. Turn after turn
// is the queue's work; this is one turn of it.
func (d *Device) Speak(ctx context.Context, text string) error {
	if !d.Paired() {
		return errNotPaired
	}
	res, err := d.client.TTSUtteranceCtx(ctx, Slot, text, requestID())
	if err != nil {
		return err
	}
	d.utterance = res.UtteranceID
	done, ok := d.client.WaitTTSUtteranceDone(ctx, res.UtteranceID, 5*time.Minute)
	if !ok {
		return fmt.Errorf("the device didn't say it was done reading")
	}
	if done != "played" && done != "done" {
		return fmt.Errorf("the device stopped reading: %s", done)
	}
	d.utterance = ""
	return nil
}

// Chat says the text to the character, who answers aloud of itself: what
// a confession's words are for. The device writes nothing of it back to
// the app, and the app keeps nothing of it.
func (d *Device) Chat(ctx context.Context, text string) error {
	if !d.Paired() {
		return errNotPaired
	}
	_, err := d.client.PlayChat(ctx, d.cc(""), text, nil)
	return err
}

// Stop stops what is being read.
func (d *Device) Stop(ctx context.Context) error {
	if !d.Paired() || d.utterance == "" {
		return nil
	}
	err := d.client.TTSUtteranceCancel(ctx, d.utterance)
	d.utterance = ""
	return err
}

// Status says whether the device is reachable.
func (d *Device) Status(ctx context.Context) error {
	if !d.Paired() {
		return errNotPaired
	}
	_, err := d.client.DeviceState(ctx)
	return err
}

// ApplyPersona dresses the device's character as a persona of the app's:
// its name and its words, its action group, and its face, where a model
// has been given.
func (d *Device) ApplyPersona(ctx context.Context, p Persona, vrm []byte) error {
	if !d.Paired() {
		return errNotPaired
	}
	cc := d.cc(p.Name)
	if _, err := d.client.SetInfo(ctx, cc, &clt01.CharacterSettings{
		Name:               p.Name,
		Description:        p.Summary,
		CharacterPrompt:    p.Prompt,
		Sex:                p.Group.Sex,
		AnimationIntensity: p.Group.Intensity,
	}, false); err != nil {
		return err
	}
	if len(vrm) > 0 {
		name := p.VRM
		if name == "" {
			name = "model.vrm"
		}
		if _, err := d.client.SetVrm(ctx, cc, name, vrm, true, nil); err != nil {
			return err
		}
	}
	return nil
}

// SetAlarm asks the device to speak at a time, as a feast's reminder does.
func (d *Device) SetAlarm(ctx context.Context, label string, at time.Time) error {
	if !d.Paired() {
		return errNotPaired
	}
	alarm := map[string]any{"name": label, "time": at.Format("15:04"), "enabled": true}
	b, err := json.Marshal(alarm)
	if err != nil {
		return err
	}
	_, err = d.client.AlarmCreate(ctx, b)
	return err
}

// errNotPaired is what the device answers while there is none.
var errNotPaired = errors.New("no device is paired")

func requestID() string {
	return fmt.Sprintf("codie-bibel-%d", time.Now().UnixNano())
}
