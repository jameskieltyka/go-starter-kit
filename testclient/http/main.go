package main

import (
	"context"
	"fmt"

	"github.com/jkieltyka/go-starter-kit/pkg/client"
	_ "github.com/jkieltyka/go-starter-kit/pkg/logger"
)

func main() {
	fmt.Println(client.
		NewHTTPClient().
		WithRequest(
			"GET",
			"http://localhost:8080/version",
			nil,
		).
		WithMiddleware(
			client.Logger{Url: "localhost:8080/version", Method: "GET"},
			client.CircuitBreaker{Name: "test"}.Configure(),
		).Run(context.Background()))

}
