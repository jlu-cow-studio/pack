package rpc

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jlu-cow-studio/common/dal/rpc/pack"
	"github.com/jlu-cow-studio/pack/handler"
	"google.golang.org/grpc"
)

var (
	ErrChan chan error
	port    = flag.Int("port", 8080, "The server port")
)

func Init() {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pack.RegisterPackServiceServer(s, &handler.Handler{})

		log.Printf("server listening at %v", lis.Addr())
		err = s.Serve(lis)
		close(ErrChan)
		if err != nil {
			panic(err)
		}
	}()
}
