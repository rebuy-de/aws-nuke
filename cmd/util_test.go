package cmd

import (
	"fmt"
	"sort"
	"testing"

	"github.com/Optum/aws-nuke/pkg/types"
)

func TestResolveResourceTypes(t *testing.T) {
	cases := []struct {
		base    types.Collection
		include []types.Collection
		exclude []types.Collection
		result  types.Collection
	}{
		{
			base:    types.Collection{"a", "b", "c", "d"},
			include: []types.Collection{types.Collection{"a", "b", "c"}},
			result:  types.Collection{"a", "b", "c"},
		},
		{
			base:    types.Collection{"a", "b", "c", "d"},
			exclude: []types.Collection{types.Collection{"b", "d"}},
			result:  types.Collection{"a", "c"},
		},
		{
			base:    types.Collection{"a", "b"},
			include: []types.Collection{types.Collection{}},
			result:  types.Collection{"a", "b"},
		},
		{
			base:    types.Collection{"c", "b"},
			exclude: []types.Collection{types.Collection{}},
			result:  types.Collection{"c", "b"},
		},
		{
			base:    types.Collection{"a", "b", "c", "d"},
			include: []types.Collection{types.Collection{"a", "b", "c"}},
			exclude: []types.Collection{types.Collection{"a"}},
			result:  types.Collection{"b", "c"},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			r := ResolveResourceTypes(tc.base, tc.include, tc.exclude)

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
		if ! IsTrue(ts) {
			t.Fatalf("IsTrue falsely returned 'false' for: %s", ts)
		}
	}
}
