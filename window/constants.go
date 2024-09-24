package window

// ModType is the type of Mod Masks
type ModType uint8

// ModShiftMask is the Shift modifier mask - provided for convenience only
const ModShiftMask ModType = 0x01

// ModAltMask is the Alt modifier mask - provided for convenience only
const ModAltMask ModType = 0x02

// ModControlMask is the Control modifier mask - provided for convenience only
const ModControlMask ModType = 0x04

const BufferTypeEglWindow = 0

const SurfaceOpaque = 0x01
const SurfaceShm = 0x02

const SurfaceHintResize = 0x10
const SurfaceHintRgb565 = 0x100

const PreferredFormatNone = 0
const PreferredFormatRgb565 = 1



const ZwpRelativePointerManagerV1Version = 1
const ZwpPointerConstraintsV1Version = 1
