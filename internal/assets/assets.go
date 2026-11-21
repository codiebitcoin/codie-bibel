// Package assets is what the app itself ships for the device's
// character: the figures and the voice its personas wear until the
// reader gives them their own. Jesus wears the app's own figure of him
// — a VRM from VRoid Hub, brought through the workshop's resource
// service — and speaks a voice cloned from the blessing the app also
// ships; the monk wears a vested cleric from VRoid Hub through the same
// service. The previews are rendered from these very models. The
// pixiv MIT sample humanoid (three-vrm) stays as the figure a persona
// wears should its own model be missing.
package assets

import _ "embed"

//go:embed jesus.vrm
var JesusFigure []byte

//go:embed monk.vrm
var MonkFigure []byte

//go:embed jesus-preview.png
var JesusPreview []byte

//go:embed monk-preview.png
var MonkPreview []byte

//go:embed default-figure.vrm
var DefaultFigure []byte

//go:embed jesus-voice.wav
var SeedVoice []byte
