package transaction

import (
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"time"
)

var TableName = "transaction"

// Database model
type Transaction struct {
	TransactionUuid gocql.UUID `db:"transaction_uuid"`
	AccountDestUuid gocql.UUID `db:"account_dest_uuid"`
	AccountSrcUuid  gocql.UUID `db:"account_src_uuid"`
	Timestamp       time.Time  `db:"timestamp"`
	Note            string     `db:"note"`
	Amount          float64    `db:"amount"`
	Currency        string     `db:"currency"`
}

func CreateTable(session *gocql.Session) error {
	queryStr := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
		"transaction_uuid uuid PRIMARY KEY,"+
		"account_dest_uuid uuid,"+
		"account_src_uuid uuid,"+
		"timestamp timestamp,"+
		"note text,"+
		"amount double,"+
		"currency varchar"+
		")", TableName)
	q := session.Query(queryStr).RetryPolicy(nil)
	defer q.Release()
	return q.Exec()
}

func GetListTransaction(session *gocql.Session, accountUuid string) ([]Transaction, error) {
	stmt, names := qb.Select(TableName).Where(qb.Eq("account_dest_uuid"), qb.Eq("account_src_uuid")).AllowFiltering().ToCql()
	if uuid, err := gocql.ParseUUID(accountUuid); err != nil {
		return nil, errors.New("invalid account id")
	} else {
		var transactions []Transaction
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(Transaction{
			AccountDestUuid: uuid,
			AccountSrcUuid:  uuid,
		})
		if err := q.SelectRelease(&transactions); err != nil {
			return nil, err
		} else {
			return transactions, nil
		}
	}
}

func AddTransaction(session *gocql.Session, src, dest string, transaction *Transaction) (*Transaction, error) {
	transaction.Timestamp = time.Now()
	transaction.TransactionUuid = gocql.TimeUUID()
	if srcUuid, err := gocql.ParseUUID(src); err != nil {
		return nil, errors.New("invalid source id")
	} else if destUuid, err := gocql.ParseUUID(dest); err != nil {
		return nil, errors.New("invalid dest id")
	} else {
		transaction.AccountSrcUuid = srcUuid
		transaction.AccountDestUuid = destUuid
		stmt, names := qb.Insert(TableName).Columns("transaction_uuid", "account_dest_uuid", "account_src_uuid",
			"timestamp", "note", "amount", "currency").ToCql()
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(transaction)
		if err := q.ExecRelease(); err != nil {
			return nil, err
		}
		return transaction, err
	}
}
