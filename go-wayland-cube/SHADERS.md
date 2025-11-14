# Shader Build System

The Vulkan cube demo now uses `go generate` to automatically compile GLSL shaders to SPIR-V and embed them in Go source files.

## Quick Start

```bash
# Install glslangValidator (one-time setup)
sudo apt install glslang-tools  # Ubuntu/Debian
# or
sudo dnf install glslang         # Fedora/RHEL

# Generate shaders
cd go-wayland-cube
go generate
```

## How It Works

1. **Source Shaders** (`shaders/*.glsl`)
   - `vertex.glsl` - Vertex shader with lighting calculations
   - `fragment.glsl` - Fragment shader for color output

2. **Generator** (`shaders/generate.go`)
   - Runs `glslangValidator` to compile GLSL → SPIR-V
   - Reads the binary SPIR-V files
   - Generates Go source files with embedded bytecode

3. **Generated Files** (`*_shader.go`)
   - `vertex_shader.go` - Contains `vs_spirv_source []uint32`
   - `fragment_shader.go` - Contains `fs_spirv_source []uint32`
   - These files are auto-generated and should be committed

## Benefits

✅ **Version Control** - Shader source (GLSL) is human-readable and diffable  
✅ **Reproducible** - Anyone can regenerate the SPIR-V from source  
✅ **Automated** - `go generate` handles compilation automatically  
✅ **No Runtime Deps** - SPIR-V is embedded, no external shader files needed  
✅ **Type Safe** - Shaders are Go variables, checked at compile time  

## Workflow

### Modifying Shaders

1. Edit `shaders/vertex.glsl` or `shaders/fragment.glsl`
2. Run `go generate` to recompile
3. Commit both the `.glsl` source and generated `*_shader.go` files

### Adding New Shaders

1. Create `shaders/newshader.glsl`
2. Add entry to `shaders/generate.go` in the `shaders` slice
3. Run `go generate`

## Technical Details

- **SPIR-V Format**: Binary format, stored as `[]uint32` (little-endian)
- **Compiler**: Uses `glslangValidator` (Khronos reference compiler)
- **GLSL Version**: 4.20 core profile
- **Output Format**: 8 uint32 values per line for readability
