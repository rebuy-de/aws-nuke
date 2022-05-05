package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codestarnotifications"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeStarNotificationRule struct {
	svc  *codestarnotifications.CodeStarNotifications
	id   *string
	name *string
	arn  *string
	tags map[string]*string
}

func init() {
	register("CodeStarNotificationRule", ListCodeStarNotificationRules)
}

func ListCodeStarNotificationRules(sess *session.Session) ([]Resource, error) {
	svc := codestarnotifications.New(sess)
	resources := []Resource{}

	params := &codestarnotifications.ListNotificationRulesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListNotificationRules(params)
		if err != nil {
			return nil, err
		}

		for _, notification := range output.NotificationRules {
			descOutput, err := svc.DescribeNotificationRule(&codestarnotifications.DescribeNotificationRuleInput{
				Arn: notification.Arn,
			})
			if err != nil {
				return nil, err
			}

			resources = append(resources, &CodeStarNotificationRule{
				svc:  svc,
				id:   notification.Id,
				name: descOutput.Name,
				arn:  notification.Arn,
				tags: descOutput.Tags,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (cn *CodeStarNotificationRule) Remove() error {

	_, err := cn.svc.DeleteNotificationRule(&codestarnotifications.DeleteNotificationRuleInput{
		Arn: cn.arn,
	})

	return err
}

func (cn *CodeStarNotificationRule) String() string {
	return fmt.Sprintf("%s (%s)", *cn.id, *cn.name)
}

func (cn *CodeStarNotificationRule) Properties() types.Properties {
	properties := types.NewProperties()
	for key, tag := range cn.tags {
		properties.SetTag(&key, tag)
	}
	properties.
		Set("Name", cn.name).
		Set("ID", cn.id)
	return properties
}
