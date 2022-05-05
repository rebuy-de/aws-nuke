package resources

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

func init() {
	register("AthenaWorkGroup", ListAthenaWorkGroups,
		mapCloudControl("AWS::Athena::WorkGroup"))
}

type AthenaWorkGroup struct {
	svc  *athena.Athena
	name *string
	arn  *string
}

func ListAthenaWorkGroups(sess *session.Session) ([]Resource, error) {
	svc := athena.New(sess)
	resources := []Resource{}

	// Lookup current account ID
	stsSvc := sts.New(sess)
	callerID, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	accountID := callerID.Account
	region := svc.Config.Region

	// List WorkGroup
	var workgroupNames []*string
	err = svc.ListWorkGroupsPages(
		&athena.ListWorkGroupsInput{},
		func(page *athena.ListWorkGroupsOutput, lastPage bool) bool {
			for _, workgroup := range page.WorkGroups {
				workgroupNames = append(workgroupNames, workgroup.Name)
			}
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	// Create AthenaWorkGroup resource objects
	for _, name := range workgroupNames {
		resources = append(resources, &AthenaWorkGroup{
			svc:  svc,
			name: name,
			// The GetWorkGroup API doesn't return an ARN,
			// so we need to construct one ourselves
			arn: aws.String(fmt.Sprintf(
				"arn:aws:athena:%s:%s:workgroup/%s",
				*region, *accountID, *name,
			)),
		})
	}

	return resources, err
}

func (a *AthenaWorkGroup) Remove() error {
	// Primary WorkGroup cannot be deleted,
	// but we can reset it to a clean state
	if *a.name == "primary" {
		logrus.Info("Primary Athena WorkGroup may not be deleted. Resetting configuration only.")

		// Reset the configuration to its default state
		_, err := a.svc.UpdateWorkGroup(&athena.UpdateWorkGroupInput{
			// See https://docs.aws.amazon.com/athena/latest/APIReference/API_WorkGroupConfigurationUpdates.html
			// for documented defaults
			ConfigurationUpdates: &athena.WorkGroupConfigurationUpdates{
				EnforceWorkGroupConfiguration:    aws.Bool(false),
				PublishCloudWatchMetricsEnabled:  aws.Bool(false),
				RemoveBytesScannedCutoffPerQuery: aws.Bool(true),
				RequesterPaysEnabled:             aws.Bool(false),
				ResultConfigurationUpdates: &athena.ResultConfigurationUpdates{
					RemoveEncryptionConfiguration: aws.Bool(true),
					RemoveOutputLocation:          aws.Bool(true),
				},
			},
			Description: aws.String(""),
			WorkGroup:   a.name,
		})

		// Remove any tags
		wgTagsRes, err := a.svc.ListTagsForResource(&athena.ListTagsForResourceInput{
			ResourceARN: a.arn,
		})
		if err != nil {
			return err
		}
		var tagKeys []*string
		for _, tag := range wgTagsRes.Tags {
			tagKeys = append(tagKeys, tag.Key)
		}
		_, err = a.svc.UntagResource(&athena.UntagResourceInput{
			ResourceARN: a.arn,
			TagKeys:     tagKeys,
		})
		if err != nil {
			return err
		}

		return nil
	}

	_, err := a.svc.DeleteWorkGroup(&athena.DeleteWorkGroupInput{
		RecursiveDeleteOption: aws.Bool(true),
		WorkGroup:             a.name,
	})

	return err
}

func (a *AthenaWorkGroup) Filter() error {
	// If this is the primary work group,
	// check if it's already had its configuration reset
	if *a.name == "primary" {
		// Get workgroup configuration
		wgConfigRes, err := a.svc.GetWorkGroup(&athena.GetWorkGroupInput{
			WorkGroup: a.name,
		})
		if err != nil {
			return err
		}

		// Get workgroup tags
		wgTagsRes, err := a.svc.ListTagsForResource(&athena.ListTagsForResourceInput{
			ResourceARN: a.arn,
		})
		if err != nil {
			return err
		}

		// If the workgroup is already in a "clean" state, then
		// don't add it to our plan
		wgConfig := wgConfigRes.WorkGroup.Configuration
		isCleanConfig := wgConfig.BytesScannedCutoffPerQuery == nil &&
			*wgConfig.EnforceWorkGroupConfiguration == false &&
			*wgConfig.PublishCloudWatchMetricsEnabled == false &&
			*wgConfig.RequesterPaysEnabled == false &&
			*wgConfig.ResultConfiguration == athena.ResultConfiguration{} &&
			len(wgTagsRes.Tags) == 0

		if isCleanConfig {
			return errors.New("cannot delete primary athena work group")
		}
	}
	return nil
}

func (a *AthenaWorkGroup) Properties() types.Properties {
	return types.NewProperties().
		Set("Name", *a.name).
		Set("ARN", *a.arn)
}

func (a *AthenaWorkGroup) String() string {
	return *a.name
}
