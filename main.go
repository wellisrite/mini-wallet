// julo/main.go
package main

import (
	"database/sql"
	"fmt"
	"julo/config"
	"julo/handlers"
	"julo/middleware"
	"julo/repository"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the configuration
	config.Initialize()
	cfg := config.AppConfig
	db := setupDatabase()
	router := setupRouter(db, cfg)
	// Access the configuration

	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(port, router)
}

func setupRouter(db *sql.DB, cfg *config.Config) *mux.Router {
	// Create a Gorilla Mux router
	router := mux.NewRouter()

	// Create an instance of newHandlers with the database connection
	walletRepository := repository.NewWalletRepository(db)
	transactionsRepository := repository.NewTransactionsRepository(db)

	handlers := handlers.NewHandlers(walletRepository, transactionsRepository)

	// Create a subrouter with the "v1" prefix
	v1Router := router.PathPrefix("/api/v1").Subrouter()
	// Define your routes within the "/api/v1" prefix
	v1Router.HandleFunc("/init", handlers.InitializeAccount).Methods("POST")

	walletRouter := v1Router.PathPrefix("/wallet").Subrouter()

	walletRouter.Use(middleware.JWTMiddleware)

	walletRouter.HandleFunc("", handlers.EnableWalletHandler).Methods("POST")
	walletRouter.HandleFunc("", handlers.DisableWalletHandler).Methods("PATCH")

	walletRouter.HandleFunc("", handlers.ViewWalletBalanceHandler).Methods("GET")

	// transactions
	walletRouter.HandleFunc("/transactions", handlers.ViewWalletTransactionsHandler).Methods("GET")
	walletRouter.HandleFunc("/deposits", handlers.DepositHandler).Methods("POST")
	walletRouter.HandleFunc("/withdrawals", handlers.WithdrawalHandler).Methods("POST")

	return router
}
