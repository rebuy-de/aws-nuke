package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/gamelift"
)

type GameLiftFleet struct {
	svc     *gamelift.GameLift
	FleetId string
}

func init() {
	register("GameLiftFleet", ListGameLiftFleets)
}

func ListGameLiftFleets(sess *session.Session) ([]Resource, error) {
	svc := gamelift.New(sess)

	resp, err := svc.ListFleets(&gamelift.ListFleetsInput{})
	if err != nil {
		return nil, err
	}

	fleets := make([]Resource, 0)
	for _, fleetId := range resp.FleetIds {
		fleet := &GameLiftFleet{
			svc:     svc,
			FleetId: *fleetId, // Dereference the fleetId pointer
		}
		fleets = append(fleets, fleet)
	}

	return fleets, nil
}

func (fleet *GameLiftFleet) Remove() error {
	params := &gamelift.DeleteFleetInput{
		FleetId: aws.String(fleet.FleetId),
	}

	_, err := fleet.svc.DeleteFleet(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *GameLiftFleet) String() string {
	return i.FleetId
}
