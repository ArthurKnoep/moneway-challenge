package v1

import (
	"context"
	"github.com/ArthurKnoep/moneway-challenge/lib/database/models/transaction"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"

	"github.com/gocql/gocql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ArthurKnoep/moneway-challenge/lib/database/models/account"
	"github.com/ArthurKnoep/moneway-challenge/pkg/api/v1"
)

type transactionServiceServer struct {
	dbSession     *gocql.Session
	balanceClient v1.BalanceServiceClient
}

func NewTransactionServiceServer(session *gocql.Session, conn *grpc.ClientConn) v1.TransactionServiceServer {
	return &transactionServiceServer{
		dbSession:     session,
		balanceClient: v1.NewBalanceServiceClient(conn),
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
					Timestamp:       &ts,
					Note:            elem.Note,
					Amount:          elem.Amount,
					Currency:        elem.Currency,
				})
			}
			return &v1.ListTransactionResponse{
				Transaction: ret,
			}, nil
		}
	}
}

func (s *transactionServiceServer) recreditUser(ctx context.Context, accountDest, currency string, amount float64) {
	_, _ = s.balanceClient.UpdateBalance(ctx, &v1.UpdateBalanceRequest{
		AccountUuid: accountDest,
		Balance:     amount,
		Currency:    currency,
	})
}

func (s *transactionServiceServer) Transaction(ctx context.Context, req *v1.TransactionRequest) (*v1.TransactionResponse, error) {
	if req.Amount < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Amount cannot be negative")
	}
	if _, err := s.balanceClient.UpdateBalance(ctx, &v1.UpdateBalanceRequest{
		AccountUuid: req.AccountSrcUuid,
		Balance:     -req.Amount,
		Currency:    req.Currency,
	}); err != nil {
		return nil, err
	} else if _, err := s.balanceClient.UpdateBalance(ctx, &v1.UpdateBalanceRequest{
		AccountUuid: req.AccountDestUuid,
		Balance:     req.Amount,
		Currency:    req.Currency,
	}); err != nil {
		s.recreditUser(ctx, req.AccountSrcUuid, req.Currency, req.Amount)
		return nil, err
	}
	if rest, err := transaction.AddTransaction(s.dbSession, req.AccountSrcUuid, req.AccountDestUuid, &transaction.Transaction{
		Note:     req.Note,
		Amount:   req.Amount,
		Currency: req.Currency,
	}); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	} else {
		return &v1.TransactionResponse{
			TransactionUuid: rest.TransactionUuid.String(),
		}, nil
	}
}

func (s *transactionServiceServer) UpdateTransaction(ctx context.Context, req *v1.UpdateTransactionRequest) (*v1.UpdateTransactionResponse, error) {
	if _, err := transaction.UpdateTransaction(s.dbSession, req.TransactionUuid, req.Note); err != nil {
		return nil, status.Errorf(codes.Unknown, err.Error())
	}
	return &v1.UpdateTransactionResponse{}, nil
}