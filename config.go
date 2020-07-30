package pgxtesting

import (
	"net/url"
	"os"
	"sync"
)

var (
	configMu   sync.Mutex
	defaultURL = "postgres://test:test@localhost:5432/test?sslmode=disable&pool_max_conns=500"
	envName    = "PGURL"
)

// sets the default url for postgres tests, panics if the url is invalid
func SetURL(pgxurl string) {
	_, err := url.Parse(pgxurl)
	if err != nil {
		panic(err)
	}
	configMu.Lock()
	defaultURL = pgxurl
	configMu.Unlock()
}

func GetURL() string {
	configMu.Lock()
	defer configMu.Unlock()
	return defaultURL
}

func GetEnvName() string {
	configMu.Lock()
	defer configMu.Unlock()
	return envName

}

func SetEnvVar(name string) {
	configMu.Lock()
	envName = name
	configMu.Unlock()
}

func GetPgxPoolURL() string {
	v := os.Getenv(GetEnvName())
	if v != "" {
		return v
	}
	return GetURL()
}
