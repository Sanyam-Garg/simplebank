// The connection to the databse is defined in this file
package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	// Must add "_" in front of this package. We do not call any functions from this, but this is a driver which needs to be imported.
	_ "github.com/lib/pq"
)

const (
	// Need to get the go driver for this.
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

// This is the main entry point for all tests in golang
func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	// If connection is successful, returns a DB or DBTX object.
	testQueries = New(conn)

	// m.Run() starts the main test and exits returning an exit code. The os.Exit then performs operations according to the code.
	os.Exit(m.Run())
}
