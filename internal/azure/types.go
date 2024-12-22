package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

type SubscriptionsMsg struct {
	Subs []armsubscription.Subscription
}

type ResourceGroupsMsg struct {
	SubscriptionID string
	Groups        []armresources.ResourceGroup
}

type ResourcesMsg struct {
	ResourceGroupName string
	Resources        []armresources.GenericResourceExpanded
}

type ErrorMsg struct {
	Error error
}

type LoadingMsg bool 