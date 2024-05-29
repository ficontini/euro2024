package transport

import (
	"context"
	"encoding/json"
	"net/http"

	matchendpoint "github.com/ficontini/euro2024/matchservice/endpoint"
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

func WriteResponse(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
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
