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

func getIntParsers(maxSize int) func(envVarValue string, v reflect.Value) error {
	return func(envVarValue string, v reflect.Value) error {
		intValue, err := strconv.Atoi(envVarValue)
		if err != nil {
			return GetError(ErrInvalidTypeConversion, err)
		}
		if intValue > maxSize {
			return GetError(ErrInvalidTypeConversion, errors.New("value greater than max available"))
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
	reflect.Int:   getIntParsers(math.MaxInt),
	reflect.Int8:  getIntParsers(math.MaxInt8),
	reflect.Int16: getIntParsers(math.MaxInt16),
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
				return GetParserNotAvailableError(fType)
			}
			if err := parser(envVarValue, f); err != nil {
				return err
			}
		}

	}
	return nil
}

func GetParserNotAvailableError(fType reflect.StructField) error {
	return fmt.Errorf("%w. kind=%s", ErrParserNotAvailable, fType.Type.Kind())
}

func GetError(customErr, baseErr error) error {
	return fmt.Errorf("%w: %s", customErr, baseErr)
}
