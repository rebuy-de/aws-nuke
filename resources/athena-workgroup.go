package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/aws/aws-sdk-go/service/athena/athenaiface"
)

type AthenaWorkGroup struct {
	svc     athenaiface.AthenaAPI
	summary *athena.WorkGroupSummary
}

func init() {
	register("AthenaWorkGroup", ListAthenaWorkGroup)
}

func ListAthenaWorkGroup(sess *session.Session) ([]Resource, error) {
	svc := athena.New(sess)
	resources := []Resource{}

	params := &athena.ListWorkGroupsInput{
		MaxResults: aws.Int64(20),
	}

	for {
		resp, err := svc.ListWorkGroups(params)
		if err != nil {
			return nil, err
		}

		for _, workGroupSummary := range resp.WorkGroups {
			resources = append(resources, &AthenaWorkGroup{
				svc:     svc,
				summary: workGroupSummary,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *AthenaWorkGroup) Remove() error {

	_, err := f.svc.DeleteWorkGroup(&athena.DeleteWorkGroupInput{
		WorkGroup: f.summary.Name,
	})

	return err
}

func (f *AthenaWorkGroup) String() string {
	return *f.summary.Name
}
