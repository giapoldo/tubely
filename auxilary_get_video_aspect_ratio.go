package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {

	args := []string{"-v", "error", "-print_format", "json", "-show_streams", filePath}
	command := exec.Command("ffprobe", args...)
	var buffer bytes.Buffer
	command.Stdout = &buffer

	err := command.Run()
	if err != nil {
		log.Fatal("Couldn't get aspect ratio")
	}

	type ffprobeResult struct {
		Streams []struct {
			Width  int64 `json:"width"`
			Height int64 `json:"height"`
		} `json:"streams"`
	}

	var aspectRatioJSON ffprobeResult

	err = json.Unmarshal(buffer.Bytes(), &aspectRatioJSON)
	if err != nil {
		return "", nil
	}

	ar169 := 16.0 / 9.0
	ar916 := 9.0 / 16.0

	aspectRatio := float64(aspectRatioJSON.Streams[0].Width) / float64(aspectRatioJSON.Streams[0].Height)

	if aspectRatio >= ar169*0.9 && aspectRatio <= ar169*1.1 {
		return "16:9", nil
	} else if aspectRatio >= ar916*0.9 && aspectRatio <= ar916*1.1 {
		return "9:16", nil
	} else {
		return "other", nil
	}
}
