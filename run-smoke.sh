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

docker run -e XDG_RUNTIME_DIR=/tmp \
           -e WAYLAND_DISPLAY=$WAYLAND_DISPLAY \
           -v $XDG_RUNTIME_DIR/$wld:/tmp/$wld  \
           --user=$(id -u):$(id -g) \
           go-wayland-demos /go/bin/go-wayland-smoke
