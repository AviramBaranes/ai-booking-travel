package fileupload

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"encore.app/internal/api_errors"
)

func TestExtractFile(t *testing.T) {
	t.Run("returns error for non-multipart content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/locations/hertz", nil)
		req.Header.Set("Content-Type", "application/json")

		_, err := ExtractFile(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, ErrInvalidContentType, err)
	})

	t.Run("returns error when no file field in form", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/locations/hertz", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		_, err := ExtractFile(req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		api_errors.AssertApiError(t, ErrGetFileFromForm, err)
	})

	t.Run("extracts file successfully", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "locations.csv")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		part.Write([]byte("test file content"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/locations/hertz", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		file, err := ExtractFile(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer file.Close()

		buf := make([]byte, 100)
		n, _ := file.Read(buf)
		if string(buf[:n]) != "test file content" {
			t.Fatalf("expected file content 'test file content', got %q", string(buf[:n]))
		}
	})
}
