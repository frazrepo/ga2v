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
