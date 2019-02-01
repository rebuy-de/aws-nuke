package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMPatchBaselineAssociation struct {
	svc           *ssm.SSM
	patchBaseline *string
	patchGroup    *string
}

func init() {
	register("SSMPatchBaselineAssociation", ListSSMPatchBaselineAssociations)
}

func ListSSMPatchGroups(svc *ssm.SSM) ([]*string, error) {
	patchGroups := make([]*string, 0)
	params := &ssm.DescribePatchGroupsInput{
		MaxResults: aws.Int64(50),
	}
	for {
		resp, err := svc.DescribePatchGroups(params)
		if err != nil {
			return nil, err
		}

		for _, mapping := range resp.Mappings {
			patchGroups = append(patchGroups, mapping.PatchGroup)
		}
		if params.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return patchGroups, nil
}

func ListSSMPatchBaselineAssociations(sess *session.Session) ([]Resource, error) {
	svc := ssm.New(sess)

	resources := make([]Resource, 0)

	patchGroups, err := ListSSMPatchGroups(svc)
	if err != nil {
		return nil, err
	}
	for _, patchGroup := range patchGroups {
		resp, err := svc.GetPatchBaselineForPatchGroup(&ssm.GetPatchBaselineForPatchGroupInput{
			PatchGroup: patchGroup,
		})
		if err != nil {
			return nil, err
		}
		resources = append(resources, &SSMPatchBaselineAssociation{
			svc:           svc,
			patchBaseline: resp.BaselineId,
			patchGroup:    patchGroup,
		})
	}

	return resources, nil
}

func (e *SSMPatchBaselineAssociation) Remove() error {
	_, err := e.svc.DeregisterPatchBaselineForPatchGroup(&ssm.DeregisterPatchBaselineForPatchGroupInput{
		BaselineId: e.patchBaseline,
		PatchGroup: e.patchGroup,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *SSMPatchBaselineAssociation) String() string {
	return fmt.Sprintf("%s -> %s", *e.patchBaseline, *e.patchGroup)
}
