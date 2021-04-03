package xkbcommon

const KeymapFormatTextV1 = 1
const ComposeFormatTextV1 = 1

const ModShiftMask uint8 = 0x01
const ModAltMask uint8 = 0x02
const ModControlMask uint8 = 0x04

const ContextNoFlags = 0
const ComposeCompileNoFlags = 0
const ComposeStateNoFlags = 0

/** Status of the Compose sequence state machine. */
type ComposeStatus uint8

const (
	/** The initial state; no sequence has started yet. */
	ComposeNothing ComposeStatus = 0
	/** In the middle of a sequence. */
	ComposeComposing ComposeStatus = 1
	/** A complete sequence has been matched. */
	ComposeComposed ComposeStatus = 2
	/** The last sequence was cancelled due to an unmatched keysym. */
	ComposeCancelled ComposeStatus = 3
)

/** The effect of a keysym fed to xkb_compose_state_feed(). */
type ComposeFeedResult uint8

const (
	/** The keysym had no effect - it did not affect the status. */
	ComposeFeedIgnored ComposeFeedResult = 0
	/** The keysym started, advanced or cancelled a sequence. */
	ComposeFeedAccepted ComposeFeedResult = 1
)
