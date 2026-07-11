package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/the-web3/mock-risk-server/protobuf/riskcontroller"
)

func parseAmount(amount string) (int64, error) {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return 0, nil
	}
	return strconv.ParseInt(amount, 10, 64)
}

func withdrawTxKey(requestID string, businessTxID string, chainID string) string {
	return fmt.Sprintf("%s:%s:%s:%s", withdrawTxKeyPrefix, requestID, businessTxID, chainID)
}

// withdrawVerifiedKey 标记某笔离线提现已校验通过的幂等键。
// 校验通过后写入（值为通过时的 hash），后续再次校验直接命中返回“重复校验”。
// 校验失败时不写入，因此允许重试。
func withdrawVerifiedKey(requestID string, businessTxID string, chainID string) string {
	return fmt.Sprintf("%s:%s:%s:%s", withdrawVerifiedKeyPrefix, requestID, businessTxID, chainID)
}

func transactionFlowKey(requestID string, userAddress string) string {
	return fmt.Sprintf("%s:%s:%s", transactionFlowKeyPrefix, requestID, userAddress)
}

func hashWithdrawTx(tx *riskcontroller.WithdrawTxList) (string, error) {
	return hashCanonicalWithdrawTx(toCanonicalWithdrawTx(tx))
}

func hashCanonicalWithdrawTx(tx canonicalWithdrawTx) (string, error) {
	value, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	return hashBytes(value), nil
}

func hashBytes(value []byte) string {
	hash := sha256.Sum256(value)
	return hex.EncodeToString(hash[:])
}

func (rss *RiskServerWireServices) getTransactionFlow() (*transactionFlowValue, error) {
	data, err := rss.levelStore.Get([]byte(RiskKey))
	if err != nil {
		if err == errors.ErrNotFound {
			return &transactionFlowValue{}, nil
		}
		return nil, err
	}
	var flow transactionFlowValue
	if err := json.Unmarshal(data, &flow); err != nil {
		return nil, err
	}
	return &flow, nil
}

func toCanonicalWithdrawTx(tx *riskcontroller.WithdrawTxList) canonicalWithdrawTx {
	return canonicalWithdrawTx{
		RequestId:       tx.GetRequestId(),
		BusinessTxId:    tx.GetBusinessId(),
		ChainId:         tx.GetChainId(),
		From:            tx.GetFrom(),
		To:              tx.GetTo(),
		Value:           tx.GetValue(),
		ContractAddress: tx.GetContractAddress(),
		TokenId:         tx.GetTokenId(),
		TokenMeta:       tx.GetTokenMeta(),
	}
}
