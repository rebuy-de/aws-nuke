package config

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Optum/aws-nuke/pkg/types"
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
		Regions:          []string{"eu-west-1", "stratoscale"},
		Accounts: map[string]Account{
			"555133742": Account{
				Presets: []string{"terraform"},
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
			Targets:  types.Collection{"DynamoDBTable", "S3Bucket", "S3Object"},
			Excludes: types.Collection{"IAMRole"},
		},
		Presets: map[string]PresetDefinitions{
			"terraform": PresetDefinitions{
				Filters: Filters{
					"S3Bucket": {
						Filter{
							Type:  FilterTypeGlob,
							Value: "my-statebucket-*",
						},
					},
				},
			},
		},
		CustomEndpoints: []*CustomRegion{
			&CustomRegion{
				Region:                "stratoscale",
				TLSInsecureSkipVerify: true,
				Services: CustomServices{
					&CustomService{
						Service: "ec2",
						URL:     "https://stratoscale.cloud.internal/api/v2/aws/ec2",
					},
					&CustomService{
						Service:               "s3",
						URL:                   "https://stratoscale.cloud.internal:1060",
						TLSInsecureSkipVerify: true,
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

func TestFilterMerge(t *testing.T) {
	config, err := Load("test-fixtures/example.yaml")
	if err != nil {
		t.Fatal(err)
	}

	filters, err := config.Filters("555133742")
	if err != nil {
		t.Fatal(err)
	}

	expect := Filters{
		"S3Bucket": []Filter{
			Filter{
				Type: "glob", Value: "my-statebucket-*",
			},
		},
		"IAMRole": []Filter{
			Filter{
				Type:  "exact",
				Value: "uber.admin",
			},
		},
		"IAMRolePolicyAttachment": []Filter{
			Filter{
				Type:  "exact",
				Value: "uber.admin -> AdministratorAccess",
			},
		},
	}

	if !reflect.DeepEqual(filters, expect) {
		t.Errorf("Read struct mismatches:")
		t.Errorf("  Got:      %#v", filters)
		t.Errorf("  Expected: %#v", expect)
	}
}

func TestGetCustomRegion(t *testing.T) {
	config, err := Load("test-fixtures/example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	stratoscale := config.CustomEndpoints.GetRegion("stratoscale")
	if stratoscale == nil {
		t.Fatal("Expected to find a set of custom endpoints for region10")
	}
	euWest1 := config.CustomEndpoints.GetRegion("eu-west-1")
	if euWest1 != nil {
		t.Fatal("Expected to euWest1 without a set of custom endpoints")
	}

	t.Run("TestGetService", func(t *testing.T) {
		ec2Service := stratoscale.Services.GetService("ec2")
		if ec2Service == nil {
			t.Fatal("Expected to find a custom ec2 service for region10")
		}
		rdsService := stratoscale.Services.GetService("rds")
		if rdsService != nil {
			t.Fatal("Expected to not find a custom rds service for region10")
		}

	})
}
