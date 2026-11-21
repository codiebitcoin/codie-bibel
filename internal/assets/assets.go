// Package assets is what the app itself ships for the device's
// character: the figures and the voice its personas wear until the
// reader gives them their own. Jesus wears the app's own figure of him
// — a VRM from VRoid Hub, brought through the workshop's resource
// service — and speaks a voice cloned from the blessing the app also
// ships; the monk wears pixiv's MIT-licensed VRM sample humanoid
// (three-vrm), a rigged figure the device's motions fit.
package assets

import _ "embed"

//go:embed jesus.vrm
var JesusFigure []byte

//go:embed jesus-preview.png
var JesusPreview []byte

//go:embed default-figure.vrm
var DefaultFigure []byte

//go:embed default-preview.png
var DefaultPreview []byte

//go:embed jesus-voice.wav
var SeedVoice []byte
