package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/bedrockagent"
)

type BedrockAgentAlias struct {
	svc            *bedrockagent.BedrockAgent
	AgentId        *string
	AgentAliasId   *string
	AgentAliasName *string
}

func init() {
	register("BedrockAgentAlias", ListBedrockAgentAliases)
}

func ListBedrockAgentAliases(sess *session.Session) ([]Resource, error) {
	svc := bedrockagent.New(sess)
	resources := []Resource{}

	agentIds, err := ListBedrockAgentIds(svc)
	if err != nil {
		return nil, err
	}

	for _, agentId := range agentIds {
		params := &bedrockagent.ListAgentAliasesInput{
			MaxResults: aws.Int64(100),
			AgentId:    aws.String(agentId),
		}
		for {
			output, err := svc.ListAgentAliases(params)
			if err != nil {
				return nil, err
			}

			for _, agentAliasInfo := range output.AgentAliasSummaries {
				resources = append(resources, &BedrockAgentAlias{
					svc:            svc,
					AgentId:        aws.String(agentId),
					AgentAliasName: agentAliasInfo.AgentAliasName,
					AgentAliasId:   agentAliasInfo.AgentAliasId,
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

func ListBedrockAgentIds(svc *bedrockagent.BedrockAgent) ([]string, error) {

	agentIds := []string{}
	params := &bedrockagent.ListAgentsInput{
		MaxResults: aws.Int64(100),
	}
	for {
		output, err := svc.ListAgents(params)
		if err != nil {
			return nil, err
		}

		for _, agent := range output.AgentSummaries {
			agentIds = append(agentIds, *agent.AgentId)
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return agentIds, nil
}

func (f *BedrockAgentAlias) Remove() error {
	_, err := f.svc.DeleteAgentAlias(&bedrockagent.DeleteAgentAliasInput{
		AgentAliasId: f.AgentAliasId,
		AgentId:      f.AgentId,
	})
	return err
}

func (f *BedrockAgentAlias) String() string {
	return *f.AgentAliasName
}
