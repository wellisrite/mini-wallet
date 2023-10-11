// julo/handlers/wallet_handlers.go
package handlers

import (
	"fmt"
	"julo/domain"
	"julo/domain/dto"

	"net/http"
)

func (h *Handlers) InitializeAccount(w http.ResponseWriter, r *http.Request) {
	// Retrieve the customer ID from the request (assuming it's in the form data)
	customerXID := r.FormValue("customer_xid")
	if customerXID == "" {
		// Customer XID is missing, return a failed response
		response := dto.WalletResponseError{
			Status: "fail",
			Data: dto.WallerErrorResponseData{
				Error: map[string][]string{
					"customer_xid": []string{"Missing data for required field."},
				},
			},
		}
		fmt.Println("missing field")
		h.respondWithJSON(w, http.StatusBadRequest, response)
		return
	}

	// If the customer exists, generate a new JWT token with the customer ID
	token, err := domain.GenerateJWTToken(customerXID)
	if err != nil {
		// Handle token generation error, e.g., return an error response
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	// Create a success response with the wallet data.
	response := dto.WalletResponse{
		Status: "success",
		Data:   &dto.InitializationResponse{Token: token},
	}

	// Check if the customer ID already exists in the database
	if exist, err := h.WalletRepository.WalletExists(customerXID); !exist {
		if err != nil {
			// Handle database insertion error, e.g., return an error response
			h.ErrorResponse(w, http.StatusInternalServerError, "Failed to initialize wallet")
			fmt.Println(err)
			return
		}

		// If the customer does not exist, insert mandatory fields into the database
		wallet := domain.InitializeWallet(customerXID)
		if err := h.WalletRepository.InsertWallet(wallet); err != nil {
			// Handle database insertion error, e.g., return an error response
			h.ErrorResponse(w, http.StatusInternalServerError, "Failed to insert wallet")
			fmt.Println(err)
			return
		}

		// Return the JWT token in the response
		h.respondWithJSON(w, http.StatusCreated, response)
		return
	}

	// Return the JWT token in the response
	h.respondWithJSON(w, http.StatusOK, response)
	return
}
