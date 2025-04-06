package main

// import (
// 	"strings"

// 	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
// )

// func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {

// 	if video.VideoURL == nil {
// 		return video, nil
// 	}

// 	fields := strings.Split(*video.VideoURL, ",")

// 	s, err := generatePresignedURL(cfg.s3Client, fields[0], fields[1], cfg.presignTimeout)
// 	if err != nil {
// 		return database.Video{}, err
// 	}
// 	video.VideoURL = &s

// 	return video, nil
// }
