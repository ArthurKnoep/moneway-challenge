package cmd

import (
	"context"

	"github.com/caarlos0/env"

	"github.com/ArthurKnoep/moneway-challenge/lib/database"
	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/server/config"
	"github.com/ArthurKnoep/moneway-challenge/pkg/protocol/grpc"
	"github.com/ArthurKnoep/moneway-challenge/pkg/service/v1"
)

func RunServer() error {
	ctx := context.Background()
	var cfg config.Config

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	session := database.Init(&cfg)
	defer session.Close()

	return grpc.RunServer(ctx,
		v1.NewAccountServiceServer(session),
		v1.NewBalanceServiceServer(session),
		v1.NewTransactionServiceServer(session),
		cfg.Port)
}
