package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rebuy-de/aws-nuke/v2/pkg/config"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("S3Bucket", ListS3Buckets,
		mapCloudControl("AWS::S3::Bucket"))
}

type S3Bucket struct {
	svc          *s3.S3
	name         string
	creationDate time.Time
	tags         []*s3.Tag
	featureFlags config.FeatureFlags
}

func ListS3Buckets(s *session.Session) ([]Resource, error) {
	svc := s3.New(s)

	buckets, err := DescribeS3Buckets(svc)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, bucket := range buckets {
		tags, err := svc.GetBucketTagging(&s3.GetBucketTaggingInput{
			Bucket: bucket.Name,
		})

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() == "NoSuchTagSet" {
					resources = append(resources, &S3Bucket{
						svc:          svc,
						name:         aws.StringValue(bucket.Name),
						creationDate: aws.TimeValue(bucket.CreationDate),
						tags:         make([]*s3.Tag, 0),
					})
				}
			}
			continue
		}

		resources = append(resources, &S3Bucket{
			svc:          svc,
			name:         aws.StringValue(bucket.Name),
			creationDate: aws.TimeValue(bucket.CreationDate),
			tags:         tags.TagSet,
		})
	}

	return resources, nil
}

func DescribeS3Buckets(svc *s3.S3) ([]s3.Bucket, error) {
	resp, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	buckets := make([]s3.Bucket, 0)
	for _, out := range resp.Buckets {
		bucketLocationResponse, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: out.Name})

		if err != nil {
			continue
		}

		location := UnPtrString(bucketLocationResponse.LocationConstraint, endpoints.UsEast1RegionID)
		region := UnPtrString(svc.Config.Region, endpoints.UsEast1RegionID)
		if location == region && out != nil {
			buckets = append(buckets, *out)
		}
	}

	return buckets, nil
}

func (e *S3Bucket) Remove() error {
	objectsThreshold := e.featureFlags.PostponeIf.S3Bucket.ObjectsThreshold
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

	totalObjects := 0
	err = e.svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: &e.name,
	}, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
		totalObjects += len(output.Contents)
		if objectsThreshold > 0 && objectsThreshold < totalObjects {
			return false
		}
		return !lastPage
	})
	if err != nil {
		return err
	}

	if objectsThreshold > 0 && objectsThreshold < totalObjects {
		_, err := e.svc.PutBucketLifecycleConfiguration(&s3.PutBucketLifecycleConfigurationInput{
			Bucket: &e.name,
			LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
				Rules: []*s3.LifecycleRule{
					{
						Expiration: &s3.LifecycleExpiration{
							Days: aws.Int64(1),
						},
						Filter: &s3.LifecycleRuleFilter{
							Prefix: aws.String(""),
						},
						ID:     aws.String("aws-nuke"),
						Status: aws.String("Enabled"),
						NoncurrentVersionExpiration: &s3.NoncurrentVersionExpiration{
							NoncurrentDays: aws.Int64(1),
						},
						AbortIncompleteMultipartUpload: &s3.AbortIncompleteMultipartUpload{
							DaysAfterInitiation: aws.Int64(1),
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		return config.NewItemPostponedError(
			fmt.Sprintf(
				"bucket %q has more than %d objects, lifecycle was enabled and postponed for a later removal",
				e.name,
				objectsThreshold,
			),
		)
	} else {
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
	return nil
}

func (e *S3Bucket) RemoveAllVersions() error {
	params := &s3.ListObjectVersionsInput{
		Bucket: &e.name,
	}

	iterator := newS3DeleteVersionListIterator(e.svc, params)
	return s3manager.NewBatchDeleteWithClient(e.svc).Delete(aws.BackgroundContext(), iterator)
}

func (e *S3Bucket) RemoveAllObjects() error {
	params := &s3.ListObjectsInput{
		Bucket: &e.name,
	}

	iterator := s3manager.NewDeleteListIterator(e.svc, params)
	return s3manager.NewBatchDeleteWithClient(e.svc).Delete(aws.BackgroundContext(), iterator)
}

func (e *S3Bucket) FeatureFlags(ff config.FeatureFlags) {
	e.featureFlags = ff
}

func (e *S3Bucket) Properties() types.Properties {
	properties := types.NewProperties().
		Set("Name", e.name).
		Set("CreationDate", e.creationDate.Format(time.RFC3339))

	for _, tag := range e.tags {
		properties.SetTag(tag.Key, tag.Value)
	}

	return properties
}

func (e *S3Bucket) String() string {
	return fmt.Sprintf("s3://%s", e.name)
}

type s3DeleteVersionListIterator struct {
	Bucket    *string
	Paginator request.Pagination
	objects   []*s3.ObjectVersion
}

func newS3DeleteVersionListIterator(svc s3iface.S3API, input *s3.ListObjectVersionsInput, opts ...func(*s3DeleteVersionListIterator)) s3manager.BatchDeleteIterator {
	iter := &s3DeleteVersionListIterator{
		Bucket: input.Bucket,
		Paginator: request.Pagination{
			NewRequest: func() (*request.Request, error) {
				var inCpy *s3.ListObjectVersionsInput
				if input != nil {
					tmp := *input
					inCpy = &tmp
				}
				req, _ := svc.ListObjectVersionsRequest(inCpy)
				return req, nil
			},
		},
	}

	for _, opt := range opts {
		opt(iter)
	}
	return iter
}

// Next will use the S3API client to iterate through a list of objects.
func (iter *s3DeleteVersionListIterator) Next() bool {
	if len(iter.objects) > 0 {
		iter.objects = iter.objects[1:]
	}

	if len(iter.objects) == 0 && iter.Paginator.Next() {
		output := iter.Paginator.Page().(*s3.ListObjectVersionsOutput)
		iter.objects = output.Versions

		for _, entry := range output.DeleteMarkers {
			iter.objects = append(iter.objects, &s3.ObjectVersion{
				Key:       entry.Key,
				VersionId: entry.VersionId,
			})
		}
	}

	return len(iter.objects) > 0
}

// Err will return the last known error from Next.
func (iter *s3DeleteVersionListIterator) Err() error {
	return iter.Paginator.Err()
}

// DeleteObject will return the current object to be deleted.
func (iter *s3DeleteVersionListIterator) DeleteObject() s3manager.BatchDeleteObject {
	return s3manager.BatchDeleteObject{
		Object: &s3.DeleteObjectInput{
			Bucket:    iter.Bucket,
			Key:       iter.objects[0].Key,
			VersionId: iter.objects[0].VersionId,
		},
	}
}
