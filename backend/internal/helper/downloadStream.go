package helper

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func StreamFileToClient(w http.ResponseWriter, reader io.Reader, filename string, size int64) error {
	w.Header().Set("Content-Disposition",fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	if size > 0 {
		w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	}
	_, err := io.Copy(w, reader)
	return err
}
