package resources

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/signer"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type SignerSigningJob struct {
	svc                 *signer.Signer
	jobId               *string
	reason              string
	isRevoked           *bool
	createdAt           time.Time
	profileName         *string
	profileVersion      *string
	platformId          *string
	platformDisplayName *string
	jobOwner            *string
	jobInvoker          *string
}

func init() {
	register("SignerSigningJob", ListSignerSigningJobs)
}

func ListSignerSigningJobs(sess *session.Session) ([]Resource, error) {
	svc := signer.New(sess)
	resources := []Resource{}
	const reason string = "Revoked by AWS Nuke"

	listJobsInput := &signer.ListSigningJobsInput{}

	err := svc.ListSigningJobsPages(listJobsInput, func(page *signer.ListSigningJobsOutput, lastPage bool) bool {
		for _, job := range page.Jobs {
			resources = append(resources, &SignerSigningJob{
				svc:                 svc,
				jobId:               job.JobId,
				reason:              reason,
				isRevoked:           job.IsRevoked,
				createdAt:           *job.CreatedAt,
				profileName:         job.ProfileName,
				profileVersion:      job.ProfileVersion,
				platformId:          job.PlatformId,
				platformDisplayName: job.PlatformDisplayName,
				jobOwner:            job.JobOwner,
				jobInvoker:          job.JobInvoker,
			})
		}
		return true // continue iterating over pages
	})
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func (j *SignerSigningJob) Filter() error {
	// Consider all non-revoked jobs
	if *j.isRevoked {
		return fmt.Errorf("job already revoked")
	}
	return nil
}

func (j *SignerSigningJob) Remove() error {
	// Signing jobs are viewable by the ListSigningJobs operation for two years after they are performed [1]
	// As a precaution we are updating Signing jobs statuses to revoked. This indicates that the signature is no longer valid.
	// [1] https://awscli.amazonaws.com/v2/documentation/api/latest/reference/signer/start-signing-job.html
	revokeInput := &signer.RevokeSignatureInput{
		JobId:  j.jobId,
		Reason: aws.String(j.reason),
	}
	_, err := j.svc.RevokeSignature(revokeInput)
	return err
}

func (j *SignerSigningJob) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("JobId", j.jobId)
	properties.Set("CreatedAt", j.createdAt.Format(time.RFC3339))
	properties.Set("ProfileName", j.profileName)
	properties.Set("ProfileVersion", j.profileVersion)
	properties.Set("PlatformId", j.platformId)
	properties.Set("PlatformDisplayName", j.platformDisplayName)
	properties.Set("JobOwner", j.jobOwner)
	properties.Set("JobInvoker", j.jobInvoker)
	return properties
}
