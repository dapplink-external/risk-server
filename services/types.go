package services

type TransactionBody struct {
	RequestId    string `json:"request_id"`
	FromAddress  string `json:"from_address"`
	ToAddress    string `json:"to_address"`
	Amount       string `json:"amount"`
	TokenAddress string `json:"token_address"`
}

type transactionFlowValue struct {
	DepositAmount  string `json:"deposit_amount"`
	WithdrawAmount string `json:"withdraw_amount"`
	PositionAmount string `json:"position_amount"`
}

type canonicalWithdrawTx struct {
	RequestId       string `json:"request_id"`
	BusinessTxId    string `json:"business_tx_id"`
	ChainId         string `json:"chain_id"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	ContractAddress string `json:"contract_address"`
	TokenId         string `json:"token_id"`
	TokenMeta       string `json:"token_meta"`
}
