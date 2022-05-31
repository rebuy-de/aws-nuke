package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/worklink"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type WorkLinkFleet struct {
	svc              *worklink.WorkLink
	fleetARN         *string
	fleetName        *string
	fleetCompanyCode *string
	fleetDisplayName *string
}

func init() {
	register("WorkLinkFleet", ListWorkLinkFleets)
}

func ListWorkLinkFleets(sess *session.Session) ([]Resource, error) {
	svc := worklink.New(sess)
	resources := []Resource{}

	params := &worklink.ListFleetsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListFleets(params)
		if err != nil {
			return nil, err
		}

		for _, fleet := range output.FleetSummaryList {
			resources = append(resources, &WorkLinkFleet{
				svc:              svc,
				fleetARN:         fleet.FleetArn,
				fleetName:        fleet.FleetName,
				fleetCompanyCode: fleet.CompanyCode,
				fleetDisplayName: fleet.DisplayName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *WorkLinkFleet) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("CompanyCode", f.fleetCompanyCode)
	properties.Set("DisplayName", f.fleetDisplayName)

	return properties
}

func (f *WorkLinkFleet) Remove() error {
	_, err := f.svc.DeleteFleet(&worklink.DeleteFleetInput{
		FleetArn: f.fleetARN,
	})

	return err
}

func (f *WorkLinkFleet) String() string {
	return *f.fleetName
}
