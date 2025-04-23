package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/minhngoc28/swift-api/handlers"
	"github.com/minhngoc28/swift-api/utils"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	fmt.Println("DB_URL is:", dbURL)

	if dbURL == "" {
		log.Fatalln("Missing DB_URL environment variable")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("Connected to PostgreSQL!")

	// Import CSV
	err = utils.ParseAndInsertCSV(db, "swift.csv")
	if err != nil {
		log.Fatalln("Failed to import CSV:", err)
	}

	r := gin.Default()
	r.GET("/v1/swift-codes/:code", handlers.GetSwiftCodeHandler(db))
	r.GET("/v1/swift-codes", handlers.GetAllSwiftCodesHandler(db))

	fmt.Println("🚀 Server running at http://localhost:8080")
	r.Run(":8080")
}
