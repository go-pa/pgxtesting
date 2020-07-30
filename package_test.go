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
	pool := CreateTestDatabase(t, testURL)

	// this should not generate any logs
	defer pool.Close()
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
	equal(defaultURL, GetDefaultURL())
	os.Setenv("PGURL", "")
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
}
