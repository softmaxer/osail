package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/softmaxer/osail/handler"
)

func main() {
	godotenv.Load(".env")
	dbPath := os.Getenv("DB_PATH")
	fmt.Println("initializing DB with: ", dbPath)
	router := handler.Router(dbPath)
	router.Run(":8080")
}
