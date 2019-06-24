package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type SecretsManagerSecret struct {
	svc *secretsmanager.SecretsManager
	ARN *string
}

func init() {
	register("SecretsManagerSecret", ListSecretsManagerSecrets)
}

func ListSecretsManagerSecrets(sess *session.Session) ([]Resource, error) {
	svc := secretsmanager.New(sess)
	resources := []Resource{}

	params := &secretsmanager.ListSecretsInput{
		MaxResults: aws.Int64(100),
	}

	for {
		output, err := svc.ListSecrets(params)
		if err != nil {
			return nil, err
		}

		for _, secrets := range output.SecretList {
			resources = append(resources, &SecretsManagerSecret{
				svc: svc,
				ARN: secrets.ARN,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

func (f *SecretsManagerSecret) Remove() error {

	_, err := f.svc.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   f.ARN,
		ForceDeleteWithoutRecovery: aws.Bool(true),
	})

	return err
}

func (f *SecretsManagerSecret) String() string {
	return *f.ARN
}
