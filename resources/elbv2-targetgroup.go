package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type ELBv2TargetGroup struct {
	svc  *elbv2.ELBV2
	name *string
	arn  *string
	tags []*elbv2.Tag
}

func init() {
	register("ELBv2TargetGroup", ListELBv2TargetGroups)
}

func ListELBv2TargetGroups(sess *session.Session) ([]Resource, error) {
	svc := elbv2.New(sess)
	var tagReqELBv2TargetGroupARNs []*string
	targetGroupArnToName := make(map[string]*string)

	err := svc.DescribeTargetGroupsPages(nil,
		func(page *elbv2.DescribeTargetGroupsOutput, lastPage bool) bool {
			for _, targetGroup := range page.TargetGroups {
				tagReqELBv2TargetGroupARNs = append(tagReqELBv2TargetGroupARNs, targetGroup.TargetGroupArn)
				targetGroupArnToName[*targetGroup.TargetGroupArn] = targetGroup.TargetGroupName
			}
			return !lastPage
		})
	if err != nil {
		return nil, err
	}

	// Tags for ELBv2 target groups need to be fetched separately
	// We can only specify up to 20 in a single call
	// See: https://github.com/aws/aws-sdk-go/blob/0e8c61841163762f870f6976775800ded4a789b0/service/elbv2/api.go#L5398
	resources := make([]Resource, 0)
	for len(tagReqELBv2TargetGroupARNs) > 0 {
		requestElements := len(tagReqELBv2TargetGroupARNs)
		if requestElements > 20 {
			requestElements = 20
		}

		tagResp, err := svc.DescribeTags(&elbv2.DescribeTagsInput{
			ResourceArns: tagReqELBv2TargetGroupARNs[:requestElements],
		})
		if err != nil {
			return nil, err
		}
		for _, tagInfo := range tagResp.TagDescriptions {
			resources = append(resources, &ELBv2TargetGroup{
				svc:  svc,
				name: targetGroupArnToName[*tagInfo.ResourceArn],
				arn:  tagInfo.ResourceArn,
				tags: tagInfo.Tags,
			})
		}

		// Remove the elements that were queried
		tagReqELBv2TargetGroupARNs = tagReqELBv2TargetGroupARNs[requestElements:]
	}
	return resources, nil
}

func (e *ELBv2TargetGroup) Remove() error {
	_, err := e.svc.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
		TargetGroupArn: e.arn,
	})

	if err != nil {
		return err
	}

	return nil
}

func (e *ELBv2TargetGroup) Properties() types.Properties {
	properties := types.NewProperties()
	for _, tagValue := range e.tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}
	return properties
}

func (e *ELBv2TargetGroup) String() string {
	return *e.name
}
