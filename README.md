# envtags

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/6904ddba8e6747559c7b4141b0f91e71)](https://www.codacy.com/gh/taciogt/envtags/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=taciogt/envtags&amp;utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/6904ddba8e6747559c7b4141b0f91e71)](https://www.codacy.com/gh/taciogt/envtags/dashboard?utm_source=github.com&utm_medium=referral&utm_content=taciogt/envtags&utm_campaign=Badge_Coverage)

🚧 (wip) 

envtags is a package to support env tags to load environment variables on structs. It is more about studying Go reflection approach than to create something that already exists, at least for now.

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

## Refs

Better test output with:
```shell
go install github.com/rakyll/gotest
```

Inspired by 
*   https://github.com/kelseyhightower/envconfig/blob/master/envconfig.go
*   https://github.com/caarlos0/env

References
*   https://go.dev/ref/spec#Numeric_types
