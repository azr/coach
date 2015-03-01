// Coach is a simple package that lets you retry operations
package coach

import (
	"time"
)

// Retry will execute operation, if operation returns an error,
// time.Sleep for the duration returned by callback or exits if callback returns an error.
//
// This is pretty convenient to allow you to check the operation's error before retrying, if any.
//
// A negative or zero duration causes Sleep to return immediately.
func Retry(operation func() error, callback func(error) (time.Duration, error)) error {
	for err := operation(); err != nil; err = operation() {
		if d, err := callback(err); err != nil {
			return err
		} else {
			time.Sleep(d)
		}
	}
	return nil
}
