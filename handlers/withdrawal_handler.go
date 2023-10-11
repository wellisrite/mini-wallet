package handlers

import (
	"fmt"
	"julo/domain"
	"julo/domain/dto"
	"net/http"
	"strconv"
)

func (h *Handlers) WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data from the request
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

	// Extract form values
	amountStr := r.FormValue("amount")
	referenceID := r.FormValue("reference_id")

	// Convert the amount to an integer
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	uniqueUuid, err := h.TransactionRepository.GenerateUniqueUUID()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to process withdrawal", http.StatusInternalServerError)
		return
	}

	// Create a new withdrawal object
	withdrawal := &domain.Withdrawal{
		ID:          uniqueUuid,
		WithdrawnBy: customerID,
		Status:      "success",
		Amount:      amount,
		ReferenceID: referenceID,
	}

	// Call the service to process the withdrawal
	err = h.TransactionRepository.AddWithdrawal(withdrawal)
	if err != nil {
		fmt.Println(err)

		h.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return the success response
	response := dto.WithdrawalResponse{
		Status: "success",
		Data: dto.WithdrawalResponseData{
			Withdrawal: withdrawal,
		},
	}

	h.respondWithJSON(w, http.StatusOK, response)
}
