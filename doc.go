/*
Package wayland is the root of the wayland sample repository.

# go-wayland-simple-shm

Simple shm demo, draws into shared memory, does not rely on the window package.

# go-wayland-smoke

Smoke demo. Reacts on mouse input, uses the window package.

# go-wayland-imageviewer

ImageViewer demo. Displays image file. Does not use window package, draws
it's own decorations. Draws fonts in the titlebar, for this it needs the
Deja Vu font fonts-dejavu.

# go-wayland-editor

Editor demo. Currently does not work.

# os

Provides basic OS functions like the creation of anonymous temporary file
(CreateAnonymousFile), Mmap, Munmap, and Socket communication.

# wlclient

Utility functions found in the wayland-client, provided for convenience.

# xkbcommon

Wrapper around the C library libxkbcommon. Used inside the window package.
Needs libxkbcommon-dev for compilation and recommends libx11-data for run time
operation.

# window

Implements a window model on top of wayland. Aims to be a lot like the original
window.c code. Uses wl.

# wlcursor

Loads an X cursors, requires to have a cursor theme installed, for instance
dmz-cursor-theme.

# cairoshim

Like cairo but not cairo. Does not depend on anything.

# external

External contains an error checking code and the swizzle function for multiple
architectures. No dependencies.

# unstable

Unstable wayland protocols. Depends on wl.

# wl

The wayland itself, does not require any external deps (except for wayland
server during runtime).

# xdg

Stable xdg protocol. Depends on wl.

*/
package wayland
