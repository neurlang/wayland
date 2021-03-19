# wayland

Golang Wayland

![Screenshot of Golang Wayland App](wayland.png?raw=true "Screenshot of Golang Wayland App")

This includes two sample apps that render into a shared memory. Tested on pc
(x64) and raspberry pi (arm64). The smoke app reacts on mouse input.

# Dependencies

This relies on CGO to compile a wrapper around the libwayland library.

```
sudo apt-get install libwayland-dev
```

# Installation

```
go get github.com/neurlang/wayland/...
```
