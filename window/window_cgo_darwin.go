//go:build darwin && cgo
// +build darwin,cgo

package window

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework QuartzCore -framework CoreVideo

#import <Cocoa/Cocoa.h>
#import <QuartzCore/QuartzCore.h>
#import <CoreVideo/CoreVideo.h>
#import <objc/runtime.h>
#import <objc/message.h>

// Forward declarations for Go callbacks
extern void goPopupCorner(void* windowPtr, int x, int y, int w, int h, void* popupWindowPtr);
extern void goMouseMotion(void* windowPtr, float x, float y);
extern void goScheduleRedraw(void* windowPtr);
extern void goWindowResize(void* windowPtr, int width, int height);
extern void goMouseButton(void* windowPtr, int button, int state);
extern void goKeyPress(void* windowPtr, unsigned short keyCode, int state, unsigned long modifiers, unsigned short unicodeChar);
extern void goMouseWheel(void* windowPtr, float dx, float dy, int discrete);

// Simple window wrapper
typedef struct {
    void* nsWindow;
    void* nsView;
    void* goWindowPtr;
    CVDisplayLinkRef displayLink;
    void* eventMonitor; // id<NSObject>
    void* buttonMonitor; // id<NSObject> for button events
    void* keyMonitor; // id<NSObject> for keyboard events
    void* wheelMonitor; // id<NSObject> for mouse wheel events
    CGImageRef currentImage;
    void* imageLock; // NSLock*
    int needsRedraw; // Atomic flag for redraw requests
    void* windowDelegate; // NSWindowDelegate
} DarwinWindow;

// CVDisplayLink callback - called on separate thread
static CVReturn displayLinkCallback(CVDisplayLinkRef displayLink,
                                   const CVTimeStamp* now,
                                   const CVTimeStamp* outputTime,
                                   CVOptionFlags flagsIn,
                                   CVOptionFlags* flagsOut,
                                   void* displayLinkContext) {
    DarwinWindow* dw = (DarwinWindow*)displayLinkContext;

    // Always call redraw on every frame - let the Go side decide if rendering is needed
    void* windowPtr = dw->goWindowPtr;
    dispatch_async(dispatch_get_main_queue(), ^{
        goScheduleRedraw(windowPtr);
    });

    return kCVReturnSuccess;
}

// Store image data in DarwinWindow struct instead of view ivars
// This avoids runtime class issues

// Helper to run block on main thread
static void runOnMainThread(void (^block)(void)) {
    if ([NSThread isMainThread]) {
        block();
    } else {
        dispatch_async(dispatch_get_main_queue(), block);
    }
}

// Window delegate resize handler
static void handleWindowResize(void* darwinWindowPtr, NSNotification* notification) {
    if (!darwinWindowPtr) return;

    DarwinWindow* dw = (DarwinWindow*)darwinWindowPtr;
    if (dw->goWindowPtr) {
        NSWindow* window = [notification object];
        NSRect contentRect = [[window contentView] frame];
        int width = (int)contentRect.size.width;
        int height = (int)contentRect.size.height;

        void* windowPtr = dw->goWindowPtr;
        dispatch_async(dispatch_get_main_queue(), ^{
            goWindowResize(windowPtr, width, height);
        });
    }
}

// Window delegate class - created once per process
static Class getDarwinWindowDelegateClass() {
    static Class delegateClass = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        // Create a new class that inherits from NSObject
        delegateClass = objc_allocateClassPair([NSObject class], "DarwinWindowDelegate", 0);

        // Add ivar for darwinWindow pointer
        class_addIvar(delegateClass, "darwinWindow", sizeof(void*), log2(sizeof(void*)), "^v");

        // Add windowDidResize: method
        IMP windowDidResizeImp = imp_implementationWithBlock(^(id self, NSNotification* notification) {
            void* darwinWindowPtr = NULL;
            object_getInstanceVariable(self, "darwinWindow", &darwinWindowPtr);
            handleWindowResize(darwinWindowPtr, notification);
        });
        class_addMethod(delegateClass, @selector(windowDidResize:), windowDidResizeImp, "v@:@");

        // Add windowShouldClose: method
        IMP windowShouldCloseImp = imp_implementationWithBlock(^BOOL(id self, NSWindow* sender) {
            return YES;
        });
        class_addMethod(delegateClass, @selector(windowShouldClose:), windowShouldCloseImp, "B@:@");

        // Register the class
        objc_registerClassPair(delegateClass);
    });
    return delegateClass;
}

// Create a simple window
static DarwinWindow* darwin_createWindow(int width, int height, const char* title, int decorated, void* goWindowPtr) {
    __block DarwinWindow* dw = NULL;
    __block NSWindow* window = nil;
    __block NSImageView* imageView = nil;

    // Create window directly - called before event loop starts
    @autoreleasepool {
        NSRect frame = NSMakeRect(100, 100, width, height);

        // Set style mask based on decorated flag
        NSWindowStyleMask styleMask;

        if (decorated) {
            styleMask = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable |
                       NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;
        } else {
            // For undecorated windows (popups), use borderless style
            styleMask = NSWindowStyleMaskBorderless;
        }

        window = [[NSWindow alloc]
            initWithContentRect:frame
            styleMask:styleMask
            backing:NSBackingStoreBuffered
            defer:NO];

        [window setTitle:[NSString stringWithUTF8String:title]];
        [window center];
        [window setAcceptsMouseMovedEvents:YES];

        // For decorated windows with transparency, hide title bar elements
        if (decorated) {
            // Normal decorated window - do nothing special
        } else {
            // Borderless window for popups
            [window setBackgroundColor:[NSColor colorWithRed:0.9 green:0.9 blue:0.9 alpha:1.0]];
            [window setOpaque:YES];
            [window setHasShadow:YES];
            [window setLevel:NSPopUpMenuWindowLevel]; // Keep popup above parent
        }

        // Use NSImageView for displaying bitmap content
        imageView = [[NSImageView alloc] initWithFrame:frame];
        [imageView setImageScaling:NSImageScaleAxesIndependently];
        [imageView setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
        [imageView setImageFrameStyle:NSImageFrameNone];
        [imageView setEditable:NO];
        [imageView setAnimates:NO];

        // Set window background
        if (decorated) {
            [window setBackgroundColor:[NSColor blackColor]];
        }
        [window setOpaque:YES];

        [window setContentView:imageView];
        [window makeKeyAndOrderFront:nil];
    }

    // Create DarwinWindow struct
    dw = malloc(sizeof(DarwinWindow));
    dw->nsWindow = (void*)CFBridgingRetain(window);
    dw->nsView = (void*)CFBridgingRetain(imageView);
    dw->goWindowPtr = goWindowPtr;
    dw->displayLink = NULL;
    dw->eventMonitor = NULL;
    dw->buttonMonitor = NULL;
    dw->keyMonitor = NULL;
    dw->currentImage = NULL;
    dw->imageLock = (void*)CFBridgingRetain([[NSLock alloc] init]);
    dw->needsRedraw = 1; // Start with redraw needed

    // Create and set window delegate for resize events
    Class delegateClass = getDarwinWindowDelegateClass();
    id delegate = [[delegateClass alloc] init];
    object_setInstanceVariable(delegate, "darwinWindow", dw);
    [window setDelegate:delegate];
    dw->windowDelegate = (void*)CFBridgingRetain(delegate);

    // Set up mouse tracking for both moved and entered events
    void* windowPtr = goWindowPtr;
    NSEventMask eventMask = NSEventMaskMouseMoved | NSEventMaskMouseEntered | NSEventMaskMouseExited;
    id eventMonitor = [NSEvent addLocalMonitorForEventsMatchingMask:eventMask
        handler:^NSEvent*(NSEvent* event) {
            if ([event window] == window) {
                NSPoint locationInWindow = [event locationInWindow];
                NSPoint locationInView = [imageView convertPoint:locationInWindow fromView:nil];

                // Flip Y coordinate to match Cairo's top-down coordinate system
                float y = [imageView bounds].size.height - locationInView.y;

                goMouseMotion(windowPtr, (float)locationInView.x, y);

            }
            return event;
        }];
    dw->eventMonitor = (void*)CFBridgingRetain(eventMonitor);

    // Set up mouse button tracking
    NSEventMask buttonMask = NSEventMaskLeftMouseDown | NSEventMaskLeftMouseUp |
                             NSEventMaskRightMouseDown | NSEventMaskRightMouseUp;
    id buttonMonitor = [NSEvent addLocalMonitorForEventsMatchingMask:buttonMask
        handler:^NSEvent*(NSEvent* event) {
            if (1) {
                NSPoint locationInWindow = [event locationInWindow];
                NSPoint locationInView = [imageView convertPoint:locationInWindow fromView:nil];

                // Flip Y coordinate to match Cairo's top-down coordinate system
                float y = [imageView bounds].size.height - locationInView.y;

                goMouseMotion(windowPtr, (float)locationInView.x, y);

                int button = 0;
                int state = 0;

                NSEventType eventType = [event type];
                if (eventType == NSEventTypeLeftMouseDown) {
                    button = 272; // Left button (BTN_LEFT in Linux)
                    state = 1;    // Pressed
                } else if (eventType == NSEventTypeLeftMouseUp) {
                    button = 272;
                    state = 0;    // Released
                } else if (eventType == NSEventTypeRightMouseDown) {
                    button = 273; // Right button (BTN_RIGHT in Linux)
                    state = 1;
                } else if (eventType == NSEventTypeRightMouseUp) {
                    button = 273;
                    state = 0;
                }

                goMouseButton(windowPtr, button, state);
            }
            return event;
        }];
    dw->buttonMonitor = (void*)CFBridgingRetain(buttonMonitor);

    // Set up keyboard event tracking
    NSEventMask keyMask = NSEventMaskKeyDown | NSEventMaskKeyUp | NSEventMaskFlagsChanged;
    id keyMonitor = [NSEvent addLocalMonitorForEventsMatchingMask:keyMask
        handler:^NSEvent*(NSEvent* event) {
            if ([event window] == window) {
                unsigned short keyCode = [event keyCode];
                NSEventType eventType = [event type];
                int state = 0;

                if (eventType == NSEventTypeKeyDown) {
                    state = 1; // Pressed
                } else if (eventType == NSEventTypeKeyUp) {
                    state = 0; // Released
                } else if (eventType == NSEventTypeFlagsChanged) {
                    // Handle modifier keys
                    // For now, we'll skip these or handle them separately
                    return event;
                }

                // Get modifier flags
                unsigned long modifiers = [event modifierFlags];

                // Get unicode character if available
                unsigned short unicodeChar = 0;
                NSString* characters = [event characters];
                if (characters && [characters length] > 0) {
                    unicodeChar = [characters characterAtIndex:0];
                }

                goKeyPress(windowPtr, keyCode, state, modifiers, unicodeChar);
            }
            return event;
        }];
    dw->keyMonitor = (void*)CFBridgingRetain(keyMonitor);

    // Set up mouse wheel tracking
    NSEventMask wheelMask = NSEventMaskScrollWheel;
    id wheelMonitor = [NSEvent addLocalMonitorForEventsMatchingMask:wheelMask
        handler:^NSEvent*(NSEvent* event) {
            if ([event window] == window) {
                NSPoint locationInWindow = [event locationInWindow];
                NSPoint locationInView = [imageView convertPoint:locationInWindow fromView:nil];

                // Flip Y coordinate to match Cairo's top-down coordinate system
                float y = [imageView bounds].size.height - locationInView.y;

                goMouseMotion(windowPtr, (float)locationInView.x, y);

                float dx = [event deltaX];
                float dy = [event deltaY];

                int discrete = 0;
                if ([event hasPreciseScrollingDeltas]) {
                    discrete = (int)dy;
                }

                goMouseWheel(windowPtr, dx, dy, discrete);
            }
            return event;
        }];
    dw->wheelMonitor = (void*)CFBridgingRetain(wheelMonitor);

    return dw;
}

// Destroy window
static void darwin_destroyWindow(DarwinWindow* dw) {
    if (dw) {
        // Check if we're on main thread
        BOOL isMainThread = [NSThread isMainThread];

        if (isMainThread) {
            // Already on main thread, execute directly
            @autoreleasepool {
                if (dw->displayLink) {
                    CVDisplayLinkStop(dw->displayLink);
                    CVDisplayLinkRelease(dw->displayLink);
                    dw->displayLink = NULL;
                }
                if (dw->eventMonitor) {
                    id monitor = (__bridge_transfer id)dw->eventMonitor;
                    [NSEvent removeMonitor:monitor];
                    dw->eventMonitor = NULL;
                }
                if (dw->buttonMonitor) {
                    id monitor = (__bridge_transfer id)dw->buttonMonitor;
                    [NSEvent removeMonitor:monitor];
                    dw->buttonMonitor = NULL;
                }
                if (dw->keyMonitor) {
                    id monitor = (__bridge_transfer id)dw->keyMonitor;
                    [NSEvent removeMonitor:monitor];
                    dw->keyMonitor = NULL;
                }
                if (dw->wheelMonitor) {
                    id monitor = (__bridge_transfer id)dw->wheelMonitor;
                    [NSEvent removeMonitor:monitor];
                    dw->wheelMonitor = NULL;
                }
                if (dw->imageLock) {
                    NSLock* lock = (__bridge NSLock*)dw->imageLock;
                    [lock lock];
                    if (dw->currentImage) {
                        CGImageRelease(dw->currentImage);
                        dw->currentImage = NULL;
                    }
                    [lock unlock];
                    CFBridgingRelease(dw->imageLock);
                    dw->imageLock = NULL;
                }
                if (dw->windowDelegate) {
                    id delegate = (__bridge_transfer id)dw->windowDelegate;
                    object_setInstanceVariable(delegate, "darwinWindow", NULL);
                    dw->windowDelegate = NULL;
                }
                if (dw->nsView) {
                    CFBridgingRelease(dw->nsView);
                    dw->nsView = NULL;
                }
                if (dw->nsWindow) {
                    NSWindow* window = (__bridge_transfer NSWindow*)dw->nsWindow;
                    [window setDelegate:nil];
                    [window close];
                }
            }
        } else {
            // Not on main thread, use dispatch_sync
            dispatch_sync(dispatch_get_main_queue(), ^{
                @autoreleasepool {
                    if (dw->displayLink) {
                        CVDisplayLinkStop(dw->displayLink);
                        CVDisplayLinkRelease(dw->displayLink);
                        dw->displayLink = NULL;
                    }
                    if (dw->eventMonitor) {
                        id monitor = (__bridge_transfer id)dw->eventMonitor;
                        [NSEvent removeMonitor:monitor];
                        dw->eventMonitor = NULL;
                    }
                    if (dw->buttonMonitor) {
                        id monitor = (__bridge_transfer id)dw->buttonMonitor;
                        [NSEvent removeMonitor:monitor];
                        dw->buttonMonitor = NULL;
                    }
                    if (dw->keyMonitor) {
                        id monitor = (__bridge_transfer id)dw->keyMonitor;
                        [NSEvent removeMonitor:monitor];
                        dw->keyMonitor = NULL;
                    }
                    if (dw->wheelMonitor) {
                        id monitor = (__bridge_transfer id)dw->wheelMonitor;
                        [NSEvent removeMonitor:monitor];
                        dw->wheelMonitor = NULL;
                    }
                    if (dw->imageLock) {
                        NSLock* lock = (__bridge NSLock*)dw->imageLock;
                        [lock lock];
                        if (dw->currentImage) {
                            CGImageRelease(dw->currentImage);
                            dw->currentImage = NULL;
                        }
                        [lock unlock];
                        CFBridgingRelease(dw->imageLock);
                        dw->imageLock = NULL;
                    }
                    if (dw->windowDelegate) {
                        id delegate = (__bridge_transfer id)dw->windowDelegate;
                        object_setInstanceVariable(delegate, "darwinWindow", NULL);
                        dw->windowDelegate = NULL;
                    }
                    if (dw->nsView) {
                        CFBridgingRelease(dw->nsView);
                        dw->nsView = NULL;
                    }
                    if (dw->nsWindow) {
                        NSWindow* window = (__bridge_transfer NSWindow*)dw->nsWindow;
                        [window setDelegate:nil];
                        [window close];
                    }
                }
            });
        }
        free(dw);
    }
}

// Set window title
static void darwin_setTitle(DarwinWindow* dw, const char* title) {
    if (dw && dw->nsWindow) {
        NSString* titleStr = [NSString stringWithUTF8String:title];
        dispatch_async(dispatch_get_main_queue(), ^{
            @autoreleasepool {
                NSWindow* window = (__bridge NSWindow*)dw->nsWindow;
                [window setTitle:titleStr];
            }
        });
    }
}

// Run main loop
static void darwin_runMainLoop() {
    @autoreleasepool {
        [NSApplication sharedApplication];
        [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
        [NSApp activateIgnoringOtherApps:YES];
        [NSApp run];
    }
}

// Stop main loop
static void darwin_stopMainLoop() {
    @autoreleasepool {
        [NSApp stop:nil];
        NSEvent* event = [NSEvent otherEventWithType:NSEventTypeApplicationDefined
                                            location:NSMakePoint(0, 0)
                                       modifierFlags:0
                                           timestamp:0
                                        windowNumber:0
                                             context:nil
                                             subtype:0
                                               data1:0
                                               data2:0];
        [NSApp postEvent:event atStart:YES];
    }
}

// Get window size
static void darwin_getWindowSize(DarwinWindow* dw, int* width, int* height) {
    if (dw && dw->nsWindow) {
        @autoreleasepool {
            NSWindow* window = (__bridge NSWindow*)dw->nsWindow;
            NSRect frame = [[window contentView] frame];
            *width = (int)frame.size.width;
            *height = (int)frame.size.height;
        }
    }
}

// Set fullscreen
static void darwin_setFullscreen(DarwinWindow* dw, int fullscreen) {
    if (dw && dw->nsWindow) {
        dispatch_async(dispatch_get_main_queue(), ^{
            @autoreleasepool {
                NSWindow* window = (__bridge NSWindow*)dw->nsWindow;
                BOOL isFullscreen = ([window styleMask] & NSWindowStyleMaskFullScreen) != 0;
                if ((fullscreen && !isFullscreen) || (!fullscreen && isFullscreen)) {
                    [window toggleFullScreen:nil];
                }
            }
        });
    }
}

// Resize window
static void darwin_resizeWindow(DarwinWindow* dw, int width, int height) {
    if (dw && dw->nsWindow) {
        dispatch_async(dispatch_get_main_queue(), ^{
            @autoreleasepool {
                NSWindow* window = (__bridge NSWindow*)dw->nsWindow;
                NSRect frame = [window frame];
                NSRect contentRect = NSMakeRect(frame.origin.x, frame.origin.y, width, height);
                NSRect newFrame = [window frameRectForContentRect:contentRect];
                [window setFrame:newFrame display:YES animate:NO];
            }
        });
    }
}

// Start display link for continuous redraw
static void darwin_startDisplayLink(DarwinWindow* dw) {
    if (dw && dw->goWindowPtr && !dw->displayLink) {
        // Create display link
        CVReturn ret = CVDisplayLinkCreateWithActiveCGDisplays(&dw->displayLink);
        if (ret != kCVReturnSuccess) {
            return;
        }

        // Set the callback
        CVDisplayLinkSetOutputCallback(dw->displayLink, &displayLinkCallback, dw);

        // Get the main display ID
        NSWindow* window = (__bridge NSWindow*)dw->nsWindow;
        NSScreen* screen = [window screen];
        NSDictionary* screenDescription = [screen deviceDescription];
        NSNumber* screenNumber = [screenDescription objectForKey:@"NSScreenNumber"];
        CGDirectDisplayID displayID = (CGDirectDisplayID)[screenNumber unsignedIntValue];
        CVDisplayLinkSetCurrentCGDisplay(dw->displayLink, displayID);

        // Start the display link
        CVDisplayLinkStart(dw->displayLink);
    }
}

// Request a redraw on next display link callback
static void darwin_requestRedraw(DarwinWindow* dw) {
    if (dw) {
        // Atomic set
        __sync_fetch_and_or(&dw->needsRedraw, 1);
    }
}

// Enable mouse tracking
static void darwin_enableMouseTracking(DarwinWindow* dw) {
    if (dw && dw->nsWindow && dw->goWindowPtr) {
        @autoreleasepool {
            NSWindow* window = (__bridge NSWindow*)dw->nsWindow;
            id view = [window contentView];

            // Get the actual class of the view instance
            Class viewClass = object_getClass(view);
            Ivar goWindowPtrIvar = class_getInstanceVariable(viewClass, "goWindowPtr");
            if (goWindowPtrIvar) {
                // Set using direct memory access for non-object types
                ptrdiff_t offset = ivar_getOffset(goWindowPtrIvar);
                void** ptr = (void**)((char*)(__bridge void*)view + offset);
                *ptr = dw->goWindowPtr;
            }

            // Call updateTrackingAreas using objc_msgSend
            ((void(*)(id, SEL))objc_msgSend)(view, @selector(updateTrackingAreas));
        }
    }
}

// Get mouse position in window
static void darwin_getMousePosition(DarwinWindow* dw, float* x, float* y) {
    if (dw && dw->nsWindow) {
        @autoreleasepool {
            NSWindow* window = (__bridge NSWindow*)dw->nsWindow;
            NSPoint mouseLocation = [window mouseLocationOutsideOfEventStream];
            *x = mouseLocation.x;
            *y = mouseLocation.y;
        }
    }
}

// Release callback for CGDataProvider
static void releaseDataCallback(void *info, const void *data, size_t size) {
    free((void*)data);
}

// Draw bitmap to window
static void darwin_drawBitmap(DarwinWindow* dw, void* data, int width, int height) {
    if (!dw || !dw->nsView || !data || width <= 0 || height <= 0) {
        return;
    }

    @autoreleasepool {
        NSImageView* imageView = (__bridge NSImageView*)dw->nsView;
        NSLock* lock = (__bridge NSLock*)dw->imageLock;

        // Create CGImage from raw BGRA data (Cairo format: BGRA premultiplied)
        CGColorSpaceRef colorSpace = CGColorSpaceCreateDeviceRGB();

        // Copy the data to ensure it stays valid
        size_t dataSize = width * height * 4;
        void* dataCopy = malloc(dataSize);
        memcpy(dataCopy, data, dataSize);

        CGDataProviderRef provider = CGDataProviderCreateWithData(
            NULL,
            dataCopy,
            dataSize,
            releaseDataCallback
        );

        CGImageRef cgImage = CGImageCreate(
            width,
            height,
            8,                  // bits per component
            32,                 // bits per pixel
            width * 4,          // bytes per row
            colorSpace,
            kCGImageAlphaPremultipliedFirst | kCGBitmapByteOrder32Little,
            provider,
            NULL,
            false,
            kCGRenderingIntentDefault
        );

        if (!cgImage) {
            CGDataProviderRelease(provider);
            CGColorSpaceRelease(colorSpace);
            return;
        }

        // Update the stored image
        [lock lock];
        if (dw->currentImage) {
            CGImageRelease(dw->currentImage);
        }
        dw->currentImage = cgImage;
        CGImageRetain(cgImage);
        [lock unlock];

        // Retain cgImage for the block
        CGImageRetain(cgImage);

        // Create NSImage and display it in the NSImageView - ALL on main thread
        // Use dispatch_sync if we're already on main thread, dispatch_async otherwise
        if ([NSThread isMainThread]) {
            @autoreleasepool {
                NSSize imageSize = NSMakeSize(width, height);
                NSImage* nsImage = [[NSImage alloc] initWithCGImage:cgImage size:imageSize];

                if (nsImage) {
                    [imageView setImage:nsImage];
                    [imageView setNeedsDisplay:YES];
                }

                CGImageRelease(cgImage);
            }
        } else {
            dispatch_sync(dispatch_get_main_queue(), ^{
                @autoreleasepool {
                    NSSize imageSize = NSMakeSize(width, height);
                    NSImage* nsImage = [[NSImage alloc] initWithCGImage:cgImage size:imageSize];

                    if (nsImage) {
                        [imageView setImage:nsImage];
                        [imageView setNeedsDisplay:YES];
                    }

                    CGImageRelease(cgImage);
                }
            });
        }

        CGImageRelease(cgImage);
        CGDataProviderRelease(provider);
        CGColorSpaceRelease(colorSpace);
    }
}

// Position popup window relative to parent window
static void darwin_positionPopup(DarwinWindow* popupDw, DarwinWindow* parentDw, int offsetX, int offsetY) {
    if (!popupDw || !popupDw->nsWindow || !parentDw || !parentDw->nsWindow) {
        printf("[DEBUG C] darwin_positionPopup: NULL pointer check failed\n");
        return;
    }

    printf("[DEBUG C] darwin_positionPopup: offsetX=%d, offsetY=%d\n", offsetX, offsetY);

    dispatch_async(dispatch_get_main_queue(), ^{
        @autoreleasepool {
            NSWindow* popupWindow = (__bridge NSWindow*)popupDw->nsWindow;
            NSWindow* parentWindow = (__bridge NSWindow*)parentDw->nsWindow;

            // Get parent window's frame
            NSRect parentFrame = [parentWindow frame];
            NSRect parentContentRect = [parentWindow contentRectForFrameRect:parentFrame];

            printf("[DEBUG C] Parent window frame: x=%.1f, y=%.1f, w=%.1f, h=%.1f\n",
                   parentContentRect.origin.x, parentContentRect.origin.y,
                   parentContentRect.size.width, parentContentRect.size.height);

            // Get popup window size
            NSRect popupFrame = [popupWindow frame];
            CGFloat popupWidth = popupFrame.size.width;
            CGFloat popupHeight = popupFrame.size.height;

            // Calculate popup position in screen coordinates
            // offsetX and offsetY are in Cairo coordinates (top-left origin, Y down)
            // We need to convert to macOS screen coordinates (bottom-left origin, Y up)

            // Convert Y from top-down to bottom-up
            CGFloat screenX = parentContentRect.origin.x + offsetX - popupWidth/2;
            CGFloat screenY = parentContentRect.origin.y + (parentContentRect.size.height - offsetY) - popupHeight/2;

            printf("[DEBUG C] Popup screen position: x=%.1f, y=%.1f (w=%.1f, h=%.1f)\n",
                   screenX, screenY, popupWidth, popupHeight);

            goPopupCorner(parentDw->goWindowPtr, (int) (offsetX - popupWidth/2) , (int) (offsetY - popupHeight/2), popupWidth, popupHeight, popupDw->goWindowPtr);

            // Set new position
            [popupWindow setFrameOrigin:NSMakePoint(screenX, screenY)];

            // Make popup window a child of parent (so it stays on top)
            [parentWindow addChildWindow:popupWindow ordered:NSWindowAbove];

            printf("[DEBUG C] Popup window positioned and shown\n");

            // Show the popup window
            [popupWindow orderFront:nil];
        }
    });
}
*/
import "C"
import (
	"sync"
	"time"
	"unsafe"

	"github.com/neurlang/wayland/wl"
)

type darwinWindowHandle struct {
	cWindow    *C.DarwinWindow
	goWindow   *Window
	windowID   uintptr
	lastMouseX float32
	lastMouseY float32
}

var (
	windowRegistry = make(map[uintptr]*darwinWindowHandle)
	windowMutex    sync.RWMutex
	nextWindowID   uintptr = 1
)

func darwin_createWindow(width, height int32, title string, decorated bool, goWindow *Window) *darwinWindowHandle {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))

	cDecorated := C.int(0)
	if decorated {
		cDecorated = C.int(1)
	}

	// Allocate a unique ID for this window
	windowMutex.Lock()
	windowID := nextWindowID
	nextWindowID++
	windowMutex.Unlock()

	// Pass the ID as a pointer (safe to pass to C)
	cWindow := C.darwin_createWindow(C.int(width), C.int(height), cTitle, cDecorated, unsafe.Pointer(windowID))

	handle := &darwinWindowHandle{
		cWindow:  cWindow,
		goWindow: goWindow,
		windowID: windowID,
	}

	windowMutex.Lock()
	windowRegistry[windowID] = handle
	windowMutex.Unlock()

	return handle
}

func goPopupGone(windowPtr, popupWindowPtr *Window) {
	window := windowPtr
	if popupWindowPtr == nil {
		return
	}
	if windowPtr == nil {
		return
	}
	popupWindowID := uintptr(unsafe.Pointer(popupWindowPtr.darwinHandle.windowID))

	var newList [][5]uintptr
	var found bool
	for _, dimensions := range window.popupList {
		if dimensions[4] != popupWindowID {
			newList = append(newList, dimensions)
		} else {
			if len(window.popupList) == 1 {
				newList = nil
			}
			found = true
		}
	}
	if found {
		window.popupList = newList
	}
}

//export goPopupCorner
func goPopupCorner(windowPtr unsafe.Pointer, x, y, w, h C.int, popupWindowPtr unsafe.Pointer) {
	windowID := uintptr(windowPtr)

	windowMutex.RLock()
	handle, ok := windowRegistry[windowID]
	windowMutex.RUnlock()

	if !ok || handle == nil || handle.goWindow == nil {
		return
	}
	window := handle.goWindow

	// Debug: log mouse motion
	println("[GO DEBUG] Window rectangle set for popup window", windowID, "at", (x), (y))

	popupWindowID := uintptr(popupWindowPtr)

	window.popupList = append(window.popupList, [5]uintptr{
		uintptr((int(x))),
		uintptr((int(y))),
		uintptr((int(w))),
		uintptr((int(h))),
		uintptr(popupWindowID),
	})

}

//export goMouseMotion
func goMouseMotion(windowPtr unsafe.Pointer, x, y C.float) {
	windowID := uintptr(windowPtr)

	windowMutex.RLock()
	handle, ok := windowRegistry[windowID]
	windowMutex.RUnlock()

	if !ok || handle == nil || handle.goWindow == nil {
		return
	}

	window := handle.goWindow

	for _, dimensions := range window.popupList {
		if int(float32(x)) > int(dimensions[0]) && int(float32(y)) > int(dimensions[1]) &&
			int(float32(x)) < int(dimensions[0]+dimensions[2]) && int(float32(y)) < int(dimensions[1]+dimensions[3]) {

			x := float32(x) - float32(dimensions[0])
			y := float32(y) - float32(dimensions[1])

			// Debug: log mouse motion on popup
			println("[GO DEBUG] Mouse motion on popup", windowID, "at", int(float32(x)), int(float32(y)), (dimensions[2]), (dimensions[3]))

			// recursing the motion
			// XXX: pointer may be bigger than int, but we don't care
			goMouseMotion(unsafe.Pointer(uintptr(dimensions[4])),
				C.float(x),
				C.float(y))
			return
		} else {
			// Debug: log mouse motion on popup
			println("[GO DEBUG] Unhandled mouse motion on popup", windowID, "at", int(float32(x)), int(float32(y)), (dimensions[0]), (dimensions[1]))
		}
	}

	{

		// Debug: log mouse motion
		println("[GO DEBUG] Mouse motion on window", windowID, "at", float32(x), float32(y))

		// Call motion handler on all widgets
		for widget := range window.widgets {
			if widget.handler != nil {
				timestamp := uint32(time.Now().UnixNano() / 1000000)
				widget.handler.Motion(widget, window.input, timestamp, float32(x), float32(y))
			}
		}
	}
}

//export goScheduleRedraw
func goScheduleRedraw(windowPtr unsafe.Pointer) {
	windowID := uintptr(windowPtr)

	windowMutex.RLock()
	handle, ok := windowRegistry[windowID]
	windowMutex.RUnlock()

	if !ok || handle == nil || handle.goWindow == nil {
		return
	}

	window := handle.goWindow

	// Schedule redraw
	window.ScheduleRedraw()
	window.Redraw()
}

//export goWindowResize
func goWindowResize(windowPtr unsafe.Pointer, width, height C.int) {
	windowID := uintptr(windowPtr)

	windowMutex.RLock()
	handle, ok := windowRegistry[windowID]
	windowMutex.RUnlock()

	if !ok || handle == nil || handle.goWindow == nil {
		return
	}

	window := handle.goWindow
	w := int32(width)
	h := int32(height)

	// Update window size
	window.width = w
	window.height = h

	// Update all widgets
	for widget := range window.widgets {
		widget.SetAllocation(0, 0, w, h)
		widget.drawnHash = 0
		widget.drawnHashes = nil
		if widget.handler != nil {
			widget.handler.Resize(widget, w, h, w, h)
		}
	}

	// Request redraw after resize
	window.ScheduleRedraw()
}

//export goMouseButton
func goMouseButton(windowPtr unsafe.Pointer, button, state C.int) {
	windowID := uintptr(windowPtr)

	windowMutex.RLock()
	handle, ok := windowRegistry[windowID]
	windowMutex.RUnlock()

	if !ok || handle == nil || handle.goWindow == nil {
		return
	}

	window := handle.goWindow

	var proceed = len(window.popupList) == 0
	for _, dimensions := range window.popupList {

		if windowID == dimensions[4] {
			{
				// Debug: log mouse button on popup
				println("[GO DEBUG] Mouse button on popup")
				proceed = true

			}
		}
	}
	if !proceed {
		return
	}

	// Call button handler on all widgets
	for widget := range window.widgets {
		if widget.handler != nil {
			timestamp := uint32(time.Now().UnixNano() / 1000000)
			// Convert state: 1 = pressed, 0 = released
			var wlState wl.PointerButtonState
			if state == 1 {
				wlState = wl.PointerButtonStatePressed
			} else {
				wlState = wl.PointerButtonStateReleased
			}
			widget.handler.Button(widget, window.input, timestamp, uint32(button), wlState, widget.handler)
		}
	}
}

//export goKeyPress
func goKeyPress(windowPtr unsafe.Pointer, keyCode C.ushort, state C.int, modifiers C.ulong, unicodeChar C.ushort) {
	windowID := uintptr(windowPtr)

	windowMutex.RLock()
	handle, ok := windowRegistry[windowID]
	windowMutex.RUnlock()

	if !ok || handle == nil || handle.goWindow == nil {
		return
	}

	window := handle.goWindow

	// Update modifier state in input
	if window.input != nil {
		window.input.updateModifiers(uint32(modifiers))
	}

	// Call keyboard handler if set
	if window.input != nil && window.input.keyboardHandler != nil {
		timestamp := uint32(time.Now().UnixNano() / 1000000)

		// Convert state: 1 = pressed, 0 = released
		var wlState wl.KeyboardKeyState
		if state == 1 {
			wlState = wl.KeyboardKeyStatePressed
		} else {
			wlState = wl.KeyboardKeyStateReleased
		}

		// Get the widget handler from the first widget (if any)
		var widgetHandler WidgetHandler
		for widget := range window.widgets {
			if widget.handler != nil {
				widgetHandler = widget.handler
				break
			}
		}

		// Call the keyboard handler
		// vKey = virtual key code (macOS keyCode)
		// code = unicode character
		window.input.keyboardHandler.Key(window, window.input, timestamp, uint32(keyCode), uint32(unicodeChar), wlState, widgetHandler)
	}
}

//export goMouseWheel
func goMouseWheel(windowPtr unsafe.Pointer, dx, dy C.float, discrete C.int) {
	windowID := uintptr(windowPtr)

	windowMutex.RLock()
	handle, ok := windowRegistry[windowID]
	windowMutex.RUnlock()

	if !ok || handle == nil || handle.goWindow == nil {
		return
	}

	window := handle.goWindow

	var proceed = len(window.popupList) == 0
	for _, dimensions := range window.popupList {

		if windowID == dimensions[4] {
			{
				println("[GO DEBUG] Mouse wheel on popup")
				proceed = true
			}
		}
	}
	if !proceed {
		return
	}

	for widget := range window.widgets {
		if widget.handler != nil {
			timestamp := uint32(time.Now().UnixNano() / 1000000)
			widget.handler.AxisDiscrete(widget, window.input, 0, C.int(dy))
		}
	}
}

func darwin_destroyWindow(handle *darwinWindowHandle) {
	if handle != nil && handle.cWindow != nil {
		windowMutex.Lock()
		delete(windowRegistry, handle.windowID)
		windowMutex.Unlock()

		C.darwin_destroyWindow(handle.cWindow)
		handle.cWindow = nil
	}
}

func darwin_setTitle(handle *darwinWindowHandle, title string) {
	if handle != nil && handle.cWindow != nil {
		cTitle := C.CString(title)
		defer C.free(unsafe.Pointer(cTitle))
		C.darwin_setTitle(handle.cWindow, cTitle)
	}
}

func darwin_runMainLoop() {
	C.darwin_runMainLoop()
}

func darwin_stopMainLoop() {
	C.darwin_stopMainLoop()
}

func darwin_getWindowSize(handle *darwinWindowHandle) (int32, int32) {
	if handle != nil && handle.cWindow != nil {
		var width, height C.int
		C.darwin_getWindowSize(handle.cWindow, &width, &height)
		return int32(width), int32(height)
	}
	return 0, 0
}

func darwin_setFullscreen(handle *darwinWindowHandle, fullscreen bool) {
	if handle != nil && handle.cWindow != nil {
		fs := C.int(0)
		if fullscreen {
			fs = 1
		}
		C.darwin_setFullscreen(handle.cWindow, fs)
	}
}

func darwin_resizeWindow(handle *darwinWindowHandle, width, height int32) {
	if handle != nil && handle.cWindow != nil {
		C.darwin_resizeWindow(handle.cWindow, C.int(width), C.int(height))
	}
}

func darwin_startDisplayLink(handle *darwinWindowHandle) {
	if handle != nil && handle.cWindow != nil {
		C.darwin_startDisplayLink(handle.cWindow)
	}
}

func darwin_requestRedraw(handle *darwinWindowHandle) {
	if handle != nil && handle.cWindow != nil {
		C.darwin_requestRedraw(handle.cWindow)
	}
}

func darwin_enableMouseTracking(handle *darwinWindowHandle) {
	if handle != nil && handle.cWindow != nil {
		C.darwin_enableMouseTracking(handle.cWindow)
	}
}

func darwin_drawBitmap(handle *darwinWindowHandle, data []byte, width, height int32) {
	if handle != nil && handle.cWindow != nil && len(data) > 0 {
		C.darwin_drawBitmap(handle.cWindow, unsafe.Pointer(&data[0]), C.int(width), C.int(height))
	}
}

func darwin_positionPopup(popupHandle, parentHandle *darwinWindowHandle, offsetX, offsetY int32) {
	if popupHandle == parentHandle {
		panic("can't be same")
	}
	if popupHandle != nil && popupHandle.cWindow != nil && parentHandle != nil && parentHandle.cWindow != nil {
		C.darwin_positionPopup(popupHandle.cWindow, parentHandle.cWindow, C.int(offsetX), C.int(offsetY))
	}
}
