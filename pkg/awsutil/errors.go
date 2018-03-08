package awsutil

import "fmt"

type ErrServiceNotInRegion struct {
	Region  string
	Service string
}

func (err ErrServiceNotInRegion) Error() string {
	return fmt.Sprintf(
		"service '%s' is not available in region '%s'",
		err.Service, err.Region)
}
