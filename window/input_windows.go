package window

import (
	xkb "github.com/neurlang/wayland/xkbcommon"
	"github.com/zzl/go-win32api/v2/win32"
	"golang.design/x/clipboard"
	"io"
)

type Input struct {
	keyboardBuf [200]byte
	unicodeBuf  [200]uint16
	last        string
	ctrl        bool
}

func (i *Input) GetRune(notUnicode *uint32, key uint32) rune {
	i.ctrl = false
	switch key {
	case xkb.KeyLeft,
		xkb.KeyRight,
		xkb.KeyUp,
		xkb.KeyDown,
		xkb.KeyHome,
		xkb.KeyEnd,
		xkb.KeyKpEnter,
		xkb.KeyReturn,
		xkb.KeyBackspace,
		xkb.KeyDelete:
		*notUnicode = key
	}
	keyboardBuf := i.keyboardBuf[:]
	_, getKeyboardStateErr := win32.GetKeyboardState(&keyboardBuf[0])
	if getKeyboardStateErr.NilOrError() != nil {
		return 0
	}
	unicodeBuf := i.unicodeBuf[:]
	win32.ToUnicode(key, *notUnicode, &keyboardBuf[0], &unicodeBuf[0], 1, 0)

	i.last = string([]rune{rune(unicodeBuf[0])})
	switch i.last {
	case "\u0016":
		*notUnicode = 'v'
		i.last = ""
		i.ctrl = true
		return 'v'
	case "\u0001":
		*notUnicode = 'a'
		i.last = ""
		i.ctrl = true
		return 'a'
	case "\u0002":
		*notUnicode = 'b'
		i.last = ""
		i.ctrl = true
		return 'b'
	case "\u0003":
		*notUnicode = 'c'
		i.last = ""
		i.ctrl = true
		return 'c'
	}
	return rune(unicodeBuf[0])
}

func (i *Input) GetModifiers() ModType {
	if i.ctrl {
		return ModControlMask
	}
	return 0
}

func init() {
	clipboard.Init()
}

func (i *Input) DeviceSetSelection(src *DataSource, i2 uint32) {
	clipboard.Write(clipboard.FmtText, []byte(src.CopyBuffer))
}

func (i *Input) ReceiveSelectionData(s string, p io.WriteCloser) error {
	go func() {

		data := clipboard.Read(clipboard.FmtText)

		p.Write(data)
		p.Close()
	}()

	return nil
}

func (i *Input) GetUtf8() string {
	return i.last
}
