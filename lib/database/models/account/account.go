package account

import (
	"errors"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

var TableName = "account"

// Database model
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

func UsernameExist(session *gocql.Session, accountName string) (bool, error) {
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

func UuidExist(session *gocql.Session, accountUuid string) (bool, error) {
	stmt, names := qb.Select(TableName).Where(qb.Eq("account_uuid")).AllowFiltering().ToCql()
	if uuid, err := gocql.ParseUUID(accountUuid); err != nil {
		return false, errors.New("invalid account id")
	} else {
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(Account{
			AccountUUID: uuid,
		})
		var accounts []Account
		if err := q.SelectRelease(&accounts); err != nil {
			return false, err
		}
		return len(accounts) != 0, nil
	}
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

func GetByUuid(session *gocql.Session, accountUuid string) (*Account, error) {
	stmt, names := qb.Select(TableName).Where(qb.Eq("account_uuid")).ToCql()
	if uuid, err := gocql.ParseUUID(accountUuid); err != nil {
		return nil, errors.New("invalid account id")
	} else {
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(Account{
			AccountUUID: uuid,
		})
		var a Account
		if err := q.GetRelease(&a); err != nil {
			return nil, err
		}
		return &a, nil
	}
}

func SetSolde(session *gocql.Session, accountUuid string, value float64) error {
	stmt, names := qb.Update(TableName).Set("balance").Where(qb.Eq("account_uuid")).ToCql()
	if uuid, err := gocql.ParseUUID(accountUuid); err != nil {
		return errors.New("invalid account id")
	} else {
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(Account{
			AccountUUID: uuid,
			Balance:     value,
		})
		return q.Exec()
	}
}

func DeleteAccount(session *gocql.Session, accountUuid string) error {
	stmt, names := qb.Delete(TableName).Where(qb.Eq("account_uuid")).ToCql()
	if uuid, err := gocql.ParseUUID(accountUuid); err != nil {
		return errors.New("invalid account id")
	} else {
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(Account{
			AccountUUID: uuid,
		})
		return q.ExecRelease()
	}
}
