package services

type TransactionBody struct {
	RequestId    string `json:"request_id"`
	FromAddress  string `json:"from_address"`
	ToAddress    string `json:"to_address"`
	Amount       string `json:"amount"`
	TokenAddress string `json:"token_address"`
}
