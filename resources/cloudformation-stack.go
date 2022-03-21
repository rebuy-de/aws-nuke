package resources

import (
	"errors"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

const CLOUDFORMATION_MAX_DELETE_ATTEMPT = 3

func init() {
	register("CloudFormationStack", ListCloudFormationStacks)
}

func ListCloudFormationStacks(sess *session.Session) ([]Resource, error) {
	svc := cloudformation.New(sess)

	params := &cloudformation.DescribeStacksInput{}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.DescribeStacks(params)
		if err != nil {
			return nil, err
		}
		for _, stack := range resp.Stacks {
			if aws.StringValue(stack.ParentId) != "" {
				continue
			}
			resources = append(resources, &CloudFormationStack{
				svc:               svc,
				stack:             stack,
				maxDeleteAttempts: CLOUDFORMATION_MAX_DELETE_ATTEMPT,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type CloudFormationStack struct {
	svc               cloudformationiface.CloudFormationAPI
	stack             *cloudformation.Stack
	maxDeleteAttempts int
	featureFlags      config.FeatureFlags
}

func (cfs *CloudFormationStack) FeatureFlags(ff config.FeatureFlags) {
	cfs.featureFlags = ff
}

func (cfs *CloudFormationStack) Remove() error {
	return cfs.removeWithAttempts(0)
}

func (cfs *CloudFormationStack) removeWithAttempts(attempt int) error {
	if err := cfs.doRemove(); err != nil {
		logrus.Errorf("CloudFormationStack stackName=%s attempt=%d maxAttempts=%d delete failed: %s", *cfs.stack.StackName, attempt, cfs.maxDeleteAttempts, err.Error())
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "ValidationError" &&
			awsErr.Message() == "Stack ["+*cfs.stack.StackName+"] cannot be deleted while TerminationProtection is enabled" {
			if cfs.featureFlags.DisableDeletionProtection.CloudformationStack {
				logrus.Infof("CloudFormationStack stackName=%s attempt=%d maxAttempts=%d updating termination protection", *cfs.stack.StackName, attempt, cfs.maxDeleteAttempts)
				_, err = cfs.svc.UpdateTerminationProtection(&cloudformation.UpdateTerminationProtectionInput{
					EnableTerminationProtection: aws.Bool(false),
					StackName:                   cfs.stack.StackName,
				})
				if err != nil {
					return err
				}
			} else {
				logrus.Warnf("CloudFormationStack stackName=%s attempt=%d maxAttempts=%d set feature flag to disable deletion protection", *cfs.stack.StackName, attempt, cfs.maxDeleteAttempts)
				return err
			}
		}
		if attempt >= cfs.maxDeleteAttempts {
			return errors.New("CFS might not be deleted after this run.")
		} else {
			return cfs.removeWithAttempts(attempt + 1)
		}
	} else {
		return nil
	}
}

func (cfs *CloudFormationStack) doRemove() error {
	o, err := cfs.svc.DescribeStacks(&cloudformation.DescribeStacksInput{
		StackName: cfs.stack.StackName,
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ValidationFailed" && strings.HasSuffix(awsErr.Message(), " does not exist") {
				logrus.Infof("CloudFormationStack stackName=%s no longer exists", *cfs.stack.StackName)
				return nil
			}
		}
		return err
	}
	stack := o.Stacks[0]

	if *stack.StackStatus == cloudformation.StackStatusDeleteComplete {
		//stack already deleted, no need to re-delete
		return nil
	} else if *stack.StackStatus == cloudformation.StackStatusDeleteInProgress {
		logrus.Infof("CloudFormationStack stackName=%s delete in progress. Waiting", *cfs.stack.StackName)
		return cfs.svc.WaitUntilStackDeleteComplete(&cloudformation.DescribeStacksInput{
			StackName: cfs.stack.StackName,
		})
	} else if *stack.StackStatus == cloudformation.StackStatusDeleteFailed {
		logrus.Infof("CloudFormationStack stackName=%s delete failed. Attempting to retain and delete stack", *cfs.stack.StackName)
		// This means the CFS has undeleteable resources.
		// In order to move on with nuking, we retain them in the deletion.
		retainableResources, err := cfs.svc.ListStackResources(&cloudformation.ListStackResourcesInput{
			StackName: cfs.stack.StackName,
		})
		if err != nil {
			return err
		}

		retain := make([]*string, 0)

		for _, r := range retainableResources.StackResourceSummaries {
			if *r.ResourceStatus != cloudformation.ResourceStatusDeleteComplete {
				retain = append(retain, r.LogicalResourceId)
			}
		}

		_, err = cfs.svc.DeleteStack(&cloudformation.DeleteStackInput{
			StackName:       cfs.stack.StackName,
			RetainResources: retain,
		})
		if err != nil {
			return err
		}
		return cfs.svc.WaitUntilStackDeleteComplete(&cloudformation.DescribeStacksInput{
			StackName: cfs.stack.StackName,
		})
	} else {
		if err := cfs.waitForStackToStabilize(*stack.StackStatus); err != nil {
			return err
		} else if _, err := cfs.svc.DeleteStack(&cloudformation.DeleteStackInput{
			StackName: cfs.stack.StackName,
		}); err != nil {
			return err
		} else if err := cfs.svc.WaitUntilStackDeleteComplete(&cloudformation.DescribeStacksInput{
			StackName: cfs.stack.StackName,
		}); err != nil {
			return err
		} else {
			return nil
		}
	}
}
func (cfs *CloudFormationStack) waitForStackToStabilize(currentStatus string) error {
	switch currentStatus {
	case cloudformation.StackStatusUpdateInProgress:
		fallthrough
	case cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress:
		fallthrough
	case cloudformation.StackStatusUpdateRollbackInProgress:
		logrus.Infof("CloudFormationStack stackName=%s update in progress. Waiting to stabalize", *cfs.stack.StackName)
		return cfs.svc.WaitUntilStackUpdateComplete(&cloudformation.DescribeStacksInput{
			StackName: cfs.stack.StackName,
		})
	case cloudformation.StackStatusCreateInProgress:
		fallthrough
	case cloudformation.StackStatusRollbackInProgress:
		logrus.Infof("CloudFormationStack stackName=%s create in progress. Waiting to stabalize", *cfs.stack.StackName)
		return cfs.svc.WaitUntilStackCreateComplete(&cloudformation.DescribeStacksInput{
			StackName: cfs.stack.StackName,
		})
	default:
		return nil
	}
}

func (cfs *CloudFormationStack) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", cfs.stack.StackName)
	properties.Set("CreationTime", cfs.stack.CreationTime.Format(time.RFC3339))
	if cfs.stack.LastUpdatedTime == nil {
		properties.Set("LastUpdatedTime", cfs.stack.CreationTime.Format(time.RFC3339))
	} else {
		properties.Set("LastUpdatedTime", cfs.stack.LastUpdatedTime.Format(time.RFC3339))
	}

	for _, tagValue := range cfs.stack.Tags {
		properties.SetTag(tagValue.Key, tagValue.Value)
	}

	return properties
}

func (cfs *CloudFormationStack) String() string {
	return *cfs.stack.StackName
}
