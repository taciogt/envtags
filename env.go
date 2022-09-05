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

var parserByKindMap = map[reflect.Kind]func(envVarValue string, v reflect.Value) error{
	reflect.String: func(envVarValue string, v reflect.Value) error {
		v.SetString(envVarValue)
		return nil
	},
	reflect.Int: func(envVarValue string, v reflect.Value) error {
		invValue, err := strconv.Atoi(envVarValue)
		if err != nil {
			return GetError(ErrInvalidTypeConversion, err)
		}
		v.SetInt(int64(invValue))
		return nil
	},
	reflect.Float32: func(envVarValue string, v reflect.Value) error {
		floatValue, err := strconv.ParseFloat(envVarValue, 32)
		if err != nil {
			return GetError(ErrInvalidTypeConversion, err)
		}
		v.SetFloat(floatValue)
		return nil
	},
}

func Set(s interface{}) error {
	value := reflect.ValueOf(s)
	elem := value.Elem()
	typeSpec := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		f := elem.Field(i)
		fType := typeSpec.Field(i)
		envVarName := fType.Tag.Get(tagName)

		if envVarValue, ok := os.LookupEnv(envVarName); ok {
			parser, parserExists := parserByKindMap[fType.Type.Kind()]
			if !parserExists {
				return errors.New("parser not available")
			}
			if err := parser(envVarValue, f); err != nil {
				return err
			}
		}

	}
	return nil
}

func GetError(customErr, baseErr error) error {
	return fmt.Errorf("%w: %s", customErr, baseErr)
}
