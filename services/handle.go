package services

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/the-web3/mock-risk-server/protobuf/common"
	"github.com/the-web3/mock-risk-server/protobuf/riskcontroller"
)

const RiskKey = "10000"

func (rss *RiskServerWireServices) CheckAmlAddress(ctx context.Context, request *riskcontroller.CheckAmlAddressRequest) (*riskcontroller.CheckAmlAddressResponse, error) {
	var retAddressList []*riskcontroller.RetAmlAddress
	for _, reqItem := range request.AmlAddress {
		// todo: 调用 chainalysis 和漫雾等平台的接口
		retAddressList = append(retAddressList, &riskcontroller.RetAmlAddress{
			Address:     reqItem.Address,
			AddressType: "white",
		})
	}
	return &riskcontroller.CheckAmlAddressResponse{
		Code:          common.ReturnCode_SUCCESS,
		Msg:           "check address success",
		RetAmlAddress: retAddressList,
	}, nil
}

func (rss *RiskServerWireServices) CheckChainTransactions(ctx context.Context, request *riskcontroller.CheckChainTransactionsRequest) (*riskcontroller.CheckChainTransactionsResponse, error) {
	blockInfo, err := rss.rpcApiClient.GetLastestBlock()
	if err != nil {
		log.Error("GetLastestBlock failed", "err", err)
		return nil, err
	}
	var retChainTxList []*riskcontroller.RetChainTransaction
	for _, reqItem := range request.ChainTxn {
		txInfo, err := rss.rpcApiClient.GetTransactionByHash(reqItem.TxHash)
		if err != nil {
			log.Info("rpcApiClient.GetTransactionByHash", "err", err)
			return nil, err
		}
		retChainTxList = append(retChainTxList, &riskcontroller.RetChainTransaction{
			FromAddress: txInfo.From[0].Address,
			ToAddress:   txInfo.To[0].Address,
			Amount:      txInfo.From[0].Amount,
			Fee:         txInfo.Fee,
			Confirms:    blockInfo.Number.Uint64() - txInfo.BlockHeight,
		})
	}
	return &riskcontroller.CheckChainTransactionsResponse{
		Code:        common.ReturnCode_SUCCESS,
		Msg:         "check transaction success",
		RetChainTxn: retChainTxList,
	}, nil
}

func (rss *RiskServerWireServices) CheckUserTransaction(ctx context.Context, request *riskcontroller.CheckUserTransactionRequest) (*riskcontroller.CheckUserTransactionResponse, error) {
	var retUserTxList []*riskcontroller.RetUserTransaction
	for _, reqItem := range request.UserTxn {
		txInfo, err := rss.db.Transactions.QueryTransactionsByRequestId(reqItem.RequestId)
		if err != nil {
			log.Error("QueryTransactionsByRequestId", "err", err)
			return nil, err
		}
		if txInfo.FromAddress == reqItem.FromAddress && txInfo.ToAddress == reqItem.ToAddress && txInfo.Amount.String() == reqItem.Amount && txInfo.TokenAddress == reqItem.TokenAddress {
			txBody := &TransactionBody{
				RequestId:    reqItem.RequestId,
				FromAddress:  reqItem.FromAddress,
				TokenAddress: reqItem.TokenAddress,
				ToAddress:    reqItem.ToAddress,
				Amount:       reqItem.Amount,
			}
			byteTxBody, _ := json.Marshal(txBody)
			byteTxBodyHash := hexutils.BytesToHex(crypto.Keccak256(byteTxBody))
			RiskKeyHashStr := hexutils.BytesToHex(crypto.Keccak256(byteTxBody, []byte(RiskKey)))
			retUserTxList = append(retUserTxList, &riskcontroller.RetUserTransaction{
				RequestId:      reqItem.RequestId,
				TxBodyHash:     byteTxBodyHash,
				TxBodyRiskHash: RiskKeyHashStr,
			})
		} else {
			log.Warn("Check user transaction fail", "reqId", reqItem.RequestId)
			continue
		}
	}
	return &riskcontroller.CheckUserTransactionResponse{
		Code:       common.ReturnCode_SUCCESS,
		Msg:        "check transaction success",
		RetUserTxn: retUserTxList,
	}, nil
}
