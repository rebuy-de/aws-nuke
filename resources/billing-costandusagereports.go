package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costandusagereportservice"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("BillingCostandUsageReport", ListBillingCostandUsageReports)
}

type BillingCostandUsageReport struct {
	svc        *costandusagereportservice.CostandUsageReportService
	reportName *string
	s3Bucket   *string
	s3Prefix   *string
	s3Region   *string
}

func ListBillingCostandUsageReports(sess *session.Session) ([]Resource, error) {
	svc := costandusagereportservice.New(sess)
	params := &costandusagereportservice.DescribeReportDefinitionsInput{
		MaxResults: aws.Int64(5),
	}

	reports := make([]*costandusagereportservice.ReportDefinition, 0)
	err := svc.DescribeReportDefinitionsPages(params, func(page *costandusagereportservice.DescribeReportDefinitionsOutput, lastPage bool) bool {
		for _, out := range page.ReportDefinitions {
			reports = append(reports, out)
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	resources := []Resource{}
	for _, report := range reports {
		resources = append(resources, &BillingCostandUsageReport{
			svc:         svc,
			reportName:  report.ReportName,
			s3Bucket:    report.S3Bucket,
			s3Prefix:    report.S3Prefix,
			s3Region:    report.S3Region,
		})
	}

	return resources, nil
}

func (r *BillingCostandUsageReport) Remove() error {
	_, err := r.svc.DeleteReportDefinition(&costandusagereportservice.DeleteReportDefinitionInput{
		ReportName: r.reportName,
	})

	return err
}

func (r *BillingCostandUsageReport) Properties() types.Properties {
	properties := types.NewProperties()
	properties.
		Set("Name", *r.reportName).
		Set("S3Bucket", *r.s3Bucket).
		Set("s3Prefix", *r.s3Prefix).
		Set("S3Region", *r.s3Region)
	return properties
}

func (r *BillingCostandUsageReport) String() string {
	return *r.reportName
}
