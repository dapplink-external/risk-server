package walletapiclient

import (
	"context"
	"math/big"

	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-external/risk-server/protobuf/common"
	"github.com/dapplink-external/risk-server/protobuf/walletapi"
)

const ConsumerToken = "DappLinkTheWeb3202402290001"

type WalletApiGateWayServiceClient struct {
	Ctx       context.Context
	ChainId   string
	RpcClient walletapi.WalletApiGateWayServiceClient
}

func NewWalletApiGateWayServiceClient(ctx context.Context, rpcClient walletapi.WalletApiGateWayServiceClient, chainId string) (*WalletApiGateWayServiceClient, error) {
	log.Info("New risk rpc server", "ChainId", chainId)
	return &WalletApiGateWayServiceClient{Ctx: ctx, RpcClient: rpcClient, ChainId: chainId}, nil
}

func (wg *WalletApiGateWayServiceClient) GetLastestBlock() (*BlockHeader, error) {
	ltBlock := &walletapi.LastestBlockRequest{
		ChainId:       wg.ChainId,
		ConsumerToken: ConsumerToken,
		Network:       "mainnet",
	}

	latestBlockInfo, err := wg.RpcClient.GetLastestBlock(wg.Ctx, ltBlock)
	if err != nil {
		log.Error("RpcClient.GetLastestBlock failed", "err", err)
		return nil, err
	}

	if latestBlockInfo.Code != common.ReturnCode_SUCCESS {
		log.Error("GetLastestBlock failed", "resultAddressInfo", latestBlockInfo)
		return nil, errors.New("get lastest block failed")
	}
	return &BlockHeader{
		Hash:       latestBlockInfo.Hash,
		ParentHash: latestBlockInfo.ParentHash,
		Number:     big.NewInt(int64(latestBlockInfo.Height)),
		Timestamp:  latestBlockInfo.Timestamp,
	}, nil
}

func (wg *WalletApiGateWayServiceClient) GetTransactionByHash(txHash string) (*walletapi.TransactionList, error) {
	blockRequest := &walletapi.TransactionByHashRequest{
		ChainId:       wg.ChainId,
		ConsumerToken: ConsumerToken,
		Network:       "mainnet",
		Hash:          txHash,
	}
	txInfo, err := wg.RpcClient.GetTransactionByHash(wg.Ctx, blockRequest)
	if err != nil {
		log.Error("RpcClient transaction by hash failed", "err", err)
		return nil, err
	}
	if txInfo.Code != common.ReturnCode_SUCCESS {
		log.Error("get transaction by hash failed", "txInfo", txInfo)
		return nil, errors.New("get transaction by hash failed")
	}
	return txInfo.Transaction, nil
}
