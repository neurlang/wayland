package main

import (
	"time"

	vk "github.com/neurlang/wayland/vulkan"
	"github.com/rkusa/gm/mat4"
	"github.com/rkusa/gm/vec3"
	vulkan "github.com/vulkan-go/vulkan"
	"reflect"
	"unsafe"
)

type ubo struct {
	modelview           mat4.Mat4
	modelviewprojection mat4.Mat4
	normal              mat4.Mat4
}

func sync(dst, src *ubo) {
	dst.modelview = src.modelview
	dst.modelviewprojection = src.modelviewprojection
	for y := 0; y < 12; y++ {
		dst.normal[y] = src.normal[y]
	}
}

func makeTimestamp() float64 {
	return float64(time.Now().UnixNano() / int64(time.Millisecond))
}

var vVertices = []float32{
	// front
	-1.0, -1.0, +1.0, // point blue
	+1.0, -1.0, +1.0, // point magenta
	-1.0, +1.0, +1.0, // point cyan
	+1.0, +1.0, +1.0, // point white
	// back
	+1.0, -1.0, -1.0, // point red
	-1.0, -1.0, -1.0, // point black
	+1.0, +1.0, -1.0, // point yellow
	-1.0, +1.0, -1.0, // point green
	// right
	+1.0, -1.0, +1.0, // point magenta
	+1.0, -1.0, -1.0, // point red
	+1.0, +1.0, +1.0, // point white
	+1.0, +1.0, -1.0, // point yellow
	// left
	-1.0, -1.0, -1.0, // point black
	-1.0, -1.0, +1.0, // point blue
	-1.0, +1.0, -1.0, // point green
	-1.0, +1.0, +1.0, // point cyan
	// top
	-1.0, +1.0, +1.0, // point cyan
	+1.0, +1.0, +1.0, // point white
	-1.0, +1.0, -1.0, // point green
	+1.0, +1.0, -1.0, // point yellow
	// bottom
	-1.0, -1.0, -1.0, // point black
	+1.0, -1.0, -1.0, // point red
	-1.0, -1.0, +1.0, // point blue
	+1.0, -1.0, +1.0, // point magenta
}

var vColors = []float32{
	// front
	0.0, 0.0, 1.0, // blue
	1.0, 0.0, 1.0, // magenta
	0.0, 1.0, 1.0, // cyan
	1.0, 1.0, 1.0, // white
	// back
	1.0, 0.0, 0.0, // red
	0.0, 0.0, 0.0, // black
	1.0, 1.0, 0.0, // yellow
	0.0, 1.0, 0.0, // green
	// right
	1.0, 0.0, 1.0, // magenta
	1.0, 0.0, 0.0, // red
	1.0, 1.0, 1.0, // white
	1.0, 1.0, 0.0, // yellow
	// left
	0.0, 0.0, 0.0, // black
	0.0, 0.0, 1.0, // blue
	0.0, 1.0, 0.0, // green
	0.0, 1.0, 1.0, // cyan
	// top
	0.0, 1.0, 1.0, // cyan
	1.0, 1.0, 1.0, // white
	0.0, 1.0, 0.0, // green
	1.0, 1.0, 0.0, // yellow
	// bottom
	0.0, 0.0, 0.0, // black
	1.0, 0.0, 0.0, // red
	0.0, 0.0, 1.0, // blue
	1.0, 0.0, 1.0, // magenta
}

var vNormals = []float32{
	// front
	+0.0, +0.0, +1.0, // forward
	+0.0, +0.0, +1.0, // forward
	+0.0, +0.0, +1.0, // forward
	+0.0, +0.0, +1.0, // forward
	// back
	+0.0, +0.0, -1.0, // backbard
	+0.0, +0.0, -1.0, // backbard
	+0.0, +0.0, -1.0, // backbard
	+0.0, +0.0, -1.0, // backbard
	// right
	+1.0, +0.0, +0.0, // right
	+1.0, +0.0, +0.0, // right
	+1.0, +0.0, +0.0, // right
	+1.0, +0.0, +0.0, // right
	// left
	-1.0, +0.0, +0.0, // left
	-1.0, +0.0, +0.0, // left
	-1.0, +0.0, +0.0, // left
	-1.0, +0.0, +0.0, // left
	// top
	+0.0, +1.0, +0.0, // up
	+0.0, +1.0, +0.0, // up
	+0.0, +1.0, +0.0, // up
	+0.0, +1.0, +0.0, // up
	// bottom
	+0.0, -1.0, +0.0, // down
	+0.0, -1.0, +0.0, // down
	+0.0, -1.0, +0.0, // down
	+0.0, -1.0, +0.0, // down
}

func find_host_coherent_memory(vc *VkCube, allowed uint32) uint32 {
	for i := uint32(0); ((1 << i) <= allowed) && (i <= vc.memory_properties.MemoryTypeCount); i++ {

		if (0 != (allowed & (1 << i))) &&
			(0 != (vc.memory_properties.MemoryTypes[i].PropertyFlags & vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyHostVisibleBit))) &&
			(0 != (vc.memory_properties.MemoryTypes[i].PropertyFlags & vulkan.MemoryPropertyFlags(vulkan.MemoryPropertyHostCoherentBit))) {
			return i
		}
	}
	return ^uint32(0)
}
func init_cube(vc *VkCube) {

	var set_layout vulkan.DescriptorSetLayout
	vulkan.CreateDescriptorSetLayout(vc.device,
		&vulkan.DescriptorSetLayoutCreateInfo{
			SType:        vulkan.StructureTypeDescriptorSetLayoutCreateInfo,
			BindingCount: 1,
			PBindings: []vulkan.DescriptorSetLayoutBinding{
				{
					DescriptorType:     vulkan.DescriptorTypeUniformBuffer,
					DescriptorCount:    1,
					StageFlags:         vulkan.ShaderStageFlags(vulkan.ShaderStageVertexBit),
					PImmutableSamplers: nil,
				},
			},
		},
		nil,
		&set_layout)

	var pipe = vc.pipeline_layout

	vulkan.CreatePipelineLayout(vc.device,
		&vulkan.PipelineLayoutCreateInfo{
			SType:          vulkan.StructureTypePipelineLayoutCreateInfo,
			SetLayoutCount: 1,
			PSetLayouts:    []vulkan.DescriptorSetLayout{set_layout},
		},
		nil,
		&pipe)

	vc.pipeline_layout = pipe

	var vi_create_info = vulkan.PipelineVertexInputStateCreateInfo{
		SType:                         vulkan.StructureTypePipelineVertexInputStateCreateInfo,
		VertexBindingDescriptionCount: 3,
		PVertexBindingDescriptions: []vulkan.VertexInputBindingDescription{
			{
				Binding:   0,
				Stride:    3 * 4,
				InputRate: vulkan.VertexInputRateVertex,
			},
			{
				Binding:   1,
				Stride:    3 * 4,
				InputRate: vulkan.VertexInputRateVertex,
			},
			{
				Binding:   2,
				Stride:    3 * 4,
				InputRate: vulkan.VertexInputRateVertex,
			},
		},
		VertexAttributeDescriptionCount: 3,
		PVertexAttributeDescriptions: []vulkan.VertexInputAttributeDescription{
			{
				Location: 0,
				Binding:  0,
				Format:   vulkan.FormatR32g32b32Sfloat,
				Offset:   0,
			},
			{
				Location: 1,
				Binding:  1,
				Format:   vulkan.FormatR32g32b32Sfloat,
				Offset:   0,
			},
			{
				Location: 2,
				Binding:  2,
				Format:   vulkan.FormatR32g32b32Sfloat,
				Offset:   0,
			},
		},
	}

	_ = vi_create_info

	var vs_module vulkan.ShaderModule
	vulkan.CreateShaderModule(vc.device,
		&vulkan.ShaderModuleCreateInfo{
			SType:    vulkan.StructureTypeShaderModuleCreateInfo,
			CodeSize: uint(4 * len(vs_spirv_source)),
			PCode:    vs_spirv_source,
		},
		nil,
		&vs_module)

	var fs_module vulkan.ShaderModule
	vulkan.CreateShaderModule(vc.device,
		&vulkan.ShaderModuleCreateInfo{
			SType:    vulkan.StructureTypeShaderModuleCreateInfo,
			CodeSize: uint(4 * len(fs_spirv_source)),
			PCode:    fs_spirv_source,
		},
		nil,
		&fs_module)

	var pipeline_stack = ([1]vulkan.Pipeline)(vc.pipeline)

	vulkan.CreateGraphicsPipelines(vc.device,
		vulkan.PipelineCache(*vk.NilPipelineCache),
		1,
		[]vulkan.GraphicsPipelineCreateInfo{{
			SType:      vulkan.StructureTypeGraphicsPipelineCreateInfo,
			StageCount: 2,
			PStages: []vulkan.PipelineShaderStageCreateInfo{
				{
					SType:  vulkan.StructureTypePipelineShaderStageCreateInfo,
					Stage:  vulkan.ShaderStageVertexBit,
					Module: vs_module,
					PName:  "main\000",
				},
				{
					SType:  vulkan.StructureTypePipelineShaderStageCreateInfo,
					Stage:  vulkan.ShaderStageFragmentBit,
					Module: fs_module,
					PName:  "main\000",
				},
			},
			PVertexInputState: &vi_create_info,
			PInputAssemblyState: &vulkan.PipelineInputAssemblyStateCreateInfo{
				SType:                  vulkan.StructureTypePipelineInputAssemblyStateCreateInfo,
				Topology:               vulkan.PrimitiveTopologyTriangleStrip,
				PrimitiveRestartEnable: vulkan.Bool32(0),
			},
			PViewportState: &vulkan.PipelineViewportStateCreateInfo{
				SType:         vulkan.StructureTypePipelineViewportStateCreateInfo,
				ViewportCount: 1,
				ScissorCount:  1,
			},
			PRasterizationState: &vulkan.PipelineRasterizationStateCreateInfo{
				SType:                   vulkan.StructureTypePipelineRasterizationStateCreateInfo,
				RasterizerDiscardEnable: vulkan.Bool32(0),
				PolygonMode:             vulkan.PolygonModeFill,
				CullMode:                vulkan.CullModeFlags(vulkan.CullModeBackBit),
				FrontFace:               vulkan.FrontFaceClockwise,
				LineWidth:               1.0,
			},
			PMultisampleState: &vulkan.PipelineMultisampleStateCreateInfo{
				SType:                vulkan.StructureTypePipelineMultisampleStateCreateInfo,
				RasterizationSamples: 1,
			},
			PDepthStencilState: &vulkan.PipelineDepthStencilStateCreateInfo{
				SType: vulkan.StructureTypePipelineDepthStencilStateCreateInfo,
			},
			PColorBlendState: &vulkan.PipelineColorBlendStateCreateInfo{
				SType:           vulkan.StructureTypePipelineColorBlendStateCreateInfo,
				AttachmentCount: 1,
				PAttachments: []vulkan.PipelineColorBlendAttachmentState{
					{ColorWriteMask: vulkan.ColorComponentFlags(
						vulkan.ColorComponentABit |
							vulkan.ColorComponentRBit |
							vulkan.ColorComponentGBit |
							vulkan.ColorComponentBBit)},
				},
			},
			PDynamicState: &vulkan.PipelineDynamicStateCreateInfo{
				SType:             vulkan.StructureTypePipelineDynamicStateCreateInfo,
				DynamicStateCount: 2,
				PDynamicStates: []vulkan.DynamicState{
					vulkan.DynamicStateViewport,
					vulkan.DynamicStateScissor,
				},
			},
			Flags:              0,
			Layout:             vc.pipeline_layout,
			RenderPass:         vc.render_pass,
			Subpass:            0,
			BasePipelineHandle: vulkan.Pipeline(*vk.NilPipeline),
			BasePipelineIndex:  0,
		}},
		nil,
		pipeline_stack[:])

	vc.pipeline = pipeline_stack

	vc.vertex_offset = vulkan.DeviceSize(unsafe.Sizeof(ubo{}))
	vc.colors_offset = vulkan.DeviceSize(vc.vertex_offset + vulkan.DeviceSize(uintptr(len(vVertices))*unsafe.Sizeof(vVertices[0])))
	vc.normals_offset = vulkan.DeviceSize(vc.colors_offset + vulkan.DeviceSize(uintptr(len(vColors))*unsafe.Sizeof(vColors[0])))
	var mem_size vulkan.DeviceSize = vulkan.DeviceSize(vc.normals_offset) + vulkan.DeviceSize(uintptr(len(vNormals))*unsafe.Sizeof(vNormals[0]))

	var buff = vc.buffer

	vulkan.CreateBuffer(vc.device,
		&vulkan.BufferCreateInfo{
			SType: vulkan.StructureTypeBufferCreateInfo,
			Size:  mem_size,
			Usage: vulkan.BufferUsageFlags(
				vulkan.BufferUsageUniformBufferBit |
					vulkan.BufferUsageVertexBufferBit),
			Flags: 0,
		},
		nil,
		&buff)
	vc.buffer = buff

	var reqs vulkan.MemoryRequirements
	vulkan.GetBufferMemoryRequirements(vc.device, vc.buffer, &reqs)

	reqs.Deref()

	var memory_type = uint32(find_host_coherent_memory(vc, reqs.MemoryTypeBits))
	if memory_type == ^uint32(0) {
		panic("find_host_coherent_memory failed")
	}

	var memory = vc.mem

	vulkan.AllocateMemory(vc.device,
		&vulkan.MemoryAllocateInfo{
			SType:           vulkan.StructureTypeMemoryAllocateInfo,
			AllocationSize:  mem_size,
			MemoryTypeIndex: memory_type,
		},
		nil,
		&memory)
	vc.mem = memory

	var mapping = vc.mapping

	var r = vulkan.MapMemory(vc.device, vc.mem, 0, mem_size, 0, (*unsafe.Pointer)(unsafe.Pointer(&mapping)))
	if r != vulkan.Success {
		panic("vkMapMemory failed")
	}

	vc.mapping = mapping

	var vertex_offset_slice []float32
	vertex_offset_slice_hdr := (*reflect.SliceHeader)(unsafe.Pointer(&vertex_offset_slice))
	vertex_offset_slice_hdr.Data = (uintptr(unsafe.Pointer(vc.mapping)) + uintptr(vc.vertex_offset))
	vertex_offset_slice_hdr.Len = len(vVertices)
	vertex_offset_slice_hdr.Cap = len(vVertices)
	var colors_offset_slice []float32
	colors_offset_slice_hdr := (*reflect.SliceHeader)(unsafe.Pointer(&colors_offset_slice))
	colors_offset_slice_hdr.Data = (uintptr(unsafe.Pointer(vc.mapping)) + uintptr(vc.colors_offset))
	colors_offset_slice_hdr.Len = len(vColors)
	colors_offset_slice_hdr.Cap = len(vColors)
	var normals_offset_slice []float32
	normals_offset_slice_hdr := (*reflect.SliceHeader)(unsafe.Pointer(&normals_offset_slice))
	normals_offset_slice_hdr.Data = (uintptr(unsafe.Pointer(vc.mapping)) + uintptr(vc.normals_offset))
	normals_offset_slice_hdr.Len = len(vNormals)
	normals_offset_slice_hdr.Cap = len(vNormals)

	for i, v := range vVertices {
		vertex_offset_slice[i] = v
	}
	for i, v := range vColors {
		colors_offset_slice[i] = v
	}
	for i, v := range vNormals {
		normals_offset_slice[i] = v
	}

	vulkan.BindBufferMemory(vc.device, vc.buffer, vc.mem, 0)

	var desc_pool vulkan.DescriptorPool
	var create_info = vulkan.DescriptorPoolCreateInfo{
		SType:         vulkan.StructureTypeDescriptorPoolCreateInfo,
		PNext:         nil,
		Flags:         0,
		MaxSets:       1,
		PoolSizeCount: 1,
		PPoolSizes: []vulkan.DescriptorPoolSize{{
			Type:            vulkan.DescriptorTypeUniformBuffer,
			DescriptorCount: 1,
		},
		},
	}

	vulkan.CreateDescriptorPool(vc.device, &create_info, nil, &desc_pool)

	var descset = vc.descriptor_set

	vulkan.AllocateDescriptorSets(vc.device,
		&vulkan.DescriptorSetAllocateInfo{
			SType:              vulkan.StructureTypeDescriptorSetAllocateInfo,
			DescriptorPool:     desc_pool,
			DescriptorSetCount: 1,
			PSetLayouts:        []vulkan.DescriptorSetLayout{set_layout},
		}, &descset)

	vc.descriptor_set = descset

	vulkan.UpdateDescriptorSets(vc.device, 1,
		[]vulkan.WriteDescriptorSet{{
			SType:           vulkan.StructureTypeWriteDescriptorSet,
			DstSet:          vc.descriptor_set,
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vulkan.DescriptorTypeUniformBuffer,
			PBufferInfo: []vulkan.DescriptorBufferInfo{{
				Buffer: vc.buffer,
				Offset: 0,
				Range:  vulkan.DeviceSize(unsafe.Sizeof(ubo{})),
			}},
		},
		},
		0, nil)
}

func render_cube(vc *VkCube, b *VkCubeBuffer, wait_semaphore_count uint8) {
	var ubo ubo

	t := float32((makeTimestamp() - vc.start) / 5.)

	var aspect = float32(vc.height) / float32(vc.width+1)
	var projection mat4.Mat4
	var mwproject = (*mat4.Mat4)((&ubo.modelviewprojection))
	var project = (*mat4.Mat4)((&projection))
	var mw = (*mat4.Mat4)((&ubo.modelview))

	*mw = *mat4.Identity()

	mw = mw.Translate(vec3.New(0, 0, -8))

	mw.Mul(mat4.Identity().Rotation(vec3.New((45.0+(0.25*t))*3.14/180., (45.0-(0.5*t))*3.14/180., (00.0+(0.15*t))*3.14/180.)))

	Frustum(project, -2.8, 2.8, -2.8*aspect, 2.8*aspect, 6, 10)

	project.Mul(mw)

	*mwproject = *project

	/* The mat3 normalMatrix is laid out as 3 vec4s. */
	ubo.normal = ubo.modelview

	sync(vc.mapping, &ubo)

	vulkan.WaitForFences(vc.device, 1, []vulkan.Fence{b.fence}, vulkan.True, ^uint64(0))
	vulkan.ResetFences(vc.device, 1, []vulkan.Fence{b.fence})

	vulkan.BeginCommandBuffer(b.cmd_buffer[0],
		&vulkan.CommandBufferBeginInfo{
			SType: vulkan.StructureTypeCommandBufferBeginInfo,
			Flags: 0,
		})

	vulkan.CmdBeginRenderPass(b.cmd_buffer[0],
		&vulkan.RenderPassBeginInfo{
			SType:           vulkan.StructureTypeRenderPassBeginInfo,
			RenderPass:      vc.render_pass,
			Framebuffer:     b.framebuffer,
			RenderArea:      vulkan.Rect2D{Offset: vulkan.Offset2D{X: 0, Y: 0}, Extent: vulkan.Extent2D{Width: uint32(vc.width), Height: uint32(vc.height)}},
			ClearValueCount: 1,
			PClearValues: []vulkan.ClearValue{
				{
					// floats 32bit
					0, 0, 0x80, 0x3e, //0.25
					0, 0, 0x80, 0x3e, //0.25
					0, 0, 0x80, 0x3e, //0.25
					0, 0, 0x80, 0x3f, //1.00
				},
			},
		},
		vulkan.SubpassContentsInline)

	vulkan.CmdBindVertexBuffers(b.cmd_buffer[0], 0, 3,
		[]vulkan.Buffer{
			vc.buffer,
			vc.buffer,
			vc.buffer,
		},
		[]vulkan.DeviceSize{
			vc.vertex_offset,
			vc.colors_offset,
			vc.normals_offset,
		})

	vulkan.CmdBindPipeline(b.cmd_buffer[0], vulkan.PipelineBindPointGraphics, vc.pipeline[0])

	vulkan.CmdBindDescriptorSets(b.cmd_buffer[0],
		vulkan.PipelineBindPointGraphics,
		vc.pipeline_layout,
		0, 1,
		[]vulkan.DescriptorSet{vc.descriptor_set}, 0, nil)

	var viewport = vulkan.Viewport{
		X:        0,
		Y:        0,
		Width:    float32(vc.width),
		Height:   float32(vc.height),
		MinDepth: 0,
		MaxDepth: 1,
	}
	vulkan.CmdSetViewport(b.cmd_buffer[0], 0, 1, []vulkan.Viewport{viewport})

	var scissor = vulkan.Rect2D{
		Offset: vulkan.Offset2D{X: 0, Y: 0},
		Extent: vulkan.Extent2D{Width: uint32(vc.width), Height: uint32(vc.height)},
	}
	vulkan.CmdSetScissor(b.cmd_buffer[0], 0, 1, []vulkan.Rect2D{scissor})

	vulkan.CmdDraw(b.cmd_buffer[0], 4, 1, 0, 0)
	vulkan.CmdDraw(b.cmd_buffer[0], 4, 1, 4, 0)
	vulkan.CmdDraw(b.cmd_buffer[0], 4, 1, 8, 0)
	vulkan.CmdDraw(b.cmd_buffer[0], 4, 1, 12, 0)
	vulkan.CmdDraw(b.cmd_buffer[0], 4, 1, 16, 0)
	vulkan.CmdDraw(b.cmd_buffer[0], 4, 1, 20, 0)

	vulkan.CmdEndRenderPass(b.cmd_buffer[0])

	vulkan.EndCommandBuffer(b.cmd_buffer[0])

	var protected_info = vulkan.ProtectedSubmitInfo{
		SType:           vulkan.StructureTypeProtectedSubmitInfo,
		ProtectedSubmit: vc.protected,
	}

	vulkan.QueueSubmit(vc.queue, 1,
		[]vulkan.SubmitInfo{{
			SType: vulkan.StructureTypeSubmitInfo,
			PNext: unsafe.Pointer(&protected_info),
			/* headless mode does not signal vc.semaphore */
			WaitSemaphoreCount: uint32(wait_semaphore_count),
			PWaitSemaphores:    []vulkan.Semaphore{vc.semaphore},
			PWaitDstStageMask: []vulkan.PipelineStageFlags{
				vulkan.PipelineStageFlags(vulkan.PipelineStageColorAttachmentOutputBit),
			},
			CommandBufferCount: 1,
			PCommandBuffers:    b.cmd_buffer[:],
		}}, b.fence)
}
