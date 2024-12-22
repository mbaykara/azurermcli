package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	tea "github.com/charmbracelet/bubbletea"
)

func FetchSubscriptions() tea.Msg {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return ErrorMsg{err}
	}

	client, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		return ErrorMsg{err}
	}

	pager := client.NewListPager(nil)
	var subs []armsubscription.Subscription

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return ErrorMsg{err}
		}
		for _, sub := range page.Value {
			subs = append(subs, *sub)
		}
	}

	return SubscriptionsMsg{Subs: subs}
}

func FetchResourceGroups(subscriptionID string) tea.Cmd {
	return func() tea.Msg {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return ErrorMsg{err}
		}

		subscriptionID = strings.TrimSpace(subscriptionID)
		if strings.HasPrefix(subscriptionID, "/subscriptions/") {
			subscriptionID = strings.TrimPrefix(subscriptionID, "/subscriptions/")
		}

		client, err := armresources.NewResourceGroupsClient(subscriptionID, cred, nil)
		if err != nil {
			return ErrorMsg{err}
		}

		pager := client.NewListPager(&armresources.ResourceGroupsClientListOptions{})
		var groups []armresources.ResourceGroup

		for pager.More() {
			page, err := pager.NextPage(context.Background())
			if err != nil {
				return ErrorMsg{err}
			}
			for _, group := range page.Value {
				groups = append(groups, *group)
			}
		}

		return ResourceGroupsMsg{
			SubscriptionID: subscriptionID,
			Groups:        groups,
		}
	}
}

func FetchResources(subscriptionID, resourceGroupName string) tea.Cmd {
	return func() tea.Msg {
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return ErrorMsg{err}
		}

		subscriptionID = strings.TrimSpace(subscriptionID)
		if strings.HasPrefix(subscriptionID, "/subscriptions/") {
			subscriptionID = strings.TrimPrefix(subscriptionID, "/subscriptions/")
		}

		client, err := armresources.NewClient(subscriptionID, cred, nil)
		if err != nil {
			return ErrorMsg{err}
		}

		pager := client.NewListByResourceGroupPager(resourceGroupName, &armresources.ClientListByResourceGroupOptions{})
		var resources []armresources.GenericResourceExpanded

		for pager.More() {
			page, err := pager.NextPage(context.Background())
			if err != nil {
				return ErrorMsg{err}
			}
			for _, resource := range page.Value {
				resources = append(resources, *resource)
			}
		}

		return ResourcesMsg{
			ResourceGroupName: resourceGroupName,
			Resources:        resources,
		}
	}
} 