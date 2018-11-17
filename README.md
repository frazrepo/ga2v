# Introduction 

A small tool to convert an audio to a video using a static background image.

`
Usage : ga2v -a audiodir/ -i bgimage.jpg -o outputdir/
`

# Installation

**Installation - Precompiled binaries**

Linux

` $ wget https://github.com/frazrepo/ga2v/releases/download/v1.0/ga2v -O ga2v && chmod +x ga2v`


Windows

` $ wget https://github.com/frazrepo/ga2v/releases/download/v1.0/ga2v.exe -O ga2v.exe`

**Software dependencies**

* [ffmpeg](https://www.ffmpeg.org/)


# Build and Test from source

`go run ga2v.go -a assets/ -i src.jpg -o assets/mp4
`
# Building binaries for Windows and Linux

Windows 

`
GOOS=windows GOARCH=386 go build -o ga2v.exe ga2v.go
`

Linux : 

`
go build
`