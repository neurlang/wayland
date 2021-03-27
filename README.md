# wayland

[![GoReportCard](https://goreportcard.com/badge/github.com/neurlang/wayland)](https://goreportcard.com/report/github.com/neurlang/wayland)

[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/neurlang/wayland)


Golang Wayland

![Screenshot of Golang Wayland App](wayland.png?raw=true "Screenshot of Golang Wayland App")

This includes two sample apps that render into a shared memory. Tested on pc
(x64) and raspberry pi (arm64). The smoke app reacts on mouse input.

# Dependencies

None, this is a pure go implementation

# Installation

First, you need a wayland-enabled Linux, if you don't have one, install
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

Finally run weston and the executables.
