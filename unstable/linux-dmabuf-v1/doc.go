package linux_dmabuf

//go:generate go run ../../cmd/go-wayland-scanner/scanner.go -pkg linux_dmabuf -i https://gitlab.freedesktop.org/wayland/wayland-protocols/-/raw/d10d18f3d49374d2e3eb96d63511f32795aab5f7/unstable/linux-dmabuf/linux-dmabuf-unstable-v1.xml -o linux_dmabuf.go
