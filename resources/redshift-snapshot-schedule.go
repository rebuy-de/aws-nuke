package resources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/redshift"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

type RedshiftSnapshotSchedule struct {
	svc                *redshift.Redshift
	scheduleID         *string
	associatedClusters []*redshift.ClusterAssociatedToSchedule
}

func init() {
	register("RedshiftSnapshotSchedule", ListRedshiftSnapshotSchedule)
}

func ListRedshiftSnapshotSchedule(sess *session.Session) ([]Resource, error) {
	svc := redshift.New(sess)
	resources := []Resource{}

	params := &redshift.DescribeSnapshotSchedulesInput{
		MaxRecords: aws.Int64(100),
	}

	for {
		output, err := svc.DescribeSnapshotSchedules(params)
		if err != nil {
			return nil, err
		}

		for _, snapshotSchedule := range output.SnapshotSchedules {
			resources = append(resources, &RedshiftSnapshotSchedule{
				svc:                svc,
				scheduleID:         snapshotSchedule.ScheduleIdentifier,
				associatedClusters: snapshotSchedule.AssociatedClusters,
			})
		}

		if output.Marker == nil {
			break
		}

		params.Marker = output.Marker
	}

	return resources, nil
}

func (f *RedshiftSnapshotSchedule) Properties() types.Properties {
	associatedClusters := make([]string, len(f.associatedClusters))
	for i, cluster := range f.associatedClusters {
		associatedClusters[i] = *cluster.ClusterIdentifier
	}
	properties := types.NewProperties()
	properties.Set("scheduleID", f.scheduleID)
	properties.Set("associatedClusters", associatedClusters)
	return properties
}

func (f *RedshiftSnapshotSchedule) Remove() error {
	for _, associatedCluster := range f.associatedClusters {
		_, disassociateErr := f.svc.ModifyClusterSnapshotSchedule(&redshift.ModifyClusterSnapshotScheduleInput{
			ScheduleIdentifier:   f.scheduleID,
			ClusterIdentifier:    associatedCluster.ClusterIdentifier,
			DisassociateSchedule: aws.Bool(true),
		})

		if disassociateErr != nil {
			return disassociateErr
		}
	}

	_, err := f.svc.DeleteSnapshotSchedule(&redshift.DeleteSnapshotScheduleInput{
		ScheduleIdentifier: f.scheduleID,
	})

	return err
}
