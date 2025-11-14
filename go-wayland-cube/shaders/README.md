# Shader Generation

This directory contains the GLSL shader sources for the Vulkan cube demo.

## Prerequisites

You need `glslangValidator` to compile GLSL to SPIR-V:

**Ubuntu/Debian:**
```bash
sudo apt install glslang-tools
```

**Fedora/RHEL:**
```bash
sudo dnf install glslang
```

**Arch Linux:**
```bash
sudo pacman -S glslang
```

## Generating Shaders

From the `go-wayland-cube` directory, run:

```bash
go generate
```

This will:
1. Compile `shaders/vertex.glsl` → `shaders/vertex.spv` → `vertex_shader.go`
2. Compile `shaders/fragment.glsl` → `shaders/fragment.spv` → `fragment_shader.go`

The generated Go files contain the SPIR-V bytecode as `[]uint32` arrays that can be directly used by Vulkan.

## Files

- `vertex.glsl` - Vertex shader source (GLSL)
- `fragment.glsl` - Fragment shader source (GLSL)
- `generate.go` - Generator program that compiles shaders and creates Go files
- `*.spv` - Compiled SPIR-V binaries (generated, not committed)

## Modifying Shaders

1. Edit `vertex.glsl` or `fragment.glsl`
2. Run `go generate` from the parent directory
3. The corresponding `*_shader.go` files will be regenerated
