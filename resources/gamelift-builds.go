package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gamelift"
)

type GameLiftBuild struct {
	svc     *gamelift.GameLift
	BuildId string
}

func init() {
	register("GameLiftBuild", ListGameLiftBuilds)
}

func ListGameLiftBuilds(sess *session.Session) ([]Resource, error) {
	svc := gamelift.New(sess)

	resp, err := svc.ListBuilds(&gamelift.ListBuildsInput{})
	if err != nil {
		return nil, err
	}

	builds := make([]Resource, 0)
	for _, build := range resp.Builds {
		builds = append(builds, &GameLiftBuild{
			svc:     svc,
			BuildId: aws.StringValue(build.BuildId),
		})
	}

	return builds, nil
}

func (build *GameLiftBuild) Remove() error {
	params := &gamelift.DeleteBuildInput{
		BuildId: aws.String(build.BuildId),
	}

	_, err := build.svc.DeleteBuild(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *GameLiftBuild) String() string {
	return i.BuildId
}
