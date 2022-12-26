package envtags_test

import (
	"fmt"
	"github.com/taciogt/envtags"
	"log"
)

func ExampleSet_simpleStruct() {
	type Config struct {
		Foo int `env:"BAR"`
	}

	var cfg Config
	if err := envtags.Set(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("cfg:", cfg) // cfg.Foo = value from environment variable $BAR
}
