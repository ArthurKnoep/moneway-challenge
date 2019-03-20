package cmd

import (
	"context"

	"github.com/caarlos0/env"

	"github.com/ArthurKnoep/moneway-challenge/pkg/protocol/grpc"
	"github.com/ArthurKnoep/moneway-challenge/pkg/service/v1"
)

type Config struct {
	Port string `env:"PORT" envDefault:"8080"`
}

func RunServer() error {
	ctx := context.Background()
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	api := v1.NewBalanceServiceServer()
	return grpc.RunServer(ctx, api, cfg.Port)
}
