FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN apt-get update
RUN apt-get --assume-yes install dmz-cursor-theme fonts-dejavu libxkbcommon-dev libx11-data libwayland-dev libvulkan-dev libglvnd-core-dev

RUN go get -d -v ./...
RUN go install -v ./go-wayland-simple-shm
RUN go install -v ./go-wayland-smoke
RUN go install -v ./go-wayland-imageviewer
RUN go install -v ./go-wayland-texteditor
RUN go install -v ./go-wayland-texteditor/editor_backend
RUN go install -tags wayland -v ./go-wayland-cube

RUN echo "#!/bin/sh" > /go/bin/texteditor
RUN echo "cd /go/src/app/go-wayland-texteditor" >> /go/bin/texteditor
RUN echo "/go/bin/editor_backend & /go/bin/go-wayland-texteditor" >> /go/bin/texteditor
RUN chmod +x /go/bin/texteditor

RUN apt-get --assume-yes install unzip build-essential autoconf libtool libxext-dev libx11-dev x11proto-gl-dev
RUN wget https://github.com/NVIDIA/libglvnd/archive/refs/tags/v1.3.4.zip
RUN unzip v1.3.4.zip
RUN cd libglvnd-1.3.4 && ./autogen.sh && ./configure && make && make install

RUN apt-get --assume-yes install unzip build-essential autoconf libtool wayland-protocols eglexternalplatform-dev libwayland-dev libwayland-egl-backend-dev libdrm-dev
RUN wget https://github.com/NVIDIA/egl-wayland/archive/refs/tags/1.1.11.zip
RUN unzip 1.1.11.zip
RUN cd egl-wayland-1.1.11 && ./autogen.sh && make && make install

