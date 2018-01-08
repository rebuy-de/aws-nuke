package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Object struct {
	svc       *s3.S3
	bucket    string
	key       string
	versionID *string
}

func (n *S3Nuke) ListObjects() ([]Resource, error) {
	resources := make([]Resource, 0)

	buckets, err := n.DescribeBuckets()
	if err != nil {
		return nil, err
	}

	for _, name := range buckets {
		params := &s3.ListObjectVersionsInput{
			Bucket: &name,
		}

		for {
			resp, err := n.Service.ListObjectVersions(params)
			if err != nil {
				return nil, err
			}

			for _, out := range resp.Versions {
				resources = append(resources, &S3Object{
					svc:       n.Service,
					bucket:    name,
					key:       *out.Key,
					versionID: out.VersionId,
				})
			}

			for _, out := range resp.DeleteMarkers {
				resources = append(resources, &S3Object{
					svc:       n.Service,
					bucket:    name,
					key:       *out.Key,
					versionID: out.VersionId,
				})
			}

			// make sure to list all with more than 1000 objects
			if *resp.IsTruncated {
				params.KeyMarker = resp.NextKeyMarker
				continue
			}

			break
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
	if e.versionID != nil && *e.versionID != "null" {
		return fmt.Sprintf("s3://%s/%s#%s", e.bucket, e.key, *e.versionID)
	}
	return fmt.Sprintf("s3://%s/%s", e.bucket, e.key)
}
