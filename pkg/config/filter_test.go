package config_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/rebuy-de/aws-nuke/pkg/config"
	yaml "gopkg.in/yaml.v2"
)

func TestUnmarshalFilter(t *testing.T) {
	past := time.Now().UTC().Add(-24 * time.Hour)
	future := time.Now().UTC().Add(24 * time.Hour)
	cases := []struct {
		yaml            string
		match, mismatch []string
	}{
		{
			yaml:     `foo`,
			match:    []string{"foo"},
			mismatch: []string{"fo", "fooo", "o", "fo"},
		},
		{
			yaml:     `{"type":"exact","value":"foo"}`,
			match:    []string{"foo"},
			mismatch: []string{"fo", "fooo", "o", "fo"},
		},
		{
			yaml:     `{"type":"glob","value":"b*sh"}`,
			match:    []string{"bish", "bash", "bosh", "bush", "boooooosh", "bsh"},
			mismatch: []string{"woooosh", "fooo", "o", "fo"},
		},
		{
			yaml:     `{"type":"glob","value":"b?sh"}`,
			match:    []string{"bish", "bash", "bosh", "bush"},
			mismatch: []string{"woooosh", "fooo", "o", "fo", "boooooosh", "bsh"},
		},
		{
			yaml:     `{"type":"regex","value":"b[iao]sh"}`,
			match:    []string{"bish", "bash", "bosh"},
			mismatch: []string{"woooosh", "fooo", "o", "fo", "boooooosh", "bsh", "bush"},
		},
		{
			yaml:     `{"type":"contains","value":"mba"}`,
			match:    []string{"bimbaz", "mba", "bi mba z"},
			mismatch: []string{"bim-baz"},
		},
		{
			yaml: `{"type":"dateOlderThan","value":"0"}`,
			match: []string{strconv.Itoa(int(future.Unix())),
				future.Format("2006-01-02"),
				future.Format("2006/01/02"),
				future.Format("2006-01-02T15:04:05Z"),
				future.Format(time.RFC3339Nano),
				future.Format(time.RFC3339),
			},
			mismatch: []string{"",
				strconv.Itoa(int(past.Unix())),
				past.Format("2006-01-02"),
				past.Format("2006/01/02"),
				past.Format("2006-01-02T15:04:05Z"),
				past.Format(time.RFC3339Nano),
				past.Format(time.RFC3339),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.yaml, func(t *testing.T) {
			var filter config.Filter

			err := yaml.Unmarshal([]byte(tc.yaml), &filter)
			if err != nil {
				t.Fatal(err)
			}

			for _, o := range tc.match {
				match, err := filter.Match(o)
				if err != nil {
					t.Fatal(err)
				}

				if !match {
					t.Fatalf("'%v' should match", o)
				}
			}

			for _, o := range tc.mismatch {
				match, err := filter.Match(o)
				if err != nil {
					t.Fatal(err)
				}

				if match {
					t.Fatalf("'%v' should not match", o)
				}
			}
		})
	}

}

func TestFilterFailureErrorsOnInput(t *testing.T) {
	cases := []struct {
		yaml            string
		fatal, nonfatal []string
	}{
		{
			yaml: `{"type":"dateOlderThan","value":"0"}`,
			fatal: []string{},
			nonfatal: []string{"IamNotADate"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.yaml, func(t *testing.T) {
			var filter config.Filter

			err := yaml.Unmarshal([]byte(tc.yaml), &filter)
			if err != nil {
				t.Fatal(err)
			}

			for _, o := range tc.fatal {
				match, err := filter.Match(o)

				if err == nil {
					t.Fatalf("'%v' should have returned an error", o)
				}

				if match {
					//is there ever a situation where a filter will fail but also return a true value?
					t.Fatalf("'%v' should not return true match", o) 
				}

				fferr, ok := err.(*config.FilterFailed)
				if ok {
					if !fferr.Fatal {
						t.Fatalf("'%v' should have returned fatal error", o)
					}
				} else {
					//even though returning a normal Error() is fatal enough, this test is about the FilterFailed object
					t.Fatalf("'%v' should have returned a FilterFailed Object", o)
				}
			}

			for _, o := range tc.nonfatal {
				match, err := filter.Match(o)

				if err == nil {
					t.Fatalf("'%v' should have returned an error", o)
				}

				if match {
					//is there ever a situation where a filter will fail but also return a true value?
					t.Fatalf("'%v' should not return true match", o) 
				}

				fferr, ok := err.(*config.FilterFailed)
				if ok {
					if fferr.Fatal {
						t.Fatalf("'%v' should have returned nonfatal error", o)
					}
				} else {
					t.Fatalf("'%v' should have returned a FilterFailed Object", o)
				}
			}
		})
	}

}
