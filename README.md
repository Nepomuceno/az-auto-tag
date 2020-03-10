# Azure auto tag

The goal of this project is to auto-tag your resources with the creation tag specifying who has created that resource. This is don looking back at the activity log.

## Required Settings

You are required to pass your azure authentication as environment variables to the container. The application will check all the subscription that the user has access to.

You can find more details on how to use the environment variables to authenticate on azure on [Azure Docs](https://docs.microsoft.com/en-us/azure/go/azure-sdk-go-authorization)

A sample configuration would have the following env variables

```
AZURE_TENANT_ID=YOUR_TENAT_ID
AZURE_CLIENT_ID=YOUR_SERVICE_PRINCIPAL_ID
AZURE_CLIENT_SECRET=YOUR_SERVICE_PRINCIPAL_PASSWORD
```

# Limitations of the current implementation

The current implementation will always look at the last 90 days of activity log independent if it has done the same 5 minutes ago this could be easily solved but was out of the scope for this first release. 

Current you can't change the tags to be created they will always be `Created-by` and `Created-by-id` The idea it is that you would be able to specify those through environment variables like the rest of the configurations.
