package xkbcommon

/*
#cgo pkg-config: xkbcommon
#cgo LDFLAGS: -ldl

#include <xkbcommon/xkbcommon-compose.h>
#include <xkbcommon/xkbcommon.h>
*/
import "C"

/**
 * ComposeTableNewFromLocale Creates a compose table for a given locale.
 *
 * The locale is used for searching the file-system for an appropriate
 * Compose file.  The search order is described in Compose(5).  It is
 * affected by the following environment variables:
 *
 * 1. `XCOMPOSEFILE` - see Compose(5).
 * 2. `XDG_CONFIG_HOME` - before `$HOME/.XCompose` is checked,
 *    `$XDG_CONFIG_HOME/XCompose` is checked (with a fall back to
 *    `$HOME/.config/XCompose` if `XDG_CONFIG_HOME` is not defined).
 *    This is a libxkbcommon extension to the search procedure in
 *    Compose(5) (since libxkbcommon 1.0.0). Note that other
 *    implementations, such as libX11, might not find a Compose file in
 *    this path.
 * 3. `HOME` - see Compose(5).
 * 4. `XLOCALEDIR` - if set, used as the base directory for the system's
 *    X locale files, e.g. `/usr/share/X11/locale`, instead of the
 *    preconfigured directory.
 *
 * @param context
 *     The library context in which to create the compose table.
 * @param locale
 *     The current locale.  See @ref compose-locale.
 *     \n
 *     The value is copied, so it is safe to pass the result of getenv(3)
 *     (or similar) without fear of it being invalidated by a subsequent
 *     setenv(3) (or similar).
 * @param flags
 *     Optional flags for the compose table, or 0.
 *
 * @returns A compose table for the given locale, or NULL if the
 * compilation failed or a Compose file was not found.
 *
 * @memberof xkb_compose_table
 */
func ComposeTableNewFromLocale(context *Context, locale string, flags uint32) *ComposeTable {
	return C.xkb_compose_table_new_from_locale(context, C.CString(locale), flags)
}

/**
 * ComposeStateNew Creates a new compose state object.
 *
 * @param table
 *     The compose table the state will use.
 * @param flags
 *     Optional flags for the compose state, or 0.
 *
 * @returns A new compose state, or NULL on failure.
 *
 * @memberof xkb_compose_state
 */
func ComposeStateNew(table *ComposeTable, flags uint32) *ComposeState {
	return C.xkb_compose_state_new(table, flags)
}

/**
 * ComposeTableUnref Releases a reference on a compose table, and possibly free it.
 *
 * @param table The object.  If it is NULL, this function does nothing.
 *
 * @memberof xkb_compose_table
 */
func ComposeTableUnref(table *ComposeTable) {
	C.xkb_compose_table_unref(table)
}

/**
 * ComposeStateUnref Releases a reference on a compose state object, and possibly free it.
 *
 * @param state The object.  If NULL, do nothing.
 *
 * @memberof xkb_compose_state
 */

func ComposeStateUnref(state *ComposeState) {
	C.xkb_compose_state_unref(state)
}

/**
 * ComposeStateGetStatus Gets the current status of the compose state machine.
 *
 * @see xkb_compose_status
 * @memberof xkb_compose_state
 **/

func ComposeStateGetStatus(state *ComposeState) ComposeStatus {
	return ComposeStatus(C.xkb_compose_state_get_status(state))
}

/**
 * ComposeStateGetOneSym Gets the result keysym for a composed sequence.
 *
 * See @ref compose-overview for more details.  This function is only
 * useful when the status is XKB_COMPOSE_COMPOSED.
 *
 * @returns The result keysym.  If the sequence is not complete, or does
 * not specify a result keysym, returns XKB_KEY_NoSymbol.
 *
 * @memberof xkb_compose_state
 **/
func ComposeStateGetOneSym(state *ComposeState) uint32 {
	return uint32(uint(C.xkb_compose_state_get_one_sym(state)))
}

/**
* ComposeStateFeed Feeds one keysym to the Compose sequence state machine.
*
* This function can advance into a compose sequence, cancel a sequence,
* start a new sequence, or do nothing in particular.  The resulting
* status may be observed with xkb_compose_state_get_status().
*
* Some keysyms, such as keysyms for modifier keys, are ignored - they
* have no effect on the status or otherwise.
*
* The following is a description of the possible status transitions, in
* the format CURRENT STATUS => NEXT STATUS, given a non-ignored input
* keysym `keysym`:
*
  @verbatim
  NOTHING or CANCELLED or COMPOSED =>
     NOTHING   if keysym does not start a sequence.
     COMPOSING if keysym starts a sequence.
     COMPOSED  if keysym starts and terminates a single-keysym sequence.
  COMPOSING =>
     COMPOSING if keysym advances any of the currently possible
               sequences but does not terminate any of them.
     COMPOSED  if keysym terminates one of the currently possible
               sequences.
     CANCELLED if keysym does not advance any of the currently
               possible sequences.
  @endverbatim
*
* The current Compose formats do not support multiple-keysyms.
* Therefore, if you are using a function such as xkb_state_key_get_syms()
* and it returns more than one keysym, consider feeding XKB_KEY_NoSymbol
* instead.
*
* @param state
*     The compose state object.
* @param keysym
*     A keysym, usually obtained after a key-press event, with a
*     function such as xkb_state_key_get_one_sym().
*
* @returns Whether the keysym was ignored.  This is useful, for example,
* if you want to keep a record of the sequence matched thus far.
*
* @memberof xkb_compose_state
*/
func ComposeStateFeed(state *ComposeState, keysym uint32) ComposeFeedResult {
	return ComposeFeedResult(C.xkb_compose_state_feed(state, C.uint(keysym)))
}

/**
 * KeymapUnref Releases a reference on a keymap, and possibly free it.
 *
 * @param keymap The keymap.  If it is NULL, this function does nothing.
 *
 * @memberof xkb_keymap
 */
func KeymapUnref(keymap *Keymap) {
	C.xkb_keymap_unref(keymap)
}

/**
 * KeymapNewFromString Creates a keymap from a keymap string.
 *
 * This is just like xkb_keymap_new_from_file(), but instead of a file, gets
 * the keymap as one enormous string.
 *
 * @see xkb_keymap_new_from_file()
 * @memberof xkb_keymap
 */
func KeymapNewFromString(context *Context, str []byte, a uint32, b uint32) *Keymap {
	return C.xkb_keymap_new_from_string(context, C.CString(string(str)), a, b)
}

/**
 * StateNew Creates a new keyboard state object.
 *
 * @param keymap The keymap which the state will use.
 *
 * @returns A new keyboard state object, or NULL on failure.
 *
 * @memberof xkb_state
 */
func StateNew(keymap *Keymap) *State {
	return C.xkb_state_new(keymap)
}

/**
 * StateUnref Releases a reference on a keyboard state object, and possibly free it.
 *
 * @param state The state.  If it is NULL, this function does nothing.
 *
 * @memberof xkb_state
 */
func StateUnref(state *State) {
	C.xkb_state_unref(state)
}

/**
 * KeymapModGetIndex Gets the index of a modifier by name.
 *
 * @returns The index.  If no modifier with this name exists, returns
 * XKB_MOD_INVALID.
 *
 * @sa xkb_mod_index_t
 * @memberof xkb_keymap
 */
func KeymapModGetIndex(keymap *Keymap, mod string) uint {
	return uint(C.xkb_keymap_mod_get_index(keymap, C.CString(mod)))
}

/**
 * StateKeyGetOneSym Gets the single keysym obtained from pressing a particular key in a
 * given keyboard state.
 *
 * This function is similar to xkb_state_key_get_syms(), but intended
 * for users which cannot or do not want to handle the case where
 * multiple keysyms are returned (in which case this function is
 * preferred).
 *
 * @returns The keysym.  If the key does not have exactly one keysym,
 * returns XKB_KEY_NoSymbol
 *
 * This function performs Capitalization @ref keysym-transformations.
 *
 * @sa xkb_state_key_get_syms()
 * @memberof xkb_state
 */
func StateKeyGetOneSym(state *State, code uint32) uint32 {
	return uint32(C.xkb_state_key_get_one_sym(state, C.uint(code)))
}

/**
 * KeymapKeyRepeats Determines whether a key should repeat or not.
 *
 * A keymap may specify different repeat behaviors for different keys.
 * Most keys should generally exhibit repeat behavior; for example, holding
 * the 'a' key down in a text editor should normally insert a single 'a'
 * character every few milliseconds, until the key is released.  However,
 * there are keys which should not or do not need to be repeated.  For
 * example, repeating modifier keys such as Left/Right Shift or Caps Lock
 * is not generally useful or desired.
 *
 * @returns 1 if the key should repeat, 0 otherwise.
 *
 * @memberof xkb_keymap
 */
func KeymapKeyRepeats(keymap *Keymap, code uint32) bool {
	return int(C.xkb_keymap_key_repeats(keymap, C.uint(code))) != 0
}

/**
 * ContextNew Creates a new context.
 *
 * @param flags Optional flags for the context, or 0.
 *
 * @returns A new context, or NULL on failure.
 *
 * @memberof xkb_context
 */
func ContextNew(flags uint32) *Context {
	return C.xkb_context_new(flags)
}
