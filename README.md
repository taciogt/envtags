# envtags

[![Go Reference](https://pkg.go.dev/badge/github.com/taciogt/envtags.svg)](https://pkg.go.dev/github.com/taciogt/envtags)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/6904ddba8e6747559c7b4141b0f91e71)](https://www.codacy.com/gh/taciogt/envtags/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=taciogt/envtags&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/6904ddba8e6747559c7b4141b0f91e71)](https://www.codacy.com/gh/taciogt/envtags/dashboard?utm_source=github.com&utm_medium=referral&utm_content=taciogt/envtags&utm_campaign=Badge_Coverage)

_envtags_ is a package to use struct tags to automatically load environment variables on them. It is more about studying Go reflection approach than to create something better than what already exists, but an useful package with a [complete documentation](https://pkg.go.dev/github.com/taciogt/envtags) may come out of this. 

## Requirements

Go >= 1.18

## Usage

Define a struct with the `env` tag on **exported** fields to bind the fields with environment variables

```go
type Config struct {
	Foo int `env:"BAR"`
}
```

On an environment with the corresponding variables set, bind the struct to these variables using the method `envtags.Set()`

```shell
  export BAR="13" 
```

```go
var config Config
if err := envtags.Set(&config); err != nil {
	log.Fatal(err)
}
```

If the environment variable value can not be parsed to the field type, an `envtags.ErrInvalidTypeConversion` error is returned.  

If the field type is not supported, an `envtags.ErrParserNotAvailable` is returned.

## Supported field types:

- Primitive types
  - `bool`
  - `string`
  - `int`, `int64`, `int32`, `int16`, `int8`
  - `uint`, `uint64`, `uint32`, `uint16`, `uint8`
  - `float32`, `float64`
  - `complex64`, `complex128`
  - `byte`

## Refs

Better test output with:
```shell
go install github.com/rakyll/gotest
```

Inspired by 
*   https://github.com/kelseyhightower/envconfig/blob/master/envconfig.go
*   https://github.com/caarlos0/env

References
-   [Go numeric types](https://go.dev/ref/spec#Numeric_types)
-   [godoc](https://go.dev/blog/godoc)
-   [docs tips](https://tip.golang.org/doc/comment)
-   [docs examples](https://go.dev/blog/examples) 
