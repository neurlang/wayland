# apicheck

A tool to verify that `window/` and `windowtrace/` packages have identical API surfaces across all platforms.

## Usage

```bash
# Check a specific platform
go run ./cmd/apicheck/ linux

# Check all platforms
go run ./cmd/apicheck/ all

# Output as JSON
go run ./cmd/apicheck/ -json all

# Build binary
go build ./cmd/apicheck/
./apicheck all
```

## What it checks

For each platform (`linux`, `windows`, `darwin`, `js`):
- **Types**: All exported types in `window/` have counterparts in `windowtrace/`
- **Constants**: All exported constants are re-exported
- **Interfaces**: Interface method signatures match
- **Functions**: Top-level functions have matching signatures
- **Methods**: Methods on exported types have matching signatures

## Known differences

The tool currently finds these categories of mismatches:
1. Internal types/constants/methods not re-exported by `windowtrace/`
2. `windowtrace/` defines additional interfaces (`CloseHandler`, `DataHandler`, etc.)
3. Platform-specific APIs present in `window/` on some platforms but not all
4. Signature differences due to type alias resolution (e.g., `zxdg.WmBase` vs `xdg.WmBase`)
