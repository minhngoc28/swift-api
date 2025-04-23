package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type SwiftRecord struct {
	SwiftCode     string `db:"swift_code"`
	BankName      string `db:"bank_name"`
	Address       string `db:"address"`
	CountryISO2   string `db:"country_iso2"`
	CountryName   string `db:"country_name"`
	IsHeadquarter bool   `db:"is_headquarter"`
}

func ParseAndInsertCSV(db *sqlx.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Đọc dòng tiêu đề
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading header: %w", err)
	}
	fmt.Println("CSV Headers:", headers)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading record: %w", err)
		}

		swiftCode := strings.ToUpper(strings.TrimSpace(record[1]))
		bankName := strings.TrimSpace(record[3])
		address := strings.TrimSpace(record[4])
		countryISO2 := strings.ToUpper(strings.TrimSpace(record[0]))
		countryName := strings.ToUpper(strings.TrimSpace(record[6]))

		isHQ := strings.HasSuffix(swiftCode, "XXX")

		swift := SwiftRecord{
			SwiftCode:     swiftCode,
			BankName:      bankName,
			Address:       address,
			CountryISO2:   countryISO2,
			CountryName:   countryName,
			IsHeadquarter: isHQ,
		}

		// 💡 Log mã SWIFT trước khi insert
		fmt.Println("Inserting SWIFT code:", swift.SwiftCode)

		_, err = db.NamedExec(`INSERT INTO swift_codes 
			(swift_code, bank_name, address, country_iso2, country_name, is_headquarter)
			VALUES (:swift_code, :bank_name, :address, :country_iso2, :country_name, :is_headquarter)
			ON CONFLICT (swift_code) DO NOTHING`, &swift)
		if err != nil {
			fmt.Println("Insert error:", err)
		}
	}

	fmt.Println("Finished importing CSV data")
	return nil
}
