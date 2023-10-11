// julo/new_handlers.go
package handlers

import (
	"encoding/json"
	"julo/repository"
	"net/http"
)

// newHandlers is a struct that holds the database connection and can be embedded in handlers.
type Handlers struct {
	WalletRepository      repository.IWalletRepository
	TransactionRepository repository.ITransactionRepository
}

// NewNewHandlers initializes a newHandlers struct with the provided database connection.
func NewHandlers(walletRepository repository.IWalletRepository, transactionRepository repository.ITransactionRepository) *Handlers {
	return &Handlers{
		WalletRepository:      walletRepository,
		TransactionRepository: transactionRepository,
	}
}

func (h *Handlers) respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// ErrorResponse generates an error response in the specified format
func (h *Handlers) ErrorResponse(w http.ResponseWriter, status int, message string) {
	response := map[string]interface{}{
		"status": "fail",
		"data": map[string]interface{}{
			"error": message,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
