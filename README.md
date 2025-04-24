# Remitly 2025 - Home Exercise

This project is a solution to the Remitly Internship 2025 Home Exercise.

The goal of the exercise is to:
- Parse SWIFT code data from a CSV file
- Store it in a PostgreSQL database
- Expose the data via a RESTful API built with Go and Gin
- Support retrieving, inserting, and deleting SWIFT code records

All services run inside Docker containers using Docker Compose.

---

## Tech Stack

- Go (Golang)
- PostgreSQL
- Docker & Docker Compose
- Go standard `net/http` + `testing` packages

---

## Project Structure

- `main.go` – App entry point, router setup
- `Dockerfile` – Go build instructions
- `docker-compose.yml` – Container setup for API + PostgreSQL
- `init.sql` – SQL script to initialize DB schema
- `swift.csv` – Source data file
- `handlers/swift_handler.go` – HTTP handler for GET endpoint
- `models/swift.go` – SWIFT code struct & DB logic
- `utils/parser.go` – Logic to parse CSV and insert into DB
- `handlers/swift_handler_test.go` – Unit tests for API handler
- `README.md` – You’re reading it

---

## How to Run

### Prerequisites

- Docker
- Docker Compose

### Quick Start

```bash
git clone https://github.com/minhngoc28/swift-api.git
cd swift-api
docker-compose up --build
```
This command will:
- Build the Go backend service
- Start a PostgreSQL database (with volume persistence)
- Parse and insert swift.csv data into the database (only on first run)
- Start a web server on http://localhost:8080

## API Endpoints
Base URL: `http://localhost:8080`

---

### GET /swift-codes
Get all SWIFT codes.

**Example:**

```bash
curl http://localhost:8080/swift-codes
```

### GET /swift-codes/:code
Get detail of a specific SWIFT code.

**Example:**
```bash
curl http://localhost:8080/swift-codes/TPEOPLPWOBP
```

### GET /swift-codes/country/:iso2
Get all SWIFT codes from a specific country using ISO2 country code.

**Example:**
```bash
curl http://localhost:8080/swift-codes/country/PL
```

### POST /swift-codes
Add a new SWIFT code.

**Example:**
```bash
curl -X POST http://localhost:8080/swift-codes \
-H "Content-Type: application/json" \
-d '{
  "swiftCode": "TESTPLPWXXX",
  "bankName": "TEST BANK",
  "address": "123 Test Street, Warsaw",
  "countryISO2": "PL",
  "countryName": "POLAND",
  "isHeadquarter": true
}'

```

### DELETE /swift-codes/:code
Delete a SWIFT code

**Example:**
```bash
curl -X DELETE http://localhost:8080/swift-codes/TESTPLPWXXX
```

## Notes & Assumptions

- `swift.csv` is parsed and inserted only on the first run (or when DB volume is reset).
- `swift_code` is used as the primary key and assumed to be globally unique.
- A SWIFT code is treated as a **headquarter** if `is_headquarter = true`, and its **branches** share the same 8-character prefix.
- `GET /swift-codes/:code` returns additional `branches` only if the code belongs to a headquarter.
- Error handling is basic and for demonstration purposes. In production, more granular handling should be added.
- The app runs in `debug` mode for easier local development. Use `GIN_MODE=release` in production.
- Volume persistence ensures data is not lost between runs unless explicitly removed using `--volumes`.

## Test Instructions
### Manual testing with curl

You can test key endpoints directly with `curl`:

```bash
# Get all codes
curl http://localhost:8080/swift-codes

# Get a specific code
curl http://localhost:8080/swift-codes/TPEOPLPWOBP

# Get all codes from a country
curl http://localhost:8080/swift-codes/country/PL

# Create a new code
curl -X POST http://localhost:8080/swift-codes \
-H "Content-Type: application/json" \
-d '{
  "swiftCode": "TESTPLPWXXX",
  "bankName": "TEST BANK",
  "address": "123 Test Street, Warsaw",
  "countryISO2": "PL",
  "countryName": "POLAND",
  "isHeadquarter": true
}'

# Delete a code
curl -X DELETE http://localhost:8080/swift-codes/TESTPLPWXXX
```

## Running Unit Tests

### 1. Ensure Docker container is running and DB `swift` is ready:
```bash
docker-compose up -d
```

### 2. Ensure the database contains required schema:
```bash
docker exec -i swift-api-swift-db-1 psql -U postgres -d swift < init.sql
```

### 3. Export DB connection URL:
```bash
export TEST_DB_URL=postgres://postgres:mysecretpassword@localhost:5432/swift?sslmode=disable
```

### 4. Run tests:
```bash
go test ./handlers
```

> Tests will automatically create & cleanup test records.

---

## Clean Up (Optional)
```bash
docker-compose down -v
```

---

## Author
**minhngoc28** – https://github.com/minhngoc28
