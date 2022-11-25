package envtags

import (
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
)

const tagName = "env"

var (
	// ErrInvalidTypeConversion is an error returned when the environment variable is not properly parsed to the expected field type
	ErrInvalidTypeConversion = errors.New("invalid type conversion")
	// ErrParserNotAvailable is returned when the value to be set has no parser for its reflect.Kind
	ErrParserNotAvailable = errors.New("parser not available")
)

func getIntParser(minSize int, maxSize int) func(envVarValue string, v reflect.Value) error {
	return func(envVarValue string, v reflect.Value) error {
		intValue, err := strconv.Atoi(envVarValue)
		if err != nil {
			return getError(ErrInvalidTypeConversion, err)
		}
		if intValue > maxSize {
			return getError(ErrInvalidTypeConversion, errors.New("value greater than max available"))
		} else if intValue < minSize {
			return getError(ErrInvalidTypeConversion, errors.New("value less than min available"))
		}
		v.SetInt(int64(intValue))
		return nil
	}
}

var parserByKindMap = map[reflect.Kind]func(envVarValue string, v reflect.Value) error{
	reflect.String: func(envVarValue string, v reflect.Value) error {
		v.SetString(envVarValue)
		return nil
	},
	reflect.Int:   getIntParser(math.MinInt, math.MaxInt),
	reflect.Int8:  getIntParser(math.MinInt8, math.MaxInt8),
	reflect.Int16: getIntParser(math.MinInt16, math.MaxInt16),
	reflect.Int32: getIntParser(math.MinInt32, math.MaxInt32),
	reflect.Int64: getIntParser(math.MinInt64, math.MaxInt64),
	reflect.Uint: func(envVarValue string, v reflect.Value) error {
		uintValue, err := strconv.ParseUint(envVarValue, 10, 64)
		if err != nil {
			return getError(ErrInvalidTypeConversion, err)
		}
		v.SetUint(uintValue)
		return nil
	},
	reflect.Float32: func(envVarValue string, v reflect.Value) error {
		floatValue, err := strconv.ParseFloat(envVarValue, 32)
		if err != nil {
			return getError(ErrInvalidTypeConversion, err)
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
				return getError(ErrParserNotAvailable, fmt.Errorf("parser for %s not found", fType.Type.Kind()))
			}
			if err := parser(envVarValue, f); err != nil {
				return err
			}
		}

	}
	return nil
}

func getError(customErr, baseErr error) error {
	return fmt.Errorf("%w: %s", customErr, baseErr)
}
