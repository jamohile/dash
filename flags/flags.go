package flags

import "fmt"

// A utility to help working with flags.
// Makes it easier to create, merge, and format them.
type Flags struct {
	Values map[string]string
}

func (f *Flags) Merge(overrides Flags) Flags {
	merged := Flags{
		Values: make(map[string]string),
	}

	// Copy in base.
	for k, v := range f.Values {
		merged.Values[k] = v
	}

	// Copy in overrides.
	for k, v := range overrides.Values {
		merged.Values[k] = v
	}

	return merged
}

func (f *Flags) ToArgs() []string {
	formatted := make([]string, 0)

	for k, v := range f.Values {
		// Special case for empty flags, which are just presence based.
		if v != "" {
			formatted = append(formatted, fmt.Sprintf("--%s=%s", k, v))
		} else {
			formatted = append(formatted, fmt.Sprintf("--%s", k))
		}
	}

	return formatted
}

func FromValues(values map[string]string) Flags {
	return Flags{
		Values: values,
	}
}
