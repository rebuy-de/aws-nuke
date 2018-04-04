package terraform

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"testing"
)

var profile string
var region string
var ec2Client *ec2.EC2

func TestMain(m *testing.M) {
	// Due to the invasive nature of these tests, an explicit AWS profile is required.
	flag.StringVar(&profile, "profile", "", "AWS Profile to use for the test")
	flag.StringVar(&region, "region", "", "AWS region to use for the test")
	flag.Parse()
	if profile == "" {
		fmt.Println("profile flag is required")
		os.Exit(1)
	}
	if region == "" {
		fmt.Println("region flag is required")
		os.Exit(1)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", profile),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ec2Client = ec2.New(sess)

	code := m.Run()
	ec2Client.DeleteVpc(&ec2.DeleteVpcInput{})
	os.Exit(code)
}

func TestInit(t *testing.T) {
	tf, err := New(region, profile, "testdata")
	Assert(t, err)
	Assert(t, tf.Init())
}

func TestCreatePlan(t *testing.T) {
	tf, err := New(region, profile, "testdata")
	Assert(t, err)
	Assert(t, tf.Init())
	plan, _, err := tf.CreatePlan([]string{"aws_vpc.vpc"})
	Assert(t, err)

	if _, err := os.Stat(plan.Name()); os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func TestCreateResource(t *testing.T) {
	tf, err := New(region, profile, "testdata")
	Assert(t, err)
	Assert(t, tf.Init())
	Assert(t, tf.CreateResource("aws_vpc.vpc"))

	res, err := ec2Client.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String("AWSNukeTerraformTest")},
		}},
	})
	Assert(t, err)

	if len(res.Vpcs) == 0 {
		t.Fatal("no VPCs with tag Name:AWSNukeTerraformTest exist")
	}
	if len(res.Vpcs) > 1 {
		t.Errorf("found %v VPCs called AWSNukeTerraformTest, expected 1", len(res.Vpcs))
	}
	for _, vpc := range res.Vpcs {
		if _, err := ec2Client.DeleteVpc(&ec2.DeleteVpcInput{VpcId: vpc.VpcId}); err != nil {
			t.Errorf("Failed to remove VPC with tag Name:AWSNukeTerraform as part of cleanup")
		}
	}
}

func TestRemoveResource(t *testing.T) {
	tf, err := New(region, profile, "testdata")
	Assert(t, err)
	Assert(t, tf.Init())
	Assert(t, tf.CreateResource("aws_vpc.vpc"))
	Assert(t, tf.RemoveResource("aws_vpc.vpc"))

	res, err := ec2Client.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String("AWSNukeTerraformTest")},
		}},
	})
	Assert(t, err)

	if len(res.Vpcs) != 0 {
		t.Fatalf("expected VPC to be removed but found: %v", res.Vpcs)
	}
}

func TestApplyPlan(t *testing.T) {
	tf, err := New(region, profile, "testdata")
	Assert(t, err)
	Assert(t, tf.Init())
	plan, _, err := tf.CreatePlan([]string{"aws_vpc.vpc"})
	Assert(t, err)

	Assert(t, tf.ApplyPlan(plan))

	res, err := ec2Client.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String("AWSNukeTerraformTest")},
		}},
	})
	if err != nil {
		t.Fatalf("failed to describe VPCs. May exit in bad state: %v", err)
	}

	if len(res.Vpcs) == 0 {
		t.Fatal("no VPCs with tag Name:AWSNukeTerraformTest exist")
	}
	if len(res.Vpcs) > 1 {
		t.Errorf("found %v VPCs called AWSNukeTerraformTest, expected 1", len(res.Vpcs))
	}

	for _, vpc := range res.Vpcs {
		if _, err := ec2Client.DeleteVpc(&ec2.DeleteVpcInput{VpcId: vpc.VpcId}); err != nil {
			t.Errorf("Failed to remove VPC with tag Name:AWSNukeTerraform as part of cleanup")
		}
	}
}

func TestResourceProperty(t *testing.T) {
	tf, err := New(region, profile, "testdata")
	Assert(t, err)
	Assert(t, tf.Init())
	plan, _, err := tf.CreatePlan([]string{"aws_vpc.vpc"})
	Assert(t, err)
	Assert(t, tf.ApplyPlan(plan))

	res, err := ec2Client.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String("AWSNukeTerraformTest")},
		}},
	})
	Assert(t, err)

	if id, err := tf.ResourceProperty("aws_vpc.vpc", "id"); err != nil {
		fmt.Println(id)
		t.Error(err)
	} else if id != *res.Vpcs[0].VpcId {
		t.Errorf("expected VPC ID to be %v but got %v", *res.Vpcs[0].VpcId, id)
	}

	for _, vpc := range res.Vpcs {
		if _, err := ec2Client.DeleteVpc(&ec2.DeleteVpcInput{VpcId: vpc.VpcId}); err != nil {
			t.Errorf("Failed to remove VPC with tag Name:AWSNukeTerraform as part of cleanup")
		}
	}
}

func Assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
