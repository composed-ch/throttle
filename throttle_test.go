package throttle

import (
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

func TestSecondEntityHasToAwaitGlobalToken(t *testing.T) {
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
