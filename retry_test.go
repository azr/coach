package coach_test

import (
	"errors"
	"testing"
	"time"

	"github.com/azr/coach"
)

func TestRetry(t *testing.T) {

	i := 0

	op := func() error {
		if i < 3 {
			i++
			return errors.New("Foo")
		}
		return nil
	}

	cb := func(error) (time.Duration, error) {
		return time.Nanosecond, nil
	}

	if err := coach.Retry(op, cb); !(err == nil && i == 3) {
		t.Fatal("Test failed")
	}
}
