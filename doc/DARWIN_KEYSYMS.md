# macOS (Darwin) Key Symbols

## Overview

The `keysyms_darwin.go` file provides macOS virtual key code mappings for the xkbcommon package. Unlike X11 keysyms (Linux) or Windows virtual key codes, macOS uses its own set of virtual key codes defined in `HIToolbox/Events.h`.

## Key Differences from Other Platforms

### Linux (X11 KeySyms)
- Uses X11 keysym values (e.g., `KeyQ = 0x0051`, `KEYq = 0x0071`)
- Uppercase and lowercase have different values
- Based on Unicode code points

### Windows (Virtual Key Codes)
- Uses Windows VK codes (e.g., `KeyQ = 81`)
- Uppercase and lowercase share the same code
- Based on ASCII values for letters

### macOS (Virtual Key Codes)
- Uses macOS virtual key codes (e.g., `KeyQ = 0x0C`)
- Uppercase and lowercase share the same code
- Based on physical keyboard layout (ANSI)
- Independent of keyboard layout or language

## Important Notes

### Case Sensitivity
On macOS, **uppercase and lowercase letters use the same virtual key code**:
```go
const KeyQ = 0x0C  // Same code
const KEYq = 0x0C  // Same code
```

The actual character produced depends on:
- Shift key state
- Caps Lock state
- Keyboard layout

### Physical Layout
macOS virtual key codes represent **physical key positions**, not characters:
- `0x0C` is always the Q key on QWERTY
- On AZERTY keyboards, `0x0C` produces 'A'
- On Dvorak keyboards, `0x0C` produces "'"

## Key Code Reference

### Letters (QWERTY Layout)
```
Row 1: Q(0x0C) W(0x0D) E(0x0E) R(0x0F) T(0x11) Y(0x10) U(0x20) I(0x22) O(0x1F) P(0x23)
Row 2: A(0x00) S(0x01) D(0x02) F(0x03) G(0x05) H(0x04) J(0x26) K(0x28) L(0x25)
Row 3: Z(0x06) X(0x07) C(0x08) V(0x09) B(0x0B) N(0x2D) M(0x2E)
```

### Navigation
```
Arrows:  Left(0x7B) Right(0x7C) Up(0x7E) Down(0x7D)
Home/End: Home(0x73) End(0x77)
Pages:   PageUp(0x74) PageDown(0x79)
```

### Modifiers
```
Shift:   Left(0x38) Right(0x3C)
Control: Left(0x3B) Right(0x3E)
Option:  Left(0x3A) Right(0x3D)
Command: Left(0x37) Right(0x36)
```

### Function Keys
```
F1-F12: 0x7A, 0x78, 0x63, 0x76, 0x60, 0x61, 0x62, 0x64, 0x65, 0x6D, 0x67, 0x6F
```

## Usage Example

```go
import "github.com/neurlang/wayland/xkbcommon"

func handleKey(keyCode uint32) {
    switch keyCode {
    case xkbcommon.KeyQ, xkbcommon.KEYq:
        fmt.Println("Q key pressed")
    case xkbcommon.KeyEscape:
        fmt.Println("Escape pressed")
    case xkbcommon.KeyReturn:
        fmt.Println("Return pressed")
    }
}
```

## Testing Key Codes

To verify key codes on macOS:

```objective-c
- (void)keyDown:(NSEvent *)event {
    unsigned short keyCode = [event keyCode];
    NSLog(@"Key code: 0x%02X (%d)", keyCode, keyCode);
}
```

## References

- [Apple Developer: NSEvent](https://developer.apple.com/documentation/appkit/nsevent)
- [HIToolbox Events.h](https://github.com/phracker/MacOSX-SDKs/blob/master/MacOSX10.13.sdk/System/Library/Frameworks/Carbon.framework/Versions/A/Frameworks/HIToolbox.framework/Versions/A/Headers/Events.h)
- [Virtual Key Codes](https://eastmanreference.com/complete-list-of-applescript-key-codes)

## Implementation Status

✅ Required keys (KeyQ, KEYq)
✅ Navigation keys (arrows, home, end, page up/down)
✅ All letter keys (A-Z, both cases)
✅ Number keys (0-9)
✅ Function keys (F1-F12)
✅ Modifier keys (shift, control, option, command)
✅ Special keys (return, backspace, delete, escape, tab, space)
✅ Keypad keys (0-9, operators)
✅ Punctuation keys

## Future Enhancements

- [ ] International keyboard layouts
- [ ] Media keys (play, pause, volume)
- [ ] Additional special keys
- [ ] Key code to character mapping utilities
