package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Bucket struct {
	svc  *s3.S3
	name string
}

func (n *S3Nuke) DescribeBuckets() ([]string, error) {
	resp, err := n.Service.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	buckets := make([]string, 0)
	for _, out := range resp.Buckets {
		buckets = append(buckets, *out.Name)
	}

	return buckets, nil
}

func (n *S3Nuke) ListBuckets() ([]Resource, error) {
	buckets, err := n.DescribeBuckets()
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, name := range buckets {
		resources = append(resources, &S3Bucket{
			svc:  n.Service,
			name: name,
		})
	}

	return resources, nil
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
