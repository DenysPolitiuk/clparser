package clparser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		in    []string
		out   map[string]string
		valid bool
	}{
		{strings.Split("-t 30", " "), map[string]string{"-t": "30"}, true},
		{strings.Split("-g", " "), map[string]string{"-g": ""}, true},
		{strings.Split("-g -i", " "), map[string]string{"-g": "", "-i": ""}, true},
		{strings.Split("--name", " "), map[string]string{"--name": ""}, true},
		{strings.Split("--name=denys", " "), map[string]string{"--name": "denys"}, true},
		{[]string{}, map[string]string{}, true},
		{strings.Split("--name=denys=pol", " "), nil, false},
		{[]string{"-"}, nil, false},
	}
	for _, c := range cases {
		got, err := Parse(c.in)
		switch c.valid {
		case true:
			if err != nil {
				t.Errorf("Parser(%q) returned %q instead of nil", c.in, err)
			}
			if !verifyMapsEqual(c.out, got) {
				t.Errorf("Parser(%q) returned %q instead of %q", c.in, got, c.out)
			}
		case false:
			if err == nil {
				t.Errorf("Parser(%q) returned nil instead of %q", c.in, got)
			}
		}
	}
}

func verifyMapsEqual(first map[string]string, second map[string]string) bool {
	for k, v := range first {
		if second[k] != v {
			return false
		}
	}
	return true
}
