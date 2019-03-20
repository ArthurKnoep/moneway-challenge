package account

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

var TableName = "account"

type Account struct {
	AccountUUID gocql.UUID `db:"account_uuid"`
	Username    string     `db:"username"`
	Balance     float64    `db:"balance"`
	Currency    string     `db:"currency"`
}

func CreateTable(session *gocql.Session) error {
	queryStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
		"account_uuid uuid PRIMARY KEY,"+
		"username varchar,"+
		"balance double,"+
		"currency varchar"+
		")", TableName)
	q := session.Query(queryStr).RetryPolicy(nil)
	defer q.Release()
	return q.Exec()
}

func AlreadyExist(session *gocql.Session, accountName string) (bool, error) {
	stmt, names := qb.Select(TableName).Where(qb.Eq("username")).AllowFiltering().ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindStruct(Account{
		Username: accountName,
	})
	var accounts []Account
	if err := q.SelectRelease(&accounts); err != nil {
		return false, err
	}
	return len(accounts) != 0, nil
}

func CreateAccount(session *gocql.Session, account, currency string) (*Account, error) {
	a := Account{
		AccountUUID: gocql.TimeUUID(),
		Username:    account,
		Balance:     0,
		Currency:    currency,
	}
	stmt, names := qb.Insert(TableName).Columns("account_uuid", "username", "balance", "currency").ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindStruct(a)
	if err := q.ExecRelease(); err != nil {
		return nil, err
	}
	return &a, nil
}
