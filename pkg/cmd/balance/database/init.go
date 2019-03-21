package database

import (
	"github.com/gocql/gocql"

	"github.com/ArthurKnoep/moneway-challenge/lib/database"
	"github.com/ArthurKnoep/moneway-challenge/lib/database/models/account"
	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/balance/config"
)

func Init(cfg *config.Config) *gocql.Session {
	session := database.Connect(cfg.DatabaseHost, cfg.DatabaseKeyspace)
	if err := account.CreateTable(session); err != nil {
		panic(err)
	}
	return session
}
