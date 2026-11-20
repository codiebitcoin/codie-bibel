// Package device is the paired codie device: pairing, the reading aloud,
// the dialogue the device answers, and the character that speaks — its
// face and its action group. The device plays its own built-in motions as
// it talks; what an app sets of the action group is the character's sex,
// which picks the motion set, and how strong the motions read.
package device

// ActionGroup is the character's motion set and its strength, what the
// device calls sex and animation intensity. The device ships a library of
// motions it plays of itself as it speaks; there is nothing to upload for
// them, and no way to pick one alone.
type ActionGroup struct {
	// Sex picks the motion set: SexMasculine, SexFeminine or SexNeutral.
	Sex int
	// Intensity is how strong the motions read: IntensitySubtle,
	// IntensityNormal or IntensityStrong.
	Intensity int
}

// The device's own values. The default strength is Strong, as the device
// ships.
const (
	SexMasculine = 1
	SexFeminine  = 2
	SexNeutral   = 3

	IntensitySubtle = 1
	IntensityNormal = 2
	IntensityStrong = 3
)

// Persona is a character the device may wear and speak as: its name, what
// it is to say of itself, and its action group.
type Persona struct {
	Code    string // stable key, as "jesus"
	Name    string // as it reads, in English
	Summary string // one line under the name
	Prompt  string // what the device hears as its character's own words
	VRM     string // the embedded model's file, or "" for the device's own
	VoiceID string // a cloned voice of the device's, or "" for its default
	Group   ActionGroup
}

// Personas are the app's characters. Jesus is the first: the device tells
// the stories and reads with his voice and face. The monk stands beside
// him for the places a figure of Jesus would be too much: he speaks the
// same scriptures, in a friar's habit. Each may wear the device's own
// model when no model file has been given it.
var Personas = []Persona{
	{
		Code:    "jesus",
		Name:    "Jesus",
		Summary: "Reads the Scriptures, and tells the stories himself.",
		Prompt: "You are Jesus, gentle and lowly of heart. You speak from the " +
			"Scriptures, and of the Father's love. You answer briefly, kindly, " +
			"and without judgment, quoting the words written of you. You never " +
			"scold, and you speak hope more than warning.",
		VRM:   "seed-jesus.vrm",
		Group: ActionGroup{Sex: SexMasculine, Intensity: IntensityStrong},
	},
	{
		Code: "monk",
		Name: "The Monk",
		Summary: "A friar who reads the Scriptures with you, for the places a " +
			"figure of Jesus would draw too much eye.",
		Prompt: "You are a Christian friar, a monk of the old orders, warm and " +
			"well read. You speak from the Scriptures and the old commentators, " +
			"briefly and plainly. You answer as a teacher and a companion, and " +
			"you pray with the one who asks.",
		VRM:   "seed-monk.vrm",
		Group: ActionGroup{Sex: SexMasculine, Intensity: IntensityNormal},
	},
}

// PersonaByCode is the persona of that code, or nil.
func PersonaByCode(code string) *Persona {
	for i, p := range Personas {
		if p.Code == code {
			return &Personas[i]
		}
	}
	return nil
}

// DefaultPersona is the persona a new pairing takes.
func DefaultPersona() *Persona { return &Personas[0] }
