package httpmiddleware

import (
	"context"
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

//Middleware can be formed using the following structure
func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("some-header", "some-value")

		ctx := context.WithValue(r.Context(), "some-key", "somevalue")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
