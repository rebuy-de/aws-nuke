package resources

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/quicksight"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/rebuy-de/aws-nuke/v2/pkg/types"
)

func init() {
	register("QuicksightSubscription", DescribeQuicksightSubscription)
}

type QuicksightSubscription struct {
	svc               *quicksight.QuickSight
	accountId         *string
	subscriptionName  *string
	notificationEmail *string
	edition           *string
}

func DescribeQuicksightSubscription(session *session.Session) ([]Resource, error) {
	const activeSubscriptionStatus = "ACCOUNT_CREATED"

	stsSvc := sts.New(session)
	callerID, err := stsSvc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}
	accountId := callerID.Account

	resources := []Resource{}

	quicksightSvc := quicksight.New(session)

	describeSubscriptionOutput, err := quicksightSvc.DescribeAccountSubscription(&quicksight.DescribeAccountSubscriptionInput{
		AwsAccountId: accountId,
	})

	if err != nil {
		var resoureceNotFoundException *quicksight.ResourceNotFoundException
		if !errors.As(err, &resoureceNotFoundException) {
			return nil, err
		}
		return resources, nil
	}

	//The account name is only available some time later after the Subscription creation.
	//Since it is an important value to identify the resource, it will wait till it is available
	if *describeSubscriptionOutput.AccountInfo.AccountSubscriptionStatus != activeSubscriptionStatus || describeSubscriptionOutput.AccountInfo.AccountName == nil {
		return resources, nil
	}

	resources = append(resources, &QuicksightSubscription{
		svc:               quicksightSvc,
		accountId:         accountId,
		subscriptionName:  describeSubscriptionOutput.AccountInfo.AccountName,
		notificationEmail: describeSubscriptionOutput.AccountInfo.NotificationEmail,
		edition:           describeSubscriptionOutput.AccountInfo.Edition,
	})

	return resources, nil
}

func (subscription *QuicksightSubscription) Remove() error {
	terminateProtectionEnabled := false

	describeSettingsOutput, err := subscription.svc.DescribeAccountSettings(&quicksight.DescribeAccountSettingsInput{
		AwsAccountId: subscription.accountId,
	})
	if err != nil {
		return err
	}

	if *describeSettingsOutput.AccountSettings.TerminationProtectionEnabled {
		updateSettingsInput := quicksight.UpdateAccountSettingsInput{
			AwsAccountId:                 subscription.accountId,
			DefaultNamespace:             describeSettingsOutput.AccountSettings.DefaultNamespace,
			NotificationEmail:            describeSettingsOutput.AccountSettings.NotificationEmail,
			TerminationProtectionEnabled: &terminateProtectionEnabled,
		}

		_, err = subscription.svc.UpdateAccountSettings(&updateSettingsInput)
		if err != nil {
			return err
		}
	}

	_, err = subscription.svc.DeleteAccountSubscription(&quicksight.DeleteAccountSubscriptionInput{
		AwsAccountId: subscription.accountId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (subscription *QuicksightSubscription) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("Edition", subscription.edition).
		Set("NotificationEmail", subscription.notificationEmail).
		Set("Name", subscription.subscriptionName)

	return properties
}

func (subscription *QuicksightSubscription) String() string {
	return *subscription.subscriptionName
}
