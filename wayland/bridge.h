int
_wl_shm_add_listener(struct wl_shm *wl_shm, void *data);


int
_wl_registry_add_listener(struct wl_registry *wl_registry, void *data);

void
_zxdg_shell_v6_add_listener(void *shell, void *data);


int
_zxdg_surface_v6_add_listener(void *zxdg_surface_v6, void *data);

int
_zxdg_toplevel_v6_add_listener(void *zxdg_toplevel_v6, void *data);

void
_zxdg_toplevel_v6_set_title(void *zxdg_toplevel_v6, void *titl, int len);

int
_wl_callback_add_listener(struct wl_callback *wl_callback, void *data);

int
_wl_buffer_add_listener(struct wl_buffer *wl_buffer, void *data);

int
_wl_surface_add_listener(struct wl_surface *wl_surface, void *data);

int
_wl_seat_add_listener(struct wl_seat *wl_seat, void *data);
int
_wl_pointer_add_listener(struct wl_pointer *wl_pointer, void *data);

int
_wl_output_add_listener(struct wl_output *wl_output, void *data);
