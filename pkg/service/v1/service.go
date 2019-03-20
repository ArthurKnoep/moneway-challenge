package v1

import (
	"context"
	"github.com/ArthurKnoep/moneway-challenge/pkg/api/v1"
)

type balanceServiceServer struct {
}

func NewBalanceServiceServer() v1.BalanceServiceServer {
	return &balanceServiceServer{}
}

func (s *balanceServiceServer) GetBalance(ctx context.Context, request *v1.BalanceRequest) (*v1.BalanceResult, error) {
	return nil, nil
}
