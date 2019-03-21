package cmd_balance

import (
	"context"

	"github.com/caarlos0/env"

	"github.com/ArthurKnoep/moneway-challenge/lib/database"
	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/server/config"
	"github.com/ArthurKnoep/moneway-challenge/pkg/protocol/grpc/balance"
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
	return balance.RunServer(ctx,
		v1.NewBalanceServiceServer(session),
		cfg.Port)
}
