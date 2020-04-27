package httpserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test defaults are valid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer()
			assert.True(t, s.Done != nil)
			assert.True(t, s.StopSignal != nil)
			assert.True(t, s.Router != nil)
			assert.True(t, s.BasePath == "/")
		})
	}
}

func TestServer_WithBaseRoutes(t *testing.T) {

	tests := []struct {
		name  string
		route Route
	}{
		{"add a version GET route", Route{"/version", "GET", func(res http.ResponseWriter, req *http.Request) {
			_, _ = res.Write([]byte("ok"))
		}}},
		{"add a version POST route", Route{"/version", "POST", func(res http.ResponseWriter, req *http.Request) {
			_, _ = res.Write([]byte("ok"))
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:     mux.NewRouter(),
				BasePath:   "/v1",
				Done:       make(chan bool),
				StopSignal: make(chan os.Signal),
			}
			s = s.WithBaseRoutes(tt.route)
			route := s.Router.Get(tt.route.Path + tt.route.Method)

			result := httptest.NewRecorder()
			route.GetHandler().ServeHTTP(result, &http.Request{})

			assert.Equal(t, result.Body.Bytes(), []byte("ok"))
			methods, _ := route.GetMethods()
			assert.Equal(t, methods, []string{tt.route.Method})
		})
	}
}

func TestServer_WithRoutes(t *testing.T) {
	tests := []struct {
		name  string
		route Route
	}{
		{"add a version GET route", Route{"/version", "GET", func(res http.ResponseWriter, req *http.Request) {
			_, _ = res.Write([]byte("ok"))
		}}},
		{"add a version POST route", Route{"/version", "POST", func(res http.ResponseWriter, req *http.Request) {
			_, _ = res.Write([]byte("ok"))
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:     mux.NewRouter(),
				BasePath:   "/v1",
				Done:       make(chan bool),
				StopSignal: make(chan os.Signal),
			}
			s = s.WithRoutes(tt.route)
			route := s.Router.Get(s.BasePath + tt.route.Path + tt.route.Method)

			result := httptest.NewRecorder()
			route.GetHandler().ServeHTTP(result, &http.Request{})

			assert.Equal(t, result.Body.Bytes(), []byte("ok"))
			methods, _ := route.GetMethods()
			assert.Equal(t, methods, []string{tt.route.Method})
		})
	}
}

func TestServer_WaitShutdown(t *testing.T) {
	t.Run("trigger stop from done signal", func(t *testing.T) {
		s := &Server{
			Router:     mux.NewRouter(),
			Done:       make(chan bool, 1),
			StopSignal: make(chan os.Signal, 1),
		}
		go s.Start(8080)
		s.Done <- true
		s.WaitShutdown()
	})

	t.Run("trigger stop for os signal", func(t *testing.T) {
		s := &Server{
			Router:     mux.NewRouter(),
			Done:       make(chan bool, 1),
			StopSignal: make(chan os.Signal, 1),
		}
		s.StopSignal <- os.Interrupt
		go s.Start(8080)
		s.WaitShutdown()
	})
}
