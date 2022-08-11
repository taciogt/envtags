package envtags

import (
	"fmt"
	"os"
	"reflect"
)

const tagName = "env"

func Set(s interface{}) bool {
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

		f.SetString(envVarValue)
	}
	return true
}
