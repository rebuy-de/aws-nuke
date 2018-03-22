package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESReceiptFilter struct {
	svc  *ses.SES
	name *string
}

func init() {
	register("SESReceiptFilter", ListSESReceiptFilters)
}

func ListSESReceiptFilters(sess *session.Session) ([]Resource, error) {
	svc := ses.New(sess)
	resources := []Resource{}

	params := &ses.ListReceiptFiltersInput{}

	output, err := svc.ListReceiptFilters(params)
	if err != nil {
		return nil, err
	}

	for _, filter := range output.Filters {
		resources = append(resources, &SESReceiptFilter{
			svc:  svc,
			name: filter.Name,
		})
	}

	return resources, nil
}

func (f *SESReceiptFilter) Remove() error {

	_, err := f.svc.DeleteReceiptFilter(&ses.DeleteReceiptFilterInput{
		FilterName: f.name,
	})

	return err
}

func (f *SESReceiptFilter) String() string {
	return *f.name
}
