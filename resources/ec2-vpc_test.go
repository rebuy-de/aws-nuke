package resources

import "testing"

func TestEC2VPC(t *testing.T) {
	if err := ResourceTypeTest("aws_vpc", "vpc", "EC2VPC", t); err != nil {
		t.Error(err)
	}
}
