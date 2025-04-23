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
			c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
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
