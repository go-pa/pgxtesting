package pgxtesting

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

// TestPool .
type TestPool struct {
	*pgxpool.Pool
	Name        string // database name
	URL         PoolURL
	originalURL PoolURL
}

// Close closes the test pool and deletes the testing database.
func (t *TestPool) Close() {
	if t.Pool != nil {
		t.Pool.Close()
	}

	pool, err := connectPostgres(t.originalURL)
	if err != nil {
		fmt.Println("pgxtesting.Pool.Cleanup error:", err)
		return
	}
	defer pool.Close()

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Println("pgxtesting.Pool.Cleanup error:", err)
		return
	}
	defer conn.Release()

	sql := fmt.Sprintf("drop database %s", t.URL.Name())
	_, err = conn.Exec(context.Background(), sql)
	if err != nil {
		fmt.Printf("pgxtesting.Pool.Cleanup error running '%s': %v\n", sql, err)
		return
	}

}

func CreateTestDatabaseEnv(tb testing.TB) *TestPool {
	tb.Helper()
	return CreateTestDatabase(tb, GetURL())
}

func CreateTestDatabase(tb testing.TB, pgxPoolURL string) *TestPool {
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
