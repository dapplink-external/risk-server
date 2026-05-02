package database

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type Transactions struct {
	Guid         string    `gorm:"primaryKey; column:guid"`
	RequestId    string    `gorm:"column:request_id"`
	FromAddress  string    `gorm:"column:from_address"`
	ToAddress    string    `gorm:"column:to_address"`
	TokenAddress string    `gorm:"column:token_address"`
	TokenId      string    `gorm:"column:token_id"`
	TokenMeta    string    `gorm:"column:token_meta"`
	Fee          *big.Int  `gorm:"column:fee"`
	Amount       *big.Int  `gorm:"column:amount"`
	Status       string    `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

type TransactionsView interface {
	QueryTransactionsByRequestId(string) (*Transactions, error)
}

type TransactionsDB interface {
	TransactionsView

	StoreTransactions(string, []*Transactions) error
}

type transactionsDB struct {
	gorm *gorm.DB
}

func NewTransactionsDB(db *gorm.DB) TransactionsDB {
	return &transactionsDB{gorm: db}
}

func (db *transactionsDB) StoreTransactions(businessId string, txs []*Transactions) error {
	result := db.gorm.Table("transaction_"+businessId).CreateInBatches(&txs, len(txs))
	return result.Error
}

func (db *transactionsDB) QueryTransactionsByRequestId(requestId string) (*Transactions, error) {
	var txInfo Transactions
	if err := db.gorm.Table("transaction").Where("requestId=?", requestId).Take(&txInfo).Error; err != nil {
		log.Error("query transactions error", "requestId", requestId, "err", err)
		return nil, err
	}
	return &txInfo, nil
}
