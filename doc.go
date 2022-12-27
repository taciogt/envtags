/*
Package envtags provides easy to use methods for setting struct fields from environment variables.

	type Config struct {
		Foo int `env:"BAR"`
	}

	var cfg Config
	Set(&cfg)

Currently, it supports the most common types in the root of the struct and parses the environment variables using the package strconv. For details on how the environment variable can be set, check [strconv.ParseBool], [strconv.ParseInt], [strconv.ParseUint]

* Unsigned integer
*/
package envtags
