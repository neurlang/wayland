package vulkan

import "github.com/vulkan-go/vulkan"

import "unsafe"

// Wayland-related

type WaylandSurfaceCreateFlags uint32

type WaylandSurfaceCreateInfo struct {
	SType   vulkan.StructureType
	PNext   unsafe.Pointer
	Flags   WaylandSurfaceCreateFlags
	Display uintptr
	Surface uintptr
}
