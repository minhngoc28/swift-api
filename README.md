# SWIFT Code API – Remitly Internship 2025

This is a RESTful API written in Go that parses SWIFT/BIC code data from a CSV file and stores it in a PostgreSQL database. The API exposes endpoints to look up bank details by SWIFT code.

This project was developed as part of the **Remitly Internship 2025 Home Exercise**.

---

## Features

- Parses `swift.csv` containing SWIFT/BIC codes and related bank information
- Loads data into a PostgreSQL database
- Exposes a REST API to query bank info by SWIFT code
- Containerized using Docker and Docker Compose
- Includes unit tests and error handling for edge cases

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
This will:
- Start a PostgreSQL container
- Run init.sql to create the swift_codes table
- Build and start the Go service on port 8080
- Parse swift.csv and insert records into the database

### Test API

```bash
curl http://localhost:8080/swift/ABCDUS33
```
