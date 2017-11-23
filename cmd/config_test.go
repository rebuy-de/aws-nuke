package cmd

import (
	"os"
	"path"
	"reflect"
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
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	config, err := LoadConfig(path.Join(cwd, "..", "config", "example.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	expect := NukeConfig{
		AccountBlacklist: []string{"1234567890"},
		Regions:          []string{"eu-west-1"},
		Accounts: map[string]NukeConfigAccount{
			"555133742": NukeConfigAccount{
				Filters: map[string][]string{
					"IamRole": []string{
						"uber.admin",
					},
					"IamRolePolicyAttachment": []string{
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
