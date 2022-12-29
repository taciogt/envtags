package envtags

import (
	"errors"
	"math"
	"os"
	"regexp"
	"strconv"
	"testing"
)

func TestSetPrimitiveFieldTypes(t *testing.T) {
	type config struct {
		Bool bool `env:"BOOL"`

		Word string `env:"FOO"`

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

		Float32 float32 `env:"FLOAT_32"`
		Float64 float64 `env:"FLOAT_64"`

		Complex64  complex64  `env:"COMPLEX_64"`
		Complex128 complex128 `env:"COMPLEX_128"`

		Byte byte `env:"BYTE"`
		Rune rune `env:"RUNE"`
	}

	tests := []struct {
		name     string
		expected config
		envVars  map[string]string
		wantErr  error
	}{
		// bool type fields
		{
			name: "set bool field with true as text",
			envVars: map[string]string{
				"BOOL": "true",
			},
			expected: config{Bool: true},
		}, {
			name: "set bool field with false as text",
			envVars: map[string]string{
				"BOOL": "false",
			},
			expected: config{Bool: false},
		}, {
			name: "set bool field with empty string",
			envVars: map[string]string{
				"BOOL": "",
			},
			expected: config{Bool: false},
		}, {
			name: "set bool field with 1",
			envVars: map[string]string{
				"BOOL": "1",
			},
			expected: config{Bool: true},
		}, {
			name: "set bool field with 0",
			envVars: map[string]string{
				"BOOL": "0",
			},
			expected: config{Bool: false},
		}, {
			name: "set bool field with invalid value",
			envVars: map[string]string{
				"BOOL": "invalid-string",
			},
			wantErr: ErrInvalidTypeConversion,
		},
		// string type field
		{
			name:     "set string field",
			expected: config{Word: "bar"},
			envVars: map[string]string{
				"FOO": "bar",
			},
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
		// unsigned integer type fields
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
		}, {
			name: "set unsigned uint8 field with negative number",
			envVars: map[string]string{
				"UINT_8": "-1",
			},
			wantErr: ErrInvalidTypeConversion,
		},
		// float type fields
		{
			name:     "set float32 field with valid value",
			expected: config{Float32: 1.23},
			envVars: map[string]string{
				"FLOAT_32": "1.23",
			},
		}, {
			name:     "set float32 field with Inf",
			expected: config{Float32: float32(math.Inf(+1))},
			envVars: map[string]string{
				"FLOAT_32": "+inf",
			},
		}, {
			name:     "set float32 field with invalid value on envvar",
			expected: config{Float32: 0},
			envVars: map[string]string{
				"FLOAT_32": "invalid.value",
			},
			wantErr: ErrInvalidTypeConversion,
		}, {
			name:     "set float64 field",
			expected: config{Float64: 123.456},
			envVars: map[string]string{
				"FLOAT_64": "123.456",
			},
		},
		// complex type fields
		{
			name:     "set complex64 field",
			expected: config{Complex64: 123.456 - 789.012i},
			envVars: map[string]string{
				"COMPLEX_64": "123.456-789.012i",
			},
		}, {
			name: "set complex64 field with invalid value",
			envVars: map[string]string{
				"COMPLEX_64": "invalid value",
			},
			wantErr: ErrInvalidTypeConversion,
		}, {
			name:     "set complex128 field",
			expected: config{Complex128: 123.45678901234567890 + 1i},
			envVars: map[string]string{
				"COMPLEX_128": "123.456789012345678901234567890+1i",
			},
		},
		// byte type fields
		{
			name:     "set byte field",
			expected: config{Byte: 255},
			envVars: map[string]string{
				"BYTE": "255",
			},
		}, {
			name: "set byte field with value larger than maximum value",
			envVars: map[string]string{
				"BYTE": "256",
			},
			wantErr: ErrInvalidTypeConversion,
		},
		// rune type fields
		//{
		//	name: "set rune field",
		//	envVars: map[string]string{
		//		"RUNE": "a",
		//	},
		//	expected: config{Rune: 'a'},
		//},
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

			if err := Set(&cfg); !errors.Is(err, tt.wantErr) {
				t.Errorf("err different than expected, want '%+v', got '%+v'", tt.wantErr, err)
				return
			}
			if cfg != tt.expected {
				t.Errorf("Set(&s), \nwant %+v,\ngot  %+v", tt.expected, cfg)
			}
		})
	}
}

func TestSetCustomTypes(t *testing.T) {
	type StringType string
	type FooType struct{}

	type config struct {
		StringField StringType `env:"STRING"`
		FooField    FooType    `env:"FOO"`
	}

	tests := []struct {
		name     string
		expected config
		envVars  map[string]string
		wantErr  error
	}{
		// bool type fields
		{
			name: "set custom string type without parser",
			envVars: map[string]string{
				"STRING": "any value",
			},
			expected: config{StringField: "any value"},
		}, {
			name: "set custom struct type without parser",
			envVars: map[string]string{
				"FOO": "any value",
			},
			wantErr: ErrParserNotAvailable,
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

			if err := Set(&cfg); !errors.Is(err, tt.wantErr) {
				t.Errorf("err different than expected, want '%+v', got '%+v'", tt.wantErr, err)
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
