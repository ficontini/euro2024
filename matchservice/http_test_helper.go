package main

import (
	"net/http/httptest"
	"testing"

	matchendpoint "github.com/ficontini/euro2024/matchservice/pkg/endpoint"
	"github.com/ficontini/euro2024/matchservice/pkg/service"
	"github.com/ficontini/euro2024/matchservice/pkg/transport"
	"github.com/ficontini/euro2024/matchservice/store"
)

func setupHTTPServer(t *testing.T, store store.Store) *httptest.Server {
	var (
		svc = service.New(store)
		eps = matchendpoint.New(svc)
		mux = transport.NewHTTPHandler(eps)
	)
	return httptest.NewServer(mux)
}

func setupHTTPClient(t *testing.T, instance string) (service.Service, error) {
	return transport.NewHTTPClient(instance)
}
