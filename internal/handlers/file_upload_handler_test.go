package handlers

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"simple-file-processor/internal/mocks/mockdb"
	"simple-file-processor/internal/mocks/mocktasks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	hKey = "FileUploadHandler"
)

func ResponseRecorder() *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	return rec
}

func MultiPartFormRequest(t *testing.T, fieldname string) *http.Request {
	// Create a multipart form file
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	// Create a file in the multipart form
	f, err := w.CreateFormFile(fieldname, "test.txt")
	if err != nil {
		t.Fatal(err)
		return nil
	}

	c := bytes.NewBufferString("This is a test file")
	_, _ = io.Copy(f, c)
	_ = w.Close() // Close the writer to finalize the multipart form

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func MultiPartFormImageRequest(t *testing.T, fieldname string) *http.Request {
	// Create a new RGBA image using image dimensions
	width := 200
	height := 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Create a new body buffer
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	// Create a file to save the image
	f, err := w.CreateFormFile(fieldname, "test.png")
	if err != nil {
		t.Fatal(err)
	}

	// Encode the image into the form file
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}

	// Close the writer to finalize the multipart form
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

// Verifies that when the file is too large, the handler returns a 413 status co
func Test_FileUploadHandler_WhenFileTooLarge_Expect413(t *testing.T) {
	rec := ResponseRecorder()
	db := new(mockdb.Database)
	ac := new(mocktasks.Client)
	req := httptest.NewRequest("POST", "/upload", nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	req.ContentLength = 1000000000   // 1GB
	req.ParseMultipartForm(10 << 20) // 10MB limit
	h := NewHandlers(log, db, ac)
	h.GetHandler(hKey)(rec, req)
	assert.Equal(t, rec.Code, 413)
	os.RemoveAll("uploads") // clean up
}

// Verifies that when field name is incorrect, the handler returns a 400 status code
func Test_FileUploadHandler_WhenFieldNameIncorrect_Expect400(t *testing.T) {
	rr := ResponseRecorder()
	fn := "test-file" // incorrect field name
	db := new(mockdb.Database)
	ac := new(mocktasks.Client)
	req := MultiPartFormRequest(t, fn)
	hand := NewHandlers(log, db, ac)
	http.HandlerFunc(hand.GetHandler(hKey)).ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 400)
	os.RemoveAll("uploads") // clean up
}

// Verifies that when the file is successfully uploaded and the database returns an error, the handler returns a 500 status code
func Test_FileUploadHandler_WhenFileUploaded_WhenErrorSavingMetadata_Expect500(t *testing.T) {
	rr := ResponseRecorder()
	db := new(mockdb.Database)
	ac := new(mocktasks.Client)
	fn := "file"
	req := MultiPartFormRequest(t, fn)
	hand := NewHandlers(log, db, ac)
	db.On("InsertFileMetadata", mock.Anything).Return(errors.New("error saving metadata"))
	http.HandlerFunc(hand.GetHandler(hKey)).ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, 500)
	os.RemoveAll("uploads") // clean up
}

// Verifies that when an image file is successfully uploaded, the handler returns a 200 status code
// and a new image resize task is created
func Test_FileUploadHandler_WhenFileImageUploaded_WhenSuccess_ExpectNewImageResizeTask(t *testing.T) {
	rr := ResponseRecorder()
	db := new(mockdb.Database)
	ac := new(mocktasks.Client)
	fn := "file"
	req := MultiPartFormImageRequest(t, fn)
	h := NewHandlers(log, db, ac)

	// Mock the database call
	db.On("InsertFileMetadata", mock.Anything).Return(nil)

	// Mock the async client enqueue call
	ac.On("Enqueue", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	// Call the handler
	http.HandlerFunc(h.GetHandler(hKey)).ServeHTTP(rr, req)

	// verify the response
	assert.Equal(t, rr.Code, 200)

	// Verify that the image resize task was created
	ac.AssertCalled(t, "Enqueue", mock.Anything, mock.Anything, mock.Anything)
	os.RemoveAll("uploads") // clean up
}
