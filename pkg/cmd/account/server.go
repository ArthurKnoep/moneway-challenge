package cmd_account

import (
	"context"

	"github.com/caarlos0/env"

	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/account/config"
	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/account/database"
	"github.com/ArthurKnoep/moneway-challenge/pkg/protocol/grpc/account"
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
	return account.RunServer(ctx,
		v1.NewAccountServiceServer(session),
		cfg.AccountPort)
}
