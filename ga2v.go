// Copyright (c) 2018, github.com/frazrepo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"io/ioutil"
	"os"

	"./gimage"
	"./gvideo"

	log "github.com/sirupsen/logrus"
)

// Extension Constants
const (
	AudioExt = ".mp3"
	VideoExt = ".mp4"
)

// Directory constants
const (
	OutDirName = "mp4"
)

func main() {

	var (
		audio           string
		backgroundImage string
	)
	log.SetLevel(log.InfoLevel)

	flag.StringVar(&audio, "a", "", "Audio file or directory")
	flag.StringVar(&backgroundImage, "i", "", "Background image")

	flag.Usage = usage

	flag.Parse()

	if audio == "" || backgroundImage == "" {
		fmt.Println("Please check options")
		flag.Usage()
		os.Exit(1)
	}

	target := audio
	bgImage := backgroundImage

	log.Info("Audio file or directory : ", target)
	log.Info("Input background image : ", bgImage)
	log.Info("Host : ", runtime.GOOS)

	tmpDir, e := ioutil.TempDir("", "gaudio2video")
	if e != nil {
		log.Error("Error creating TempDir")
		panic(e)
	}

	defer os.RemoveAll(tmpDir)

	log.Info("TempDir : ", tmpDir)

	if stat, err := os.Stat(target); err == nil && stat.IsDir() {
		files, _ := ioutil.ReadDir(target)
		for _, fn := range files {

			//Process only mp3 files
			ext := path.Ext(fn.Name())
			if ext != AudioExt {
				continue
			}

			abs, _ := filepath.Abs(target)

			aPath := path.Join(abs, fn.Name())

			log.Info("Processing " + aPath + " ...")
			processFile(aPath, bgImage, tmpDir)
		}
	} else {
		log.Info("Processing " + target + " ...")
		processFile(target, bgImage, tmpDir)
	}

	log.Info("Done!")
}

func processFile(fullAudioPath string, bgImage string, tmpDir string) {

	title := filenameWithoutExtension(filepath.Base(fullAudioPath))
	bgImageFileName := gimage.GenerateImageWithText(bgImage, title, tmpDir)

	audioFileName := fullAudioPath
	log.Debug("Full audio path : ", audioFileName)

	dir, _ := filepath.Split(fullAudioPath)
	outDir := createOutDir(dir)

	outputFileName := path.Join(outDir, title+VideoExt)
	err := gvideo.Convert(bgImageFileName, audioFileName, outputFileName)
	if err != nil {
		log.Error("Error processing ", fullAudioPath, err)
	}
}

func filenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func createOutDir(in string) string {
	oPath := path.Join(in, OutDirName)

	if _, err := os.Stat(oPath); os.IsNotExist(err) {
		os.Mkdir(oPath, os.ModePerm)
	}

	return oPath
}

var usageStr = `

██████╗  █████╗ ██████╗ ██╗   ██╗
██╔════╝ ██╔══██╗╚════██╗██║   ██║
██║  ███╗███████║ █████╔╝██║   ██║
██║   ██║██╔══██║██╔═══╝ ╚██╗ ██╔╝
╚██████╔╝██║  ██║███████╗ ╚████╔╝ 
 ╚═════╝ ╚═╝  ╚═╝╚══════╝  ╚═══╝  
                                  
Usage: ga2v [options]
Required Options:
    -a, --audio <path>               Audio file or audio direcory
    -i, --image <path>               Static background image
Common Options:
    -h, --help                       Show this message
`

// usage will print out the flag options for the server.
func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}
