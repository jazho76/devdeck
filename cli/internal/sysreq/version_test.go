package sysreq

import "testing"

func TestMeetsMinimum(t *testing.T) {
	cases := []struct {
		current, min string
		want         bool
	}{
		{"3.6b", "3.5a", true},
		{"3.5a", "3.5a", true},
		{"3.10", "3.5a", true},
		{"3.4", "3.5a", false},
		{"3.5", "3.5a", false},
		{"4.0", "3.5a", true},
		{"3.5b", "3.5a", true},
	}
	for _, c := range cases {
		if got := meetsMinimum(c.current, c.min); got != c.want {
			t.Errorf("meetsMinimum(%q, %q) = %v, want %v", c.current, c.min, got, c.want)
		}
	}
}
