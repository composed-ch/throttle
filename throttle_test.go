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

	resultA, err := th.Await(entity, &valueA, 100*time.Millisecond)
	if err != nil {
		t.Error("timeout A")
	}
	if resultA != &valueA {
		t.Error("wrong result A")
	}

	resultB, err := th.Await(entity, &valueB, time.Second)
	if err != nil {
		t.Error("timeout B")
	}
	if resultB != &valueB {
		t.Error("wrong result B")
	}

	// timeout < entity spawn interval
	_, err = th.Await(entity, &valueB, 400*time.Millisecond)
	if !errors.Is(err, TimeoutError) {
		t.Errorf("expected timeout, but succeeded")
	}
}

func TestRequestDelayedGlobally(t *testing.T) {
	entity, valueA, valueB := "somebody", "something", "other thing"
	th := New[string, string](500*time.Millisecond, 200*time.Millisecond)

	resultA, err := th.Await(entity, &valueA, 100*time.Millisecond)
	if err != nil {
		t.Error("timeout A")
	}
	if resultA != &valueA {
		t.Error("wrong result A")
	}

	resultB, err := th.Await(entity, &valueB, time.Second)
	if err != nil {
		t.Error("timeout B")
	}
	if resultB != &valueB {
		t.Error("wrong result B")
	}

	// timeout < global spawn interval
	_, err = th.Await(entity, &valueB, 400*time.Millisecond)
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

func TestSecondEntityDelayedForEntity(t *testing.T) {
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

	_, err = th.Await(entityB, &valueB, time.Second-100*time.Millisecond)
	if !errors.Is(err, TimeoutError) {
		t.Errorf("expected timeout, but succeeded")
	}
}

// TODO: spammer must not influence other user beyond global rate
func TestSpammerMustNotAffectUserBeyondGlobalRate(t *testing.T) {
	entityA, entityB := "Spammer", "Egger"
	valueA, valueB := "spam", "eggs"
	th := New[string, string](10*time.Millisecond, 300*time.Millisecond)
	done := make(chan struct{})

	// spam, spam, spam
	go func() {
		th.Await(entityA, &valueA, 100*time.Millisecond)
		th.Await(entityA, &valueA, 100*time.Millisecond)
		th.Await(entityA, &valueA, 100*time.Millisecond)
		done <- struct{}{}
	}()

	// give me some ham
	go func() {
		result, err := th.Await(entityB, &valueB, 50*time.Millisecond)
		if err != nil {
			t.Errorf("regular user ran into timeout")
		}
		if result != &valueB {
			t.Errorf("wrong result B")
		}
		done <- struct{}{}
	}()

	<-done
	<-done
}
