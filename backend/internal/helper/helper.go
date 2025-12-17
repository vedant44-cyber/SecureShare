package helper

import (
	"encoding/json"
	"errors"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func ErrorHandler(err error) {
	if err != nil {
		log.Printf("%v\n", err.Error())
	}
}

var filenameSafeRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

func SanitizeFilename(name string) string {
	// Remove any path (../, /, \)
	name = filepath.Base(name)

	// Trim spaces
	name = strings.TrimSpace(name)

	// Replace unsafe chars
	name = filenameSafeRegex.ReplaceAllString(name, "_")

	// Prevent empty filename
	if name == "" || name == "." || name == ".." {
		return "file"
	}

	// Enforce max length
	if len(name) > 255 {
		name = name[:255]
	}

	return name
}

func SanitizeFileID(id string) (string, error) {
	// Basic validation for UUID format
	u, err := uuid.Parse(id)
	if err != nil {
		return "", err
	}
	if u.Version() != 4 {
		return "", errors.New("unsupported uuid version")
	}
	return u.String(), nil
}

func ParseFileMetaJSON(metaJSON string) (FileMeta, error) {
	var meta FileMeta
	err := json.Unmarshal([]byte(metaJSON), &meta)
	if err != nil {
		return FileMeta{}, err
	}
	return meta, nil
}
