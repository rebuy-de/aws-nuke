package resources

import (
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
)

type EC2LaunchTemplate struct {
    svc  *ec2.EC2
    name *string
}

func init() {
    register("EC2LaunchTemplate", ListEC2LaunchTemplates)
}

func ListEC2LaunchTemplates(sess *session.Session) ([]Resource, error) {
    svc := ec2.New(sess)

    resp, err := svc.DescribeLaunchTemplates(nil)
    if err != nil {
        return nil, err
    }

    resources := make([]Resource, 0)
    for _, template := range resp.LaunchTemplates {
        resources = append(resources, &EC2LaunchTemplate{
            svc:  svc,
            name: template.LaunchTemplateName,
        })
    }
    return resources, nil
}

func (template *EC2LaunchTemplate) Remove() error {
    _, err := template.svc.DeleteLaunchTemplate(&ec2.DeleteLaunchTemplateInput{
        LaunchTemplateName: template.name,
    })
    return err
}

func (template *EC2LaunchTemplate) String() string {
    return *template.name
}
