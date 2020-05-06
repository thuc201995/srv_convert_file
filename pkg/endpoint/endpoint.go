package endpoint

import (
	"context"
	"io"

	endpoint "github.com/go-kit/kit/endpoint"
	service "github.com/thuc201995/srv_convert_file/pkg/service"
)

// WordToPDFRequest collects the request parameters for the WordToPDF method.
type WordToPDFRequest struct {
	File     io.Reader
	Name     string
	MimeType string
}

// WordToPDFResponse collects the response parameters for the WordToPDF method.
type WordToPDFResponse struct {
	Link string `json:"link"`
	Err  error  `json:"err"`
}

// MakeWordToPDFEndpoint returns an endpoint that invokes WordToPDF on the service.
func MakeWordToPDFEndpoint(s service.SrvConvertFileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(WordToPDFRequest)

		rs, err := s.WordToPDF(ctx, req.File, req.Name, req.MimeType)
		return WordToPDFResponse{
			Err:  err,
			Link: rs,
		}, nil
	}
}

// Failed implements Failer.
func (r WordToPDFResponse) Failed() error {
	return r.Err
}

// Failure is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failure interface {
	Failed() error
}
