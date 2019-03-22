package transaction

import (
	"errors"
	"fmt"
	"sort"
	"time"
	
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
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

func getListTransacOut(session *gocql.Session, accountUuid string) ([]Transaction, error) {
	stmt, names := qb.Select(TableName).Where(qb.Eq("account_src_uuid")).AllowFiltering().ToCql()
	if uuid, err := gocql.ParseUUID(accountUuid); err != nil {
		return nil, errors.New("invalid account id")
	} else {
		var transactions []Transaction
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(Transaction{
			AccountSrcUuid:  uuid,
		})
		if err := q.SelectRelease(&transactions); err != nil {
			return nil, err
		} else {
			var ret []Transaction
			for _, elem := range transactions {
				elem.Amount = -elem.Amount
				ret = append(ret, elem)
			}
			return ret, nil
		}
	}
}

func getListTransacIn(session *gocql.Session, accountUuid string) ([]Transaction, error) {
	stmt, names := qb.Select(TableName).Where(qb.Eq("account_dest_uuid")).AllowFiltering().ToCql()
	if uuid, err := gocql.ParseUUID(accountUuid); err != nil {
		return nil, errors.New("invalid account id")
	} else {
		var transactions []Transaction
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(Transaction{
			AccountDestUuid: uuid,
		})
		if err := q.SelectRelease(&transactions); err != nil {
			return nil, err
		} else {
			return transactions, nil
		}
	}
}

type ByTimestamp []Transaction

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp.Second() < a[j].Timestamp.Second() }

func GetListTransaction(session *gocql.Session, accountUuid string) ([]Transaction, error) {
	if out, err := getListTransacOut(session, accountUuid); err != nil {
		return nil, err
	} else if in, err := getListTransacIn(session, accountUuid); err != nil {
		return nil, err
	} else {
		var ret []Transaction
		ret = append(out, in...)
		sort.Slice(ret, func(a, b int) bool {
			return ret[a].Timestamp.Second() > ret[b].Timestamp.Second()
		})
		return ret, nil
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

func UpdateTransaction(session *gocql.Session, transactionUuid, note string) (*Transaction, error) {
	stmt, names := qb.Select(TableName).Where(qb.Eq("transaction_uuid")).ToCql()
	if uuid, err := gocql.ParseUUID(transactionUuid); err != nil {
		return nil, errors.New("invalid transaction id")
	} else {
		var transaction Transaction
		q := gocqlx.Query(session.Query(stmt), names).BindStruct(Transaction{
			TransactionUuid: uuid,
		})
		if err := q.GetRelease(&transaction); err != nil {
			return nil, err
		} else {
			transaction.Note = note
			stmt, names := qb.Update(TableName).Set("note").Where(qb.Eq("transaction_uuid")).ToCql()
			q := gocqlx.Query(session.Query(stmt), names).BindStruct(Transaction{
				TransactionUuid: uuid,
				Note: note,
			})
			if err := q.ExecRelease(); err != nil {
				return nil, err
			} else {
				return &transaction, nil
			}
		}
	}
}
