#!/bin/bash

if [ ! -v XDG_RUNTIME_DIR ]
then
	echo "XDG_RUNTIME_DIR not set!"
	exit -1
fi

wld="$WAYLAND_DISPLAY"

if [ ! -f $XDG_RUNTIME_DIR/$WAYLAND_DISPLAY  ]
then
	for i in {9..0}
	do
		if [ -f "$XDG_RUNTIME_DIR/wayland-$i" ]
		then
			wld="wayland-$i"
			break
		fi
	done
fi

docker run --gpus='all,"capabilities=compute,utility,graphics,display"' \
           -e DISPLAY='' \
           -e XDG_RUNTIME_DIR=/tmp \
           -e WAYLAND_DISPLAY=$WAYLAND_DISPLAY \
           -v $XDG_RUNTIME_DIR/$wld:/tmp/$wld  \
           --user=$(id -u):$(id -g) \
           --group-add=122 \
           --device /dev/dri:/dev/dri \
           --device /dev/nvidiactl:/dev/nvidiactl \
           --device /dev/nvidia0:/dev/nvidia0 \
           --mount type=bind,source=/usr/share/glvnd/egl_vendor.d,target=/usr/share/glvnd/egl_vendor.d \
           --mount type=bind,source=/usr/share/egl/egl_external_platform.d,target=/usr/share/egl/egl_external_platform.d \
           --mount type=bind,source=/usr/share/vulkan/icd.d,target=/usr/share/vulkan/icd.d \
           --mount type=bind,source=/usr/share/drirc.d,target=/usr/share/drirc.d \
           --mount type=bind,source=/usr/share/nvidia,target=/usr/share/nvidia \
           --mount type=bind,source=/usr/share/vulkan/implicit_layer.d,target=/usr/share/vulkan/implicit_layer.d \
           --mount type=bind,source=/usr/share/vulkan/explicit_layer.d,target=/usr/share/vulkan/explicit_layer.d \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libnvidia-egl-gbm.so.1,target=/lib/x86_64-linux-gnu/libnvidia-egl-gbm.so.1 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libEGL_nvidia.so.0,target=/lib/x86_64-linux-gnu/libEGL_nvidia.so.0 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libGLX_nvidia.so.0,target=/lib/x86_64-linux-gnu/libGLX_nvidia.so.0 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libnvidia-glsi.so.515.65.01,target=/lib/x86_64-linux-gnu/libnvidia-glsi.so.515.65.01 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libnvidia-tls.so.515.65.01,target=/lib/x86_64-linux-gnu/libnvidia-tls.so.515.65.01 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libnvidia-glcore.so.515.65.01,target=/lib/x86_64-linux-gnu/libnvidia-glcore.so.515.65.01 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libnvidia-eglcore.so.515.65.01,target=/lib/x86_64-linux-gnu/libnvidia-eglcore.so.515.65.01 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libXext.so.6,target=/lib/x86_64-linux-gnu/libXext.so.6 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libEGL.so.1,target=/lib/x86_64-linux-gnu/libEGL.so.1 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libGLdispatch.so.0,target=/lib/x86_64-linux-gnu/libGLdispatch.so.0 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libnvidia-glvkspirv.so.515.65.01,target=/lib/x86_64-linux-gnu/libnvidia-glvkspirv.so.515.65.01 \
           --mount type=bind,source=/lib/x86_64-linux-gnu/libnvidia-vulkan-producer.so,target=/lib/x86_64-linux-gnu/libnvidia-vulkan-producer.so \
           go-wayland-demos /go/bin/go-wayland-cube

exit 0


