package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/bedrockagent"
)

type BedrockFlowAlias struct {
	svc           *bedrockagent.BedrockAgent
	FlowId        *string
	FlowAliasId   *string
	FlowAliasName *string
}

func init() {
	register("BedrockFlowAlias", ListBedrockFlowAliases)
}

func ListBedrockFlowAliases(sess *session.Session) ([]Resource, error) {
	svc := bedrockagent.New(sess)
	resources := []Resource{}

	flowIds, err := ListBedrockFlowIds(svc)
	if err != nil {
		return nil, err
	}

	for _, flowId := range flowIds {
		params := &bedrockagent.ListFlowAliasesInput{
			MaxResults:     aws.Int64(100),
			FlowIdentifier: aws.String(flowId),
		}
		for {
			output, err := svc.ListFlowAliases(params)
			if err != nil {
				return nil, err
			}

			for _, flowAliasInfo := range output.FlowAliasSummaries {
				resources = append(resources, &BedrockFlowAlias{
					svc:           svc,
					FlowId:        flowAliasInfo.FlowId,
					FlowAliasId:   flowAliasInfo.Id,
					FlowAliasName: flowAliasInfo.Name,
				})
			}

			if output.NextToken == nil {
				break
			}
			params.NextToken = output.NextToken
		}

	}

	return resources, nil
}

func ListBedrockFlowIds(svc *bedrockagent.BedrockAgent) ([]string, error) {

	flowIds := []string{}
	params := &bedrockagent.ListFlowsInput{
		MaxResults: aws.Int64(100),
	}
	for {
		output, err := svc.ListFlows(params)
		if err != nil {
			return nil, err
		}

		for _, flow := range output.FlowSummaries {
			flowIds = append(flowIds, *flow.Id)
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return flowIds, nil
}

func (f *BedrockFlowAlias) Remove() error {
	_, err := f.svc.DeleteFlowAlias(&bedrockagent.DeleteFlowAliasInput{
		AliasIdentifier: f.FlowAliasId,
		FlowIdentifier:  f.FlowId,
	})
	return err
}

func (f *BedrockFlowAlias) String() string {
	return *f.FlowAliasName
}
