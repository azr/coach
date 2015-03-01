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

	cb := func(_ error) (time.Duration, error) {
		return time.Nanosecond, nil
	}

	if err := coach.Retry(op, cb); !(err == nil && i == 3) {
		t.Fatal("Test failed")
	}
}

var (
	ExampleErrorOne = errors.New("ExampleError one !")
)

type ExampleErrorTwo interface {
	Duration() time.Duration
}

type e struct{}

func (_ e) Three() time.Duration {
	return time.Nanosecond
}

func (_ e) Error() string {
	return "Two !"
}

func Example_retry(t *testing.T) {
	i := 0
	op := func() error {
		j := i
		i++
		switch j {
		case 0:
			return ExampleErrorOne
		case 1:
			return e{}
		default:
			return nil
		}
	}

	cb := func(err error) (time.Duration, error) {
		switch err {
		case ExampleErrorOne:
			return 0, nil
		default:
		}

		switch e := err.(type) {
		case ExampleErrorTwo:
			return e.Duration(), err
		default:
		}
		return 0, nil
	}

	err := coach.Retry(op, cb)
	if _, ok := err.(ExampleErrorTwo); !ok {
		t.Fatal("err should have been of ExampleErrorTwo type !?")
	}
}
