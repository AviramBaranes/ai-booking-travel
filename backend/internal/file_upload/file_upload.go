// Package fileupload holds helpers for reading files out of multipart/form-data
// requests, shared by the raw endpoints that accept spreadsheet uploads.
package fileupload

import (
	"mime/multipart"
	"net/http"
	"strings"

	"encore.app/internal/api_errors"
)

var (
	ErrInvalidContentType = api_errors.NewValidationError("invalid content type: expected multipart/form-data")
	ErrParseMultipartForm = api_errors.NewValidationError("failed to parse multipart form")
	ErrGetFileFromForm    = api_errors.NewValidationError("failed to get file from form data")
)

// maxMemory is how much of the upload is buffered in memory; the remainder is spilled to
// temporary files by ParseMultipartForm.
const maxMemory = 2 << 20 // 2 MB

// ExtractFile parses a multipart/form-data request and returns the file uploaded under the
// "file" form field. The caller owns the returned file and must close it.
func ExtractFile(req *http.Request) (multipart.File, error) {
	ct := req.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "multipart/form-data") {
		return nil, ErrInvalidContentType
	}

	if err := req.ParseMultipartForm(maxMemory); err != nil {
		return nil, ErrParseMultipartForm
	}

	file, _, err := req.FormFile("file")
	if err != nil {
		return nil, ErrGetFileFromForm
	}

	return file, nil
}
