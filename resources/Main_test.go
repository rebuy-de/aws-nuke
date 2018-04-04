package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rebuy-de/aws-nuke/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/pkg/config"
	"github.com/rebuy-de/aws-nuke/pkg/terraform"
	"os"
	"testing"
)

var testConfigPath = "./test-config.yaml"
var awsProfile string
var tf *terraform.Terraform
var awsSess *session.Session

func SetUp() error {
	cfg, err := config.Load(testConfigPath)
	if err != nil {
		return err
	}

	// TODO: Make a flag out of it
	awsProfile = "AWSNuke"

	tf, err = terraform.New(cfg.Regions[0], awsProfile, ".")
	if err != nil {
		return err
	}

	if err := tf.Init(); err != nil {
		return err
	}

	account, err := awsutil.NewAccount(awsutil.Credentials{Profile: awsProfile})
	if err != nil {
		return err
	}

	// For safety reasons the tests require a configuration to make sure the
	// correct AWS account is used.
	if err := cfg.ValidateAccount(account.ID(), account.Aliases()); err != nil {
		return err
	}
	awsSess, err = account.NewSession(cfg.Regions[0])
	return err
}

func TestMain(m *testing.M) {
	if err := SetUp(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
