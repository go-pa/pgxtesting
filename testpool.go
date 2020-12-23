package pgxtesting

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type TB interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	TempDir() string
}

// TestPool .
type TestPool struct {
	*pgxpool.Pool
	Name        string // database name
	URL         PoolURL
	originalURL PoolURL
	closed      bool
}

// Close closes the test pool and deletes the testing database.
func (t *TestPool) Close() {
	if t.Pool != nil && !t.closed {
		t.closed = true
		t.Pool.Close()
	}

	err := t.dropTestDB()
	if err != nil {
		fmt.Println(err)
	}

}

func (t *TestPool) dropTestDB() error {
	ctx := context.Background()
	pool, err := connectPostgres(t.originalURL)
	if err != nil {
		return fmt.Errorf("pgxtesting.Pool.Cleanup error: %v", err)
	}
	defer pool.Close()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("pgxtesting.Pool.Cleanup error: %v", err)

	}
	defer conn.Release()

	sql := fmt.Sprintf("drop database %s", t.URL.Name())
	_, err = conn.Exec(ctx, sql)
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code != "3D000" {
				return fmt.Errorf("pgxtesting.Pool.Cleanup error running '%s': %v\n", sql, err)
			}
		} else {
			return fmt.Errorf("pgxtesting.Pool.Cleanup error running '%s': %v\n", sql, err)
		}
	}

	return nil

}

func CreateTestDatabaseEnv(tb TB) *TestPool {
	tb.Helper()
	return CreateTestDatabase(tb, GetURL())
}

func CreateTestDatabase(tb TB, pgxPoolURL string) *TestPool {
	tb.Helper()
	pu, err := NewPoolURL(pgxPoolURL)
	if err != nil {
		tb.Fatal(err)
	}

	dbName := getRandomDBName()

	if err := createDB(pu, dbName); err != nil {
		tb.Fatalf("error creating database %v: %v", dbName, err)
	}

	URL := pu.SetName(dbName)

	pool, err := connectPostgres(URL)
	if err != nil {
		tb.Fatal(err)
	}

	tp := &TestPool{
		Pool:        pool,
		Name:        dbName,
		URL:         URL,
		originalURL: pu,
	}
	tb.Cleanup(tp.Close)
	return tp
}

func createDB(pu PoolURL, dbName string) error {
	ctx := context.Background()
	pool, err := connectPostgres(pu)
	if err != nil {
		return fmt.Errorf("cannot connect to dburl %s: %v", pu, err)
	}
	defer pool.Close()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, fmt.Sprintf("create database %s", dbName))
	if err != nil {
		return fmt.Errorf("error %v", err)
	}

	return nil
}

func connectPostgres(pu PoolURL) (*pgxpool.Pool, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(pu.String())
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	return pool, err
}
