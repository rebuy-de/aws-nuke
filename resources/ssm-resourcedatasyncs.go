package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMResourceDataSync struct {
	svc  *ssm.SSM
	name *string
}

func init() {
	register("SSMResourceDataSync", ListSSMResourceDataSyncs)
}

func ListSSMResourceDataSyncs(sess *session.Session) ([]Resource, error) {
	svc := ssm.New(sess)
	resources := []Resource{}

	params := &ssm.ListResourceDataSyncInput{
		MaxResults: aws.Int64(50),
	}

	for {
		output, err := svc.ListResourceDataSync(params)
		if err != nil {
			return nil, err
		}

		for _, resourceDataSyncItem := range output.ResourceDataSyncItems {
			resources = append(resources, &SSMResourceDataSync{
				svc:  svc,
				name: resourceDataSyncItem.SyncName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SSMResourceDataSync) Remove() error {

	_, err := f.svc.DeleteResourceDataSync(&ssm.DeleteResourceDataSyncInput{
		SyncName: f.name,
	})

	return err
}

func (f *SSMResourceDataSync) String() string {
	return *f.name
}
