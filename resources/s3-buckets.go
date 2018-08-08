package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rebuy-de/aws-nuke/pkg/types"
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
		tags, _ := retrieveTags(svc, name)

		resources = append(resources, &S3Bucket{
			svc:  svc,
			name: name,
			tags: tags,
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
	tags []*s3.Tag
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

func retrieveTags(svc *s3.S3, bucketName string) ([]*s3.Tag, error) {
	input := &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketName),
	}

	result, err := svc.GetBucketTagging(input)
	if err != nil {
		return make([]*s3.Tag, 0), err
	}

	return result.TagSet, nil
}

func (e *S3Bucket) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Name", e.name)

	for _, tag := range e.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (e *S3Bucket) String() string {
	return fmt.Sprintf("s3://%s", e.name)
}
