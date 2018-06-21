package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/neptune"
)

type NeptuneInstance struct {
	svc *neptune.Neptune
	ID  *string
}

func init() {
	register("NeptuneInstance", ListNeptuneInstances)
}

func ListNeptuneInstances(sess *session.Session) ([]Resource, error) {
	svc := neptune.New(sess)
	resources := []Resource{}

	params := &neptune.DescribeDBInstancesInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeDBInstances(params)
		if err != nil {
			return nil, err
		}

		for _, dbInstance := range output.DBInstances {
			resources = append(resources, &NeptuneInstance{
				svc: svc,
				ID:  dbInstance.DBInstanceIdentifier,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *NeptuneInstance) Remove() error {

	_, err := f.svc.DeleteDBInstance(&neptune.DeleteDBInstanceInput{
		DBInstanceIdentifier: f.ID,
		SkipFinalSnapshot:    aws.Bool(true),
	})

	return err
}

func (f *NeptuneInstance) String() string {
	return *f.ID
}
