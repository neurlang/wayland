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
extern void goMouseMotion(void* windowPtr, float x, float y);
extern void goScheduleRedraw(void* windowPtr);
extern void goWindowResize(void* windowPtr, int width, int height);
extern void goMouseButton(void* windowPtr, int button, int state, float x, float y);

// Simple window wrapper
typedef struct {
    void* nsWindow;
    void* nsView;
    void* goWindowPtr;
    CVDisplayLinkRef displayLink;
    void* eventMonitor; // id<NSObject>
    void* buttonMonitor; // id<NSObject> for button events
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
        NSWindowStyleMask styleMask = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | 
                                      NSWindowStyleMaskMiniaturizable | NSWindowStyleMaskResizable;
        
        window = [[NSWindow alloc]
            initWithContentRect:frame
            styleMask:styleMask
            backing:NSBackingStoreBuffered
            defer:NO];
        
        [window setTitle:[NSString stringWithUTF8String:title]];
        [window center];
        [window setAcceptsMouseMovedEvents:YES];
        
        // For undecorated windows, hide the title bar
        if (!decorated) {
            [window setTitlebarAppearsTransparent:YES];
            [window setTitleVisibility:NSWindowTitleHidden];
            window.styleMask |= NSWindowStyleMaskFullSizeContentView;
            
            // Move the window buttons off-screen instead of hiding them
            NSButton *closeButton = [window standardWindowButton:NSWindowCloseButton];
            NSButton *miniaturizeButton = [window standardWindowButton:NSWindowMiniaturizeButton];
            NSButton *zoomButton = [window standardWindowButton:NSWindowZoomButton];
            
            [closeButton setFrameOrigin:NSMakePoint(-100, -100)];
            [miniaturizeButton setFrameOrigin:NSMakePoint(-100, -100)];
            [zoomButton setFrameOrigin:NSMakePoint(-100, -100)];
        }
        
        // Use NSImageView for displaying bitmap content
        imageView = [[NSImageView alloc] initWithFrame:frame];
        [imageView setImageScaling:NSImageScaleAxesIndependently];
        [imageView setAutoresizingMask:NSViewWidthSizable | NSViewHeightSizable];
        [imageView setImageFrameStyle:NSImageFrameNone];
        [imageView setEditable:NO];
        [imageView setAnimates:NO];
        
        // Set window background
        [window setBackgroundColor:[NSColor blackColor]];
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
            if ([event window] == window) {
                NSPoint locationInWindow = [event locationInWindow];
                NSPoint locationInView = [imageView convertPoint:locationInWindow fromView:nil];
                
                // Flip Y coordinate to match Cairo's top-down coordinate system
                float y = [imageView bounds].size.height - locationInView.y;
                
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
                
                goMouseButton(windowPtr, button, state, (float)locationInView.x, y);
            }
            return event;
        }];
    dw->buttonMonitor = (void*)CFBridgingRetain(buttonMonitor);
    
    return dw;
}

// Destroy window
static void darwin_destroyWindow(DarwinWindow* dw) {
    if (dw) {
        // Must be on main thread for NSWindow operations
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
	
	// Call motion handler on all widgets
	for widget := range window.widgets {
		if widget.handler != nil {
			timestamp := uint32(time.Now().UnixNano() / 1000000)
			widget.handler.Motion(widget, window.input, timestamp, float32(x), float32(y))
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
func goMouseButton(windowPtr unsafe.Pointer, button, state C.int, x, y C.float) {
	windowID := uintptr(windowPtr)
	
	windowMutex.RLock()
	handle, ok := windowRegistry[windowID]
	windowMutex.RUnlock()
	
	if !ok || handle == nil || handle.goWindow == nil {
		return
	}
	
	window := handle.goWindow
	
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
