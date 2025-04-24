package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minhngoc28/swift-api/handlers"
	"github.com/minhngoc28/swift-api/utils"
)

func main() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("Environment variable DB_URL is missing")
	}

	var db *sql.DB
	var err error

	for i := 0; i < 20; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Failed to open DB: %v", err)
		} else {
			err = db.Ping()
			if err == nil {
				log.Println("Connected to DB")
				break
			}
			log.Printf("Failed to ping DB: %v", err)
		}
		log.Printf("Waiting for DB to be ready... attempt %d/20", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")

	if err := utils.ParseAndInsertCSV(dbx, "swift.csv"); err != nil {
		log.Fatalf("Failed to import CSV: %v", err)
	}

	r := gin.Default()
	r.GET("/swift-codes", handlers.GetAllSwiftCodesHandler(dbx))
	r.GET("/swift-codes/:code", handlers.GetSwiftCodeHandler(dbx))
	r.GET("/swift-codes/country/:iso2", handlers.GetSwiftCodesByCountryHandler(dbx))
	r.POST("/swift-codes", handlers.CreateSwiftCodeHandler(dbx))
	r.DELETE("/swift-codes/:code", handlers.DeleteSwiftCodeHandler(dbx))

	log.Println("Server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
