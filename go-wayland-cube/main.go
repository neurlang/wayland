package main

import (
	"fmt"
	wayland "github.com/neurlang/wayland/libwayland"
	vulkan "github.com/neurlang/wayland/vulkan"
	"golang.org/x/sys/unix"
	"os"
	"syscall"
	"unsafe"
)

const WL_PROXY_FLAG_ID_DELETED = 1

const MAX_NUM_IMAGES = 4

type Wland struct {
	wait_for_configure bool

	ndisplay      *wayland.Display
	nseat         *wayland.Seat
	nkeyboard     *wayland.Keyboard
	nshell        *wayland.XdgWmBase
	ncompositor   *wayland.Compositor
	nsurface      *wayland.Surface
	nxdg_toplevel *wayland.XdgToplevel
	nxdg_surface  *wayland.XdgSurface
	nregistry     *wayland.Registry
}

func (w *Wland) RegistryGlobal(wl_registry *wayland.Registry, name uint32, iface string, version uint32) {

	switch iface {
	case "wl_compositor":
		w.ncompositor = wayland.RegistryBindCompositorInterface(w.nregistry, name, 1)

	case "xdg_wm_base":
		w.nshell = wayland.RegistryBindXdgWmBaseInterface(w.nregistry, name, 1)

	case "wl_seat":

		w.nseat = wayland.RegistryBindSeatInterface(w.nregistry, name, 1)

	}

}

func (w *Wland) RegistryGlobalRemove(*wayland.Registry, uint32) {
}

func (w *Wland) XdgSurfaceConfigure(*wayland.XdgSurface, uint32) {
}
func (w *Wland) XdgToplevelClose(*wayland.XdgToplevel) {
	os.Exit(0)
}

func (w *Wland) XdgToplevelConfigure(*wayland.XdgToplevel, int32, int32, []int32) {
}

var wland Wland

type VkCube struct {
	width  int
	height int

	wl *Wland

	device     vulkan.Device
	swap_chain [1]vulkan.Swapchain
	queue      vulkan.Queue
	semaphore  vulkan.Semaphore

	buffers [MAX_NUM_IMAGES]VkCubeBuffer

	mapping *ubo

	buffer          vulkan.Buffer
	render_pass     vulkan.RenderPass
	vertex_offset   vulkan.DeviceSize
	colors_offset   vulkan.DeviceSize
	normals_offset  vulkan.DeviceSize
	pipeline        [1]vulkan.Pipeline
	pipeline_layout vulkan.PipelineLayout
	descriptor_set  vulkan.DescriptorSet

	protected vulkan.Bool32

	start float64

	physical_device vulkan.PhysicalDevice
	surface         vulkan.Surface
	image_count     uint32
	image_format    vulkan.Format
	cmd_pool        vulkan.CommandPool

	mem vulkan.DeviceMemory

	memory_properties vulkan.PhysicalDeviceMemoryProperties
	instance          vulkan.Instance

	protected_chain bool
}

func (vc *VkCube) Render(buf *VkCubeBuffer, wait_semaphore uint8) {
	render_cube(vc, buf, wait_semaphore)
}
func (vc *VkCube) Init() {
	init_cube(vc)
}

type VkCubeBuffer struct {
	mem         vulkan.DeviceMemory
	image       vulkan.Image
	view        [1]vulkan.ImageView
	framebuffer vulkan.Framebuffer
	fence       vulkan.Fence
	cmd_buffer  [1]vulkan.CommandBuffer

	fb     uint32
	stride uint32
}

func init_vk_objects(vc *VkCube) {

	var pass = vc.render_pass

	vulkan.CreateRenderPass(vc.device,
		&vulkan.RenderPassCreateInfo{
			SType:           vulkan.StructureTypeRenderPassCreateInfo,
			AttachmentCount: 1,
			PAttachments: []vulkan.AttachmentDescription{
				{
					Format:        vc.image_format,
					Samples:       1,
					LoadOp:        vulkan.AttachmentLoadOpClear,
					StoreOp:       vulkan.AttachmentStoreOpStore,
					InitialLayout: vulkan.ImageLayoutUndefined,
					FinalLayout:   vulkan.ImageLayoutPresentSrc,
				},
			},
			SubpassCount: 1,
			PSubpasses: []vulkan.SubpassDescription{
				{
					PipelineBindPoint:    vulkan.PipelineBindPointGraphics,
					InputAttachmentCount: 0,
					ColorAttachmentCount: 1,
					PColorAttachments: []vulkan.AttachmentReference{
						{
							Attachment: 0,
							Layout:     vulkan.ImageLayoutColorAttachmentOptimal,
						},
					},
					PResolveAttachments: []vulkan.AttachmentReference{
						{
							Attachment: vulkan.AttachmentUnused,
							Layout:     vulkan.ImageLayoutColorAttachmentOptimal,
						},
					},
					PDepthStencilAttachment: nil,
					PreserveAttachmentCount: 0,
					PPreserveAttachments:    nil,
				},
			},
			DependencyCount: 0,
		},
		nil,
		&pass)

	vc.render_pass = pass

	vc.Init()

	var flags vulkan.CommandPoolCreateFlags
	if 0 != uint32(vc.protected) {
		flags = vulkan.CommandPoolCreateFlags(vulkan.CommandPoolCreateProtectedBit)
	}

	var cmdpool = vc.cmd_pool

	vulkan.CreateCommandPool(vc.device,
		&vulkan.CommandPoolCreateInfo{
			SType:            vulkan.StructureTypeCommandPoolCreateInfo,
			QueueFamilyIndex: 0,
			Flags:            vulkan.CommandPoolCreateFlags(vulkan.CommandPoolCreateResetCommandBufferBit) | flags,
		},
		nil,
		&cmdpool)
	vc.cmd_pool = cmdpool

	var sema = vc.semaphore

	vulkan.CreateSemaphore(vc.device,
		&vulkan.SemaphoreCreateInfo{
			SType: vulkan.StructureTypeSemaphoreCreateInfo,
		},
		nil,
		&sema)

	vc.semaphore = sema
}

func choose_surface_format(vc *VkCube) vulkan.Format {
	var num_formats uint32 = 0

	vulkan.GetPhysicalDeviceSurfaceFormats(vc.physical_device, vc.surface,
		&num_formats, nil)
	if !(num_formats > 0) {
		panic("assert")
	}

	var formats = make([]vulkan.SurfaceFormat, num_formats)

	vulkan.GetPhysicalDeviceSurfaceFormats(vc.physical_device, vc.surface,
		&num_formats, formats)

	var format vulkan.Format = vulkan.FormatUndefined

	for i := uint32(0); i < num_formats; i++ {
		switch formats[i].Format {
		case vulkan.FormatR8g8b8a8Srgb:
			fallthrough
		case vulkan.FormatB8g8r8a8Srgb:
			/* These formats are all fine */
			format = formats[i].Format
			continue
		case vulkan.FormatR8g8b8Srgb:
			fallthrough
		case vulkan.FormatB8g8r8Srgb:
			fallthrough
		case vulkan.FormatR5g6b5UnormPack16:
			fallthrough
		case vulkan.FormatB5g6r5UnormPack16:
			fallthrough
			/* We would like to support these but they don't seem to work. */
		default:
			continue
		}
	}

	if !(format != vulkan.FormatUndefined) {
		panic("assert")
	}

	return format
}

func mainloop(vc *VkCube) {

	fds := []unix.PollFd{{Fd: int32(wayland.DisplayGetFd(vc.wl.ndisplay)), Events: unix.POLLIN}}

	for {
		var index uint32

		for wayland.DisplayPrepareRead(vc.wl.ndisplay) != 0 {
			wayland.DisplayDispatchPending(vc.wl.ndisplay)
		}
		n, err := wayland.DisplayFlush(vc.wl.ndisplay)
		if errno, ok := err.(syscall.Errno); n < 0 && ok {
			if int(errno) != wayland.ErrAgain {
				wayland.DisplayCancelRead(vc.wl.ndisplay)
				return
			}
		}
		n, err = unix.Poll(fds, 0)
		if err != nil && err != syscall.EINTR {
			panic(err)
		}
		if err == nil && n > 0 {
			wayland.DisplayReadEvents(vc.wl.ndisplay)
			wayland.DisplayDispatchPending(vc.wl.ndisplay)
		} else {
			wayland.DisplayCancelRead(vc.wl.ndisplay)
		}

		var result = [1]vulkan.Result{vulkan.AcquireNextImage(vc.device, vc.swap_chain[0], 60,
			vc.semaphore, nil, &index)}
		if result[0] != vulkan.Success {
			return
		}

		if !(index <= MAX_NUM_IMAGES) {
			panic("assert")
		}

		vc.Render(&vc.buffers[index], 1)

		vulkan.QueuePresent(vc.queue,
			&vulkan.PresentInfo{
				SType:          vulkan.StructureTypePresentInfo,
				SwapchainCount: 1,
				PSwapchains:    vc.swap_chain[:],
				PImageIndices:  []uint32{index},
				PResults:       result[:],
			})
		if result[0] != vulkan.Success {
			return
		}

		vulkan.QueueWaitIdle(vc.queue)
	}

}

func init_buffer(vc *VkCube, b *VkCubeBuffer) {

	var iview = b.view[0]

	vulkan.CreateImageView(vc.device,
		&vulkan.ImageViewCreateInfo{
			SType:    vulkan.StructureTypeImageViewCreateInfo,
			Image:    b.image,
			ViewType: vulkan.ImageViewType2d,
			Format:   vc.image_format,
			Components: vulkan.ComponentMapping{
				R: vulkan.ComponentSwizzleR,
				G: vulkan.ComponentSwizzleG,
				B: vulkan.ComponentSwizzleB,
				A: vulkan.ComponentSwizzleA,
			},
			SubresourceRange: vulkan.ImageSubresourceRange{
				AspectMask:     vulkan.ImageAspectFlags(vulkan.ImageAspectColorBit),
				BaseMipLevel:   0,
				LevelCount:     1,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
		},
		nil,
		&iview)
	b.view[0] = iview

	var fbuffer = b.framebuffer

	vulkan.CreateFramebuffer(vc.device,
		&vulkan.FramebufferCreateInfo{
			SType:           vulkan.StructureTypeFramebufferCreateInfo,
			RenderPass:      vc.render_pass,
			AttachmentCount: 1,
			PAttachments:    b.view[:],
			Width:           uint32(vc.width),
			Height:          uint32(vc.height),
			Layers:          1,
		},
		nil,
		&fbuffer)

	b.framebuffer = fbuffer

	var fence = b.fence

	vulkan.CreateFence(vc.device,
		&vulkan.FenceCreateInfo{
			SType: vulkan.StructureTypeFenceCreateInfo,
			Flags: vulkan.FenceCreateFlags(vulkan.FenceCreateSignaledBit),
		},
		nil,
		&fence)

	b.fence = fence

	var cmd_buffer = b.cmd_buffer

	vulkan.AllocateCommandBuffers(vc.device,
		&vulkan.CommandBufferAllocateInfo{
			SType:              vulkan.StructureTypeCommandBufferAllocateInfo,
			CommandPool:        vc.cmd_pool,
			Level:              vulkan.CommandBufferLevelPrimary,
			CommandBufferCount: 1,
		},
		cmd_buffer[:])

	b.cmd_buffer = cmd_buffer
}

func create_swapchain(vc *VkCube) {
	var surface_caps vulkan.SurfaceCapabilities

	vulkan.GetPhysicalDeviceSurfaceCapabilities(vc.physical_device, vc.surface,
		&surface_caps)

	if 0 == (surface_caps.SupportedCompositeAlpha &
		vulkan.CompositeAlphaFlags(vulkan.CompositeAlphaOpaqueBit)) {
		panic("assert")
	}

	var supported vulkan.Bool32
	vulkan.GetPhysicalDeviceSurfaceSupport(vc.physical_device, 0, vc.surface,
		&supported)
	if !(0 != uint32(supported)) {
		panic("assert")
	}

	var count uint32
	vulkan.GetPhysicalDeviceSurfacePresentModes(vc.physical_device, vc.surface,
		&count, nil)
	var present_modes = make([]vulkan.PresentMode, count)
	vulkan.GetPhysicalDeviceSurfacePresentModes(vc.physical_device, vc.surface,
		&count, present_modes)

	var present_mode vulkan.PresentMode = vulkan.PresentModeMailbox
	for i := uint32(0); i < count; i++ {
		if present_modes[i] == vulkan.PresentModeFifo {
			present_mode = vulkan.PresentModeFifo
			break
		}
	}

	var minImageCount uint32 = 2
	if minImageCount < surface_caps.MinImageCount {
		if surface_caps.MinImageCount > MAX_NUM_IMAGES {
			panic(fmt.Errorf("surface_caps.MinImageCount is too large (is: %d, max: %d)",
				surface_caps.MinImageCount, MAX_NUM_IMAGES))
		}
		minImageCount = surface_caps.MinImageCount
	}

	if surface_caps.MaxImageCount > 0 &&
		minImageCount > surface_caps.MaxImageCount {
		minImageCount = surface_caps.MaxImageCount
	}

	var flags vulkan.SwapchainCreateFlags
	if 0 != uint32(vc.protected) {
		flags = vulkan.SwapchainCreateFlags(vulkan.SwapchainCreateProtectedBit)
	}

	var swpchain = vc.swap_chain[0]

	vulkan.CreateSwapchain(vc.device,
		&vulkan.SwapchainCreateInfo{
			SType:                 vulkan.StructureTypeSwapchainCreateInfo,
			Flags:                 flags,
			Surface:               vc.surface,
			MinImageCount:         minImageCount,
			ImageFormat:           vc.image_format,
			ImageColorSpace:       vulkan.ColorSpaceSrgbNonlinear,
			ImageExtent:           vulkan.Extent2D{Width: uint32(vc.width), Height: uint32(vc.height)},
			ImageArrayLayers:      1,
			ImageUsage:            vulkan.ImageUsageFlags(vulkan.ImageUsageColorAttachmentBit),
			ImageSharingMode:      vulkan.SharingModeExclusive,
			QueueFamilyIndexCount: 1,
			PQueueFamilyIndices:   []uint32{0},
			PreTransform:          vulkan.SurfaceTransformIdentityBit,
			CompositeAlpha:        vulkan.CompositeAlphaOpaqueBit,
			PresentMode:           present_mode,
		}, nil, &swpchain)
	vc.swap_chain[0] = swpchain

	vulkan.GetSwapchainImages(vc.device, vc.swap_chain[0],
		&vc.image_count, nil)
	if !(vc.image_count > 0) {
		panic("assert")
	}
	var swap_chain_images = make([]vulkan.Image, vc.image_count)
	vulkan.GetSwapchainImages(vc.device, vc.swap_chain[0],
		&vc.image_count, swap_chain_images)

	if !(vc.image_count <= MAX_NUM_IMAGES) {
		panic("assert")
	}
	for i := uint32(0); i < vc.image_count; i++ {
		vc.buffers[i].image = swap_chain_images[i]
		init_buffer(vc, &vc.buffers[i])
	}
}

var que vulkan.Queue
var dev vulkan.Device

func init_vk(vc *VkCube) {
	const extension = "VK_KHR_wayland_surface\000"

	var inst = vc.instance

	vulkan.CreateInstance(&vulkan.InstanceCreateInfo{
		SType: vulkan.StructureTypeInstanceCreateInfo,
		PApplicationInfo: &vulkan.ApplicationInfo{
			SType:            vulkan.StructureTypeApplicationInfo,
			PApplicationName: "vkcube",
			ApiVersion:       uint32(1<<22 | 1<<12 | 0),
		},
		EnabledExtensionCount: 2,
		PpEnabledExtensionNames: []string{
			vulkan.KhrSurfaceExtensionName,
			extension,
		},
	},
		nil,
		&inst)
	vc.instance = inst

	var count uint32
	var res = vulkan.EnumeratePhysicalDevices(vc.instance, &count, nil)
	if (res != vulkan.Success) || (count == 0) {
		panic("No Vulkan devices found.\n")
	}
	var pd = make([]vulkan.PhysicalDevice, count)
	vulkan.EnumeratePhysicalDevices(vc.instance, &count, pd)
	vc.physical_device = pd[0]
	fmt.Printf("%d physical devices\n", count)

	var protected_features = vulkan.PhysicalDeviceProtectedMemoryFeatures{
		SType: vulkan.StructureTypePhysicalDeviceProtectedMemoryFeatures,
	}
	var features = vulkan.PhysicalDeviceFeatures2{
		SType: vulkan.StructureTypePhysicalDeviceFeatures2,
		PNext: unsafe.Pointer(&protected_features),
	}
	vulkan.GetPhysicalDeviceFeatures2(vc.physical_device, &features)

	if vc.protected_chain && 0 == protected_features.ProtectedMemory {
		fmt.Print("Requested protected memory but not supported by device, dropping...\n")
	}

	if vc.protected_chain && protected_features.ProtectedMemory != 0 {
		vc.protected = 1
	} else {
		vc.protected = 0
	}

	var properties vulkan.PhysicalDeviceProperties
	vulkan.GetPhysicalDeviceProperties(vc.physical_device, &properties)
	fmt.Printf("vendor id %04x, device name %s\n",
		properties.VendorID, properties.DeviceName)

	vulkan.GetPhysicalDeviceMemoryProperties(vc.physical_device, &vc.memory_properties)

	vulkan.GetPhysicalDeviceQueueFamilyProperties(vc.physical_device, &count, nil)
	if !(count > 0) {
		panic("assert")
	}
	var props = make([]vulkan.QueueFamilyProperties, count)
	vulkan.GetPhysicalDeviceQueueFamilyProperties(vc.physical_device, &count, props)
	if 0 == (props[0].QueueFlags & vulkan.QueueFlags(vulkan.QueueGraphicsBit)) {
		panic("assert")
	}

	var flag vulkan.DeviceQueueCreateFlags = 0
	if vc.protected != 0 {
		flag = vulkan.DeviceQueueCreateFlags(vulkan.DeviceQueueCreateProtectedBit)
	}

	vulkan.CreateDevice(vc.physical_device,
		&vulkan.DeviceCreateInfo{
			SType:                vulkan.StructureTypeDeviceCreateInfo,
			QueueCreateInfoCount: 1,
			PQueueCreateInfos: []vulkan.DeviceQueueCreateInfo{{
				SType:            vulkan.StructureTypeDeviceQueueCreateInfo,
				QueueFamilyIndex: 0,
				QueueCount:       1,
				Flags:            flag,
				PQueuePriorities: []float32{1.0},
			}},
			EnabledExtensionCount: 1,
			PpEnabledExtensionNames: []string{
				vulkan.KhrSwapchainExtensionName,
			},
		},
		nil,
		&dev)

	vulkan.GetDeviceQueue(dev, 0, 0, &que)

	vc.device = dev
	vc.queue = que

}

func main() {

	var vc VkCube
	vc.wl = &wland
	vc.width = 1024
	vc.height = 768
	vc.start = makeTimestamp()

	var err error

	vc.wl.ndisplay, err = wayland.DisplayConnect(nil)
	if err != nil {
		panic(err)
	}

	vc.wl.nregistry, err = wayland.DisplayGetRegistry(vc.wl.ndisplay)
	if err != nil {
		panic(err)
	}

	wayland.RegistryAddListener(vc.wl.nregistry, vc.wl)

	// Round-trip to get globals
	wayland.DisplayRoundtrip(vc.wl.ndisplay)

	// We don't need this anymore
	wayland.RegistryDestroy(vc.wl.nregistry)
	vc.wl.nregistry = nil

	vc.wl.nsurface, err = vc.wl.ncompositor.CreateSurface()
	if err != nil {
		panic(err)
	}

	if nil == vc.wl.nshell {
		panic("Compositor is missing xdg_wm_base protocol support")
	}

	vc.wl.nxdg_surface, err = vc.wl.nshell.GetSurface(vc.wl.nsurface)
	if err != nil {
		panic(err)
	}
	vc.wl.nxdg_surface.AddListener(vc.wl)

	vc.wl.nxdg_toplevel, err = vc.wl.nxdg_surface.GetToplevel()
	if err != nil {
		panic(err)
	}
	wayland.XdgToplevelAddListener(vc.wl.nxdg_toplevel, vc.wl)

	vc.wl.nxdg_toplevel.SetTitle("vkcube")
	vc.wl.wait_for_configure = true

	vc.wl.nsurface.Commit()

	init_vk(&vc)

	if !vulkan.GetPhysicalDeviceWaylandPresentationSupport((vc.instance),
		(vc.physical_device), 0, uintptr(unsafe.Pointer(vc.wl.ndisplay))) {

		panic("no wl support on physical device")
	}

	var inst = vc.instance
	var surf = vc.surface

	vulkan.CreateWaylandSurface(inst,
		&vulkan.WaylandSurfaceCreateInfo{
			SType:   vulkan.StructureTypeWaylandSurfaceCreateInfo,
			Display: uintptr(unsafe.Pointer(vc.wl.ndisplay)),
			Surface: uintptr(unsafe.Pointer(vc.wl.nsurface)),
		}, nil, &surf)

	vc.surface = surf

	vc.image_format = choose_surface_format(&vc)

	init_vk_objects(&vc)

	create_swapchain(&vc)

	mainloop(&vc)

}

// sudo dnf install libdrm-devel
// sudo dnf install libpng-devel
// sudo dnf install libxcb-devel
// sudo dnf install libwayland-client-devel
// sudo dnf install vulkan-headers
// sudo dnf install mesa-libgbm-devel
// sudo dnf install vulkan-loader-devel
