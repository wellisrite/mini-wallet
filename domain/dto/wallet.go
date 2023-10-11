package dto

import "julo/domain"

type WalletResponseData struct {
	Wallet *domain.Wallet `json:"wallet"`
}

// EnableWalletResponse represents the response when enabling a wallet.
type EnableWalletResponse struct {
	Status string              `json:"status"`
	Data   *WalletResponseData `json:"data"`
}

// InitializationResponse represents the response when initializing a wallet.
type InitializationResponse struct {
	Token string `json:"token"`
}

// InitializeWalletResponseError represents the error response for initializing the wallet.
type WalletResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

// InitializeWalletResponseError represents the error response for initializing the wallet.
type WalletResponseError struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type WallerErrorResponseData struct {
	Error map[string][]string `json:"error"`
}

type WallerErrorResponseString struct {
	Error string `json:"error"`
}

// ViewWalletBalanceResponse represents the success response for viewing the wallet balance.
type ViewWalletBalanceResponse struct {
	Status string              `json:"status"`
	Data   *WalletResponseData `json:"data"`
}

// Request structure for disabling the wallet
type DisableWalletRequest struct {
	IsDisabled bool `form:"is_disabled" binding:"required"`
}
