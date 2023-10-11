package dto

import "julo/domain"

type TransactionResponseData struct {
	Transactions []*domain.Transaction `json:"transactions"`
}

// EnableWalletResponse represents the response when enabling a wallet.
type TransactionResponse struct {
	Status string                   `json:"status"`
	Data   *TransactionResponseData `json:"data"`
}

// DTO for the request body
type DepositRequest struct {
	Amount      int    `json:"amount"`
	ReferenceID string `json:"reference_id"`
}

// DTO for the success response
type DepositResponse struct {
	Status string              `json:"status"`
	Data   DepositResponseData `json:"data"`
}

type DepositResponseData struct {
	Deposit *domain.Deposit `json:"deposit"`
}

// DTO for the success response
type WithdrawalResponse struct {
	Status string                 `json:"status"`
	Data   WithdrawalResponseData `json:"data"`
}

type WithdrawalResponseData struct {
	Withdrawal *domain.Withdrawal `json:"withdrawal"`
}
