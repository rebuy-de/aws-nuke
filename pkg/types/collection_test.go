package types_test

import (
	"fmt"
	"testing"

	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func TestSetInterset(t *testing.T) {
	s1 := types.Collection{"a", "b", "c"}
	s2 := types.Collection{"b", "a", "d"}

	r := s1.Intersect(s2)

	want := fmt.Sprint([]string{"a", "b"})
	have := fmt.Sprint(r)

	if want != have {
		t.Errorf("Wrong result. Want: %s. Have: %s", want, have)
	}
}

func TestSetRemove(t *testing.T) {
	s1 := types.Collection{"a", "b", "c"}
	s2 := types.Collection{"b", "a", "d"}

	r := s1.Remove(s2)

	want := fmt.Sprint([]string{"c"})
	have := fmt.Sprint(r)

	if want != have {
		t.Errorf("Wrong result. Want: %s. Have: %s", want, have)
	}
}

func TestSetUnion(t *testing.T) {
	s1 := types.Collection{"a", "b", "c"}
	s2 := types.Collection{"b", "a", "d"}

	r := s1.Union(s2)

	want := fmt.Sprint([]string{"a", "b", "c", "d"})
	have := fmt.Sprint(r)

	if want != have {
		t.Errorf("Wrong result. Want: %s. Have: %s", want, have)
	}
}
