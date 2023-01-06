package envtags

import (
	"errors"
	"fmt"
)

// Specific errors returned by envtags package.
//
// Errors returned by the [Set] method can be tested against these variables using [errors.Is]
var (
	ErrInvalidTypeConversion = errors.New("invalid type conversion") // returned when the environment variable is not properly parsed to the expected field type
	ErrParserNotAvailable    = errors.New("parser not available")    // the field to be set has no parser for its reflect.Kind
)

func getError(customErr, baseErr error) error {
	return fmt.Errorf("%w: %s", customErr, baseErr)
}
