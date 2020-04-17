package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jkieltyka/go-starter-kit/internal/config"
	"github.com/jkieltyka/go-starter-kit/pkg/httpserver"
	"github.com/jkieltyka/go-starter-kit/pkg/middleware"
)

func SetupHTTPServer(c *config.AppConfig) *httpserver.Server {
	versionPath := fmt.Sprintf("/%v", c.Version)
	server := httpserver.NewServer().
		WithBasePath(versionPath).
		WithMiddleware(middleware.ExampleMiddleware).
		WithBaseRoutes(VersionRoute(c)).
		WithRoutes(
			VersionRoute(c),
		)
	return server
}

func VersionRoute(c *config.AppConfig) httpserver.Route {
	type Version struct {
		Version string `json:"version,omitempty"`
	}

	return httpserver.Route{
		Path:   "/version",
		Method: http.MethodGet,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "application/json")

			data, err := json.Marshal(Version{
				Version: c.Version,
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(data)
		}),
	}
}
