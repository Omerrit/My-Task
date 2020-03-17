package filter

import (
	"fmt"
	"path/filepath"
	"regexp"
)

type FilterFunc = func(imp string) (bool, error)

type FilterSet []FilterFunc

func (fset *FilterSet) Match(s string) (bool, error) {
	if fset == nil {
		return false, nil
	}
	for _, subf := range *fset {
		ok, err := subf(s)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func (fset *FilterSet) IsEmpty() bool {
	return fset == nil || len(*fset) == 0
}

func NewFilter(pattern string) (FilterFunc, error) {
	if pattern == "" {
		f := func(_ string) (bool, error) {
			return false, nil
		}
		return f, nil
	}
	if pattern == "*" { // match all
		pattern = ".*"
	}
	// '^' and '$' need to match case where a string equal to the pattern
	reg, err := regexp.Compile("^" + pattern + "$")
	if err == nil {
		f := func(s string) (bool, error) {
			return reg.MatchString(s), nil
		}
		return f, nil
	}
	_, err = filepath.Match(pattern, "ahaha")
	if err == nil {
		f := func(s string) (bool, error) {
			return filepath.Match(pattern, s)
		}
		return f, nil
	}
	return nil, fmt.Errorf("Pattern <%s> is invalid", pattern)
}

func NewFilterSet(patterns []string) (FilterSet, error) {
	fs := FilterSet{}
	for _, pattern := range patterns {
		subf, err := NewFilter(pattern)
		if err != nil {
			return nil, fmt.Errorf("pattern(%v):%w", pattern, err)
		}
		fs = append(fs, subf)
	}
	return fs, nil
}

func (fset *FilterSet) Append(other FilterSet) {
	if *fset == nil {
		*fset = make(FilterSet, 0)
	}
	*fset = append(*fset, other...)
}
