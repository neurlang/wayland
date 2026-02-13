# Window Decoration Interaction Implementation

## Overview
This document describes the implementation of interactive window decorations for Linux/Wayland, including window dragging and button actions.

## Features Implemented

### 1. Window Dragging
- Users can click and drag the titlebar to move the window
- Uses `xdg_toplevel.Move()` protocol method
- Requires seat and serial from pointer enter event

### 2. Button Actions

#### Close Button
- Calls the window's `closeHandler.Close()` if set
- Located at the rightmost position in the titlebar

#### Maximize Button
- Toggles between maximized and normal window states
- Uses `xdg_toplevel.SetMaximized()` and `UnsetMaximized()`
- Button icon changes to show current state (single square vs. overlapping squares)
- Located in the middle of the three buttons

#### Minimize Button
- Minimizes the window using `xdg_toplevel.SetMinimized()`
- Located at the leftmost position of the three buttons

### 3. Visual Feedback
- Buttons change color on hover
- Active state: darker button background with lighter icon
- Inactive state: same color as titlebar with dimmed icon
- Hover detection updates in real-time as pointer moves

## Implementation Details

### Pointer Event Routing
Input events are routed through the window's input handlers:

1. **PointerEnter**: Detects when pointer enters titlebar or shadow surface
   - Titlebar: Initializes hover tracking
   - Shadow: Ignores input (decorative only)

2. **PointerLeave**: Clears hover state when pointer leaves

3. **PointerMotion**: Updates hover button based on pointer position
   - Calculates which button (if any) is under the pointer
   - Updates visual state accordingly

4. **PointerButton**: Handles click actions
   - Left click (button 272) triggers actions
   - Routes to appropriate handler based on hover button

### Surface Registration
- Decoration surfaces are registered in `Display.surface2window` map
- Allows input system to route events to correct window
- Surfaces are unregistered when destroyed

### Coordinate System
- Pointer coordinates are relative to the surface they're over
- Titlebar surface: (0,0) is top-left of titlebar
- Button positions calculated from right edge of titlebar

### Button Layout
```
[Title Text]                    [Min] [Max] [Close]
                                 ^     ^     ^
                                 |     |     |
                          width-3*32  -2*32  -32
```

## Code Structure

### Files Modified
- `window/decoration_linux.go`: Core decoration and interaction logic
- `window/window_linux.go`: Input event routing

### Key Methods in WindowDecoration

#### Input Handling
- `HandlePointerEnter(serial, x, y)`: Track pointer entering titlebar
- `HandlePointerLeave()`: Clear hover state
- `HandlePointerMotion(x, y)`: Update hover button
- `HandlePointerButton(serial, button, state)`: Handle clicks

#### Action Handlers
- `handleDragStart(serial)`: Initiate window move
- `handleClose()`: Close window
- `handleMaximize()`: Toggle maximize state
- `handleMinimize()`: Minimize window

#### Helper Methods
- `updateHoverButton()`: Determine which button is under pointer
- `SetHoverButton(btn)`: Update hover state and redraw

## Testing
Build and run any of the example applications:
```bash
cd go-wayland-smoke
go build -o smoke main.go
./smoke
```

Try the following interactions:
1. Click and drag the titlebar to move the window
2. Click the minimize button (leftmost) to minimize
3. Click the maximize button (middle) to toggle maximize
4. Click the close button (rightmost) to close
5. Hover over buttons to see visual feedback

## Future Enhancements
- Cursor changes on hover (hand cursor over buttons)
- Double-click titlebar to maximize
- Right-click titlebar for window menu
- Resize handles on window edges
- Snap-to-edge functionality during drag
