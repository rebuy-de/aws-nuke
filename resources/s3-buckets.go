package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
			return nil, err
		}

		if EqualStringPtr(bucketLocationResponse.LocationConstraint, svc.Config.Region) {
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
	params := &s3.DeleteBucketInput{
		Bucket: &e.name,
	}

	_, err := e.svc.DeleteBucket(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *S3Bucket) String() string {
	return fmt.Sprintf("s3://%s", e.name)
}
