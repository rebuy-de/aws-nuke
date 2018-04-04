package resources

import "testing"

func TestEC2Volume(t *testing.T) {
	if err := ResourceTypeTest("aws_ebs_volume", "volume", "EC2Volume", t); err != nil {
		t.Error(err)
	}
}
