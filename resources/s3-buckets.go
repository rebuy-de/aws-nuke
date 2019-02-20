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

type S3Bucket struct {
	svc  *s3.S3
	name string
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

func (e *S3Bucket) Remove() error {
	_, err := e.svc.DeleteBucketPolicy(&s3.DeleteBucketPolicyInput{
		Bucket: &e.name,
	})
	if err != nil {
		return err
	}

	_, err = e.svc.PutBucketLogging(&s3.PutBucketLoggingInput{
		Bucket:              &e.name,
		BucketLoggingStatus: &s3.BucketLoggingStatus{},
	})
	if err != nil {
		return err
	}

	err = e.RemoveAllVersions()
	if err != nil {
		return err
	}

	err = e.RemoveAllObjects()
	if err != nil {
		return err
	}

	_, err = e.svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: &e.name,
	})

	return err
}

func (e *S3Bucket) RemoveAllVersions() error {
	params := &s3.ListObjectVersionsInput{
		Bucket: &e.name,
	}

	iterator := NewDeleteVersionListIterator(e.svc, params)
	return s3manager.NewBatchDeleteWithClient(e.svc).Delete(aws.BackgroundContext(), iterator)
}

func (e *S3Bucket) RemoveAllObjects() error {
	params := &s3.ListObjectsInput{
		Bucket: &e.name,
	}

	iterator := s3manager.NewDeleteListIterator(e.svc, params)
	return s3manager.NewBatchDeleteWithClient(e.svc).Delete(aws.BackgroundContext(), iterator)
}

func (e *S3Bucket) String() string {
	return fmt.Sprintf("s3://%s", e.name)
}
