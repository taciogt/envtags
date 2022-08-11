package envtags

import (
	"os"
	"testing"
)

func TestSet(t *testing.T) {
	//tests := []struct {
	//	name string
	//	input interface{}
	//	want bool
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if got := Set(); got != tt.want {
	//			t.Errorf("Set() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}

	t.Run("first test", func(t *testing.T) {
		if err := os.Setenv("FOO", "bar"); err != nil {
			t.Error(err)
			return
		}

		type S struct {
			Foo string `env:"FOO"`
		}

		var s S

		//Set(s)  TODO: what happens here?
		Set(&s)

		if s.Foo != "bar b" {
			t.Errorf("unexpected env var set. expected=\"bar\". got=\"%s\"", s.Foo)
		}
	})

}
