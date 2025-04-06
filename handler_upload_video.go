package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {

	maxMemory := int64(10 << 30)
	r.Body = http.MaxBytesReader(w, r.Body, maxMemory)

	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Video metadata fetch failed", err)
		return
	}
	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "User is not the video's owner", err)
		return
	}

	// TODO: implement the upload here

	videoFile, multipartHeader, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to parse multipart request", err)
		return
	}
	defer videoFile.Close()

	mediaType := multipartHeader.Header.Get("Content-Type")

	mType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid MIME type", err)
		return
	}

	if mType != "video/mp4" {
		respondWithError(w, http.StatusBadRequest, "File type not allowed", err)
		return
	}

	tempFilename := "tubely-upload.mp4"
	localUploadFile, err := os.CreateTemp("", tempFilename)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file", err)
		return
	}

	_, err = io.Copy(localUploadFile, videoFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "File copy failed", err)
		return
	}
	localUploadFile.Seek(0, io.SeekStart)

	aspectRatio, err := getVideoAspectRatio(localUploadFile.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch aspect ratio", err)
		return
	}

	procesedVideoPath, err := processVideoForFastStart(localUploadFile.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't process video for fast start", err)
		return

	}

	localUploadFile.Close()
	os.Remove(localUploadFile.Name())

	processedVideo, err := os.Open(procesedVideoPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't open video for upload", err)
		return
	}

	key := make([]byte, 32)
	rand.Read(key)

	var objectKey string
	switch aspectRatio {
	case "16:9":
		objectKey = fmt.Sprintf("landscape/%s.mp4", hex.EncodeToString(key))
	case "9:16":
		objectKey = fmt.Sprintf("portrait/%s.mp4", hex.EncodeToString(key))
	case "other":
		objectKey = fmt.Sprintf("other/%s.mp4", hex.EncodeToString(key))
	}

	putObjectInput := s3.PutObjectInput{
		Bucket:      &cfg.s3Bucket,
		Key:         &objectKey,
		Body:        processedVideo,
		ContentType: &mType,
	}

	_, err = cfg.s3Client.PutObject(context.TODO(), &putObjectInput)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to upload video to fileserver", err)
		return
	}

	s := fmt.Sprintf("%s/%s", cfg.s3CfDistribution, objectKey)
	video.VideoURL = &s

	err = cfg.db.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update video database", err)
		return
	}

	// presignedVideo, err := cfg.dbVideoToSignedVideo(video)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Unable to presign video url", err)
	// 	return
	// }

	respondWithJSON(w, http.StatusOK, video)
}
