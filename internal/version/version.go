package version

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jkieltyka/go-starter-kit/internal/config"
	"github.com/jkieltyka/go-starter-kit/pkg/httpserver"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	GoVersion    string
	GitCommit    string
	BuildTime    string
	BuildVersion string
)

type Version struct {
	Version   string `json:"version,omitempty"`
	Hash      string `json:"hash,omitempty"`
	Buildtime string `json:"buildtime,omitempty"`
}

func VersionRoute() httpserver.Route {

	return httpserver.Route{
		Path:   "/version",
		Method: http.MethodGet,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "application/json")

			data, err := json.Marshal(Version{
				Version:   BuildVersion,
				Hash:      GitCommit,
				Buildtime: BuildTime,
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

type VersionGRPCServer struct{}

func NewVersionGRPCServer(cfg *config.AppConfig) *VersionGRPCServer {
	return &VersionGRPCServer{}
}

func (v *VersionGRPCServer) Version(ctx context.Context, empty *emptypb.Empty) (*VersionReply, error) {
	return &VersionReply{
		Version:   BuildVersion,
		Hash:      GitCommit,
		Buildtime: BuildTime,
	}, nil
}
