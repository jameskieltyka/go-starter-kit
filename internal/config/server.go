package config

// import (
// 	"context"
// 	"net/http"
// )

// type Route struct {
// 	Path    String
// 	Method  String
// 	Handler http.HandlerFunc
// }

// type Middleware func(next http.Handler) http.Handler

// type Server gorilla.router

// func ExampleMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		requestId := r.Header.Get("some-header")
// 		w.Header().Set("some-header", "some-value")

// 		ctx := context.WithValue(r.Context(), "some-key", "somevalue")
// 		next.ServerHttp(w, r.WithContext(ctx))
// 	})
// }

// func (server *Server) WithMiddleware(m ...Middleware) *server {
// 	//function on router to use middleware chain
// 	s.Router.User(m...)
// 	return s
// }

// func (server *Server) WithRoutes(basePath string, routes ...Route) *server {
// 	sub := s.router.PathPrefix(basePath).Subrouter()
// }
