package envtags

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const tagName = "env"

// Specific errors returned by envtags package.
//
// Errors returned by the [Set] method can be tested against these variables using [errors.Is]
var (
	ErrInvalidTypeConversion = errors.New("invalid type conversion") // returned when the environment variable is not properly parsed to the expected field type
	ErrParserNotAvailable    = errors.New("parser not available")    // the field to be set has no parser for its reflect.Kind
)

func getIntParser(bitSize int) func(envVarValue string, v reflect.Value) error {
	return func(envVarValue string, v reflect.Value) error {
		intValue, err := strconv.ParseInt(envVarValue, 0, bitSize)
		if err != nil {
			return getError(ErrInvalidTypeConversion, err)
		}
		v.SetInt(intValue)
		return nil
	}
}

func getUIntParser(bitSize int) func(envVarValue string, v reflect.Value) error {
	return func(envVarValue string, v reflect.Value) error {
		uintValue, err := strconv.ParseUint(envVarValue, 0, bitSize)
		if err != nil {
			return getError(ErrInvalidTypeConversion, err)
		}
		v.SetUint(uintValue)
		return nil
	}
}

func getFloatParser(bitSize int) func(envVarValue string, v reflect.Value) error {
	return func(envVarValue string, v reflect.Value) error {
		floatValue, err := strconv.ParseFloat(envVarValue, bitSize)
		if err != nil {
			return getError(ErrInvalidTypeConversion, err)
		}
		v.SetFloat(floatValue)
		return nil
	}
}

func getComplexParser(bitSize int) func(envVarValue string, v reflect.Value) error {
	return func(envVarValue string, v reflect.Value) error {
		complexValue, err := strconv.ParseComplex(envVarValue, bitSize)
		if err != nil {
			return getError(ErrInvalidTypeConversion, err)
		}
		v.SetComplex(complexValue)
		return nil
	}
}

var parserByKindMap = map[reflect.Kind]func(envVarValue string, v reflect.Value) error{
	reflect.Bool: func(envVarValue string, v reflect.Value) error {
		if envVarValue == "" {
			v.SetBool(false)
			return nil
		}
		boolValue, err := strconv.ParseBool(envVarValue)
		if err != nil {
			return getError(ErrInvalidTypeConversion, err)
		}
		v.SetBool(boolValue)
		return nil
	},
	reflect.String: func(envVarValue string, v reflect.Value) error {
		v.SetString(envVarValue)
		return nil
	},
	reflect.Int:        getIntParser(64),
	reflect.Int8:       getIntParser(8),
	reflect.Int16:      getIntParser(16),
	reflect.Int32:      getIntParser(32),
	reflect.Int64:      getIntParser(64),
	reflect.Uint:       getUIntParser(64),
	reflect.Uint8:      getUIntParser(8),
	reflect.Uint16:     getUIntParser(16),
	reflect.Uint32:     getUIntParser(32),
	reflect.Uint64:     getUIntParser(32),
	reflect.Float32:    getFloatParser(32),
	reflect.Float64:    getFloatParser(64),
	reflect.Complex64:  getComplexParser(64),
	reflect.Complex128: getComplexParser(128),
}

/*
Set receives a struct pointer and sets its fields using the value from environment variables defined in the struct tag "env".
*/
func Set(s interface{}) error {
	return set(s, tagDetails{})
}

func set(s interface{}, options tagDetails) error {
	value := reflect.ValueOf(s)
	elem := value.Elem()
	typeSpec := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		f := elem.Field(i)
		fType := typeSpec.Field(i)
		tagValue := fType.Tag.Get(tagName)

		details := parseTagValue(tagValue).Update(options)

		k := fType.Type.Kind()
		if k == reflect.Struct {
			if err := set(f.Addr().Interface(), details); err != nil {
				return err
			}
			continue
		}

		if envVarValue, ok := os.LookupEnv(details.GetEnvVar()); ok {
			parser, parserExists := parserByKindMap[fType.Type.Kind()]
			if !parserExists {
				return getError(ErrParserNotAvailable, fmt.Errorf("parser for %s not found", fType.Type.Kind()))
			}
			if err := parser(envVarValue, f); err != nil {
				return fmt.Errorf("failed to parse value for field %s: %w", elem.Type().Field(i).Name, err)
			}
		}

	}
	return nil
}

func getError(customErr, baseErr error) error {
	return fmt.Errorf("%w: %s", customErr, baseErr)
}
