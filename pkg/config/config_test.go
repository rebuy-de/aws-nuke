package config

import (
	"fmt"
	"github.com/rebuy-de/aws-nuke/pkg/types"
	"reflect"
	"strings"
	"testing"
)

func TestConfigBlacklist(t *testing.T) {
	config := new(Nuke)

	if config.HasBlacklist() {
		t.Errorf("HasBlacklist() returned true on a nil backlist.")
	}

	if config.InBlacklist("blubber") {
		t.Errorf("InBlacklist() returned true on a nil backlist.")
	}

	config.AccountBlacklist = []string{}

	if config.HasBlacklist() {
		t.Errorf("HasBlacklist() returned true on a empty backlist.")
	}

	if config.InBlacklist("foobar") {
		t.Errorf("InBlacklist() returned true on a empty backlist.")
	}

	config.AccountBlacklist = append(config.AccountBlacklist, "bim")

	if !config.HasBlacklist() {
		t.Errorf("HasBlacklist() returned false on a backlist with one element.")
	}

	if !config.InBlacklist("bim") {
		t.Errorf("InBlacklist() returned false on looking up an existing value.")
	}

	if config.InBlacklist("baz") {
		t.Errorf("InBlacklist() returned true on looking up an non existing value.")
	}
}

func TestLoadExampleConfig(t *testing.T) {
	config, err := Load("test-fixtures/example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	expect := Nuke{
		AccountBlacklist: []string{"1234567890"},
		Regions:          []string{"eu-west-1"},
		Accounts: map[string]Account{
			"555133742": Account{
				Filters: Filters{
					"IAMRole": {
						NewExactFilter("uber.admin"),
					},
					"IAMRolePolicyAttachment": {
						NewExactFilter("uber.admin -> AdministratorAccess"),
					},
				},
				ResourceTypes: ResourceTypes{
					types.Collection{"S3Bucket"},
					nil,
				},
			},
		},
		ResourceTypes: ResourceTypes{
			Targets: types.Collection{"S3Object", "S3Bucket"},
			Excludes: types.Collection{"IAMRole"},
		},
	}

	if !reflect.DeepEqual(*config, expect) {
		t.Errorf("Read struct mismatches:")
		t.Errorf("  Got:      %#v", *config)
		t.Errorf("  Expected: %#v", expect)
	}
}

func TestResolveDeprecations(t *testing.T) {
	config := Nuke{
		AccountBlacklist: []string{"1234567890"},
		Regions:          []string{"eu-west-1"},
		Accounts: map[string]Account{
			"555133742": {
				Filters: Filters{
					"IamRole": {
						NewExactFilter("uber.admin"),
						NewExactFilter("foo.bar"),
					},
					"IAMRolePolicyAttachment": {
						NewExactFilter("uber.admin -> AdministratorAccess"),
					},
				},
			},
			"2345678901": {
				Filters: Filters{
					"ECRrepository": {
						NewExactFilter("foo:bar"),
						NewExactFilter("bar:foo"),
					},
					"IAMRolePolicyAttachment": {
						NewExactFilter("uber.admin -> AdministratorAccess"),
					},
				},
			},
		},
	}

	expect := map[string]Account{
		"555133742": {
			Filters: Filters{
				"IAMRole": {
					NewExactFilter("uber.admin"),
					NewExactFilter("foo.bar"),
				},
				"IAMRolePolicyAttachment": {
					NewExactFilter("uber.admin -> AdministratorAccess"),
				},
			},
		},
		"2345678901": {
			Filters: Filters{
				"ECRRepository": {
					NewExactFilter("foo:bar"),
					NewExactFilter("bar:foo"),
				},
				"IAMRolePolicyAttachment": {
					NewExactFilter("uber.admin -> AdministratorAccess"),
				},
			},
		},
	}

	err := config.resolveDeprecations()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(config.Accounts, expect) {
		t.Errorf("Read struct mismatches:")
		t.Errorf("  Got:      %#v", config.Accounts)
		t.Errorf("  Expected: %#v", expect)
	}

	invalidConfig := Nuke{
		AccountBlacklist: []string{"1234567890"},
		Regions:          []string{"eu-west-1"},
		Accounts: map[string]Account{
			"555133742": {
				Filters: Filters{
					"IamUserAccessKeys": {
						NewExactFilter("X")},
					"IAMUserAccessKey": {
						NewExactFilter("Y")},
				},
			},
		},
	}

	err = invalidConfig.resolveDeprecations()
	if err == nil || !strings.Contains(err.Error(), "using deprecated resource type and replacement") {
		t.Fatal("invalid config did not cause correct error")
	}
}

func TestConfigValidation(t *testing.T) {
	config, err := Load("test-fixtures/example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		ID         string
		Aliases    []string
		ShouldFail bool
	}{
		{ID: "555133742", Aliases: []string{"staging"}, ShouldFail: false},
		{ID: "1234567890", Aliases: []string{"staging"}, ShouldFail: true},
		{ID: "1111111111", Aliases: []string{"staging"}, ShouldFail: true},
		{ID: "555133742", Aliases: []string{"production"}, ShouldFail: true},
		{ID: "555133742", Aliases: []string{}, ShouldFail: true},
		{ID: "555133742", Aliases: []string{"staging", "prod"}, ShouldFail: true},
	}

	for i, tc := range cases {
		name := fmt.Sprintf("%d_%s/%v/%t", i, tc.ID, tc.Aliases, tc.ShouldFail)
		t.Run(name, func(t *testing.T) {
			err := config.ValidateAccount(tc.ID, tc.Aliases)
			if tc.ShouldFail && err == nil {
				t.Fatal("Expected an error but didn't get one.")
			}
			if !tc.ShouldFail && err != nil {
				t.Fatalf("Didn't excpect an error, but got one: %v", err)
			}
		})
	}
}
