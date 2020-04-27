package httpserver

import (
	"net/http"
	"os"
	"reflect"
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
	type fields struct {
		Server     http.Server
		Router     *mux.Router
		BasePath   string
		Done       chan bool
		StopSignal chan os.Signal
	}
	type args struct {
		routes []Route
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Server
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Server:     tt.fields.Server,
				Router:     tt.fields.Router,
				BasePath:   tt.fields.BasePath,
				Done:       tt.fields.Done,
				StopSignal: tt.fields.StopSignal,
			}
			if got := s.WithBaseRoutes(tt.args.routes...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.WithBaseRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_WithRoutes(t *testing.T) {
	type fields struct {
		Server     http.Server
		Router     *mux.Router
		BasePath   string
		Done       chan bool
		StopSignal chan os.Signal
	}
	type args struct {
		routes []Route
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Server:     tt.fields.Server,
				Router:     tt.fields.Router,
				BasePath:   tt.fields.BasePath,
				Done:       tt.fields.Done,
				StopSignal: tt.fields.StopSignal,
			}
			// s.Router.GetRoute()
			if got := s.WithRoutes(tt.args.routes...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.WithRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_LogRoutes(t *testing.T) {
	type fields struct {
		Server     http.Server
		Router     *mux.Router
		BasePath   string
		Done       chan bool
		StopSignal chan os.Signal
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Server:     tt.fields.Server,
				Router:     tt.fields.Router,
				BasePath:   tt.fields.BasePath,
				Done:       tt.fields.Done,
				StopSignal: tt.fields.StopSignal,
			}
			s.LogRoutes()
		})
	}
}

func TestServer_WaitShutdown(t *testing.T) {
	type fields struct {
		Server     http.Server
		Router     *mux.Router
		BasePath   string
		Done       chan bool
		StopSignal chan os.Signal
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Server:     tt.fields.Server,
				Router:     tt.fields.Router,
				BasePath:   tt.fields.BasePath,
				Done:       tt.fields.Done,
				StopSignal: tt.fields.StopSignal,
			}
			s.WaitShutdown()
		})
	}
}
