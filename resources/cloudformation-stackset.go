package resources

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
	"github.com/sirupsen/logrus"
)

func init() {
	register("CloudFormationStackSet", ListCloudFormationStackSets,
		mapCloudControl("AWS::CloudFormation::StackSet"))
}

func ListCloudFormationStackSets(sess *session.Session) ([]Resource, error) {
	svc := cloudformation.New(sess)

	params := &cloudformation.ListStackSetsInput{
		Status: aws.String(cloudformation.StackSetStatusActive),
	}
	resources := make([]Resource, 0)

	for {
		resp, err := svc.ListStackSets(params)
		if err != nil {
			return nil, err
		}
		for _, stackSetSummary := range resp.Summaries {
			resources = append(resources, &CloudFormationStackSet{
				svc:             svc,
				stackSetSummary: stackSetSummary,
				sleepDuration:   10 * time.Second,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type CloudFormationStackSet struct {
	svc             cloudformationiface.CloudFormationAPI
	stackSetSummary *cloudformation.StackSetSummary
	sleepDuration   time.Duration
}

func (cfs *CloudFormationStackSet) findStackInstances() (map[string][]string, error) {
	accounts := make(map[string][]string)

	input := &cloudformation.ListStackInstancesInput{
		StackSetName: cfs.stackSetSummary.StackSetName,
	}

	for {
		resp, err := cfs.svc.ListStackInstances(input)
		if err != nil {
			return nil, err
		}
		for _, stackInstanceSummary := range resp.Summaries {
			if regions, ok := accounts[*stackInstanceSummary.Account]; !ok {
				accounts[*stackInstanceSummary.Account] = []string{*stackInstanceSummary.Region}
			} else {
				accounts[*stackInstanceSummary.Account] = append(regions, *stackInstanceSummary.Region)
			}
		}

		if resp.NextToken == nil {
			break
		}

		input.NextToken = resp.NextToken
	}

	return accounts, nil
}

func (cfs *CloudFormationStackSet) waitForStackSetOperation(operationId string) error {
	for {
		result, err := cfs.svc.DescribeStackSetOperation(&cloudformation.DescribeStackSetOperationInput{
			StackSetName: cfs.stackSetSummary.StackSetName,
			OperationId:  &operationId,
		})
		if err != nil {
			return err
		}
		logrus.Infof("Got stackInstance operation status on stackSet=%s operationId=%s status=%s", *cfs.stackSetSummary.StackSetName, operationId, *result.StackSetOperation.Status)
		if *result.StackSetOperation.Status == cloudformation.StackSetOperationResultStatusSucceeded {
			return nil
		} else if *result.StackSetOperation.Status == cloudformation.StackSetOperationResultStatusFailed || *result.StackSetOperation.Status == cloudformation.StackSetOperationResultStatusCancelled {
			return fmt.Errorf("unable to delete stackSet=%s operationId=%s status=%s", *cfs.stackSetSummary.StackSetName, operationId, *result.StackSetOperation.Status)
		} else {
			logrus.Infof("Waiting on stackSet=%s operationId=%s status=%s", *cfs.stackSetSummary.StackSetName, operationId, *result.StackSetOperation.Status)
			time.Sleep(cfs.sleepDuration)
		}
	}
}

func (cfs *CloudFormationStackSet) deleteStackInstances(accountId string, regions []string) error {
	logrus.Infof("Deleting stack instance accountId=%s regions=%s", accountId, strings.Join(regions, ","))
	regionsInput := make([]*string, len(regions))
	for i, region := range regions {
		regionsInput[i] = aws.String(region)
		fmt.Printf("region=%s i=%d\n", region, i)
	}
	result, err := cfs.svc.DeleteStackInstances(&cloudformation.DeleteStackInstancesInput{
		StackSetName: cfs.stackSetSummary.StackSetName,
		Accounts:     []*string{&accountId},
		Regions:      regionsInput,
		RetainStacks: aws.Bool(true), //this will remove the stack set instance from the stackset, but will leave the stack in the account/region it was deployed to
	})

	fmt.Printf("got result=%v err=%v\n", result, err)

	if result == nil {
		return fmt.Errorf("got null result")
	}
	if err != nil {
		return err
	}

	return cfs.waitForStackSetOperation(*result.OperationId)
}

func (cfs *CloudFormationStackSet) Remove() error {
	accounts, err := cfs.findStackInstances()
	if err != nil {
		return err
	}
	for accountId, regions := range accounts {
		err := cfs.deleteStackInstances(accountId, regions)
		if err != nil {
			return err
		}
	}
	_, err = cfs.svc.DeleteStackSet(&cloudformation.DeleteStackSetInput{
		StackSetName: cfs.stackSetSummary.StackSetName,
	})
	return err
}

func (cfs *CloudFormationStackSet) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", cfs.stackSetSummary.StackSetName)
	properties.Set("StackSetId", cfs.stackSetSummary.StackSetId)

	return properties
}

func (cfs *CloudFormationStackSet) String() string {
	return *cfs.stackSetSummary.StackSetName
}
