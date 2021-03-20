# wayland

Golang Wayland

![Screenshot of Golang Wayland App](wayland.png?raw=true "Screenshot of Golang Wayland App")

This includes two sample apps that render into a shared memory. Tested on pc
(x64) and raspberry pi (arm64). The smoke app reacts on mouse input.

# Dependencies

None, this is a pure go implementation

# Installation

First of all you need a wayland-enabled Linux, if you don't have one, install
the Wayland compositor weston that is useful for testing:

```
sudo apt-get install weston
```

Next install the demos:

```
go get github.com/neurlang/wayland/...
go install github.com/neurlang/wayland/go-wayland-simple-shm
go install github.com/neurlang/wayland/go-wayland-smoke
```

And finally run weston and the executables.
