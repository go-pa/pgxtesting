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

// SetDefaultURL sets the default URL for postgres tests, panics if the url is invalid
func SetDefaultURL(pgxurl string) {
	_, err := url.Parse(pgxurl)
	if err != nil {
		panic(err)
	}
	configMu.Lock()
	defaultURL = pgxurl
	configMu.Unlock()
}

// SetDefaultURL returns the default URL used to connect to the postgres server
// while creating and destorying test databases.
func GetDefaultURL() string {
	configMu.Lock()
	defer configMu.Unlock()
	return defaultURL
}


// SetEnvName sets the enviroment variable name used to fetch the pgxpool URL,
// defaults name is PGURL.
func SetEnvName(name string) {
	configMu.Lock()
	envName = name
	configMu.Unlock()
}

// GetEnvName gets the enviroment variable name used to fetch the pgxpool URL,
// defaults name is PGURL.
func GetEnvName() string {
	configMu.Lock()
	defer configMu.Unlock()
	return envName

}

// GetURL returns the pgx pool URL from environment or the default value.
func GetURL() string {
	v := os.Getenv(GetEnvName())
	if v != "" {
		return v
	}
	return GetDefaultURL()
}
