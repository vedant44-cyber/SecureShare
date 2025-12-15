package helper

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"
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
