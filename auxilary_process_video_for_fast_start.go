package main

import (
	"fmt"
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {

	processedFilePath := fmt.Sprintf("%s.processing", filePath)

	args := []string{"-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", processedFilePath}
	command := exec.Command("ffmpeg", args...)

	err := command.Run()
	if err != nil {
		return "", nil
	}

	return processedFilePath, nil
}
