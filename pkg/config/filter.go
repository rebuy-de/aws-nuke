package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mb0/glob"
	"errors"
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type FilterType string

const (
	FilterTypeEmpty    FilterType = ""
	FilterTypeExact               = "exact"
	FilterTypeGlob                = "glob"
	FilterTypeRegex               = "regex"
	FilterTypeContains            = "contains"
	FilterTypePcre                = "pcre"
)

type Filters map[string][]Filter

type Filter struct {
	Property string
	Type     FilterType
	Value    string
}

func (f Filter) Match(o string) (bool, error) {
	switch f.Type {
	case FilterTypeEmpty:
		fallthrough

	case FilterTypeExact:
		return f.Value == o, nil

	case FilterTypeContains:
		return strings.Contains(o, string(f.Value)), nil

	case FilterTypeGlob:
		return glob.Match(string(f.Value), o)

	case FilterTypeRegex:
		re, err := regexp.Compile(string(f.Value))
		if err != nil {
			return false, err
		}
		return re.MatchString(o), nil

	case FilterTypePcre:
		re, err := pcre.Compile(string(f.Value), 0)
		if err != nil {
			return false, errors.New(err.Message)
		}
		m := re.MatcherString(o, 0)
		return m.Matches(), nil

	default:
		return false, fmt.Errorf("unknown type %s", f.Type)
	}
}

func (f *Filter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string

	if unmarshal(&value) == nil {
		f.Type = FilterTypeExact
		f.Value = value
		return nil
	}

	m := map[string]string{}
	err := unmarshal(m)
	if err != nil {
		return err
	}

	f.Type = FilterType(m["type"])
	f.Value = m["value"]
	f.Property = m["property"]
	return nil
}

func NewExactFilter(value string) Filter {
	return Filter{
		Type:  FilterTypeExact,
		Value: value,
	}
}
