package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kendra"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type KendraIndex struct {
	svc  *kendra.Kendra
	name string
	id   string
}

func init() {
	register("KendraIndex", ListKendraIndexes)
}

func ListKendraIndexes(sess *session.Session) ([]Resource, error) {
	svc := kendra.New(sess)
	resources := []Resource{}

	params := &kendra.ListIndicesInput{
		MaxResults: aws.Int64(100),
	}

	for {
		resp, err := svc.ListIndices(params)
		if err != nil {
			return nil, err
		}
		for _, index := range resp.IndexConfigurationSummaryItems {
			resources = append(resources, &KendraIndex{
				svc:  svc,
				id:   *index.Id,
				name: *index.Name,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

func (i *KendraIndex) Remove() error {
	_, err := i.svc.DeleteIndex(&kendra.DeleteIndexInput{
		Id: &i.id,
	})
	return err
}

func (i *KendraIndex) String() string {
	return i.id
}

func (i *KendraIndex) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", i.name)

	return properties
}
