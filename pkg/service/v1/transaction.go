package v1

import (
	"context"
	"github.com/ArthurKnoep/moneway-challenge/lib/database/models/transaction"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/gocql/gocql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ArthurKnoep/moneway-challenge/lib/database/models/account"
	"github.com/ArthurKnoep/moneway-challenge/pkg/api/v1"
)

type transactionServiceServer struct {
	dbSession *gocql.Session
}

func NewTransactionServiceServer(session *gocql.Session) v1.TransactionServiceServer {
	return &transactionServiceServer{
		dbSession: session,
	}
}

func (s *transactionServiceServer) ListTransaction(ctx context.Context, req *v1.ListTransactionRequest) (*v1.ListTransactionResponse, error) {
	if len(req.AccountUuid) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid account_uuid")
	}
	if exist, err := account.UuidExist(s.dbSession, req.AccountUuid); err != nil {
		return nil, status.Errorf(codes.Unknown, "Check exist: %v", err)
	} else if exist == false {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	} else {
		if list, err := transaction.GetListTransaction(s.dbSession, req.AccountUuid); err != nil {
			return nil, status.Errorf(codes.Unknown, "Error while listing")
		} else {
			var ret []*v1.ListTransactionResponse_Transaction
			for _, elem := range list {
				var ts timestamp.Timestamp
				ts.Seconds = elem.Timestamp.Unix()
				ts.Nanos = int32(elem.Timestamp.UnixNano())
				ret = append(ret, &v1.ListTransactionResponse_Transaction{
					TransactionUuid: elem.TransactionUuid.String(),
					Timestamp: &ts,
					Note: elem.Note,
					Amount: elem.Amount,
					Currency: elem.Currency,
				})
			}
			return &v1.ListTransactionResponse{
				Transaction: ret,
			}, nil
		}
	}
}