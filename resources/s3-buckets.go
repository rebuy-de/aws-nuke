package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func init() {
	register("S3Bucket", ListS3Buckets)
}

func ListS3Buckets(s *session.Session) ([]Resource, error) {
	svc := s3.New(s)

	buckets, err := DescribeS3Buckets(svc)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, name := range buckets {
		resources = append(resources, &S3Bucket{
			svc:  svc,
			name: name,
		})
	}

	return resources, nil
}

func DescribeS3Buckets(svc *s3.S3) ([]string, error) {
	resp, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	buckets := make([]string, 0)
	for _, out := range resp.Buckets {
		bucketLocationResponse, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: out.Name})

		if err != nil {
			continue
		}

		location := UnPtrString(bucketLocationResponse.LocationConstraint, endpoints.UsEast1RegionID)
		region := UnPtrString(svc.Config.Region, endpoints.UsEast1RegionID)
		if location == region {
			buckets = append(buckets, *out.Name)
		}

	}

	return buckets, nil
}

type S3Bucket struct {
	svc  *s3.S3
	name string
}

func (e *S3Bucket) Remove() error {
	_, err := e.svc.DeleteBucketPolicy(&s3.DeleteBucketPolicyInput{
		Bucket: &e.name,
	})

	if err != nil {
		return err
	}

	err = e.RemoveAllObjects()
	if err != nil {
		return err
	}

	params := &s3.DeleteBucketInput{
		Bucket: &e.name,
	}

	_, err = e.svc.DeleteBucket(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *S3Bucket) RemoveAllObjects() error {
	params := &s3.ListObjectsInput{
		Bucket: &e.name,
	}

	iterator := s3manager.NewDeleteListIterator(e.svc, params)

	err := s3manager.NewBatchDeleteWithClient(e.svc).Delete(aws.BackgroundContext(), iterator)
	if err != nil {
		return err
	}

	return nil
}

func (e *S3Bucket) String() string {
	return fmt.Sprintf("s3://%s", e.name)
}
