package types_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func TestPropertiesEquals(t *testing.T) {
	cases := []struct {
		p1, p2 types.Properties
		result bool
	}{
		{
			p1:     nil,
			p2:     nil,
			result: true,
		},
		{
			p1:     nil,
			p2:     types.NewProperties(),
			result: false,
		},
		{
			p1:     types.NewProperties(),
			p2:     types.NewProperties(),
			result: true,
		},
		{
			p1:     types.NewProperties().Set("blub", "blubber"),
			p2:     types.NewProperties().Set("blub", "blubber"),
			result: true,
		},
		{
			p1:     types.NewProperties().Set("blub", "foo"),
			p2:     types.NewProperties().Set("blub", "bar"),
			result: false,
		},
		{
			p1:     types.NewProperties().Set("bim", "baz").Set("blub", "blubber"),
			p2:     types.NewProperties().Set("bim", "baz").Set("blub", "blubber"),
			result: true,
		},
		{
			p1:     types.NewProperties().Set("bim", "baz").Set("blub", "foo"),
			p2:     types.NewProperties().Set("bim", "baz").Set("blub", "bar"),
			result: false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if tc.p1.Equals(tc.p2) != tc.result {
				t.Errorf("Test Case failed. Want %t. Got %t.", !tc.result, tc.result)
				t.Errorf("p1: %s", tc.p1.String())
				t.Errorf("p2: %s", tc.p2.String())
			} else if tc.p2.Equals(tc.p1) != tc.result {
				t.Errorf("Test Case reverse check failed. Want %t. Got %t.", !tc.result, tc.result)
				t.Errorf("p1: %s", tc.p1.String())
				t.Errorf("p2: %s", tc.p2.String())
			}
		})
	}
}

func TestPropertiesSetTag(t *testing.T) {
	cases := []struct {
		name  string
		key   *string
		value interface{}
		want  string
	}{
		{
			name:  "string",
			key:   aws.String("name"),
			value: "blubber",
			want:  `[tag:name: "blubber"]`,
		},
		{
			name:  "string_ptr",
			key:   aws.String("name"),
			value: aws.String("blubber"),
			want:  `[tag:name: "blubber"]`,
		},
		{
			name:  "int",
			key:   aws.String("int"),
			value: 42,
			want:  `[tag:int: "42"]`,
		},
		{
			name:  "nil",
			key:   aws.String("nothing"),
			value: nil,
			want:  `[]`,
		},
		{
			name:  "empty_key",
			key:   aws.String(""),
			value: "empty",
			want:  `[]`,
		},
		{
			name:  "nil_key",
			key:   nil,
			value: "empty",
			want:  `[]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := types.NewProperties()

			p.SetTag(tc.key, tc.value)
			have := p.String()

			if tc.want != have {
				t.Errorf("'%s' != '%s'", tc.want, have)
			}
		})
	}
}

func TestPropertiesSetTagWithPrefix(t *testing.T) {
	cases := []struct {
		name   string
		prefix string
		key    *string
		value  interface{}
		want   string
	}{
		{
			name:   "empty",
			prefix: "",
			key:    aws.String("name"),
			value:  "blubber",
			want:   `[tag:name: "blubber"]`,
		},
		{
			name:   "nonempty",
			prefix: "bish",
			key:    aws.String("bash"),
			value:  "bosh",
			want:   `[tag:bish:bash: "bosh"]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := types.NewProperties()

			p.SetTagWithPrefix(tc.prefix, tc.key, tc.value)
			have := p.String()

			if tc.want != have {
				t.Errorf("'%s' != '%s'", tc.want, have)
			}
		})
	}
}
