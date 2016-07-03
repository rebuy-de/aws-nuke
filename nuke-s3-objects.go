package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Object struct {
	svc    *s3.S3
	bucket string
	key    string
}

func (n *S3Nuke) ListObjects() ([]Resource, error) {
	resources := make([]Resource, 0)

	buckets, err := n.DescribeBuckets()
	if err != nil {
		return nil, err
	}

	for _, name := range buckets {
		params := &s3.ListObjectsInput{
			Bucket: &name,
		}

		resp, err := n.svc.ListObjects(params)
		if err != nil {
			return nil, err
		}

		for _, out := range resp.Contents {
			resources = append(resources, &S3Object{
				svc:    n.svc,
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
