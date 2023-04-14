package envtags

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

func parseRune(envVarValue string, v reflect.Value) error {
	if len(envVarValue) != 1 {
		return getError(ErrInvalidTypeConversion, errors.New("environment variable length different than rune size"))
	}
	for _, letter := range envVarValue {
		v.SetInt(int64(letter))
	}
	return nil
}

type parserFunc func(string, reflect.Value) error

var parserByKindMap = map[reflect.Kind]parserFunc{
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

func getParser(k reflect.Kind, d tagDetails) (parserFunc, error) {
	if k == reflect.Int32 && d.IsRune {
		return parseRune, nil
	}

	if parser, parserExists := parserByKindMap[k]; parserExists {
		return parser, nil
	}
	return nil, getError(ErrParserNotAvailable, fmt.Errorf("parser for %s not found", k))

}
