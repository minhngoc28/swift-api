package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupRouter(db *sqlx.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/swift-codes/:code", GetSwiftCodeHandler(db))
	r.POST("/swift-codes", CreateSwiftCodeHandler(db))
	r.DELETE("/swift-codes/:code", DeleteSwiftCodeHandler(db))
	return r
}

func connectTestDB(t *testing.T) *sqlx.DB {
	dsn := os.Getenv("TEST_DB_URL")
	if dsn == "" {
		dsn = "postgres://postgres:mysecretpassword@localhost:5432/swift?sslmode=disable"
	}
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func TestGetSwiftCodeHandler_ExistingCode(t *testing.T) {
	db := connectTestDB(t)

	_, _ = db.Exec(`
		INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter)
		VALUES ('TPEOPLPWOBP', 'PEKAO TOWARZYSTWO FUNDUSZY INWESTYCYJNYCH SPOLKA AKCYJNA', '123 Warsaw', 'PL', 'POLAND', false)
		ON CONFLICT DO NOTHING
	`)
	t.Cleanup(func() {
		db.Exec("DELETE FROM swift_codes WHERE swift_code = 'TPEOPLPWOBP'")
	})

	router := setupRouter(db)
	req, _ := http.NewRequest("GET", "/swift-codes/TPEOPLPWOBP", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "PEKAO TOWARZYSTWO FUNDUSZY")
}

func TestGetSwiftCodeHandler_NotFound(t *testing.T) {
	db := connectTestDB(t)
	router := setupRouter(db)

	req, _ := http.NewRequest("GET", "/swift-codes/INVALID123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
}

func TestCreateSwiftCodeHandler(t *testing.T) {
	db := connectTestDB(t)
	router := setupRouter(db)

	body := `{
		"swiftCode": "UNITTEST01",
		"bankName": "Unit Test Bank",
		"address": "Test Address",
		"countryISO2": "PL",
		"countryName": "POLAND",
		"isHeadquarter": false
	}`
	t.Cleanup(func() {
		db.Exec("DELETE FROM swift_codes WHERE swift_code = 'UNITTEST01'")
	})

	req, _ := http.NewRequest("POST", "/swift-codes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Contains(t, w.Body.String(), "UNITTEST01")
}

func TestDeleteSwiftCodeHandler(t *testing.T) {
	db := connectTestDB(t)
	router := setupRouter(db)

	_, _ = db.Exec(`
		INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter)
		VALUES ('UNITTEST01', 'Test Bank', 'Somewhere', 'PL', 'POLAND', false)
		ON CONFLICT DO NOTHING
	`)
	t.Cleanup(func() {
		db.Exec("DELETE FROM swift_codes WHERE swift_code = 'UNITTEST01'")
	})

	req, _ := http.NewRequest("DELETE", "/swift-codes/UNITTEST01", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "UNITTEST01")
}
