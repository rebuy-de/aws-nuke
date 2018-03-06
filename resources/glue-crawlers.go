package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

type GlueCrawler struct {
	svc  *glue.Glue
	name *string
}

func init() {
	register("GlueCrawler", ListGlueCrawlers)
}

func ListGlueCrawlers(sess *session.Session) ([]Resource, error) {
	svc := glue.New(sess)
	resources := []Resource{}

	params := &glue.GetCrawlersInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.GetCrawlers(params)
		if err != nil {
			return nil, err
		}

		for _, crawler := range output.Crawlers {
			resources = append(resources, &GlueCrawler{
				svc:  svc,
				name: crawler.Name,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *GlueCrawler) Remove() error {

	_, err := f.svc.DeleteCrawler(&glue.DeleteCrawlerInput{
		Name: f.name,
	})

	return err
}

func (f *GlueCrawler) String() string {
	return *f.name
}
