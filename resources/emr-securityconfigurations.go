package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/emr"
)

type EMRSecurityConfiguration struct {
	svc  *emr.EMR
	name *string
}

func init() {
	register("EMRSecurityConfiguration", ListEMRSecurityConfiguration)
}

func ListEMRSecurityConfiguration(sess *session.Session) ([]Resource, error) {
	svc := emr.New(sess)
	resources := []Resource{}

	params := &emr.ListSecurityConfigurationsInput{}

	for {
		resp, err := svc.ListSecurityConfigurations(params)
		if err != nil {
			return nil, err
		}

		for _, securityConfiguration := range resp.SecurityConfigurations {
			resources = append(resources, &EMRSecurityConfiguration{
				svc:  svc,
				name: securityConfiguration.Name,
			})
		}

		if resp.Marker == nil {
			break
		}

		params.Marker = resp.Marker
	}

	return resources, nil
}

func (f *EMRSecurityConfiguration) Remove() error {

	//Call names are inconsistent in the SDK
	_, err := f.svc.DeleteSecurityConfiguration(&emr.DeleteSecurityConfigurationInput{
		Name: f.name,
	})

	return err
}

func (f *EMRSecurityConfiguration) String() string {
	return *f.name
}
