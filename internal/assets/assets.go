// Package assets is what the app itself ships for the device's
// character: the default figure both personas wear until the reader
// gives them a model of their own. It is pixiv's MIT-licensed VRM sample
// humanoid (three-vrm), a rigged figure the device's motions fit.
package assets

import _ "embed"

//go:embed default-figure.vrm
var DefaultFigure []byte
