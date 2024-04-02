package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/softmaxer/localflow/api"
)

func main() {
	godotenv.Load(".env")
	dbPath := os.Getenv("DB_PATH")
	fmt.Println("initializing DB with: ", dbPath)
	router := api.Router(dbPath)
	router.Run(":8080")
}
