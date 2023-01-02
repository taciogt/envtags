package envtags

import (
	"errors"
	"os"
	"testing"
)

func TestStructTypes(t *testing.T) {
	type SimpleStruct struct {
		Int    int    `env:"INT"`
		String string `env:"STRING"`
	}
	type SecondStruct struct {
		InnerStruct SimpleStruct `env:",prefix:INNER_"`
	}

	type config struct {
		SimpleStructValue SimpleStruct
		WithPrefix        SimpleStruct `env:",prefix:PREFIX_"`
		DeepStruct        SecondStruct `env:",prefix:DEEP_"`
		//SimpleStructPointer *SimpleStruct
	}

	tests := []struct {
		name     string
		expected config
		envVars  map[string]string
		wantErr  error
	}{
		{
			name: "set struct fields for simple struct",
			envVars: map[string]string{
				"INT":    "123",
				"STRING": "a word",
			},
			expected: config{
				SimpleStructValue: SimpleStruct{
					Int:    123,
					String: "a word"},
			},
		},
		{
			name: "set struct fields with prefix option",
			envVars: map[string]string{
				"INT":           "123",
				"STRING":        "a word",
				"PREFIX_INT":    "456",
				"PREFIX_STRING": "another sentence",
			},
			expected: config{
				SimpleStructValue: SimpleStruct{
					Int:    123,
					String: "a word"},
				WithPrefix: SimpleStruct{Int: 456,
					String: "another sentence"},
			},
		},
		{
			name: "set deep struct fields",
			envVars: map[string]string{
				"DEEP_INNER_INT":    "1234",
				"DEEP_INNER_STRING": "deep word",
			},
			expected: config{
				DeepStruct: SecondStruct{
					SimpleStruct{
						Int:    1234,
						String: "deep word",
					}}},
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
