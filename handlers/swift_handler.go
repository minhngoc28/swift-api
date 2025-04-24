package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/minhngoc28/swift-api/models"
)

func GetSwiftCodeHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.ToUpper(strings.TrimSpace(c.Param("code")))
		fmt.Println("Looking for SWIFT code:", code)

		var swift models.SwiftCode
		err := db.Get(&swift, `SELECT * FROM swift_codes WHERE swift_code = $1`, code)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("SWIFT code %s not found", code)})
			return
		}

		// If it's a headquarter, also return a list of its branches
		if swift.IsHeadquarter {
			var branches []models.SwiftCode
			err = db.Select(&branches, `
                SELECT * FROM swift_codes
                WHERE is_headquarter = false AND LEFT(swift_code, 8) = LEFT($1, 8)
            `, code)
			if err == nil {
				c.JSON(http.StatusOK, gin.H{
					"address":       swift.Address,
					"bankName":      swift.BankName,
					"countryISO2":   swift.CountryISO2,
					"countryName":   swift.CountryName,
					"isHeadquarter": true,
					"swiftCode":     swift.SwiftCode,
					"branches":      branches,
				})
				return
			}
		}

		// If it's not a headquarter, return the data as-is
		c.JSON(http.StatusOK, swift)
	}
}

func GetAllSwiftCodesHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var codes []models.SwiftCode
		err := db.Select(&codes, "SELECT * FROM swift_codes")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SWIFT codes"})
			return
		}
		c.JSON(http.StatusOK, codes)
	}
}

func GetSwiftCodesByCountryHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		countryISO2 := strings.ToUpper(strings.TrimSpace(c.Param("iso2")))

		var codes []models.SwiftCode
		err := db.Select(&codes, "SELECT * FROM swift_codes WHERE country_iso2 = $1", countryISO2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SWIFT codes"})
			return
		}

		if len(codes) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("No SWIFT codes found for country %s", countryISO2)})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"countryISO2": countryISO2,
			"countryName": codes[0].CountryName,
			"swiftCodes":  codes,
		})
	}
}

func CreateSwiftCodeHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.SwiftCode

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
			return
		}

		query := `
			INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter)
			VALUES (:swift_code, :bank_name, :address, :country_iso2, :country_name, :is_headquarter)
		`

		_, err := db.NamedExec(query, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Insert failed: %v", err)})
			return
		}

		// ✅ Fix: Trả về đúng key theo JSON field
		c.JSON(http.StatusCreated, gin.H{
			"message":   "SWIFT code created successfully",
			"swiftCode": input.SwiftCode,
		})
	}
}

func DeleteSwiftCodeHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := strings.ToUpper(strings.TrimSpace(c.Param("code")))

		res, err := db.Exec("DELETE FROM swift_codes WHERE swift_code = $1", code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Delete failed: %v", err)})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "SWIFT code not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "SWIFT code deleted successfully", "swiftCode": code})
	}
}
