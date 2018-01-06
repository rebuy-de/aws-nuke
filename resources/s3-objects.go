package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Object struct {
	svc    *s3.S3
	bucket string
	key    string
}

func init() {
	register("S3Object", ListS3Objects)
}

func ListS3Objects(sess *session.Session) ([]Resource, error) {
	svc := s3.New(sess)

	resources := make([]Resource, 0)

	buckets, err := DescribeS3Buckets(svc)
	if err != nil {
		return nil, err
	}

	for _, name := range buckets {
		params := &s3.ListObjectsInput{
			Bucket: &name,
		}

		resp, err := svc.ListObjects(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.Contents {
			resources = append(resources, &S3Object{
				svc:    svc,
				bucket: name,
				key:    *out.Key,
			})
		}
	}

	return resources, nil
}

func (e *S3Object) Remove() error {
	params := &s3.DeleteObjectInput{
		Bucket: &e.bucket,
		Key:    &e.key,
	}

	_, err := e.svc.DeleteObject(params)
	if err != nil {
		return err
	}

	return nil
}

func (e *S3Object) String() string {
	return fmt.Sprintf("s3://%s/%s", e.bucket, e.key)
}
