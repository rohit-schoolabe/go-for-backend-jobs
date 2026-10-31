package port

import "testing"

func TestAddr(t *testing.T) {
	cases := []struct {
		env  map[string]string
		want string
	}{
		{map[string]string{}, ":8080"},
		{map[string]string{"PORT": "9000"}, ":9000"},
		{map[string]string{"PORT": "9000",
			"ADDR": "127.0.0.1:7000"}, "127.0.0.1:7000"},
	}
	for _, tc := range cases {
		got := Addr(func(k string) string { return tc.env[k] })
		if got != tc.want {
			t.Errorf("%v: got %q, want %q",
				tc.env, got, tc.want)
		}
	}
}
