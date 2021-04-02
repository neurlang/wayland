FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN apt-get update
RUN apt-get --assume-yes install dmz-cursor-theme fonts-dejavu libxkbcommon-dev libx11-data
RUN go get -d -v ./...
RUN go install -v ./go-wayland-simple-shm
RUN go install -v ./go-wayland-smoke
RUN go install -v ./go-wayland-imageviewer

