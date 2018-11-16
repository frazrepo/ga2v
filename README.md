# Introduction 

A small tool to convert an audio to a video using a static background image.

# Getting Started

1.	Installation process

* Just download the [precompiled binary](https://github.com/frazrepo/ga2v/releases).

2.	Software dependencies

* ffmpeg

3.	Latest releases

# Build and Test

`go run ga2v.go -d assets/ -i src.jpg -o assets/mp4
`
# Compiling for Windows and Linux

Windows 

`
GOOS=windows GOARCH=386 go build -o ga2v.exe ga2v.go
`

Linux : 

`
go build
`