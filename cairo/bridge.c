#include <cairo.h>

#include "bridge.h"
#include "_cgo_export.h"

cairo_user_data_key_t key;

void _cairo_destroy_func (void *data)
{
	cairocallback_cairo_destroy_func(data);
}


cairo_status_t
_cairo_surface_set_user_data(cairo_surface_t		 *surface,
			     void			 *user_data)
{
	return cairo_surface_set_user_data(surface, &key, user_data, &_cairo_destroy_func);
}

