package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/simpledb"
)

type SimpleDBDomain struct {
	svc        *simpledb.SimpleDB
	domainName *string
}

func init() {
	register("SimpleDBDomain", ListSimpleDBDomains)
}

func ListSimpleDBDomains(sess *session.Session) ([]Resource, error) {
	svc := simpledb.New(sess)
	resources := []Resource{}

	params := &simpledb.ListDomainsInput{
		MaxNumberOfDomains: aws.Int64(100),
	}

	for {
		output, err := svc.ListDomains(params)
		if err != nil {
			return nil, err
		}

		for _, domainName := range output.DomainNames {
			resources = append(resources, &SimpleDBDomain{
				svc:        svc,
				domainName: domainName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SimpleDBDomain) Remove() error {

	_, err := f.svc.DeleteDomain(&simpledb.DeleteDomainInput{
		DomainName: f.domainName,
	})

	return err
}

func (f *SimpleDBDomain) String() string {
	return *f.domainName
}
