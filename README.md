# wayland

[![GoReportCard](https://goreportcard.com/badge/github.com/neurlang/wayland)](https://goreportcard.com/report/github.com/neurlang/wayland)

[![GoDoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/neurlang/wayland)


Golang Wayland

![Screenshot of Golang Wayland App](wayland.png?raw=true "Screenshot of Golang Wayland App")

This includes two sample apps that render into a shared memory. Tested on pc
(x64) and raspberry pi (arm64). The smoke app reacts on mouse input.

![Screenshot of Golang Wayland Text Editor](gophers.png?raw=true "Screenshot of Golang Wayland Text Editor")

This text editor is one of the sample apps. It is a Linux-only programming text editor.


# Dependencies

None, this is a pure go implementation. But for a meaningful keyboard
support, you need the C libxkbcommon library for apps that require keyboard.

Golang install:

```
sudo apt-get install golang
```

or

```
sudo dnf install golang
```

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

Install libxkbcommon:

```
sudo apt-get install libxkbcommon-dev
```

or in Fedora:

```
sudo dnf install libxkbcommon-devel
```


Next, get the demos:

```
go get github.com/neurlang/wayland/...
```

Then, install them:

```
go get github.com/neurlang/wayland/go-wayland-simple-shm@latest
go get github.com/neurlang/wayland/go-wayland-smoke@latest
go get github.com/neurlang/wayland/go-wayland-imageviewer@latest
go get github.com/neurlang/wayland/go-wayland-texteditor@latest
```

Using golang version < 1.16 (there is support for golang >= 1.09):

```
go install github.com/neurlang/wayland/go-wayland-simple-shm
go install github.com/neurlang/wayland/go-wayland-smoke
go install github.com/neurlang/wayland/go-wayland-imageviewer
go install github.com/neurlang/wayland/go-wayland-texteditor
```

Finally, run weston and the executables.
