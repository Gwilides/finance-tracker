package main

import (
	"github.com/Gwilides/finance-tracker/configs"
	"github.com/Gwilides/finance-tracker/internal/account"
	"github.com/Gwilides/finance-tracker/internal/category"
	"github.com/Gwilides/finance-tracker/internal/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	config := configs.LoadConfig()
	db, err := gorm.Open(postgres.Open(config.Db.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&user.User{}, &account.Account{}, &category.Category{})
}
