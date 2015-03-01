// Coach is a simple package that lets you retry operations
package coach

import (
	"time"
)

// Retry if operation fails, wait for the duration returned by callback or exits if the error is too big.
// Just return nil if operation is successfull.
//
// This is pretty convenient to allow you to check the operation's error before retrying, if any.
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
