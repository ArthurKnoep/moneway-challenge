package cmd_transaction

import (
	"context"

	"github.com/caarlos0/env"

	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/transaction/config"
	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/transaction/database"
	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/transaction/grpc"
	"github.com/ArthurKnoep/moneway-challenge/pkg/protocol/grpc/transaction"
	"github.com/ArthurKnoep/moneway-challenge/pkg/service/v1"
)

func RunServer() error {
	ctx := context.Background()
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	conn := grpc.Init(&cfg)
	session := database.Init(&cfg)

	defer conn.Close()
	defer session.Close()
	return transaction.RunServer(ctx,
		v1.NewTransactionServiceServer(session, conn),
		cfg.TransactionPort)
}
