package utils

import (
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/spf13/viper"
	"net/url"
)

func CreateMetadataClient(c client.Config) *ec2metadata.EC2Metadata {
	client := ec2metadata.NewClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion,
		// Additional functions to modify the client
		func(client *client.Client) {
			if (viper.IsSet("metadata.version")) {
				client.APIVersion = viper.GetString("metadata.version")
			}
		}, func(client *client.Client) {
			if (viper.IsSet("metadata.host")) {
				endpoint, err := url.Parse(viper.GetString("metadata.host"))
				if err != nil {
					return
				}
				api, _ := url.Parse(client.APIVersion)
				client.Endpoint = endpoint.ResolveReference(api).String()
			}
		}, func(client *client.Client) {
			if (viper.IsSet("metadata.timeout")) {
				client.Config.HTTPClient.Timeout = viper.GetDuration("metadata.timeout")
			}
		})
	return client
}
