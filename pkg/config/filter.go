package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mb0/glob"
)

type FilterType string

const (
	FilterTypeEmpty         FilterType = ""
	FilterTypeExact                    = "exact"
	FilterTypeGlob                     = "glob"
	FilterTypeRegex                    = "regex"
	FilterTypeContains                 = "contains"
	FilterTypeDateOlderThan            = "dateOlderThan"
)

type Filters map[string][]Filter

func (f Filters) Merge(f2 Filters) {
	for resourceType, filter := range f2 {
		f[resourceType] = append(f[resourceType], filter...)
	}
}

type Filter struct {
	Property string
	Type     FilterType
	Value    string
	Invert   string
}

func (f Filter) Match(o string) (bool, error) {
	switch f.Type {
	case FilterTypeEmpty:
		fallthrough

	case FilterTypeExact:
		return f.Value == o, nil

	case FilterTypeContains:
		return strings.Contains(o, f.Value), nil

	case FilterTypeGlob:
		return glob.Match(f.Value, o)

	case FilterTypeRegex:
		re, err := regexp.Compile(f.Value)
		if err != nil {
			return false, err
		}
		return re.MatchString(o), nil

	case FilterTypeDateOlderThan:
		if o == "" {
			return false, nil
		}
		duration, err := time.ParseDuration(f.Value)
		if err != nil {
			return false, err
		}
		fieldTime, err := parseDate(o)
		if err != nil {
			return false, err
		}
		fieldTimeWithOffset := fieldTime.Add(duration)

		return fieldTimeWithOffset.After(time.Now()), nil

	default:
		return false, fmt.Errorf("unknown type %s", f.Type)
	}
}

func parseDate(input string) (time.Time, error) {
	if i, err := strconv.ParseInt(input, 10, 64); err == nil {
		t := time.Unix(i, 0)
		return t, nil
	}

	formats := []string{"2006-01-02",
		"2006/01/02",
		"2006-01-02T15:04:05Z",
		time.RFC3339Nano, // Format of t.MarshalText() and t.MarshalJSON()
		time.RFC3339,
	}
	for _, f := range formats {
		t, err := time.Parse(f, input)
		if err == nil {
			return t, nil
		}
	}
	return time.Now(), fmt.Errorf("unable to parse time %s", input)
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
	f.Invert = m["invert"]
	return nil
}

func NewExactFilter(value string) Filter {
	return Filter{
		Type:  FilterTypeExact,
		Value: value,
	}
}
