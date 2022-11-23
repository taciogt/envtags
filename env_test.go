package envtags

import (
	"errors"
	"os"
	"testing"
)

func TestSetFieldTypes(t *testing.T) {
	type Config struct {
		Word string `env:"FOO"`

		Int   int   `env:"INT"`
		Int8  int8  `env:"INT_8"`
		Int16 int16 `env:"INT_16"`

		Float32 float32 `env:"FLOAT_NUMBER"`
	}

	tests := []struct {
		name     string
		expected Config
		envVars  map[string]string
		wantErr  error
	}{{
		name:     "set string field",
		expected: Config{Word: "bar"},
		envVars: map[string]string{
			"FOO": "bar",
		},
	}, {
		name:     "set integer field",
		expected: Config{Int: 123},
		envVars: map[string]string{
			"INT": "123",
		},
	}, {
		name:     "set integer field with big value",
		expected: Config{Int: 21474836},
		envVars: map[string]string{
			"INT": "21474836", // value bigger than a int16
		},
	}, {
		name:    "set integer field with invalid env var",
		wantErr: ErrInvalidTypeConversion,
		envVars: map[string]string{
			"INT": "abc",
		},
	}, {
		name:     "set int8 field",
		expected: Config{Int8: 19},
		envVars: map[string]string{
			"INT_8": "19",
		},
	}, {
		name:     "set int8 field with negative value",
		expected: Config{Int8: -13},
		envVars: map[string]string{
			"INT_8": "-13",
		},
	}, {
		name: "set int8 field with value greater than max size",
		envVars: map[string]string{
			"INT_8": "130", // max size is 127
		},
		wantErr: ErrInvalidTypeConversion,
	}, {
		name:     "set int16 field",
		expected: Config{Int16: 32767},
		envVars: map[string]string{
			"INT_16": "32767",
		},
	}, {
		name:     "set int16 field with negative value",
		expected: Config{Int16: -32768},
		envVars: map[string]string{
			"INT_16": "-32768",
		},
	}, {
		name:     "set float field",
		expected: Config{Float32: 1.23},
		envVars: map[string]string{
			"FLOAT_NUMBER": "1.23",
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				if err := os.Setenv(k, v); err != nil {
					t.Error(err)
					return
				}
			}
			defer os.Clearenv()

			var cfg Config

			if err := Set(&cfg); err != tt.wantErr && !errors.Is(err, tt.wantErr) {
				t.Errorf("err different than expected, want %+v, got %+v", tt.wantErr, err)
				return
			}
			if cfg != tt.expected {
				t.Errorf("Set(&s), want %+v, got %+v", tt.expected, cfg)
			}
		})
	}
}
