package database

import (
	"github.com/ArthurKnoep/moneway-challenge/lib/database/models/account"
	"github.com/gocql/gocql"

	"github.com/ArthurKnoep/moneway-challenge/pkg/cmd/server/config"
)

func connect(cfg *config.Config) *gocql.Session {
	cluster := gocql.NewCluster(cfg.DatabaseHost)
	cluster.Keyspace = cfg.DatabaseKeyspace
	cluster.Consistency = gocql.Quorum
	if session, err := cluster.CreateSession(); err != nil {
		panic(err)
	} else {
		return session
	}
}

func Init(cfg *config.Config) *gocql.Session {
	session := connect(cfg)
	if err := account.CreateTable(session); err != nil {
		panic(err)
	}
	return session
}
