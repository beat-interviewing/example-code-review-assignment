package beatly

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/endpoint"
	transport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

// Handler returns an http.Handler capable of multiplexing the various endpoints
// made available by the BEAT.ly service.
func Handler(s Service) http.Handler {

	r := mux.NewRouter()

	r.Handle("/link", transport.NewServer(
		CreateEndpoint(s),
		DecodeCreateRequest(),
		EncodeCreateResponse(),
		transport.ServerErrorEncoder(ErrorEncoder()),
	)).Methods("POST")

	r.Handle("/{id}", transport.NewServer(
		VisitEndpoint(s),
		DecodeVisitRequest(),
		EncodeVisitResponse(),
		transport.ServerErrorEncoder(ErrorEncoder()),
	)).Methods("GET")

	l := log.NewLogfmtLogger(os.Stdout)

	return LoggingMiddleware(l)(r)
}

func ErrorEncoder() transport.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			E string `json:"error"`
		}{
			err.Error(),
		})
	}
}

// CreateEndpoint encapsulates the business logic of creating a shortened link.
//
// This function will be passed a CreateRequest as an argument and will return a
// CreateResponse as a result.
func CreateEndpoint(service Service) endpoint.Endpoint {
	return func(_ context.Context, v interface{}) (interface{}, error) {

		req := v.(*CreateRequest)

		link := &Link{
			Target:   req.Target,
			Redirect: req.Redirect,
		}

		err := service.Create(link)
		if err != nil {
			return nil, err
		}

		return &CreateResponse{
			ID:       link.IDHash,
			URL:      fmt.Sprintf("https://beat.ly/%s", link.IDHash),
			Target:   link.Target,
			Redirect: link.Redirect,
		}, nil
	}
}

type CreateRequest struct {
	Target   string `json:"target"`
	Redirect int    `json:"redirect"`
}

func DecodeCreateRequest() transport.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		v := &CreateRequest{}
		err := json.NewDecoder(r.Body).Decode(v)
		return v, err
	}
}

type CreateResponse struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Target   string `json:"target"`
	Redirect int    `json:"redirect"`
}

func (r *CreateResponse) StatusCode() int {
	return http.StatusCreated
}

func EncodeCreateResponse() transport.EncodeResponseFunc {
	return transport.EncodeJSONResponse
}

// VisitEndpoint encapsulates the business logic of visiting a shortened link.
// Each visit is counted for analytics purposes.
//
// This function will be passed a VisitRequest as an argument and will return a
// VisitResponse as a result.
func VisitEndpoint(service Service) endpoint.Endpoint {
	return func(_ context.Context, req interface{}) (interface{}, error) {

		link, err := service.Visit(req.(*VisitRequest).ID)
		if err != nil {
			return nil, err
		}

		return &VisitResponse{
			Target:   link.Target,
			Redirect: link.Redirect,
		}, nil
	}
}

type VisitRequest struct {
	ID string
}

func DecodeVisitRequest() transport.DecodeRequestFunc {
	return func(_ context.Context, req *http.Request) (interface{}, error) {
		m := mux.Vars(req)
		if id, ok := m["id"]; ok {
			return &VisitRequest{ID: id}, nil
		}
		return nil, fmt.Errorf("invalid request")
	}
}

type VisitResponse struct {
	Target   string
	Redirect int
}

func EncodeVisitResponse() transport.EncodeResponseFunc {
	return func(_ context.Context, w http.ResponseWriter, req interface{}) error {
		r := req.(*VisitResponse)
		w.Header().Set("Location", r.Target)
		w.WriteHeader(r.Redirect)
		return nil
	}
}
