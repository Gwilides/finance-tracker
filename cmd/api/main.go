package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gwilides/finance-tracker/configs"
	"github.com/Gwilides/finance-tracker/internal/account"
	"github.com/Gwilides/finance-tracker/internal/auth"
	"github.com/Gwilides/finance-tracker/internal/user"
	"github.com/Gwilides/finance-tracker/pkg/db"
)

func App() http.Handler {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	db := db.NewDb(&config.Db)

	//Repositories
	userRepository := user.NewUserRepository(db)
	accountRepository := account.NewAccountRepository(db)

	//Services
	authService := auth.NewAuthService(userRepository)
	accountService := account.NewAccountService(&account.AccountServiceDeps{
		UserRepository:    userRepository,
		AccountRepository: accountRepository,
	})

	//Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		AuthService: authService,
		Config:      &config.Auth,
	})
	account.NewAccountHandler(router, &account.AccountHandlerDeps{
		Service: accountService,
		Config:  &config.Auth,
	})
	return router
}

func main() {
	server := http.Server{
		Addr:    "localhost:8090",
		Handler: App(),
	}
	fmt.Println("Server is listening on port 8090")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
