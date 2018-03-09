package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
)

type ServiceDiscoveryService struct {
	svc *servicediscovery.ServiceDiscovery
	ID  *string
}

func init() {
	register("ServiceDiscoveryService", ListServiceDiscoveryServices)
}

func ListServiceDiscoveryServices(sess *session.Session) ([]Resource, error) {
	svc := servicediscovery.New(sess)
	resources := []Resource{}

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
			resources = append(resources, &ServiceDiscoveryService{
				svc: svc,
				ID:  service.Id,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ServiceDiscoveryService) Remove() error {

	_, err := f.svc.DeleteService(&servicediscovery.DeleteServiceInput{
		Id: f.ID,
	})

	return err
}

func (f *ServiceDiscoveryService) String() string {
	return *f.ID
}
