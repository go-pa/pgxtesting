package pgxtesting

import (
	"context"
	"os"
	"strings"
	"testing"
)

var testURL = "postgres://test:test@localhost:45432/test?sslmode=disable&pool_max_conns=1000"

func TestPoolURL(t *testing.T) {
	const testu = "postgres://user:pass@host:500/dbname?pool_max_conns=100"
	u, err := NewPoolURL(testu)
	if err != nil {
		t.Error(err)
	}

	equal := func(a, b string) {
		t.Helper()
		if a != b {
			t.Errorf("not equal: '%v' '%v'", a, b)

		}
	}

	equal(u.Name(), "dbname")
	equal(u.String(), testu)
	equal(u.ConnURL().String(), "postgres://user:pass@host:500/dbname")

}

func TestCreateTestDB(t *testing.T) {
	pool := CreateTestDatabase(t, testURL)

	row := pool.QueryRow(context.Background(), "SELECT current_database()")
	var v string
	err := row.Scan(&v)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(v, "go_test_") {
		t.Errorf("%s does not start with go_test_", v)
	}
	if v != pool.URL.Name() {
		t.Errorf("'%s' != '%s'", v, pool.URL.Name())
	}
}

func TestClosePool(t *testing.T) {
	t.Parallel()
	pool := CreateTestDatabase(t, testURL)

	// this should not generate any logs
	defer pool.Close()
	if pool.closed {
		t.Error("pool is closed")
	}
	pool.Close()
	if !pool.closed {
		t.Error("pool is not closed")
	}
	pool.Close()
	if !pool.closed {
		t.Error("pool is not closed")
	}
}

func TestClosePoolStillOpenConnections(t *testing.T) {
	t.Parallel()

	pool := CreateTestDatabase(t, testURL)

	pool2, err := connectPostgres(pool.URL)
	if err != nil {
		t.Fatal(err)
	}

	pool.Pool.Close()
	pool.closed = true

	err = pool.dropTestDB() // this tames some time to run

	pool2.Close()
	if err == nil {
		t.Fatal("expected error due to connection still being open")
	}
	if !strings.HasPrefix(err.Error(), "pgxtesting.Pool.Cleanup error running 'drop database ") ||
		!strings.Contains(err.Error(), "SQLSTATE 55006") {
		t.Fatalf("got wrong error message: %v", err)
	}

	err = pool.dropTestDB()
	if err != nil {
		t.Fatal(err)
	}

	pool.Close()
	pool.Close()

}

func TestGlobalConfig(t *testing.T) {
	t.Parallel()
	equal := func(a, b string) {
		t.Helper()
		if a != b {
			t.Errorf("not equal: '%v' '%v'", a, b)

		}
	}
	os.Setenv("PGURL", "")
	equal(defaultURL, GetDefaultURL())
	equal(GetURL(), GetDefaultURL())
	os.Setenv("PGURL", "postgres://foo")
	equal(GetURL(), "postgres://foo")
	SetDefaultURL("postgres://foo")
	equal(GetURL(), "postgres://foo")
	equal(defaultURL, GetDefaultURL())

	SetEnvName("DBURL")
	os.Setenv("DBURL", "")
	equal(GetURL(), GetDefaultURL())
	os.Setenv("DBURL", "postgres://bar")
	equal(GetURL(), "postgres://bar")
	SetEnvName("PGURL")
}
