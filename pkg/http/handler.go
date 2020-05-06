package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	http1 "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	endpoint "github.com/thuc201995/srv_convert_file/pkg/endpoint"
)

const (
	MB = 1 << 20
)

var (
	errEmptyFile        = errors.New("file is required")
	errMaxFileSize      = errors.New("max file upload size is 5MB")
	errFileNotSupported = errors.New("file not supported")
	errInvalidFileName  = errors.New("invalid file name")
	acceptedFileType    = []string{
		"application/doc",
		"application/ms-doc",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}
)

// makeWordToPDFHandler creates the handler logic
func makeWordToPDFHandler(m *mux.Router, endpoints endpoint.Endpoints, options []http1.ServerOption) {
	m.Methods("POST").Path("/word-to-pdf").Handler(http1.NewServer(endpoints.WordToPDFEndpoint, decodeWordToPDFRequest, encodeWordToPDFResponse, options...))
	m.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./public/pdf"))))
}

// decodeWordToPDFRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body.
func decodeWordToPDFRequest(_ context.Context, r *http.Request) (interface{}, error) {
	req := endpoint.WordToPDFRequest{}
	file, handler, err := r.FormFile("file")
	var maxFileSize int64 = 5 * MB
	r.ParseMultipartForm(maxFileSize)

	if err != nil {
		return req, errEmptyFile
	}

	if handler.Size > maxFileSize {
		return req, errMaxFileSize
	}
	var mimeType = handler.Header.Get("Content-Type")
	_, found := Find(acceptedFileType, mimeType)

	if found == false {
		return req, errFileNotSupported
	}
	req = endpoint.WordToPDFRequest{
		File:     file,
		Name:     handler.Filename,
		MimeType: mimeType,
	}

	return req, err
}

// encodeWordToPDFResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func encodeWordToPDFResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(response)
	return nil
}
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
func ErrorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

// This is used to set the http status, see an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/pkg/addtransport/http.go#L133
func err2code(err error) int {
	switch err {
	case errMaxFileSize, errEmptyFile, errFileNotSupported, errInvalidFileName:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
