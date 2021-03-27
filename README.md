# wayland

[![GoReportCard](https://goreportcard.com/badge/github.com/neurlang/wayland)](https://goreportcard.com/report/github.com/neurlang/wayland)

[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/neurlang/wayland)


Golang Wayland

![Screenshot of Golang Wayland App](wayland.png?raw=true "Screenshot of Golang Wayland App")

This includes two sample apps that render into a shared memory. Tested on pc
(x64) and raspberry pi (arm64). The smoke app reacts on mouse input.

# Dependencies

None, this is a pure go implementation

# Docker Installation

Run the docker build command in the provided build.sh script.

Next, start your Wayland compositor, and run demos using the docker run
commands that can be found in the run-shm.sh or run-smoke.sh scripts.

*Important:* If your docker requires root privileges, use sudo -E to start
the programs. This is because *$XDG_RUNTIME_DIR* env variable is required.

# Installation

First, you need a wayland-enabled Linux, if you don't have one, install
the Wayland compositor weston that is useful for testing:

```
sudo apt-get install weston
```

Next, get the demos:

```
go get github.com/neurlang/wayland/...
```

Then, install them:

```
go install github.com/neurlang/wayland/go-wayland-simple-shm@latest
go install github.com/neurlang/wayland/go-wayland-smoke@latest
go install github.com/neurlang/wayland/go-wayland-imageviewer@latest
```

Using golang version < 1.16 (there is support for golang >= 1.09):

```
go install github.com/neurlang/wayland/go-wayland-simple-shm
go install github.com/neurlang/wayland/go-wayland-smoke
go install github.com/neurlang/wayland/go-wayland-imageviewer
```

Finally, run weston and the executables.
