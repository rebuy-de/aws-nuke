package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
)

type ServiceDiscoveryInstance struct {
	svc        *servicediscovery.ServiceDiscovery
	serviceID  *string
	instanceID *string
}

func init() {
	register("ServiceDiscoveryInstance", ListServiceDiscoveryInstances)
}

func ListServiceDiscoveryInstances(sess *session.Session) ([]Resource, error) {
	svc := servicediscovery.New(sess)
	resources := []Resource{}
	services := []*servicediscovery.ServiceSummary{}

	params := &servicediscovery.ListServicesInput{
		MaxResults: aws.Int64(100),
	}

	// Collect all services, using separate for loop
	// due to multi-service pagination issues
	for {
		output, err := svc.ListServices(params)
		if err != nil {
			return nil, err
		}

		for _, service := range output.Services {
			services = append(services, service)
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	//Collect Instances for de-registration
	for _, service := range services {
		instanceParams := &servicediscovery.ListInstancesInput{
			ServiceId:  service.Id,
			MaxResults: aws.Int64(100),
		}

		output, err := svc.ListInstances(instanceParams)
		if err != nil {
			return nil, err
		}

		for _, instance := range output.Instances {
			resources = append(resources, &ServiceDiscoveryInstance{
				svc:        svc,
				serviceID:  service.Id,
				instanceID: instance.Id,
			})
		}
		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ServiceDiscoveryInstance) Remove() error {

	_, err := f.svc.DeregisterInstance(&servicediscovery.DeregisterInstanceInput{
		InstanceId: f.instanceID,
		ServiceId:  f.serviceID,
	})

	return err
}

func (f *ServiceDiscoveryInstance) String() string {
	return fmt.Sprintf("%s -> %s", *f.instanceID, *f.serviceID)
}
