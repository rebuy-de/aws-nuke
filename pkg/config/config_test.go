package config

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestConfigBlacklist(t *testing.T) {
	config := new(NukeConfig)

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
	config, err := LoadConfig("test-fixtures/example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	expect := NukeConfig{
		AccountBlacklist: []string{"1234567890"},
		Regions:          []string{"eu-west-1"},
		Accounts: map[string]NukeConfigAccount{
			"555133742": NukeConfigAccount{
				Filters: map[string][]string{
					"IAMRole": []string{
						"uber.admin",
					},
					"IAMRolePolicyAttachment": []string{
						"uber.admin -> AdministratorAccess",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(*config, expect) {
		t.Errorf("Read struct mismatches:")
		t.Errorf("  Got:      %#v", *config)
		t.Errorf("  Expected: %#v", expect)
	}
}

func TestResolveDeprecations(t *testing.T) {
	config := NukeConfig{
		AccountBlacklist: []string{"1234567890"},
		Regions:          []string{"eu-west-1"},
		Accounts: map[string]NukeConfigAccount{
			"555133742": {
				Filters: map[string][]string{
					"IamRole":                 {"uber.admin", "foo.bar"},
					"IAMRolePolicyAttachment": {"uber.admin -> AdministratorAccess"},
				},
			},
			"2345678901": {
				Filters: map[string][]string{
					"ECRrepository":           {"foo:bar", "bar:foo"},
					"IAMRolePolicyAttachment": {"uber.admin -> AdministratorAccess"},
				},
			},
		},
	}

	expect := map[string]NukeConfigAccount{
		"555133742": {
			Filters: map[string][]string{
				"IAMRole":                 {"uber.admin", "foo.bar"},
				"IAMRolePolicyAttachment": {"uber.admin -> AdministratorAccess"},
			},
		},
		"2345678901": {
			Filters: map[string][]string{
				"ECRRepository":           {"foo:bar", "bar:foo"},
				"IAMRolePolicyAttachment": {"uber.admin -> AdministratorAccess"},
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

	invalidConfig := NukeConfig{
		AccountBlacklist: []string{"1234567890"},
		Regions:          []string{"eu-west-1"},
		Accounts: map[string]NukeConfigAccount{
			"555133742": {
				Filters: map[string][]string{
					"IamUserAccessKeys": {"X"},
					"IAMUserAccessKey":  {"Y"},
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
	config, err := LoadConfig("test-fixtures/example.yaml")
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
