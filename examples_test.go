package envtags_test

import (
	"fmt"
	"github.com/taciogt/envtags"
	"log"
	"os"
)

func ExampleSet_simpleStruct() {
	const envKey = "BAR"
	defer func(value string) { _ = os.Setenv(envKey, value) }(os.Getenv(envKey))

	if err := os.Setenv(envKey, "123"); err != nil {
		log.Fatal(err)
	}

	type Config struct {
		Foo int `env:"BAR"`
	}

	var cfg Config
	if err := envtags.Set(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("cfg: %+v", cfg)
	// Output: cfg: {Foo:123}
}
