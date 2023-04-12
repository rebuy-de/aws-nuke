package resources

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/rebuy-de/aws-nuke/v2/pkg/awsutil"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type AppRunnerVPCConnector struct {
	svc              *apprunner.AppRunner
	vpcConnectorArn  string
	vpcConnectorName string
	status           string
	createDate       *time.Time
}

func init() {
	register("AppRunnerVPCConnector", ListAppRunnerVPCConnectors)
}

func ListAppRunnerVPCConnectors(sess *session.Session) ([]Resource, error) {
	svc := apprunner.New(sess)
	resources := make([]Resource, 0)

	params := &apprunner.ListVpcConnectorsInput{
		MaxResults: aws.Int64(20),
	}

	for {
		resp, err := svc.ListVpcConnectors(params)
		if err != nil {
			// The ErrUnknownEndpoint occurs when the region doesn't support AppRunner so we will
			// skip those regions
			if _, ok := err.(awsutil.ErrUnknownEndpoint); ok {
				return resources, nil
			}
			return nil, err
		}

		for _, vpcConnector := range resp.VpcConnectors {
			resources = append(resources, &AppRunnerVPCConnector{
				svc:              svc,
				vpcConnectorArn:  *vpcConnector.VpcConnectorArn,
				vpcConnectorName: *vpcConnector.VpcConnectorName,
				status:           *vpcConnector.Status,
				createDate:       vpcConnector.CreatedAt,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *AppRunnerVPCConnector) Remove() error {

	_, err := f.svc.DeleteVpcConnector(&apprunner.DeleteVpcConnectorInput{
		VpcConnectorArn: &f.vpcConnectorArn,
	})

	return err
}

func (f *AppRunnerVPCConnector) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("VpcConnectorName", f.vpcConnectorName)
	properties.Set("CreateDate", f.createDate.Format(time.RFC3339))
	properties.Set("Status", f.status)

	return properties
}

func (f *AppRunnerVPCConnector) String() string {
	return f.vpcConnectorArn
}
