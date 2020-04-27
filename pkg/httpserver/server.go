package httpserver

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type Route struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

type Server struct {
	http.Server
	Router     *mux.Router
	BasePath   string
	Done       chan bool
	StopSignal chan os.Signal
}

func NewServer() *Server {
	return &Server{
		Router:     mux.NewRouter(),
		BasePath:   "/",
		Done:       make(chan bool),
		StopSignal: make(chan os.Signal),
	}
}

func (s *Server) WithBasePath(basepath string) *Server {
	s.BasePath = basepath
	return s
}

func (s *Server) WithMiddleware(m ...mux.MiddlewareFunc) *Server {
	//function on router to use middleware chain
	s.Router.Use(m...)
	return s
}

//Base routes are for endpoints that should not be versioned, for example /health or /version
func (s *Server) WithBaseRoutes(routes ...Route) *Server {
	sub := s.Router.PathPrefix("/").Subrouter()
	for _, route := range routes {
		sub.HandleFunc(route.Path, route.Handler).Methods(route.Method).Name(fmt.Sprintf("%s%s", route.Path, route.Method))
	}

	return s
}

func (s *Server) WithRoutes(routes ...Route) *Server {
	sub := s.Router.PathPrefix(s.BasePath).Subrouter()
	for _, route := range routes {
		sub.HandleFunc(route.Path, route.Handler).Methods(route.Method).Name(fmt.Sprintf("%s%s%s", s.BasePath, route.Path, route.Method))
	}

	return s
}

func (s *Server) LogRoutes() {
	err := s.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.URLPath()
		if err != nil {
			return err
		}

		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}

		zap.S().Infof("added path %s, methods: %s", path, methods)
		return nil
	})

	if err != nil {
		zap.S().Infof("unable to walk server routes, err %v", err)
	}
}

func (s *Server) WaitShutdown() {
	signal.Notify(s.StopSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-s.Done:
		zap.S().Info("server stopping due to finish request")
	case sig := <-s.StopSignal:
		zap.S().Infof("server stopped due to signal %v", sig)
	}

	//give the server 10 seconds to finish any outstanding requests
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		zap.S().Infof("server failed to gracefully shut down, err %v", err)
	}
}

// Start starts the server on the defined port
func (s *Server) Start(port int) {
	// Set server handler and listen port
	s.Handler = s.Router
	s.Addr = fmt.Sprintf(":%v", port)
	zap.S().Infof("Starting server %v", s.Addr)

	s.LogRoutes()

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			zap.S().Errorf("server closed due to error %v", err)
			s.Done <- true
		}
	}()

	// wait shutdown
	s.WaitShutdown()
}
