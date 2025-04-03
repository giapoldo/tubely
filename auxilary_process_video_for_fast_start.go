package main

import (
	"fmt"
	"log"
	"os/exec"
)

func processVideoForFastStart(filePath string) (string, error) {

	processedFilePath := fmt.Sprintf("%s.processing", filePath)

	args := []string{"-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", processedFilePath}
	command := exec.Command("ffmpeg", args...)

	err := command.Run()
	if err != nil {
		log.Fatal("Couldn't put moov atom at the start")
	}

	fmt.Println(processedFilePath)

	return processedFilePath, nil
}
