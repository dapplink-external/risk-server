package services

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/log"
	"github.com/syndtr/goleveldb/leveldb/errors"

	"github.com/the-web3/mock-risk-server/protobuf/common"
	"github.com/the-web3/mock-risk-server/protobuf/riskcontroller"
)

const RiskKey = "10000"
const withdrawTxKeyPrefix = "withdraw_tx"
const withdrawVerifiedKeyPrefix = "withdraw_verified"
const transactionFlowKeyPrefix = "transaction_flow"

func (rss *RiskServerWireServices) SubmitWithdraw(ctx context.Context, request *riskcontroller.RiskWithdrawTransactionRequest) (*riskcontroller.RiskWithdrawTransactionResponse, error) {
	if rss.AccessToken != request.GetConsumerToken() {
		return &riskcontroller.RiskWithdrawTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "invalid consumer token",
		}, nil
	}
	for _, tx := range request.GetWithdrawTxn() {
		if tx == nil {
			log.Error("nil transaction")
			continue
		}
		key := withdrawTxKey(tx.GetRequestId(), tx.GetBusinessId(), tx.GetChainId())
		value, err := json.Marshal(toCanonicalWithdrawTx(tx))
		if err != nil {
			log.Error("marshal withdraw tx failed", "err", err)
			return &riskcontroller.RiskWithdrawTransactionResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "marshal withdraw transaction failed",
			}, nil
		}
		if err := rss.levelStore.Put([]byte(key), value); err != nil {
			log.Error("store withdraw tx failed", "err", err, "key", key)
			return &riskcontroller.RiskWithdrawTransactionResponse{
				Code: common.ReturnCode_ERROR,
				Msg:  "store withdraw transaction failed",
			}, nil
		}
	}
	return &riskcontroller.RiskWithdrawTransactionResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "submit withdraw transaction success",
	}, nil
}

func (rss *RiskServerWireServices) CheckOfflineWithdraw(ctx context.Context, request *riskcontroller.CheckOfflineTransactionRequest) (*riskcontroller.CheckOfflineTransactionResponse, error) {
	if rss.AccessToken != request.GetConsumerToken() {
		return &riskcontroller.CheckOfflineTransactionResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "invalid consumer token",
		}, nil
	}
	var txResults []*riskcontroller.CheckOfflineTxResult
	for _, tx := range request.GetCheckOfflineTxn() {
		if tx == nil {
			continue
		}
		verifiedKey := withdrawVerifiedKey(tx.GetRequestId(), tx.GetBusinessId(), tx.GetChainId())

		// 1. 幂等检查：该笔是否已校验通过。命中则直接返回“重复校验”，不再走后续逻辑。
		verifiedHash, err := rss.levelStore.Get([]byte(verifiedKey))
		if err != nil && err != errors.ErrNotFound {
			log.Error("get withdraw verified state failed", "err", err, "key", verifiedKey)
			return &riskcontroller.CheckOfflineTransactionResponse{
				Code:     common.ReturnCode_ERROR,
				Msg:      "get withdraw verified state failed",
				TxResult: txResults,
			}, nil
		}
		if err == nil {
			txResults = append(txResults, &riskcontroller.CheckOfflineTxResult{
				BusinessTxId: tx.GetRequestId(),
				Status:       riskcontroller.CheckOfflineTxStatus_CHECK_DUPLICATE,
				RiskKeyHash:  string(verifiedHash),
			})
			continue
		}

		// 2. 未校验过：读取已提交的原始交易，比对 hash。
		key := withdrawTxKey(tx.GetRequestId(), tx.GetBusinessId(), tx.GetChainId())
		storedValue, err := rss.levelStore.Get([]byte(key))
		if err != nil {
			if err == errors.ErrNotFound {
				// 未提交，校验失败，不写入幂等键，允许后续重试。
				txResults = append(txResults, &riskcontroller.CheckOfflineTxResult{
					BusinessTxId: tx.GetRequestId(),
					Status:       riskcontroller.CheckOfflineTxStatus_CHECK_FAILED,
				})
				continue
			}
			log.Error("get withdraw tx failed", "err", err, "key", key)
			return &riskcontroller.CheckOfflineTransactionResponse{
				Code:     common.ReturnCode_ERROR,
				Msg:      "get withdraw transaction failed",
				TxResult: txResults,
			}, nil
		}
		requestHash, err := hashWithdrawTx(tx)
		if err != nil {
			log.Error("hash request withdraw tx failed", "err", err)
			return &riskcontroller.CheckOfflineTransactionResponse{
				Code:     common.ReturnCode_ERROR,
				Msg:      "hash withdraw transaction failed",
				TxResult: txResults,
			}, nil
		}
		storedHash := hashBytes(storedValue)

		// 3. hash 不匹配：校验失败，不写入幂等键，允许重试。
		if requestHash != storedHash {
			txResults = append(txResults, &riskcontroller.CheckOfflineTxResult{
				BusinessTxId: tx.GetRequestId(),
				Status:       riskcontroller.CheckOfflineTxStatus_CHECK_FAILED,
				RiskKeyHash:  storedHash,
			})
			continue
		}

		// 4. 校验通过：落最终状态（幂等键），后续再次校验将返回“重复校验”。
		if err := rss.levelStore.Put([]byte(verifiedKey), []byte(storedHash)); err != nil {
			log.Error("store withdraw verified state failed", "err", err, "key", verifiedKey)
			return &riskcontroller.CheckOfflineTransactionResponse{
				Code:     common.ReturnCode_ERROR,
				Msg:      "store withdraw verified state failed",
				TxResult: txResults,
			}, nil
		}
		txResults = append(txResults, &riskcontroller.CheckOfflineTxResult{
			BusinessTxId: tx.GetRequestId(),
			Status:       riskcontroller.CheckOfflineTxStatus_CHECK_SUCCESS,
			RiskKeyHash:  storedHash,
		})
	}
	return &riskcontroller.CheckOfflineTransactionResponse{
		Code:     common.ReturnCode_SUCCESS,
		Msg:      "check offline withdraw transaction success",
		TxResult: txResults,
	}, nil
}

func (rss *RiskServerWireServices) CheckAmlAddress(ctx context.Context, request *riskcontroller.CheckAmlAddressRequest) (*riskcontroller.CheckAmlAddressResponse, error) {
	if rss.AccessToken != request.GetConsumerToken() {
		return &riskcontroller.CheckAmlAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "invalid consumer token",
		}, nil
	}
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

func (rss *RiskServerWireServices) SubmitTransactionFlow(ctx context.Context, request *riskcontroller.TransactionFlowRequest) (*riskcontroller.TransactionFlowResponse, error) {
	if rss.AccessToken != request.GetConsumerToken() {
		return &riskcontroller.TransactionFlowResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "invalid consumer token",
		}, nil
	}
	flow := &transactionFlowValue{
		DepositAmount:  request.GetDepositAmount(),
		WithdrawAmount: request.GetWithdrawAmount(),
		PositionAmount: request.GetPositionAmount(),
	}
	value, err := json.Marshal(flow)
	if err != nil {
		log.Error("marshal transaction flow failed", "err", err)
		return &riskcontroller.TransactionFlowResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "marshal transaction flow failed",
		}, nil
	}
	key := transactionFlowKey(request.GetRequestId(), request.GetUserAddress())
	if err := rss.levelStore.Put([]byte(key), value); err != nil {
		log.Error("store transaction flow failed", "err", err)
		return &riskcontroller.TransactionFlowResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "store transaction flow failed",
		}, nil
	}
	if err := rss.levelStore.Put([]byte(RiskKey), value); err != nil {
		log.Error("store latest transaction flow failed", "err", err)
		return &riskcontroller.TransactionFlowResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "store latest transaction flow failed",
		}, nil
	}
	return &riskcontroller.TransactionFlowResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "submit transaction flow success",
	}, nil
}

func (rss *RiskServerWireServices) CheckedTransactionFlow(ctx context.Context, request *riskcontroller.TransactionFlowCheckedRequest) (*riskcontroller.TransactionFlowCheckedResponse, error) {
	if rss.AccessToken != request.GetConsumerToken() {
		return &riskcontroller.TransactionFlowCheckedResponse{
			Code:     common.ReturnCode_ERROR,
			Msg:      "invalid consumer token",
			IsPassed: false,
		}, nil
	}
	currentFlow, err := rss.getTransactionFlow()
	if err != nil {
		log.Error("get transaction flow failed", "err", err)
		return &riskcontroller.TransactionFlowCheckedResponse{
			Code:     common.ReturnCode_ERROR,
			Msg:      "get transaction flow failed",
			IsPassed: false,
		}, nil
	}
	withdrawAmount, err := parseAmount(request.GetWithdrawAmount())
	if err != nil {
		log.Error("parse withdraw amount failed", "err", err, "withdraw_amount", request.GetWithdrawAmount())
		return &riskcontroller.TransactionFlowCheckedResponse{
			Code:     common.ReturnCode_ERROR,
			Msg:      "parse withdraw amount failed",
			IsPassed: false,
		}, nil
	}
	depositAmount, err := parseAmount(currentFlow.DepositAmount)
	if err != nil {
		log.Error("parse deposit amount failed", "err", err, "deposit_amount", currentFlow.DepositAmount)
		return &riskcontroller.TransactionFlowCheckedResponse{
			Code:     common.ReturnCode_ERROR,
			Msg:      "parse deposit amount failed",
			IsPassed: false,
		}, nil
	}
	positionAmount, err := parseAmount(currentFlow.PositionAmount)
	if err != nil {
		log.Error("parse position amount failed", "err", err, "position_amount", currentFlow.PositionAmount)
		return &riskcontroller.TransactionFlowCheckedResponse{
			Code:     common.ReturnCode_ERROR,
			Msg:      "parse position amount failed",
			IsPassed: false,
		}, nil
	}
	storedWithdrawAmount, err := parseAmount(currentFlow.WithdrawAmount)
	if err != nil {
		log.Error("parse stored withdraw amount failed", "err", err, "withdraw_amount", currentFlow.WithdrawAmount)
		return &riskcontroller.TransactionFlowCheckedResponse{
			Code:     common.ReturnCode_ERROR,
			Msg:      "parse stored withdraw amount failed",
			IsPassed: false,
		}, nil
	}
	isPassed := depositAmount+storedWithdrawAmount+positionAmount >= withdrawAmount
	return &riskcontroller.TransactionFlowCheckedResponse{
		Code:     common.ReturnCode_SUCCESS,
		Msg:      "check transaction flow success",
		IsPassed: isPassed,
	}, nil
}

func (rss *RiskServerWireServices) CheckChainTransactions(ctx context.Context, request *riskcontroller.CheckChainTransactionsRequest) (*riskcontroller.CheckChainTransactionsResponse, error) {
	if rss.AccessToken != request.GetConsumerToken() {
		return &riskcontroller.CheckChainTransactionsResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "invalid consumer token",
		}, nil
	}
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
