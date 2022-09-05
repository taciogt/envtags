package envtags

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const tagName = "env"

var (
	// ErrInvalidTypeConversion is an error returned when the environment variable is not properly parsed to the expected field type
	ErrInvalidTypeConversion = errors.New("invalid type conversion")
)

func Set(s interface{}) error {
	value := reflect.ValueOf(s)
	elem := value.Elem()
	typeSpec := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		f := elem.Field(i)
		fType := typeSpec.Field(i)

		envVarName := fType.Tag.Get(tagName)
		envVarValue, ok := os.LookupEnv(envVarName)

		if ok {
			if fType.Type.Kind() == reflect.String {
				f.SetString(envVarValue)
			} else if fType.Type.Kind() == reflect.Int {
				invValue, err := strconv.Atoi(envVarValue)
				if err != nil {
					return GetError(ErrInvalidTypeConversion, err)
				}
				f.SetInt(int64(invValue))
			}
		}

	}
	return nil
}

func GetError(customErr, baseErr error) error {
	return fmt.Errorf("%w: %s", customErr, baseErr)
}
