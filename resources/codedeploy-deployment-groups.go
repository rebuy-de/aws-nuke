package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type CodeDeployDeploymentGroup struct {
	svc                 *codedeploy.CodeDeploy
	deploymentGroupName *string
	applicationName     *string
}

func init() {
	register("CodeDeployDeploymentGroup", ListCodeDeployDeploymentGroups)
}

func ListCodeDeployDeploymentGroups(sess *session.Session) ([]Resource, error) {
	svc := codedeploy.New(sess)
	resources := []Resource{}

	appParams := &codedeploy.ListApplicationsInput{}
	appResp, err := svc.ListApplications(appParams)
	if err != nil {
		return nil, err
	}

	for _, appName := range appResp.Applications {
		// For each application, list deployment groups
		deploymentGroupParams := &codedeploy.ListDeploymentGroupsInput{
			ApplicationName: appName,
		}
		deploymentGroupResp, err := svc.ListDeploymentGroups(deploymentGroupParams)
		if err != nil {
			return nil, err
		}

		for _, group := range deploymentGroupResp.DeploymentGroups {
			resources = append(resources, &CodeDeployDeploymentGroup{
				svc:                 svc,
				deploymentGroupName: group,
				applicationName:     appName,
			})
		}
	}

	return resources, nil
}

func (f *CodeDeployDeploymentGroup) Remove() error {
	_, err := f.svc.DeleteDeploymentGroup(&codedeploy.DeleteDeploymentGroupInput{
		ApplicationName:     f.applicationName,
		DeploymentGroupName: f.deploymentGroupName,
	})

	return err
}

func (f *CodeDeployDeploymentGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("DeploymentGroupName", f.deploymentGroupName)
	properties.Set("ApplicationName", f.applicationName)
	return properties
}

func (f *CodeDeployDeploymentGroup) String() string {
	return *f.deploymentGroupName
}
