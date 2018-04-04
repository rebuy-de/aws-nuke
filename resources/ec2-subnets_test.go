package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
)

func TestListEC2Subnets(t *testing.T) {
	rs, err := ListEC2Subnets(awsSess)
	if err != nil {
		t.Fatal(err)
	}

	if len(rs) != 0 {
		t.Fatalf("expected %v, found %v", 0, len(rs))
	}

	if err := tf.CreateResource("aws_subnet.subnet"); err != nil {
		t.Fatal(err)
	}

	rs, err = ListEC2Subnets(awsSess)
	if err != nil {
		t.Error(err)
	}

	if len(rs) != 1 {
		t.Errorf("expected %v, found %v", 1, len(rs))
	}

	if err := tf.RemoveAllResources(); err != nil {
		t.Fatal(err)
	}
}

func TestEC2Subnet_Remove(t *testing.T) {
	if err := tf.CreateResource("aws_subnet.subnet"); err != nil {
		t.Fatal(err)
	}

	id, err := tf.ResourceProperty("aws_subnet.subnet", "id")
	if err != nil {
		t.Fatal(err)
	}

	subnet := EC2Subnet{ec2.New(awsSess), aws.String(id)}
	if err := subnet.Remove(); err != nil {
		t.Error(err)
	}

	if err := tf.RemoveAllResources(); err != nil {
		t.Fatal(err)
	}
}
