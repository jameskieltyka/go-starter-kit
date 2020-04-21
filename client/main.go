package main

import (
	"context"
	"fmt"
	"time"

	versioner "github.com/jkieltyka/go-starter-kit/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	c := versioner.NewVersionerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("making request")
	fmt.Println(c.Version(ctx, &emptypb.Empty{}))
}
