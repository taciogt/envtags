package envtags

import (
	"fmt"
	"os"
	"reflect"
)

const tagName = "env"

/*
Set receives a struct pointer and sets its fields using the value from
environment variables defined in the struct tag `env`.
*/
func Set(s interface{}) error {
	return set(s, tagDetails{})
}

func set(s interface{}, details tagDetails) error {
	value := reflect.ValueOf(s)
	elem := value.Elem()
	typeSpec := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		f := elem.Field(i)
		fType := typeSpec.Field(i)
		tagValue := fType.Tag.Get(tagName)

		details := parseTagValue(tagValue).Update(details)

		k := fType.Type.Kind()
		if k == reflect.Struct {
			if err := set(f.Addr().Interface(), details); err != nil {
				return err
			}
			continue
		}

		if envVarValue, ok := os.LookupEnv(details.GetEnvVar()); ok {
			parser, err := getParser(k, details)
			if err != nil {
				return err
			}
			if err := parser(envVarValue, f); err != nil {
				return fmt.Errorf("failed to parse value for field %s: %w", elem.Type().Field(i).Name, err)
			}
		}

	}
	return nil
}
