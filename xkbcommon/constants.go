package xkbcommon

const KEYMAP_FORMAT_TEXT_V1 = 1
const COMPOSE_FORMAT_TEXT_V1 = 1

const MOD_SHIFT_MASK uint8 = 0x01
const MOD_ALT_MASK uint8 = 0x02
const MOD_CONTROL_MASK uint8 = 0x04

const CONTEXT_NO_FLAGS = 0
const COMPOSE_COMPILE_NO_FLAGS = 0
const COMPOSE_STATE_NO_FLAGS = 0

/** Status of the Compose sequence state machine. */
type ComposeStatus uint8

const (
	/** The initial state; no sequence has started yet. */
	COMPOSE_NOTHING ComposeStatus = 0
	/** In the middle of a sequence. */
	COMPOSE_COMPOSING ComposeStatus = 1
	/** A complete sequence has been matched. */
	COMPOSE_COMPOSED ComposeStatus = 2
	/** The last sequence was cancelled due to an unmatched keysym. */
	COMPOSE_CANCELLED ComposeStatus = 3
)

/** The effect of a keysym fed to xkb_compose_state_feed(). */
type ComposeFeedResult uint8

const (
	/** The keysym had no effect - it did not affect the status. */
	COMPOSE_FEED_IGNORED ComposeFeedResult = 0
	/** The keysym started, advanced or cancelled a sequence. */
	COMPOSE_FEED_ACCEPTED ComposeFeedResult = 1
)
