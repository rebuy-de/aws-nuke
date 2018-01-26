package types_test

import (
	"sort"
	"testing"

	"github.com/rebuy-de/aws-nuke/pkg/types"
)

func TestSetRetain(t *testing.T) {
	s1 := types.Set{"a", "b", "c"}
	s2 := types.Set{"b", "a", "d"}

	r := s1.Intersect(s2)
	sort.Strings(r)

	if len(r) != 2 || r[0] != "a" || r[1] != "b" {
		t.Errorf("Wrong result. Want: [a, b]. Got: %v", r)
	}
}
