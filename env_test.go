package envtags

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"testing"
)

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
		},
		//{
		//	name: "set custom struct type without parser",
		//	envVars: map[string]string{
		//		"FOO": "any value",
		//	},
		//	wantErr: ErrParserNotAvailable,
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
