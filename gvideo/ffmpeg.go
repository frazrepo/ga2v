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

package gvideo

import (
	"fmt"
	os_exec "os/exec"
	"runtime"

	log "github.com/sirupsen/logrus"
)

const isDebug = true

// Convert an audio file to video using a static background image
func Convert(image string, audio string, out string) error {
	template := "ffmpeg -loop 1 -i %s -y -i %s -c:v libx264 -c:a copy -shortest %s"

	if runtime.GOOS == "windows" {
		return exec(template, image, audio, out)
	}

	//Add additional formatting for non-windows host
	return exec(template, "\""+image+"\"", "\""+audio+"\"", "\""+out+"\"")
}
func exec(template string, params ...interface{}) error {
	cmd := fmt.Sprintf(template, params...)
	log.Debug("Running command : " + cmd)

	interpreter := ""
	interpreterArgs := ""

	if runtime.GOOS == "windows" {
		interpreter = "cmd"
		interpreterArgs = "/C"
	} else {
		interpreter = "sh"
		interpreterArgs = "-c"
	}

	c, err := os_exec.Command(interpreter, interpreterArgs, cmd).CombinedOutput()

	log.Debug(fmt.Sprintf("ffmpeg out:\n%s\n", string(c)))

	return err
}
