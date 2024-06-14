package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gamelift"
)

type GameLiftMatchmakingConfiguration struct {
	svc  *gamelift.GameLift
	Name string
}

func init() {
	register("GameLiftMatchmakingConfiguration", ListMatchmakingConfigurations)
}

func ListMatchmakingConfigurations(sess *session.Session) ([]Resource, error) {
	svc := gamelift.New(sess)

	resp, err := svc.DescribeMatchmakingConfigurations(&gamelift.DescribeMatchmakingConfigurationsInput{})
	if err != nil {
		return nil, err
	}

	configs := make([]Resource, 0)
	for _, config := range resp.Configurations {
		q := &GameLiftMatchmakingConfiguration{
			svc:  svc,
			Name: *config.Name,
		}
		configs = append(configs, q)
	}

	return configs, nil
}

func (config *GameLiftMatchmakingConfiguration) Remove() error {
	params := &gamelift.DeleteMatchmakingConfigurationInput{
		Name: aws.String(config.Name),
	}

	_, err := config.svc.DeleteMatchmakingConfiguration(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *GameLiftMatchmakingConfiguration) String() string {
	return i.Name
}
