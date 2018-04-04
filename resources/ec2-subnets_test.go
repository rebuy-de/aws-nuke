package resources

import "testing"

func TestEC2Subnet(t *testing.T) {
	if err := ResourceTypeTest("aws_subnet", "subnet", "EC2Subnet", t); err != nil {
		t.Error(err)
	}
}
