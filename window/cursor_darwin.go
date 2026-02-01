package window

// Cursor type constants for macOS
// These map to NSCursor types in Cocoa

const CursorHand1 = 1        // NSCursor.pointingHandCursor
const CursorLeftPtr = 2      // NSCursor.arrowCursor (default)
const CursorIbeam = 3        // NSCursor.IBeamCursor (text selection)
const CursorBottomLeft = 4   // Resize bottom-left
const CursorBottomRight = 5  // Resize bottom-right
const CursorBottom = 6       // Resize bottom
const CursorDragging = 7     // NSCursor.openHandCursor
const CursorLeft = 8         // Resize left
const CursorRight = 9        // Resize right
const CursorTopLeft = 10     // Resize top-left
const CursorTopRight = 11    // Resize top-right
const CursorTop = 12         // Resize top
const CursorWatch = 13       // NSCursor.operationNotAllowedCursor
const CursorDndMove = 14     // Drag and drop move
const CursorDndCopy = 15     // Drag and drop copy
const CursorDndForbidden = 16 // NSCursor.operationNotAllowedCursor
const CursorBlank = 17       // Hidden cursor

// Additional macOS-specific cursors
const CursorCrosshair = 18   // NSCursor.crosshairCursor
const CursorClosedHand = 19  // NSCursor.closedHandCursor
const CursorDisappearingItem = 20 // NSCursor.disappearingItemCursor
const CursorResizeLeftRight = 21  // NSCursor.resizeLeftRightCursor
const CursorResizeUpDown = 22     // NSCursor.resizeUpDownCursor
