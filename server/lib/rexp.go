package lib

import (
	"regexp"
)

// StringMap -
type StringMap map[string]string

// Rexp -
type Rexp struct {
	Exp *regexp.Regexp
}

// Match -
func (r *Rexp) Match(s string) StringMap {
	match := r.Exp.FindStringSubmatch(s)
	if match == nil {
		return nil
	}

	captures := make(map[string]string)

	for i, name := range r.Exp.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]
	}
	return captures
}
