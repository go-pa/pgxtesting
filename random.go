package pgxtesting

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	random   = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	randomMu sync.Mutex
)

func getRandomDBName() string {
	randomMu.Lock()
	defer randomMu.Unlock()
	return fmt.Sprintf("pgxtesting_%v", random.Int63())
}
