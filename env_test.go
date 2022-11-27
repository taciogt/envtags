package envtags

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"testing"
)

func TestSetFieldTypes(t *testing.T) {
	type config struct {
		Word      string    `env:"FOO"`
		Float32   float32   `env:"FLOAT_NUMBER"`
		Complex64 complex64 `env:"COMPLEX_64"`

		Int   int   `env:"INT"`
		Int8  int8  `env:"INT_8"`
		Int16 int16 `env:"INT_16"`
		Int32 int32 `env:"INT_32"`
		Int64 int64 `env:"INT_64"`

		UInt   uint   `env:"UINT"`
		UInt8  uint8  `env:"UINT_8"`
		UInt16 uint16 `env:"UINT_16"`
		UInt32 uint32 `env:"UINT_32"`
		UInt64 uint64 `env:"UINT_64"`
	}

	tests := []struct {
		name     string
		expected config
		envVars  map[string]string
		wantErr  error
	}{{
		name:     "set string field",
		expected: config{Word: "bar"},
		envVars: map[string]string{
			"FOO": "bar",
		},
	}, {
		name:     "set float field",
		expected: config{Float32: 1.23},
		envVars: map[string]string{
			"FLOAT_NUMBER": "1.23",
		},
	}, {
		name: "set float field",
		envVars: map[string]string{
			"COMPLEX_64": "-",
		},
		wantErr: ErrParserNotAvailable,
	},
		// int type fields
		{
			name:     "set integer field",
			expected: config{Int: 123},
			envVars: map[string]string{
				"INT": "123",
			},
		}, {
			name:     "set integer field with big value",
			expected: config{Int: 21474836},
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
			expected: config{Int8: 19},
			envVars: map[string]string{
				"INT_8": "19",
			},
		}, {
			name:     "set int8 field with negative value",
			expected: config{Int8: -13},
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
			name: "set int8 field with value less than min size",
			envVars: map[string]string{
				"INT_8": "-129", // max size is 127
			},
			wantErr: ErrInvalidTypeConversion,
		}, {
			name:     "set int16 field",
			expected: config{Int16: 32767},
			envVars: map[string]string{
				"INT_16": "32767",
			},
		}, {
			name:     "set int16 field with negative value",
			expected: config{Int16: -32768},
			envVars: map[string]string{
				"INT_16": "-32768",
			},
		}, {
			name:     "set int32 field",
			expected: config{Int32: 2147483647},
			envVars: map[string]string{
				"INT_32": "2147483647",
			},
		}, {
			name:     "set int64 field",
			expected: config{Int64: 9223372036854775807},
			envVars: map[string]string{
				"INT_64": "9223372036854775807",
			},
		},
		// unsigned integer types
		{
			name:     "set unsigned integer field",
			expected: config{UInt: 123},
			envVars: map[string]string{
				"UINT": "123",
			},
		}, {
			name: "set unsigned integer field for string bigger than max size",
			envVars: map[string]string{
				"UINT": "184467440737095516150",
			},
			wantErr: ErrInvalidTypeConversion,
		}, {
			name: "try to set unsigned integer field with negative envvar",
			envVars: map[string]string{
				"UINT": "-1",
			},
			wantErr: ErrInvalidTypeConversion,
		}, {
			name:     "set unsigned uint64 field",
			expected: config{UInt64: 123},
			envVars: map[string]string{
				"UINT_64": "123",
			},
		}, {
			name:     "set unsigned uint32 field",
			expected: config{UInt32: 123},
			envVars: map[string]string{
				"UINT_32": "123",
			},
		}, {
			name: "try to set unsigned uint32 field with value bigger than max size",
			envVars: map[string]string{
				"UINT_32": "4294967296",
			},
			wantErr: ErrInvalidTypeConversion,
		}, {
			name:     "set unsigned uint16 field",
			expected: config{UInt16: 123},
			envVars: map[string]string{
				"UINT_16": "123",
			},
		}, {
			name:     "set unsigned uint8 field",
			expected: config{UInt8: 123},
			envVars: map[string]string{
				"UINT_8": "123",
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

			var cfg config

			if err := Set(&cfg); err != tt.wantErr && !errors.Is(err, tt.wantErr) {
				t.Errorf("err different than expected, want %+v, got %+v", tt.wantErr, err)
				return
			}
			if cfg != tt.expected {
				t.Errorf("Set(&s), \nwant %+v,\ngot  %+v", tt.expected, cfg)
			}
		})
	}
}

func FuzzSetUint(f *testing.F) {
	type config struct {
		UInt8 uint8 `env:"UINT_8"`
	}
	f.Fuzz(func(t *testing.T, s string) {
		ignoredEntryRegex, err := regexp.Compile("(^\\s+$)|(^0.*$)")
		if err != nil {
			t.Error(err)
		}
		if ignoredEntryRegex.Match([]byte(s)) {
			t.Skip()
		}

		envVarName := "UINT_8"
		if err := os.Setenv(envVarName, s); err != nil {
			t.Skip()
		}
		var cfg config
		if err := Set(&cfg); err != nil && !errors.Is(err, ErrInvalidTypeConversion) {
			t.Error(err)
		} else if errors.Is(err, ErrInvalidTypeConversion) {
			t.Skip()
		}

		_ = Set(&cfg)
		if os.Getenv(envVarName) != strconv.Itoa(int(cfg.UInt8)) {
			t.Errorf("cfg field no set as expected. got=\"%d\", want=\"%s\"", cfg.UInt8, s)
		}

	})
}
