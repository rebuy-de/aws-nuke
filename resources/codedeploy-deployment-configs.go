package resources

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeDeployDeploymentConfig struct {
	svc                  *codedeploy.CodeDeploy
	deploymentConfigName *string
}

func init() {
	register("CodeDeployDeploymentConfig", ListCodeDeployDeploymentConfigs, mapCloudControl("AWS::CodeDeploy::DeploymentConfig"))
}

func ListCodeDeployDeploymentConfigs(sess *session.Session) ([]Resource, error) {
	svc := codedeploy.New(sess)
	resources := []Resource{}

	params := &codedeploy.ListDeploymentConfigsInput{}

	for {
		resp, err := svc.ListDeploymentConfigs(params)
		if err != nil {
			return nil, err
		}

		for _, config := range resp.DeploymentConfigsList {
			resources = append(resources, &CodeDeployDeploymentConfig{
				svc:                  svc,
				deploymentConfigName: config,
			})
		}

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return resources, nil
}

func (f *CodeDeployDeploymentConfig) Filter() error {
	if strings.HasPrefix(*f.deploymentConfigName, "CodeDeployDefault") {
		return fmt.Errorf("cannot delete default codedeploy config")
	}
	return nil
}

func (f *CodeDeployDeploymentConfig) Remove() error {
	_, err := f.svc.DeleteDeploymentConfig(&codedeploy.DeleteDeploymentConfigInput{
		DeploymentConfigName: f.deploymentConfigName,
	})

	return err
}

func (f *CodeDeployDeploymentConfig) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("DeploymentConfigName", f.deploymentConfigName)
	return properties
}

func (f *CodeDeployDeploymentConfig) String() string {
	return *f.deploymentConfigName
}
