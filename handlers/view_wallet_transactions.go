package handlers

import (
	"julo/domain"
	"julo/domain/dto" // Import the relevant DTO package and other dependencies.
	"net/http"
)

// ViewWalletTransactionsHandler retrieves and returns the wallet transactions for the authenticated customer.
func (h *Handlers) ViewWalletTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract customerID from the request context (using your authentication middleware).
	customerID := r.Context().Value("customerXID").(string)

	// Fetch the wallet information including balance from the repository.
	wallet, err := h.WalletRepository.GetWallet(customerID)
	if err != nil {
		h.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Check if the wallet is disabled.
	if wallet.Status == domain.STATUS_DISABLED {
		// Wallet is disabled, return a specific error response
		response := dto.WalletResponseError{
			Status: "fail",
			Data:   dto.WallerErrorResponseString{Error: "Wallet disabled"},
		}

		h.respondWithJSON(w, http.StatusNotFound, response)
		return
	}

	// Fetch the wallet transactions from the repository (you will need to implement this function).
	transactions, err := h.TransactionRepository.GetTransactionsByCustomerID(customerID)
	if err != nil {
		h.ErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Create a success response with the transactions data.
	response := dto.TransactionResponse{
		Status: "success",
		Data: &dto.TransactionResponseData{
			Transactions: transactions, // Set the transactions data.
		},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}
