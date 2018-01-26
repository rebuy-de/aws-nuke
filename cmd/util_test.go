package cmd

import (
	"fmt"
	"sort"
	"testing"

	"github.com/rebuy-de/aws-nuke/pkg/types"
)

func TestResolveResourceTypes(t *testing.T) {
	cases := []struct {
		base    types.Set
		include []types.Set
		exclude []types.Set
		result  types.Set
	}{
		{
			base:    types.Set{"a", "b", "c", "d"},
			include: []types.Set{types.Set{"a", "b", "c"}},
			result:  types.Set{"a", "b", "c"},
		},
		{
			base:    types.Set{"a", "b", "c", "d"},
			exclude: []types.Set{types.Set{"b", "d"}},
			result:  types.Set{"a", "c"},
		},
		{
			base:    types.Set{"a", "b"},
			include: []types.Set{types.Set{}},
			result:  types.Set{"a", "b"},
		},
		{
			base:    types.Set{"c", "b"},
			exclude: []types.Set{types.Set{}},
			result:  types.Set{"c", "b"},
		},
		{
			base:    types.Set{"a", "b", "c", "d"},
			include: []types.Set{types.Set{"a", "b", "c"}},
			exclude: []types.Set{types.Set{"a"}},
			result:  types.Set{"b", "c"},
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
