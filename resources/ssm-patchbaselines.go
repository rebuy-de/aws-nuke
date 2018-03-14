package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMPatchBaseline struct {
	svc             *ssm.SSM
	ID              *string
	defaultBaseline *bool
}

func init() {
	register("SSMPatchBaseline", ListSSMPatchBaselines)
}

func ListSSMPatchBaselines(sess *session.Session) ([]Resource, error) {
	svc := ssm.New(sess)
	resources := []Resource{}

	patchBaselineFilter := []*ssm.PatchOrchestratorFilter{
		{
			Key:    aws.String("OWNER"),
			Values: []*string{aws.String("Self")},
		},
	}

	params := &ssm.DescribePatchBaselinesInput{
		MaxResults: aws.Int64(50),
		Filters:    patchBaselineFilter,
	}

	for {
		output, err := svc.DescribePatchBaselines(params)
		if err != nil {
			return nil, err
		}

		for _, baselineIdentity := range output.BaselineIdentities {
			resources = append(resources, &SSMPatchBaseline{
				svc:             svc,
				ID:              baselineIdentity.BaselineId,
				defaultBaseline: baselineIdentity.DefaultBaseline,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SSMPatchBaseline) Remove() error {

	_, err := f.svc.DeletePatchBaseline(&ssm.DeletePatchBaselineInput{
		BaselineId: f.ID,
	})

	return err
}

func (f *SSMPatchBaseline) String() string {
	return *f.ID
}

func (f *SSMPatchBaseline) Filter() error {
	if *f.defaultBaseline {
		return fmt.Errorf("cannot delete default patch baseline, reassign default first")
	}
	return nil
}
