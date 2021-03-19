#include <wayland-client.h>

#include "bridge.h"
#include "_cgo_export.h"



void
_go_shm_listener(void *data, struct wl_shm *wl_shm, uint32_t format)
{
	wlcallback_shm_format((uintptr_t)(data), wl_shm, format);
}

struct wl_shm_listener _shm_listener = {
    _go_shm_listener
};

int
_wl_shm_add_listener(struct wl_shm *wl_shm, void *data)
{
	wl_shm_add_listener(wl_shm, &_shm_listener, data);
}

void
_registry_handle_global(void *data, struct wl_registry *registry,
		       uint32_t id, const char *interface, uint32_t version)
{
	char *iface = interface; 

	wlcallback_registry_global((uintptr_t)data, registry, id, iface, version);


}

void
_registry_handle_global_remove(void *data, struct wl_registry *registry,
			      uint32_t name)
{
}

const struct wl_registry_listener _registry_listener = {
	_registry_handle_global,
	_registry_handle_global_remove
};

int
_wl_registry_add_listener(struct wl_registry *wl_registry, void *data)
{
    wl_registry_add_listener(wl_registry, &_registry_listener,
        data);
}


static void
_xdg_shell_ping(void *data, struct zxdg_shell_v6 *shell, uint32_t serial)
{
	zxdgcallback_handle_shell_ping((uintptr_t)(data), shell, serial);
}

static const struct zxdg_shell_v6_listener _xdg_shell_listener = {
	_xdg_shell_ping,
};

void
_zxdg_shell_v6_add_listener(void *shell, void *data)
{
    zxdg_shell_v6_add_listener((struct zxdg_shell_v6 *)shell, &_xdg_shell_listener, data);
}

static void
_handle_xdg_surface_configure(void *data, struct zxdg_surface_v6 *surface,
			     uint32_t serial)
{
	zxdgcallback_handle_surface_configure((uintptr_t)data, surface, serial);

}

static const struct zxdg_surface_v6_listener _xdg_surface_listener = {
	_handle_xdg_surface_configure,
};


int
_zxdg_surface_v6_add_listener(void *zxdg_surface_v6, void *data)
{
    return zxdg_surface_v6_add_listener((struct zxdg_surface_v6 *)zxdg_surface_v6, &_xdg_surface_listener, data);
}

void
_handle_xdg_toplevel_configure(void *data,
			  struct zxdg_toplevel_v6 *zxdg_toplevel_v6,
			  int32_t width,
			  int32_t height,
			  struct wl_array *states)
{
	int32_t *pos = states->data;
	int len = states->size;

	zxdgcallback_handle_toplevel_configure((uintptr_t)(data), zxdg_toplevel_v6, width, height, pos,len);
	

}

void
_handle_xdg_toplevel_close(void *data,
		      struct zxdg_toplevel_v6 *zxdg_toplevel_v6)
{
	zxdgcallback_handle_toplevel_close((uintptr_t)(data), zxdg_toplevel_v6);

//      goHandleClose(data);
}

static const struct zxdg_toplevel_v6_listener _zxdg_toplevel_v6_listener = {
_handle_xdg_toplevel_configure,
_handle_xdg_toplevel_close,

};


int
_zxdg_toplevel_v6_add_listener(void *zxdg_toplevel_v6, void *data)
{
	return zxdg_toplevel_v6_add_listener((struct zxdg_toplevel_v6 *)zxdg_toplevel_v6,
			      &_zxdg_toplevel_v6_listener, data);
}

void
_zxdg_toplevel_v6_set_title(void *zxdg_toplevel_v6, void *titl, int len)
{
	char *title = (char*) titl;
	if ((len > 0) && (title[len-1] == 0)) {
		// now it is null terminated
		wl_proxy_marshal((struct wl_proxy *) zxdg_toplevel_v6,
				 ZXDG_TOPLEVEL_V6_SET_TITLE, &title[0]);
	}
}

static void
_handle_frame(void *data, struct wl_callback *callback, uint32_t time)
{
	wlcallback_handle_frame((uintptr_t)(data), callback, time);
}

static const struct wl_callback_listener _frame_listener = {
_handle_frame
};


int
_wl_callback_add_listener(struct wl_callback *wl_callback, void *data)
{
	return wl_callback_add_listener(wl_callback, &_frame_listener, data);
}

void _handle_buffer_release(void *data,
			struct wl_buffer *wl_buffer)
{
	wlcallback_handle_buffer_release((uintptr_t)(data), wl_buffer);
}

static const struct wl_buffer_listener _buffer_listener = {
	_handle_buffer_release

};

int
_wl_buffer_add_listener(struct wl_buffer *wl_buffer, void *data)
{
	return wl_buffer_add_listener(wl_buffer, &_buffer_listener, data);
}



static void
_handle_surface_enter(void *data,
	      struct wl_surface *wl_surface, struct wl_output *wl_output)
{
	wlcallback_handle_surface_enter((uintptr_t)(data), wl_surface, wl_output);
}

static void
_handle_surface_leave(void *data,
	      struct wl_surface *wl_surface, struct wl_output *wl_output)
{
	wlcallback_handle_surface_leave((uintptr_t)(data), wl_surface, wl_output);
}

static const struct wl_surface_listener _surface_listener = {
	_handle_surface_enter,
	_handle_surface_leave
};


int
_wl_surface_add_listener(struct wl_surface *wl_surface, void *data)
{
	return wl_surface_add_listener(wl_surface, &_surface_listener, data);
}

static void
_handle_seat_handle_capabilities(void *data, struct wl_seat *seat,
			 enum wl_seat_capability caps)
{
	uint32_t capabilities = caps;

	wlcallback_handle_seat_capabilities((uintptr_t)(data), seat, capabilities);
}

static void
_handle_seat_handle_name(void *data, struct wl_seat *seat,
		 const char *name)
{
	char *seatname = name; 


	wlcallback_handle_seat_name((uintptr_t)(data), seat, seatname);
}

static const struct wl_seat_listener _seat_listener = {
	_handle_seat_handle_capabilities,
	_handle_seat_handle_name
};

int
_wl_seat_add_listener(struct wl_seat *wl_seat, void *data)
{
	return wl_seat_add_listener(wl_seat, &_seat_listener, data);
}





static void _handle_pointer_enter(void *data, struct wl_pointer *wl_pointer,
		      uint32_t serial,
		      struct wl_surface *surface,
		      wl_fixed_t surface_x,
		      wl_fixed_t surface_y)
{
	wlcallback_handle_pointer_enter((uintptr_t)(data), wl_pointer, serial, surface, surface_x, surface_y);
}
static void _handle_pointer_leave(void *data, struct wl_pointer *wl_pointer,
		      uint32_t serial,
		      struct wl_surface *surface)
{
	wlcallback_handle_pointer_leave((uintptr_t)(data), wl_pointer, serial, surface);
}
static void _handle_pointer_motion(void *data, struct wl_pointer *wl_pointer,
		       uint32_t time,
		       wl_fixed_t surface_x,
		       wl_fixed_t surface_y)
{
	wlcallback_handle_pointer_motion((uintptr_t)(data), wl_pointer, time, surface_x, surface_y);
}
static void _handle_pointer_button(void *data, struct wl_pointer *wl_pointer,
		       uint32_t serial,
		       uint32_t time,
		       uint32_t button,
		       uint32_t state)
{
	wlcallback_handle_pointer_button((uintptr_t)(data), wl_pointer, serial, time, button, state);
}
static void _handle_pointer_axis(void *data, struct wl_pointer *wl_pointer,
		     uint32_t time,
		     uint32_t axis,
		     wl_fixed_t value)
{
	wlcallback_handle_pointer_axis((uintptr_t)(data), wl_pointer, time, axis, value);
}
static void _handle_pointer_frame(void *data, struct wl_pointer *wl_pointer)
{
	wlcallback_handle_pointer_frame((uintptr_t)(data), wl_pointer);
}
static void _handle_pointer_axis_source(void *data, struct wl_pointer *wl_pointer,
			    uint32_t axis_source)
{
	wlcallback_handle_pointer_axis_source((uintptr_t)(data), wl_pointer, axis_source);
}
static void _handle_pointer_axis_stop(void *data, struct wl_pointer *wl_pointer,
			  uint32_t time,
			  uint32_t axis)
{
	wlcallback_handle_pointer_axis_stop((uintptr_t)(data), wl_pointer, time, axis);
}
static void _handle_pointer_axis_discrete(void *data,
			      struct wl_pointer *wl_pointer,
			      uint32_t axis,
			      int32_t discrete)
{
	wlcallback_handle_pointer_axis_discrete((uintptr_t)(data), wl_pointer, axis, discrete);
}


static const struct wl_pointer_listener _pointer_listener = {
	_handle_pointer_enter,
	_handle_pointer_leave,
	_handle_pointer_motion,
	_handle_pointer_button,
	_handle_pointer_axis,
	_handle_pointer_frame,
	_handle_pointer_axis_source,
	_handle_pointer_axis_stop,
	_handle_pointer_axis_discrete
};

int
_wl_pointer_add_listener(struct wl_pointer *wl_pointer, void *data)
{
	return wl_pointer_add_listener(wl_pointer, &_pointer_listener,
					data);
}

static void
_handle_output_geometry(void *data,
			struct wl_output *wl_output,
			int x, int y,
			int physical_width,
			int physical_height,
			int subpixel,
			const char *make,
			const char *model,
			int transform)
{
	char *omake = make;
	char *omodel = model;

	wlcallback_handle_output_geometry((uintptr_t)(data), wl_output, x, y, physical_width, physical_height, subpixel, omake, omodel, transform);
}

static void
_handle_output_done(void *data,
		     struct wl_output *wl_output)
{
	wlcallback_handle_output_done((uintptr_t)(data), wl_output);
}

static void
_handle_output_scale(void *data,
		     struct wl_output *wl_output,
		     int32_t scale)
{
	wlcallback_handle_output_scale((uintptr_t)(data), wl_output, scale);
}

static void
_handle_output_mode(void *data,
		    struct wl_output *wl_output,
		    uint32_t flags,
		    int width,
		    int height,
		    int refresh)
{
	wlcallback_handle_output_mode((uintptr_t)(data), wl_output, flags, width, height, refresh);
}

static const struct wl_output_listener _output_listener = {
	_handle_output_geometry,
	_handle_output_mode,
	_handle_output_done,
	_handle_output_scale
};

int
_wl_output_add_listener(struct wl_output *wl_output, void *data)
{
	return wl_output_add_listener(wl_output, &_output_listener, data);
}
