package pgxtesting_test

import (
	"context"
	"log"
	"testing"

	"github.com/go-pa/pgxtesting"
	"github.com/golang-migrate/migrate"
)

func Example() {
	// the t argument in your test function
	var t *testing.T

	{

		// You can create a test database from a database url.
		// It will automatically clean up (drop the database) after the test has been run.
		// The CreateTestDatabase function will fatal the test if the
		// connection/creation to the datavbase fails there are no errors to
		// handle.
		_ = pgxtesting.CreateTestDatabase(t, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&pool_max_conns=500")

	}

	{
		// You can also set a global pgxtesting default postgres url in a
		// init() function in one of your test files.
		pgxtesting.SetDefaultURL("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&pool_max_conns=500")

		// and then only use this function to get a new pool for each test.
		_ = pgxtesting.CreateTestDatabaseEnv(t)

	}

	{

		pool := pgxtesting.CreateTestDatabaseEnv(t)

		// And you can ten do something with the database...
		row := pool.QueryRow(context.Background(), "SELECT current_database()")
		var v string
		err := row.Scan(&v)
		if err != nil {
			t.Fatal(err)
		}

		// You can close the pool early if you want to free up database
		// connections before the test has run and the testing package does it
		// automatically for you. Close() drops the database and closes the
		// connection.
		// If you have a lot of subtests and each creates it's own test
		// database this might be a good idea because the testing package does
		// not clean up before all sub tests are done so you might exceed
		// maximum number of connections if you have many subtests.
		pool.Close()

	}

	{
		pool := pgxtesting.CreateTestDatabaseEnv(t)

		// If you need an URL for the postgres test database you use pool.URL.
		// If you need to strip pgxpool specific parameters from it to use it
		// with the std library sql package or go-migrate, use the ConnURL()
		// method.
		m, _ := migrate.New(
			"file://migrations",
			pool.URL.ConnURL().String())

		if err := m.Up(); err != nil {
			log.Fatal(err)
		}

	}
}
