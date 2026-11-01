package wl

import (
	"errors"
	"net"
	"path/filepath"
	"testing"
)

func TestSocketPath(t *testing.T) {
	t.Setenv("XDG_RUNTIME_DIR", "/run/user/1000")
	t.Setenv("WAYLAND_DISPLAY", "")
	for _, c := range []struct{ addr, want string }{
		{"", "/run/user/1000/wayland-0"},
		{"wayland-1", "/run/user/1000/wayland-1"},
		{"/tmp/compositor/wayland-0", "/tmp/compositor/wayland-0"},
	} {
		if got, err := socketPath(c.addr); err != nil || got != c.want {
			t.Errorf("socketPath(%q) = %q, %v; want %q", c.addr, got, err, c.want)
		}
	}
	t.Setenv("WAYLAND_DISPLAY", "/tmp/compositor/wayland-0")
	if got, _ := socketPath(""); got != "/tmp/compositor/wayland-0" {
		t.Errorf("an absolute WAYLAND_DISPLAY resolved to %q", got)
	}

	t.Setenv("XDG_RUNTIME_DIR", "")
	if got, err := socketPath("/tmp/compositor/wayland-0"); err != nil || got != "/tmp/compositor/wayland-0" {
		t.Errorf("an absolute addr without XDG_RUNTIME_DIR = %q, %v", got, err)
	}
	if _, err := socketPath("wayland-0"); !errors.Is(err, ErrXdgRuntimeDirNotSet) {
		t.Errorf("a relative addr without XDG_RUNTIME_DIR: err = %v, want ErrXdgRuntimeDirNotSet", err)
	}
}

// Connect dials an absolute WAYLAND_DISPLAY as is, wherever
// XDG_RUNTIME_DIR points.
func TestConnectAbsoluteDisplay(t *testing.T) {
	sock := filepath.Join(t.TempDir(), "wayland-test")
	l, err := net.Listen("unix", sock)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	t.Setenv("XDG_RUNTIME_DIR", t.TempDir())
	t.Setenv("WAYLAND_DISPLAY", sock)
	d, err := Connect("")
	if err != nil {
		t.Fatalf("Connect with WAYLAND_DISPLAY=%s: %v", sock, err)
	}
	d.Context().Close()
}
