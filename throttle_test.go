package throttle

import (
	"testing"
	"time"
)

func TestInitial(t *testing.T) {
	entity, value := "somebody", "something"
	th := New[string, string](time.Hour, time.Hour)
	result, err := th.Await(entity, &value, 100*time.Millisecond)
	if result != &value {
		t.Error("wrong result")
	}
	if err != nil {
		t.Error("timeout")
	}
}
