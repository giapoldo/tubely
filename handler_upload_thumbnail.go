package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here

	const maxMemory = 10 << 20
	err = r.ParseMultipartForm(maxMemory)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some parse error", err)
		return
	}

	thumbnailFile, multipartHeader, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some parse error", err)
		return
	}
	defer thumbnailFile.Close()

	mediaType := multipartHeader.Header.Get("Content-Type")

	// fileData, err := io.ReadAll(file)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Some parse error", err)
	// 	return
	// }

	videoDbEntry, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Video metadata fetch failed", err)
		return
	}
	if videoDbEntry.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "User is not the video's owner", err)
		return
	}

	mType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Bad mime type", err)
		return
	}
	if mType != "image/png" && mType != "image/jpeg" {
		respondWithError(w, http.StatusBadRequest, "File type not allowed", err)
		return
	}

	ext, err := mime.ExtensionsByType(mediaType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad file type", err)
		return
	}

	key := make([]byte, 32)
	rand.Read(key)

	fileName := base64.RawURLEncoding.EncodeToString(key)

	thumbnail_fileName := fmt.Sprintf("%s%s", fileName, ext[0])
	thumbnail_filePath := filepath.Join(cfg.assetsRoot, thumbnail_fileName)
	file, err := os.Create(thumbnail_filePath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't upload file", err)
		return
	}
	_, err = io.Copy(file, thumbnailFile)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create thumbnail file", err)
		return
	}

	// s := base64.StdEncoding.EncodeToString(fileData)

	// thisThumbnail := thumbnail{
	// 	data:      fileData,
	// 	mediaType: mediaType,
	// }
	// videoThumbnails[videoID] = thisThumbnail
	// s = fmt.Sprintf("http://localhost:%s/api/thumbnails/%s", cfg.port, videoIDString)

	// s = fmt.Sprintf("data:%s;base64,%s", mediaType, s)

	s := fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, thumbnail_fileName)

	videoDbEntry.ThumbnailURL = &s

	err = cfg.db.UpdateVideo(videoDbEntry)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User is not the video's owner", err)
	}

	respondWithJSON(w, http.StatusOK, videoDbEntry)
}
