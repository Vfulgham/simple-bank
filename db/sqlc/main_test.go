package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// global variable so it can be used in all unit tests
// use this to connect to db (it has DBTX object)
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	// get connection to db
	testDB, err = sql.Open(getEnvVar("DBDRIVER"), getEnvVar("DBSOURCE"))
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// create Queries object
	testQueries = New(testDB)

	// m.Run runs tests
	// tells test runner via Exit command if test failed or passed
	os.Exit(m.Run())
}

// returns the value of the key from env file
func getEnvVar(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
