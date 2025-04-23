package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minhngoc28/swift-api/handlers"
	"github.com/minhngoc28/swift-api/utils"
)

func main() {
	dsn := "postgres://postgres:mysecretpassword@swift-db:5432/swift?sslmode=disable"
	if dsn == "" {
		log.Fatal("Environment variable DB_URL is missing")
	}

	var db *sql.DB
	var err error

	// Retry connection
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Failed to open DB: %v", err)
		} else {
			err = db.Ping()
			if err == nil {
				log.Println("✅ Connected to DB")
				break
			}
			log.Printf("Failed to ping DB: %v", err)
		}

		log.Printf("⏳ Waiting for DB to be ready... attempt %d/10", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	dbx := sqlx.NewDb(db, "postgres")

	// Insert CSV data
	if err := utils.ParseAndInsertCSV(dbx, "swift.csv"); err != nil {
		log.Fatalf("Failed to import CSV: %v", err)
	}

	// Start API server
	r := gin.Default()
	r.GET("/swift-codes", handlers.GetAllSwiftCodesHandler(dbx))
	r.GET("/swift-codes/:code", handlers.GetSwiftCodeHandler(dbx))

	log.Println("🚀 Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
