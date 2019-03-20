package v1

import (
	"context"
	"github.com/ArthurKnoep/moneway-challenge/lib/database/models/account"
	"github.com/ArthurKnoep/moneway-challenge/pkg/api/v1"
	"github.com/gocql/gocql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type accountServiceServer struct {
	dbSession *gocql.Session
}

func NewAccountServiceServer(session *gocql.Session) v1.AccountServiceServer {
	return &accountServiceServer{
		dbSession: session,
	}
}

func (s *accountServiceServer) CreateAccount(ctx context.Context, req *v1.CreateRequest) (*v1.CreateResponse, error) {
	if len(req.Username) == 0 || len(req.Currency) < 3 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid username or currency")
	}
	if exist, err := account.AlreadyExist(s.dbSession, req.Username); err != nil {
		return nil, status.Errorf(codes.Unknown, "Unable to check the account")
	} else if exist == true {
		return nil, status.Errorf(codes.AlreadyExists, "Account already exist")
	} else {
		if rst, err := account.CreateAccount(s.dbSession, req.Username, req.Currency); err != nil {
			return nil, status.Errorf(codes.Unknown, "Unable to create the account")
		} else {
			return &v1.CreateResponse{
				AccountUuid: rst.AccountUUID.String(),
			}, nil
		}
	}
}