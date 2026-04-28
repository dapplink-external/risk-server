package services

import (
	"context"
	
	"github.com/the-web3/mock-risk-server/protobuf/riskcontroller"
)

const BatchAddressNum = 10000

func (rss *RiskServerWireServices) CheckAmlAddress(ctx context.Context, request *riskcontroller.CheckAmlAddressRequest) (*riskcontroller.CheckAmlAddressResponse, error) {
	return nil, nil
}

func (rss *RiskServerWireServices) CheckChainTransactions(ctx context.Context, request *riskcontroller.CheckChainTransactionsRequest) (*riskcontroller.CheckChainTransactionsResponse, error) {
	return nil, nil
}

func (rss *RiskServerWireServices) CheckUserTransaction(ctx context.Context, request *riskcontroller.CheckUserTransactionRequest) (*riskcontroller.CheckUserTransactionResponse, error) {
	return nil, nil
}
