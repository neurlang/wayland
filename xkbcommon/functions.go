package xkbcommon

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
func ComposeTableNewFromLocale(*Context, string, int) *ComposeTable {
	return new(ComposeTable)
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
func ComposeStateNew(*ComposeTable, int) *ComposeState {
	return new(ComposeState)
}

/**
 * ComposeTableUnref Releases a reference on a compose table, and possibly free it.
 *
 * @param table The object.  If it is NULL, this function does nothing.
 *
 * @memberof xkb_compose_table
 */
func ComposeTableUnref(*ComposeTable) {

}

/**
 * ComposeStateUnref Releases a reference on a compose state object, and possibly free it.
 *
 * @param state The object.  If NULL, do nothing.
 *
 * @memberof xkb_compose_state
 */

func ComposeStateUnref(*ComposeState) {

}

/**
 * ComposeStateGetStatus Gets the current status of the compose state machine.
 *
 * @see xkb_compose_status
 * @memberof xkb_compose_state
 **/

func ComposeStateGetStatus(*ComposeState) ComposeStatus {
	return 0
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
func ComposeStateGetOneSym(*ComposeState) uint32 {
	return 0
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
func ComposeStateFeed(*ComposeState, uint32) ComposeFeedResult {
	return 0
}

/**
 * KeymapUnref Releases a reference on a keymap, and possibly free it.
 *
 * @param keymap The keymap.  If it is NULL, this function does nothing.
 *
 * @memberof xkb_keymap
 */
func KeymapUnref(*Keymap) {

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
func KeymapNewFromString(*Context, []byte, int, int) *Keymap {
	return new(Keymap)
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
func StateNew(*Keymap) *State {
	return new(State)
}

/**
 * StateUnref Releases a reference on a keybaord state object, and possibly free it.
 *
 * @param state The state.  If it is NULL, this function does nothing.
 *
 * @memberof xkb_state
 */
func StateUnref(*State) {

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
func KeymapModGetIndex(*Keymap, string) uint {
	return 0
}

/**
 * StateKeyGetSyms Gets the keysyms obtained from pressing a particular key in a given
 * keyboard state.
 *
 * Get the keysyms for a key according to the current active layout,
 * modifiers and shift level for the key, as determined by a keyboard
 * state.
 *
 * @param[in]  state    The keyboard state object.
 * @param[in]  key      The keycode of the key.
 * @param[out] syms_out An immutable array of keysyms corresponding the
 * key in the given keyboard state.
 *
 * As an extension to XKB, this function can return more than one keysym.
 * If you do not want to handle this case, you can use
 * xkb_state_key_get_one_sym() for a simpler interface.
 *
 * This function does not perform any @ref keysym-transformations.
 * (This might change).
 *
 * @returns The number of keysyms in the syms_out array.  If no keysyms
 * are produced by the key in the given keyboard state, returns 0 and sets
 * syms_out to NULL.
 *
 * @memberof xkb_state
 */
func StateKeyGetSyms(state *State, code uint32) []uint32 {
	return []uint32{code}
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
	return true
}
