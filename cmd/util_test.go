package cmd

import (
	"fmt"
	"sort"
	"testing"

	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func TestResolveResourceTypes(t *testing.T) {
	cases := []struct {
		name         string
		base         types.Collection
		mapping      map[string]string
		include      []types.Collection
		exclude      []types.Collection
		cloudControl []types.Collection
		result       types.Collection
	}{
		{
			base:    types.Collection{"a", "b", "c", "d"},
			include: []types.Collection{{"a", "b", "c"}},
			result:  types.Collection{"a", "b", "c"},
		},
		{
			base:    types.Collection{"a", "b", "c", "d"},
			exclude: []types.Collection{{"b", "d"}},
			result:  types.Collection{"a", "c"},
		},
		{
			base:    types.Collection{"a", "b"},
			include: []types.Collection{{}},
			result:  types.Collection{"a", "b"},
		},
		{
			base:    types.Collection{"c", "b"},
			exclude: []types.Collection{{}},
			result:  types.Collection{"c", "b"},
		},
		{
			base:    types.Collection{"a", "b", "c", "d"},
			include: []types.Collection{{"a", "b", "c"}},
			exclude: []types.Collection{{"a"}},
			result:  types.Collection{"b", "c"},
		},
		{
			name:         "CloudControlAdd",
			base:         types.Collection{"a", "b"},
			cloudControl: []types.Collection{{"x"}},
			result:       types.Collection{"a", "b", "x"},
		},
		{
			name:         "CloudControlReplaceOldStyle",
			base:         types.Collection{"a", "b", "c"},
			mapping:      map[string]string{"z": "b"},
			cloudControl: []types.Collection{{"z"}},
			result:       types.Collection{"a", "z", "c"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := ResolveResourceTypes(tc.base, tc.mapping, tc.include, tc.exclude, tc.cloudControl)

			sort.Strings(r)
			sort.Strings(tc.result)

			var (
				want = fmt.Sprint(tc.result)
				have = fmt.Sprint(r)
			)

			if want != have {
				t.Fatalf("Wrong result. Want: %s. Have: %s", want, have)
			}
		})
	}
}

func TestIsTrue(t *testing.T) {
	falseStrings := []string{"", "false", "treu", "foo"}
	for _, fs := range falseStrings {
		if IsTrue(fs) {
			t.Fatalf("IsTrue falsely returned 'true' for: %s", fs)
		}
	}

	trueStrings := []string{"true", " true", "true ", " TrUe "}
	for _, ts := range trueStrings {
		if !IsTrue(ts) {
			t.Fatalf("IsTrue falsely returned 'false' for: %s", ts)
		}
	}
}
