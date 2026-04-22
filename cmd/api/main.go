package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Gwilides/finance-tracker/configs"
	"github.com/Gwilides/finance-tracker/internal/account"
	"github.com/Gwilides/finance-tracker/internal/auth"
	"github.com/Gwilides/finance-tracker/internal/category"
	"github.com/Gwilides/finance-tracker/internal/user"
	"github.com/Gwilides/finance-tracker/pkg/db"
	"github.com/Gwilides/finance-tracker/pkg/middleware"
)

func App() http.Handler {
	config := configs.LoadConfig()
	router := http.NewServeMux()
	db := db.NewDb(&config.Db)

	//Repositories
	userRepository := user.NewUserRepository(db)
	accountRepository := account.NewAccountRepository(db)
	categoryRepository := category.NewCategoryRepository(db)

	//Services
	authService := auth.NewAuthService(userRepository)
	accountService := account.NewAccountService(&account.AccountServiceDeps{
		UserRepository:    userRepository,
		AccountRepository: accountRepository,
	})
	categoryService := category.NewCategoryService(&category.CategoryServiceDeps{
		UserRepository:     userRepository,
		CategoryRepository: categoryRepository,
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
	category.NewCategoryHandler(router, &category.CategoryHandlerDeps{
		Service: categoryService,
		Config:  &config.Auth,
	})
	return router
}

func main() {
	server := http.Server{
		Addr:    "localhost:8090",
		Handler: middleware.Logger(App()),
	}
	fmt.Println("Server is listening on port 8090")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
