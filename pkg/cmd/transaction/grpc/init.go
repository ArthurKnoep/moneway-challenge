package grpc

import (
	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/transaction/config"
	"google.golang.org/grpc"
)

func Init(cfg *config.Config) *grpc.ClientConn {
	conn, err := grpc.Dial(cfg.BalanceHost+":"+cfg.BalancePort, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}
