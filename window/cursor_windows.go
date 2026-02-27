package window

import "github.com/neurlang/winc/w32"

const (
	CursorBottomLeft   = w32.IDC_SIZENESW // ↙↗
	CursorBottomRight  = w32.IDC_SIZENWSE // ↘↖
	CursorBottom       = w32.IDC_SIZENS   // ↓↑
	CursorDragging     = w32.IDC_SIZEALL  // move (4 arrows)
	CursorLeftPtr      = w32.IDC_ARROW    // default arrow
	CursorLeft         = w32.IDC_SIZEWE   // ←→
	CursorRight        = w32.IDC_SIZEWE   // ←→
	CursorTopLeft      = w32.IDC_SIZENWSE // ↖↘
	CursorTopRight     = w32.IDC_SIZENESW // ↗↙
	CursorTop          = w32.IDC_SIZENS   // ↑↓
	CursorIbeam        = w32.IDC_IBEAM    // text caret
	CursorHand1        = w32.IDC_HAND     // hand pointer
	CursorWatch        = w32.IDC_WAIT     // hourglass / busy
	CursorDndMove      = w32.IDC_SIZEALL  // reasonable move fallback
	CursorDndCopy      = w32.IDC_ARROW    // no native copy cursor in Win32
	CursorDndForbidden = w32.IDC_NO       // forbidden
	CursorBlank        = 0                // requires custom blank cursor if truly invisible
)
