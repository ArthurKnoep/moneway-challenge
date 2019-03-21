package database

import (
	"time"

	"github.com/gocql/gocql"
)

func Connect(host, keyspace string) *gocql.Session {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	cluster.Timeout = 1 * time.Minute
	cluster.Consistency = gocql.Quorum
	if session, err := cluster.CreateSession(); err != nil {
		panic(err)
	} else {
		return session
	}
}
