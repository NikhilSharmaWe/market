package main

import (
	"log"
	"os"

	"github.com/NikhilSharmaWe/market/api"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	if err := godotenv.Load("vars.env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	db := setupDB()

	app := api.NewApplication(db)

	mux := app.Router()

	log.Fatal(mux.Start(os.Getenv("ADDR")))
}

func setupDB() *gorm.DB {
	dbAddress := os.Getenv("DBADDRESS")
	db, err := gorm.Open(postgres.Open(dbAddress), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
