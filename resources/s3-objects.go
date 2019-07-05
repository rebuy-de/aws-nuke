package resources

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type S3Object struct {
	svc       *s3.S3
	bucket    string
	key       string
	versionID *string
	latest    bool
	tags      []*s3.Tag
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
		params := &s3.ListObjectVersionsInput{
			Bucket: &name,
		}

		for {
			resp, err := svc.ListObjectVersions(params)
			if err != nil {
				return nil, err
			}

			for _, out := range resp.Versions {
				if out.Key == nil {
					continue
				}
				tags, err := retrieveObjectTags(svc, name, *out.Key, *out.VersionId)

				if err != nil {
					continue
				}

				resources = append(resources, &S3Object{
					svc:       svc,
					bucket:    name,
					key:       *out.Key,
					versionID: out.VersionId,
					latest:    UnPtrBool(out.IsLatest, false),
					tags:      tags,
				})
			}

			for _, out := range resp.DeleteMarkers {
				if out.Key == nil {
					continue
				}

				resources = append(resources, &S3Object{
					svc:       svc,
					bucket:    name,
					key:       *out.Key,
					versionID: out.VersionId,
					latest:    UnPtrBool(out.IsLatest, false),
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
		Bucket:    &e.bucket,
		Key:       &e.key,
		VersionId: e.versionID,
	}

	_, err := e.svc.DeleteObject(params)
	if err != nil {
		return err
	}

	return nil
}

func retrieveObjectTags(svc *s3.S3, bucketName string, keyName string, versionId string) ([]*s3.Tag, error) {
	input := &s3.GetObjectTaggingInput{
		Bucket:    aws.String(bucketName),
		Key:       aws.String(keyName),
		VersionId: aws.String(versionId),
	}

	result, err := svc.GetObjectTagging(input)
	if err != nil {
		return make([]*s3.Tag, 0), err
	}

	return result.TagSet, nil
}

func (e *S3Object) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Bucket", e.bucket)
	properties.Set("Key", e.key)
	properties.Set("VersionID", e.versionID)
	properties.Set("IsLatest", e.latest)

	for _, tag := range e.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (e *S3Object) String() string {
	if e.versionID != nil && *e.versionID != "null" && !e.latest {
		return fmt.Sprintf("s3://%s/%s#%s", e.bucket, e.key, *e.versionID)
	}
	return fmt.Sprintf("s3://%s/%s", e.bucket, e.key)
}
