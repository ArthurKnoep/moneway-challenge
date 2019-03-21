package v1

import (
	"context"
	"fmt"
	"github.com/ArthurKnoep/moneway-challenge/lib/database/models/account"
	"github.com/ArthurKnoep/moneway-challenge/pkg/api/v1"
	"github.com/gocql/gocql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type balanceServiceServer struct {
	dbSession *gocql.Session
}

func NewBalanceServiceServer(session *gocql.Session) v1.BalanceServiceServer {
	return &balanceServiceServer{
		dbSession: session,
	}
}

func (s *balanceServiceServer) GetBalance(ctx context.Context, req *v1.GetBalanceRequest) (*v1.GetBalanceResponse, error) {
	if len(req.AccountUuid) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid account_uuid")
	}
	if exist, err := account.UuidExist(s.dbSession, req.AccountUuid); err != nil {
		return nil, status.Errorf(codes.Unknown, "Check exist: %v", err)
	} else if exist == false {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	} else {
		if a, err := account.GetByUuid(s.dbSession, req.AccountUuid); err != nil {
			return nil, status.Errorf(codes.Unknown, "Get: %v", err)
		} else {
			return &v1.GetBalanceResponse{
				Balance:  a.Balance,
				Currency: a.Currency,
			}, nil
		}
	}
}

func (s *balanceServiceServer) UpdateBalance(ctx context.Context, req *v1.UpdateBalanceRequest) (*v1.UpdateBalanceResponse, error) {
	if len(req.AccountUuid) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid account_uuid")
	}
	if exist, err := account.UuidExist(s.dbSession, req.AccountUuid); err != nil {
		return nil, status.Errorf(codes.Unknown, "Check exist: %v", err)
	} else if exist == false {
		return nil, status.Errorf(codes.NotFound, "Account not found")
	} else {
		if a, err := account.GetByUuid(s.dbSession, req.AccountUuid); err != nil {
			return nil, status.Errorf(codes.Unknown, "Get: %v", err)
		} else if req.Currency != a.Currency {
			return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Invalid currency, expected \"%s\"", a.Currency))
		} else if a.Balance+req.Balance < 0 {
			return nil, status.Errorf(codes.InvalidArgument, "Balance cannot be negative")
		} else if err := account.SetSolde(s.dbSession, req.AccountUuid, a.Balance+req.Balance); err != nil {
			return nil, status.Errorf(codes.Unknown, "Unable to update the balance")
		} else {
			return &v1.UpdateBalanceResponse{
				Balance:  a.Balance + req.Balance,
				Currency: req.Currency,
			}, nil
		}
	}
}
