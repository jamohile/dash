package flags

import (
	"reflect"
	"sort"
	"testing"
)

func TestFormatAsArgs(t *testing.T) {
	f := Flags{
		Values: map[string]string{
			"foo": "1",
			"bar": "abc",
		},
	}

	args := f.ToArgs()
	expected := []string{"--foo=1", "--bar=abc"}

	sort.Strings(args)
	sort.Strings(expected)

	if !reflect.DeepEqual(args, expected) {
		t.Fatalf("Expected %v, got %v", expected, args)
	}
}

func TestMerge(t *testing.T) {
	f := Flags{
		Values: map[string]string{
			"foo": "1",
			"bar": "abc",
		},
	}

	merged := f.Merge(Flags{
		Values: map[string]string{
			"bar":   "def",
			"hello": "world",
		},
	})

	args := merged.ToArgs()
	expected := []string{"--foo=1", "--bar=def", "--hello=world"}

	sort.Strings(args)
	sort.Strings(expected)

	if !reflect.DeepEqual(args, expected) {
		t.Fatalf("Expected %v, got %v", expected, args)
	}
}
