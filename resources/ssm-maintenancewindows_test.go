package resources

import "testing"

func TestSSMMaintenanceWindow(t *testing.T) {
	if err := ResourceTypeTest("aws_ssm_maintenance_window", "window", "SSMMaintenanceWindow", t); err != nil {
		t.Error(err)
	}
}
