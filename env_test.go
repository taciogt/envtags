package envtags

import (
	"errors"
	"os"
	"testing"
)

func TestSetFieldTypes(t *testing.T) {
	type Config struct {
		Word string `env:"FOO"`

		Number int `env:"NUMBER"`
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
		expected: Config{Number: 123},
		envVars: map[string]string{
			"NUMBER": "123",
		},
	},
		{
			name:    "set integer field with invalid env var",
			wantErr: ErrInvalidTypeConversion,
			envVars: map[string]string{
				"NUMBER": "abc",
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
			}
			if cfg != tt.expected {
				t.Errorf("Set(&s), want %+v, got %+v", tt.expected, cfg)
			}
		})
	}
}
