package grpc

import (
	"context"
	"github.com/ArthurKnoep/moneway-challenge/pkg/api/v1"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func RunServer(ctx context.Context, account v1.AccountServiceServer, balance v1.BalanceServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	v1.RegisterAccountServiceServer(server, account)
	v1.RegisterBalanceServiceServer(server, balance)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("shutting down server....")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()
	log.Println("starting server....")
	return server.Serve(listen)
}
