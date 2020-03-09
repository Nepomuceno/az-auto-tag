package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	writeOperation      = regexp.MustCompile(`.*\/write$`)
	unsuportedProviders = map[string]bool{
		"microsoft.eventhub/namespaces/eventhubs/authorizationrules": true,
		"microsoft.eventhub/namespaces/eventhubs/consumergroups":     true,
		"microsoft.network/privatednszones/a":                        true,
		"microsoft.network/networksecuritygroups/securityrules":      true,
		"microsoft.sql/servers/databases/securityalertpolicies":      true,
		"microsoft.support/supporttickets":                           true,
		"microsoft.insights/diagnosticsettings":                      true,
		"microsoft.storage/storageaccounts/blobservices/containers":  true,
		"microsoft.network/virtualnetworks/subnets":                  true,
		"microsoft.network/routetables/routes":                       true,
		"microsoft.network/virtualnetworks/virtualnetworkpeerings":   true,
		"microsoft.sql/servers/firewallrules":                        true,
		"microsoft.authorization/roledefinitions":                    true,
		"microsoft.authorization/roleassignments":                    true,
	}
	resouceIDRegex = regexp.MustCompile(`(?m)\/subscriptions\/(?P<subscription>[^\/]+)\/resourceGroups\/(?P<resourceGroup>[^\/]+)\/providers\/(?P<resourceProvider>[^\/]+)\/(?P<resouceType>[^\/]+)\/(?P<resourceName>[^\/]+)(\/)?(?P<resourceSubtype>[^\/]+)?(\/)?(?P<resourceSubtypeName>[^\/]+)?$`)
)

func getResource(resource string) *AzureResource {
	matches := resouceIDRegex.FindStringSubmatch(resource)
	result := &AzureResource{}
	if len(matches) > 1 {
		result.Subscription = matches[1]
		result.ResourceGroup = matches[2]
		result.Provider = matches[3]
		result.Type = matches[4]
		result.Name = matches[5]
		result.SubType = matches[7]
		result.SubName = matches[9]
	}
	return result
}

func getAppName(appID *string, authGraph autorest.Authorizer) string {
	appClient := graphrbac.NewServicePrincipalsClient(tenant)
	appClient.Authorizer = authGraph
	servPrincipal, err := appClient.Get(context.Background(), *appID)
	if err != nil {
		bodyBytes, err := ioutil.ReadAll(servPrincipal.Response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}
	return *servPrincipal.DisplayName
}

func newAuthorizer() (*autorest.Authorizer, error) {
	// Carry out env var lookups
	_, clientIDExists := os.LookupEnv("AZURE_CLIENT_ID")
	_, tenantIDExists := os.LookupEnv("AZURE_TENANT_ID")
	_, fileAuthSet := os.LookupEnv("AZURE_AUTH_LOCATION")

	// Execute logic to return an authorizer from the correct method
	if clientIDExists && tenantIDExists {
		log.Println("Logging from environment")
		authorizer, err := auth.NewAuthorizerFromEnvironment()
		return &authorizer, err
	} else if fileAuthSet {
		log.Println("Logging from file")
		authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
		return &authorizer, err
	} else {
		log.Println("Logging from CLI")
		authorizer, err := auth.NewAuthorizerFromCLI()
		return &authorizer, err
	}
}

func newGraphAuthorizer() (*autorest.Authorizer, error) {
	// Carry out env var lookups
	_, clientIDExists := os.LookupEnv("AZURE_CLIENT_ID")
	_, tenantIDExists := os.LookupEnv("AZURE_TENANT_ID")
	_, fileAuthSet := os.LookupEnv("AZURE_AUTH_LOCATION")

	// Execute logic to return an authorizer from the correct method
	if clientIDExists && tenantIDExists {
		log.Println("Logging from environment")
		authorizer, err := auth.NewAuthorizerFromEnvironmentWithResource(azure.PublicCloud.GraphEndpoint)
		return &authorizer, err
	} else if fileAuthSet {
		log.Println("Logging from file")
		authorizer, err := auth.NewAuthorizerFromFile(azure.PublicCloud.GraphEndpoint)
		return &authorizer, err
	} else {
		log.Println("Logging from CLI")
		authorizer, err := auth.NewAuthorizerFromCLIWithResource(azure.PublicCloud.GraphEndpoint)
		return &authorizer, err
	}
}
