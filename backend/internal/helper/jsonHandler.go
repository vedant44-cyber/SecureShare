package helper

import "encoding/json"

type FileMeta struct {
	S3Key         string `json:"s3_key"`
	Size          int64  `json:"size"`
	UploadedAt    int64  `json:"uploaded_at"`
	TTLHours      int    `json:"ttl_hours"`
	DownloadLimit int    `json:"download_limit"`
	Filename      string `json:"filename"`
}
type UploadResponse struct {
	FileID        string `json:"file_id"`
	Size          int64  `json:"size"`
	TTLHours      int    `json:"ttl_hours"`
	DownloadLimit int    `json:"download_limit"`
	Filename      string `json:"filename"`
	Message       string `json:"message"`
}

func BuildFileMetaJSON(s3Key string, Size int64, uploadedAt int64, ttl int, downloadLimit int, filename string) (string, error) {
	metaObj := FileMeta{
		S3Key:         s3Key,
		Size:          Size,
		UploadedAt:    uploadedAt,
		TTLHours:      ttl,
		DownloadLimit: downloadLimit,
		Filename:      filename,
	}
	metaBytes, err := json.Marshal(metaObj)
	if err != nil {
		return "", err
	}
	return string(metaBytes), nil
}

func BuildUploadResponse(fileID string, size int64, ttl int, downloadLimit int, filename string) UploadResponse {
	return UploadResponse{
		FileID:        fileID,
		Size:          size,
		TTLHours:      ttl,
		DownloadLimit: downloadLimit,
		Filename:      filename,
		Message:       "file uploaded successfully",
	}
}
