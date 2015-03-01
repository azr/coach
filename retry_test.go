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

// Example_retry

var (
	ErrorOne = errors.New("Error one !")
)

type ErrorTwoInterface interface {
	Duration() time.Duration
}
type ErrorTooBigInterface interface {
	TooBig()
}

type MyErrorTwo struct{}

func (_ MyErrorTwo) Three() time.Duration {
	return time.Nanosecond
}

func (_ MyErrorTwo) Error() string {
	return "Two !"
}

type TooBigError struct{}

func (_ TooBigError) TooBig() {}

func (_ TooBigError) Error() string {
	return "Too BIG !"
}

func ExampleRetry(t *testing.T) {
	i := 0
	op := func() error {
		j := i
		i++
		switch j {
		case 0:
			return ErrorOne // an errors.New("") error
		case 1:
			return MyErrorTwo{} // a special type of error
		case 2:
			return TooBigError{} // another special type of error
		default:
			return nil
		}
	}

	cb := func(err error) (time.Duration, error) {
		if err == ErrorOne {
			return 0, nil // hah almost false alert !
		}

		switch e := err.(type) {
		case ErrorTwoInterface: //Just check if err has a .Duration() method !
			return e.Duration(), nil
		case ErrorTooBigInterface: //Just check if err has a .TooBig() method !
			return 0, err // damn, we can't make it back, let's return the error !
		default:
		}
		return 0, nil
	}

	err := coach.Retry(op, cb)
	if _, ok := err.(ErrorTooBigInterface); !ok {
		panic("err should have been of ErrorTooBigInterface type !?")
	}
}
