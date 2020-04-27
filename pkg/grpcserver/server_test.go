package grpcserver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test valid settings returned"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewServer()
			assert.True(t, len(server.StreamMiddleware) == 0)
			assert.True(t, len(server.UnaryMiddleware) == 0)
			assert.True(t, server.Done != nil)
			assert.True(t, server.StopSignal != nil)
		})
	}
}

func TestServer_WaitShutdown(t *testing.T) {

	t.Run("trigger stop from done signal", func(t *testing.T) {
		s := &Server{
			Server:     grpc.NewServer(),
			Done:       make(chan bool, 1),
			StopSignal: make(chan os.Signal, 1),
		}
		lis := bufconn.Listen(10)
		go s.Serve(lis)
		s.Done <- true
		s.WaitShutdown()
	})

	t.Run("trigger stop for os signal", func(t *testing.T) {
		s := &Server{
			Server:     grpc.NewServer(),
			Done:       make(chan bool, 1),
			StopSignal: make(chan os.Signal, 1),
		}
		s.StopSignal <- os.Interrupt
		lis := bufconn.Listen(10)
		go s.Serve(lis)
		s.WaitShutdown()
	})

}
