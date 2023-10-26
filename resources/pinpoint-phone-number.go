package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pinpointsmsvoicev2"
)

type PinpointPhoneNumber struct {
	svc   *pinpointsmsvoicev2.PinpointSMSVoiceV2
	phone string
}

func init() {
	register("PinpointPhoneNumber", ListPinpointPhoneNumbers)
}

func ListPinpointPhoneNumbers(sess *session.Session) ([]Resource, error) {
	svc := pinpointsmsvoicev2.New(sess)

	resp, err := svc.DescribePhoneNumbers(&pinpointsmsvoicev2.DescribePhoneNumbersInput{})
	if err != nil {
		return nil, err
	}

	numbers := make([]Resource, 0)
	for _, number := range resp.PhoneNumbers {
		numbers = append(numbers, &PinpointPhoneNumber{
			svc:   svc,
			phone: aws.StringValue(number.PhoneNumberId),
		})
	}

	return numbers, nil
}

func (p *PinpointPhoneNumber) Remove() error {
	params := &pinpointsmsvoicev2.ReleasePhoneNumberInput{
		PhoneNumberId: aws.String(p.phone),
	}

	_, err := p.svc.ReleasePhoneNumber(params)
	if err != nil {
		return err
	}

	return nil
}

func (p *PinpointPhoneNumber) String() string {
	return p.phone
}
