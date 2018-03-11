package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
)

type ServiceDiscoveryNamespace struct {
	svc *servicediscovery.ServiceDiscovery
	ID  *string
}

func init() {
	register("ServiceDiscoveryNamespace", ListServiceDiscoveryNamespaces)
}

func ListServiceDiscoveryNamespaces(sess *session.Session) ([]Resource, error) {
	svc := servicediscovery.New(sess)
	resources := []Resource{}

	params := &servicediscovery.ListNamespacesInput{
		MaxResults: aws.Int64(100),
	}

	// Collect all services, using separate for loop
	// due to multi-service pagination issues
	for {
		output, err := svc.ListNamespaces(params)
		if err != nil {
			return nil, err
		}

		for _, namespace := range output.Namespaces {
			resources = append(resources, &ServiceDiscoveryNamespace{
				svc: svc,
				ID:  namespace.Id,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *ServiceDiscoveryNamespace) Remove() error {

	_, err := f.svc.DeleteNamespace(&servicediscovery.DeleteNamespaceInput{
		Id: f.ID,
	})

	return err
}

func (f *ServiceDiscoveryNamespace) String() string {
	return *f.ID
}
