package throttle

import (
	"errors"
	"sync"
	"time"
)

var TimeoutError error

func init() {
	TimeoutError = errors.New("timeout")
}

// Throttle allows to throttles the usage of values of type V identified by entities of type T.
type Throttle[T comparable, V any] struct {
	// GlobalSpawn is a channel that spawns tokens at GlobalRate.
	GlobalSpawn *time.Ticker

	// EntitySpawnMux is a mutex over EntitySpawn.
	EntitySpawnMux sync.RWMutex

	// EntitySpawn maps entities to their individual channel that spawns tokens at EntityRate.
	EntitySpawn map[T]*time.Ticker

	// EntityRate is the rate at which entity-specific tokens are spanwed.
	EntityRate time.Duration
}

// New creates a new Throttle by defining its global and entity token spawn rate.
func New[T comparable, V any](globalRate, entityRate time.Duration) *Throttle[T, V] {
	return &Throttle[T, V]{
		GlobalSpawn: time.NewTicker(globalRate),
		EntitySpawn: make(map[T]*time.Ticker),
		EntityRate:  entityRate,
	}
}

// Await waits for the throttle object to spawn both a token for the individual
// entity and globally. If those tokens are spawned within the given timeout,
// the given value is returned back. Otherwise, TimeoutError is returned.
// The first token (globally, per entity) is spawned immediately.
func (t *Throttle[T, V]) Await(entity T, value *V, timeout time.Duration) (*V, error) {
	timeoutChan := time.NewTimer(timeout)

	awaitGlobalToken := func() (*V, error) {
		t.EntitySpawnMux.RLock()
		empty := len(t.EntitySpawn) == 0
		t.EntitySpawnMux.RUnlock()
		if empty {
			// first token spawns immediately
			return value, nil
		}
		select {
		case <-timeoutChan.C:
			return nil, TimeoutError
		case <-t.GlobalSpawn.C:
			return value, nil
		}
	}

	t.EntitySpawnMux.RLock()
	var entityTicker *time.Ticker
	entityTicker, ok := t.EntitySpawn[entity]
	t.EntitySpawnMux.RUnlock()
	if !ok {
		// first token spawns immediately
		result, err := awaitGlobalToken()
		if err != nil {
			return nil, err
		}
		entityTicker = time.NewTicker(t.EntityRate)
		t.EntitySpawnMux.Lock()
		t.EntitySpawn[entity] = entityTicker
		t.EntitySpawnMux.Unlock()
		return result, nil
	}

	select {
	case <-timeoutChan.C:
		return nil, TimeoutError
	case <-entityTicker.C:
		return awaitGlobalToken()
	}
}
