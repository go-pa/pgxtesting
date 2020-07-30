package pgxtesting

import (
	"net/url"
	"strings"
)

// PoolURL makes for quick manipulation of some paramets that are important for
// tests. The url is expected to always be a valid URL, utility functions will
// panic if it's not.
type PoolURL string

// NewPoolURL validates the input yurl
func NewPoolURL(poolURL string) (PoolURL, error) {
	_, err := url.Parse(poolURL)
	if err != nil {
		return "", err
	}
	return PoolURL(poolURL), nil
}

func (pu PoolURL) String() string {
	return string(pu)
}

// SetName sets the database name
func (pu PoolURL) SetName(dbName string) PoolURL {
	u, err := url.Parse(string(pu))
	if err != nil {
		panic(err)
	}
	u.Path = dbName
	return PoolURL(u.String())
}

// Name returns the database name
func (pu PoolURL) Name() string {
	u, err := url.Parse(string(pu))
	if err != nil {
		panic(err)
	}
	return strings.TrimPrefix(u.Path, "/")
}

// ConnURL returns a valid non pool connection url (can be used with go-migrate
// or sql.db)
func (pu PoolURL) ConnURL() PoolURL {
	u, err := url.Parse(string(pu))
	if err != nil {
		panic(err)
	}
	q := u.Query()
	for _, k := range []string{
		"pool_max_conns",
		"pool_min_conns",
		"pool_max_conn_lifetime",
		"pool_max_conn_idle_time",
		"pool_health_check_period",
	} {
		q.Del(k)
	}
	u.RawQuery = q.Encode()
	return PoolURL(u.String())
}
