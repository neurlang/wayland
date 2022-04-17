package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	wayland "github.com/neurlang/wayland/libwayland"
	vulkan "github.com/neurlang/wayland/vulkan"
	vk "github.com/vulkan-go/vulkan"
	"golang.org/x/sys/unix"
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
		wayland.XdgWmBaseAddListener(w.nshell, w)
	case "wl_seat":

		w.nseat = wayland.RegistryBindSeatInterface(w.nregistry, name, 1)

	}

}

func (w *Wland) RegistryGlobalRemove(*wayland.Registry, uint32) {
}

func (w *Wland) XdgSurfaceConfigure(surface *wayland.XdgSurface, serial uint32) {

	wayland.XdgSurfaceAckConfigure(surface, serial)

	if w.wait_for_configure {
		// redraw
		w.wait_for_configure = false
	}

}
func (w *Wland) XdgWmBasePing(shell *wayland.XdgWmBase, serial uint32) {
	wayland.XdgWmBasePong(shell, serial)
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

	device     vk.Device
	swap_chain [1]vk.Swapchain
	queue      vk.Queue
	semaphore  vk.Semaphore

	buffers [MAX_NUM_IMAGES]VkCubeBuffer

	mapping *ubo

	buffer          vk.Buffer
	render_pass     vk.RenderPass
	vertex_offset   vk.DeviceSize
	colors_offset   vk.DeviceSize
	normals_offset  vk.DeviceSize
	pipeline        [1]vk.Pipeline
	pipeline_layout vk.PipelineLayout
	descriptor_set  vk.DescriptorSet

	protected vk.Bool32

	start float64

	physical_device vk.PhysicalDevice
	surface         vk.Surface
	image_count     uint32
	image_format    vk.Format
	cmd_pool        vk.CommandPool

	mem vk.DeviceMemory

	memory_properties vk.PhysicalDeviceMemoryProperties
	instance          vk.Instance

	protected_chain bool
}

func (vc *VkCube) Render(buf *VkCubeBuffer, wait_semaphore uint8) {
	render_cube(vc, buf, wait_semaphore)
}
func (vc *VkCube) Init() {
	init_cube(vc)
}

type VkCubeBuffer struct {
	mem         vk.DeviceMemory
	image       vk.Image
	view        [1]vk.ImageView
	framebuffer vk.Framebuffer
	fence       vk.Fence
	cmd_buffer  [1]vk.CommandBuffer

	fb     uint32
	stride uint32
}

func init_vk_objects(vc *VkCube) {

	var pass = vc.render_pass

	vk.CreateRenderPass(vc.device,
		&vk.RenderPassCreateInfo{
			SType:           vk.StructureTypeRenderPassCreateInfo,
			AttachmentCount: 1,
			PAttachments: []vk.AttachmentDescription{
				{
					Format:        vc.image_format,
					Samples:       1,
					LoadOp:        vk.AttachmentLoadOpClear,
					StoreOp:       vk.AttachmentStoreOpStore,
					InitialLayout: vk.ImageLayoutUndefined,
					FinalLayout:   vk.ImageLayoutPresentSrc,
				},
			},
			SubpassCount: 1,
			PSubpasses: []vk.SubpassDescription{
				{
					PipelineBindPoint:    vk.PipelineBindPointGraphics,
					InputAttachmentCount: 0,
					ColorAttachmentCount: 1,
					PColorAttachments: []vk.AttachmentReference{
						{
							Attachment: 0,
							Layout:     vk.ImageLayoutColorAttachmentOptimal,
						},
					},
					PResolveAttachments: []vk.AttachmentReference{
						{
							Attachment: vk.AttachmentUnused,
							Layout:     vk.ImageLayoutColorAttachmentOptimal,
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

	var flags vk.CommandPoolCreateFlags
	if 0 != uint32(vc.protected) {
		flags = vk.CommandPoolCreateFlags(vk.CommandPoolCreateProtectedBit)
	}

	var cmdpool = vc.cmd_pool

	vk.CreateCommandPool(vc.device,
		&vk.CommandPoolCreateInfo{
			SType:            vk.StructureTypeCommandPoolCreateInfo,
			QueueFamilyIndex: 0,
			Flags:            vk.CommandPoolCreateFlags(vk.CommandPoolCreateResetCommandBufferBit) | flags,
		},
		nil,
		&cmdpool)
	vc.cmd_pool = cmdpool

	var sema = vc.semaphore

	vk.CreateSemaphore(vc.device,
		&vk.SemaphoreCreateInfo{
			SType: vk.StructureTypeSemaphoreCreateInfo,
		},
		nil,
		&sema)

	vc.semaphore = sema
}

func choose_surface_format(vc *VkCube) vk.Format {
	var num_formats uint32 = 0

	vk.GetPhysicalDeviceSurfaceFormats(vc.physical_device, vc.surface,
		&num_formats, nil)
	if !(num_formats > 0) {
		panic("assert")
	}

	var formats = make([]vk.SurfaceFormat, num_formats)

	vk.GetPhysicalDeviceSurfaceFormats(vc.physical_device, vc.surface,
		&num_formats, formats)

	for i := range formats {
		formats[i].Deref()
	}

	var format vk.Format = vk.FormatUndefined

	for i := uint32(0); i < num_formats; i++ {
		switch formats[i].Format {
		case vk.FormatR8g8b8a8Srgb:
			fallthrough
		case vk.FormatB8g8r8a8Srgb:
			/* These formats are all fine */
			format = formats[i].Format
			continue
		case vk.FormatR8g8b8Srgb:
			fallthrough
		case vk.FormatB8g8r8Srgb:
			fallthrough
		case vk.FormatR5g6b5UnormPack16:
			fallthrough
		case vk.FormatB5g6r5UnormPack16:
			fallthrough
			/* We would like to support these but they don't seem to work. */
		default:
			println(formats[i].Format)
			continue
		}
	}

	if !(format != vk.FormatUndefined) {
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

		var result = [1]vk.Result{vk.AcquireNextImage(vc.device, vc.swap_chain[0], 60,
			vc.semaphore, nil, &index)}
		if result[0] != vk.Success {
			return
		}

		if !(index <= MAX_NUM_IMAGES) {
			panic("assert")
		}

		vc.Render(&vc.buffers[index], 1)

		vk.QueuePresent(vc.queue,
			&vk.PresentInfo{
				SType:          vk.StructureTypePresentInfo,
				SwapchainCount: 1,
				PSwapchains:    vc.swap_chain[:],
				PImageIndices:  []uint32{index},
				PResults:       result[:],
			})
		if result[0] != vk.Success {
			return
		}

		vk.QueueWaitIdle(vc.queue)
	}

}

func init_buffer(vc *VkCube, b *VkCubeBuffer) {

	var iview = b.view[0]

	vk.CreateImageView(vc.device,
		&vk.ImageViewCreateInfo{
			SType:    vk.StructureTypeImageViewCreateInfo,
			Image:    b.image,
			ViewType: vk.ImageViewType2d,
			Format:   vc.image_format,
			Components: vk.ComponentMapping{
				R: vk.ComponentSwizzleR,
				G: vk.ComponentSwizzleG,
				B: vk.ComponentSwizzleB,
				A: vk.ComponentSwizzleA,
			},
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
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

	vk.CreateFramebuffer(vc.device,
		&vk.FramebufferCreateInfo{
			SType:           vk.StructureTypeFramebufferCreateInfo,
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

	vk.CreateFence(vc.device,
		&vk.FenceCreateInfo{
			SType: vk.StructureTypeFenceCreateInfo,
			Flags: vk.FenceCreateFlags(vk.FenceCreateSignaledBit),
		},
		nil,
		&fence)

	b.fence = fence

	var cmd_buffer = b.cmd_buffer

	vk.AllocateCommandBuffers(vc.device,
		&vk.CommandBufferAllocateInfo{
			SType:              vk.StructureTypeCommandBufferAllocateInfo,
			CommandPool:        vc.cmd_pool,
			Level:              vk.CommandBufferLevelPrimary,
			CommandBufferCount: 1,
		},
		cmd_buffer[:])

	b.cmd_buffer = cmd_buffer
}

func create_swapchain(vc *VkCube) {
	var surface_caps vk.SurfaceCapabilities

	vk.GetPhysicalDeviceSurfaceCapabilities(vc.physical_device, vc.surface,
		&surface_caps)

	surface_caps.Deref()

	if 0 == (surface_caps.SupportedCompositeAlpha &
		vk.CompositeAlphaFlags(vk.CompositeAlphaOpaqueBit)) {
		panic("assert")
	}

	var supported vk.Bool32
	vk.GetPhysicalDeviceSurfaceSupport(vc.physical_device, 0, vc.surface,
		&supported)
	if !(0 != uint32(supported)) {
		panic("assert")
	}

	var count uint32
	vk.GetPhysicalDeviceSurfacePresentModes(vc.physical_device, vc.surface,
		&count, nil)
	var present_modes = make([]vk.PresentMode, count)
	vk.GetPhysicalDeviceSurfacePresentModes(vc.physical_device, vc.surface,
		&count, present_modes)

	var present_mode vk.PresentMode = vk.PresentModeMailbox
	for i := uint32(0); i < count; i++ {
		if present_modes[i] == vk.PresentModeFifo {
			present_mode = vk.PresentModeFifo
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

	var flags vk.SwapchainCreateFlags
	if 0 != uint32(vc.protected) {
		flags = vk.SwapchainCreateFlags(vk.SwapchainCreateProtectedBit)
	}

	var swpchain = vc.swap_chain[0]

	vk.CreateSwapchain(vc.device,
		&vk.SwapchainCreateInfo{
			SType:                 vk.StructureTypeSwapchainCreateInfo,
			Flags:                 flags,
			Surface:               vc.surface,
			MinImageCount:         minImageCount,
			ImageFormat:           vc.image_format,
			ImageColorSpace:       vk.ColorSpaceSrgbNonlinear,
			ImageExtent:           vk.Extent2D{Width: uint32(vc.width), Height: uint32(vc.height)},
			ImageArrayLayers:      1,
			ImageUsage:            vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit),
			ImageSharingMode:      vk.SharingModeExclusive,
			QueueFamilyIndexCount: 1,
			PQueueFamilyIndices:   []uint32{0},
			PreTransform:          vk.SurfaceTransformIdentityBit,
			CompositeAlpha:        vk.CompositeAlphaOpaqueBit,
			PresentMode:           present_mode,
		}, nil, &swpchain)
	vc.swap_chain[0] = swpchain

	vk.GetSwapchainImages(vc.device, vc.swap_chain[0],
		&vc.image_count, nil)
	if !(vc.image_count > 0) {
		panic("assert")
	}
	var swap_chain_images = make([]vk.Image, vc.image_count)
	vk.GetSwapchainImages(vc.device, vc.swap_chain[0],
		&vc.image_count, swap_chain_images)

	if !(vc.image_count <= MAX_NUM_IMAGES) {
		panic("assert")
	}
	for i := uint32(0); i < vc.image_count; i++ {
		vc.buffers[i].image = swap_chain_images[i]
		init_buffer(vc, &vc.buffers[i])
	}
}

func init_vk(vc *VkCube) {
	const extension = "VK_KHR_wayland_surface\000"

	// OR without using a windowing library (Linux only, recommended for compute-only tasks)
	if err := vk.SetDefaultGetInstanceProcAddr(); err != nil {
		panic(err)
	}

	var inst vk.Instance

	vulkan.CreateInstance(&vk.InstanceCreateInfo{
		SType: vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo: &vk.ApplicationInfo{
			SType:            vk.StructureTypeApplicationInfo,
			PApplicationName: "vkcube",
			ApiVersion:       uint32(1<<22 | 1<<12 | 0),
		},
		EnabledExtensionCount: 2,
		PpEnabledExtensionNames: []string{
			vk.KhrSurfaceExtensionName + "\000",
			extension,
		},
	},
		nil,
		&inst)
	vc.instance = vk.Instance(unsafe.Pointer(inst))

	if err := vk.InitInstance(vk.Instance(unsafe.Pointer(inst))); err != nil {
		panic(err)
	}

	var count uint32
	var res = vk.EnumeratePhysicalDevices(vk.Instance(unsafe.Pointer(vc.instance)), &count, nil)
	if (res != vk.Success) || (count == 0) {
		panic("No Vulkan devices found.\n")
	}
	var pd = make([]vk.PhysicalDevice, count)
	vk.EnumeratePhysicalDevices(vc.instance, &count, pd)
	vc.physical_device = pd[0]
	fmt.Printf("%d physical devices\n", count)

	var protected_features = vk.PhysicalDeviceProtectedMemoryFeatures{
		SType: vk.StructureTypePhysicalDeviceProtectedMemoryFeatures,
	}
	var features = vk.PhysicalDeviceFeatures2{
		SType: vk.StructureTypePhysicalDeviceFeatures2,
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

	var properties vk.PhysicalDeviceProperties
	vk.GetPhysicalDeviceProperties(vc.physical_device, &properties)

	properties.Deref()

	fmt.Printf("vendor id %04x, device name %s\n",
		properties.VendorID, properties.DeviceName)

	vk.GetPhysicalDeviceMemoryProperties(vc.physical_device, &vc.memory_properties)

	vc.memory_properties.Deref()

	for i := range vc.memory_properties.MemoryTypes {
		vc.memory_properties.MemoryTypes[i].Deref()
	}

	vk.GetPhysicalDeviceQueueFamilyProperties(vc.physical_device, &count, nil)
	if !(count > 0) {
		panic("assert")
	}
	var props = make([]vk.QueueFamilyProperties, count)
	vk.GetPhysicalDeviceQueueFamilyProperties(vc.physical_device, &count, props)

	for i := range props {
		props[i].Deref()
	}

	if 0 == (props[0].QueueFlags & vk.QueueFlags(vk.QueueGraphicsBit)) {
		panic("assert")
	}

	var flag vk.DeviceQueueCreateFlags = 0
	if vc.protected != 0 {
		flag = vk.DeviceQueueCreateFlags(vk.DeviceQueueCreateProtectedBit)
	}

	var que vk.Queue
	var dev vk.Device

	vk.CreateDevice(vc.physical_device,
		&vk.DeviceCreateInfo{
			SType:                vk.StructureTypeDeviceCreateInfo,
			QueueCreateInfoCount: 1,
			PQueueCreateInfos: []vk.DeviceQueueCreateInfo{{
				SType:            vk.StructureTypeDeviceQueueCreateInfo,
				QueueFamilyIndex: 0,
				QueueCount:       1,
				Flags:            flag,
				PQueuePriorities: []float32{1.0},
			}},
			EnabledExtensionCount: 1,
			PpEnabledExtensionNames: []string{
				vk.KhrSwapchainExtensionName + "\000",
			},
		},
		nil,
		&dev)

	vk.GetDeviceQueue(dev, 0, 0, &que)

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

	wayland.DisplayRoundtrip(vc.wl.ndisplay)

	init_vk(&vc)

	if !vulkan.GetPhysicalDeviceWaylandPresentationSupport(
		(vc.physical_device), 0, uintptr(unsafe.Pointer(vc.wl.ndisplay))) {

		panic("no wl support on physical device")
	}

	var inst = vc.instance
	var surf = vc.surface

	vulkan.CreateWaylandSurface(inst,
		&vulkan.WaylandSurfaceCreateInfo{
			SType:   vk.StructureTypeWaylandSurfaceCreateInfo,
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
