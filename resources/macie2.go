package resources

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/macie2"
)

type Macie struct {
	svc *macie2.Macie2
}

func init() {
	register("Macie", CheckMacieStatus)
}

func CheckMacieStatus(sess *session.Session) ([]Resource, error) {
	svc := macie2.New(sess)

	status, err := svc.GetMacieSession(&macie2.GetMacieSessionInput{})
	if err != nil {
		if err.Error() == "AccessDeniedException: Macie is not enabled" {
			return nil, nil
		} else {
			return nil, err
		}
	}

	resources := make([]Resource, 0)
	if *status.Status == macie2.AdminStatusEnabled {
		resources = append(resources, &Macie{
			svc: svc,
		})
	}

	return resources, nil
}

func (b *Macie) Remove() error {
	_, err := b.svc.DisableMacie(&macie2.DisableMacieInput{})
	if err != nil {
		return err
	}
	return nil
}

func (b *Macie) String() string {
	return "macie"
}
