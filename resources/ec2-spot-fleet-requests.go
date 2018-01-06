package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2SpotFleetRequest struct {
	svc   *ec2.EC2
	id    string
	state string
}

func init() {
	register("EC2SpotFleetRequest", ListEC2SpotFleetRequests)
}

func ListEC2SpotFleetRequests(sess *session.Session) ([]Resource, error) {
	svc := ec2.New(sess)

	resp, err := svc.DescribeSpotFleetRequests(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, config := range resp.SpotFleetRequestConfigs {
		resources = append(resources, &EC2SpotFleetRequest{
			svc:   svc,
			id:    *config.SpotFleetRequestId,
			state: *config.SpotFleetRequestState,
		})
	}

	return resources, nil
}

func (i *EC2SpotFleetRequest) Filter() error {
	if i.state == "cancelled" {
		return fmt.Errorf("already cancelled")
	}
	return nil
}

func (i *EC2SpotFleetRequest) Remove() error {
	params := &ec2.CancelSpotFleetRequestsInput{
		TerminateInstances: aws.Bool(true),
		SpotFleetRequestIds: []*string{
			&i.id,
		},
	}

	_, err := i.svc.CancelSpotFleetRequests(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *EC2SpotFleetRequest) String() string {
	return i.id
}
