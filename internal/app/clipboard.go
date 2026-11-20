package app

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/transfer"
)

// copyText puts text on the clipboard; a test puts a recorder here.
var copyText = func(text string) error { return mygo.Clipboard.Write(transfer.TextData(text)) }
