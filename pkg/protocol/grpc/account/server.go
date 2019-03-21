package account

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"github.com/ArthurKnoep/moneway-challenge/pkg/api/v1"
)

func RunServer(ctx context.Context, account v1.AccountServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterAccountServiceServer(server, account)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down account server....")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()
	log.Printf("starting account server (port: %s)....", port)
	return server.Serve(listen)
}
