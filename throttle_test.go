package throttle

import (
	"errors"
	"testing"
	"time"
)

func TestInitialRequestServedImmediately(t *testing.T) {
	entity, value := "somebody", "something"
	th := New[string, string](time.Hour, time.Hour)

	result, err := th.Await(entity, &value, 100*time.Millisecond)
	if err != nil {
		t.Error("timeout")
	}
	if result != &value {
		t.Error("wrong result")
	}
}

func TestRequestDelayedForEntity(t *testing.T) {
	entity, valueA, valueB := "somebody", "something", "other thing"
	th := New[string, string](200*time.Millisecond, 500*time.Millisecond)
	minimumWaitTime := 500 * time.Millisecond

	resultA, err := th.Await(entity, &valueA, 100*time.Millisecond)
	if err != nil {
		t.Error("timeout A")
	}
	if resultA != &valueA {
		t.Error("wrong result A")
	}

	resultB, err := th.Await(entity, &valueB, 2*minimumWaitTime)
	if err != nil {
		t.Error("timeout B")
	}
	if resultB != &valueB {
		t.Error("wrong result B")
	}

	_, err = th.Await(entity, &valueB, minimumWaitTime-100*time.Millisecond)
	if !errors.Is(err, TimeoutError) {
		t.Errorf("expected timeout, but succeeded")
	}
}

func TestSecondEntityDelayedGlobally(t *testing.T) {
	entityA, entityB := "Alice", "Bob"
	valueA, valueB := "Apple", "Pear"
	th := New[string, string](time.Second, 3*time.Second)

	resultA, err := th.Await(entityA, &valueA, 100*time.Millisecond)
	if err != nil {
		t.Error("timeout A")
	}
	if resultA != &valueA {
		t.Error("wrong result A")
	}

	resultB, err := th.Await(entityB, &valueB, 2*time.Second)
	if err != nil {
		t.Error("timeout B")
	}
	if resultB != &valueB {
		t.Error("wrong result B")
	}
}
