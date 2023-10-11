package handlers

import (
	"fmt"
	"julo/domain"
	"julo/domain/dto"
	"net/http"
	"strconv"
)

func (h *Handlers) DepositHandler(w http.ResponseWriter, r *http.Request) {
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

	// Parse and validate the request body
	var request dto.DepositRequest
	request.Amount, err = strconv.Atoi(r.FormValue("amount"))
	if err != nil {
		// Handle repository error and return a failure response
		fmt.Println(err)
		h.ErrorResponse(w, http.StatusInternalServerError, "Error on server")
		return
	}
	request.ReferenceID = r.FormValue("reference_id")

	// Call the repository function to add the deposit
	deposit, err := h.TransactionRepository.AddDeposit(customerID, request.ReferenceID, request.Amount)
	if err != nil {
		h.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return a success response
	response := dto.DepositResponse{
		Status: "success",
		Data: dto.DepositResponseData{
			Deposit: deposit,
		},
	}
	h.respondWithJSON(w, http.StatusCreated, response)
}
