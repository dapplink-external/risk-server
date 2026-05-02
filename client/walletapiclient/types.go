package walletapiclient

import "math/big"

type Base64TxWithSignature struct {
	Base64Tx  string
	Signature string
	PublicKey string
}

type BlockHeader struct {
	Hash       string
	ParentHash string
	Number     *big.Int
	Timestamp  uint64
}
