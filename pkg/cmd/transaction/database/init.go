package database

import (
	"github.com/gocql/gocql"

	"github.com/ArthurKnoep/moneway-challenge/lib/database"
	"github.com/ArthurKnoep/moneway-challenge/lib/database/models/transaction"
	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/transaction/config"
)

func Init(cfg *config.Config) *gocql.Session {
	session := database.Connect(cfg.DatabaseHost, cfg.DatabaseKeyspace)
	if err := transaction.CreateTable(session); err != nil {
		panic(err)
	}
	return session
}
