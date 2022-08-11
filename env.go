package envtags

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const tagName = "env"

func Set(s interface{}) error {
	value := reflect.ValueOf(s)
	//if value.Type().Kind() == reflect.Ptr {
	//
	//	elem := value.Elem()
	//}
	elem := value.Elem()

	//typeSpec := value.Type()
	typeSpec := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		f := elem.Field(i)
		fmt.Println("struct field", f)
		fType := typeSpec.Field(i)
		fmt.Printf("struct type tag: %+v\n", fType.Tag)
		fmt.Printf("struct type tag.env: %+v\n", fType.Tag.Get(tagName))

		envVarName := fType.Tag.Get(tagName)
		envVarValue, _ := os.LookupEnv(envVarName)

		if fType.Type.Kind() == reflect.String {
			f.SetString(envVarValue)
		} else if fType.Type.Kind() == reflect.Int {
			invValue, err := strconv.Atoi(envVarValue)
			if err != nil {
				return err
			}
			f.SetInt(int64(invValue))
		}
	}
	return nil
}
