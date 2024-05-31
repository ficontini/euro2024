package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	matchendpoint "github.com/ficontini/euro2024/matchservice/pkg/endpoint"
	"github.com/ficontini/euro2024/matchservice/pkg/service"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(endpoints matchendpoint.Set) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeHTTPError),
	}
	m := http.NewServeMux()
	m.Handle("/matches/upcoming", httptransport.NewServer(
		endpoints.GetUpcomingMatchesEndpoint,
		decodeHTTPRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/matches/live", httptransport.NewServer(
		endpoints.GetLiveMatchesEndpoint,
		decodeHTTPRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return m
}

func NewHTTPClient(instance string) (service.Service, error) {
	var (
		options          = []httptransport.ClientOption{}
		upcomingEndpoint endpoint.Endpoint
	)
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	{
		upcomingEndpoint = httptransport.NewClient(
			http.MethodGet,
			copyURL(u, "/matches/upcoming"),
			encodeHTTPGenericRequest,
			decodeHTTPUpcomingResponse,
			options...,
		).Endpoint()
	}
	return matchendpoint.Set{
		GetUpcomingMatchesEndpoint: upcomingEndpoint,
	}, nil
}

func decodeHTTPRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	if r.Method != http.MethodGet {
		return nil, ErrNotAllowed()
	}
	return struct{}{}, nil
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, v interface{}) error {
	if err, ok := v.(Error); ok {
		return WriteResponse(w, err.Code, map[string]string{"err": err.Err})
	}
	return WriteResponse(w, http.StatusOK, v)
}
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = io.NopCloser(&buf)
	return nil
}
func decodeHTTPUpcomingResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp matchendpoint.MatchResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func WriteResponse(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

type Error struct {
	Code int
	Err  string
}

func encodeHTTPError(ctx context.Context, e error, w http.ResponseWriter) {
	if err, ok := e.(Error); ok {
		WriteResponse(w, err.Code, map[string]string{"err": err.Err})
		return
	}
	WriteResponse(w, http.StatusInternalServerError, map[string]string{})
}
func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}
func (e Error) Error() string {
	return e.Err
}
func ErrBadRequest() Error {
	return NewError(http.StatusBadRequest, "invalid JSON request")
}
func ErrNotAllowed() Error {
	return NewError(http.StatusMethodNotAllowed, "method not allowed")
}
