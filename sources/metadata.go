package sources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	//"github.com/davepgreene/propsd/api"
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/aws/aws-sdk-go/aws/client"
	"fmt"
)

var METADATA_PATHS = []string{
	//"ami-id", GET FROM EC2InstanceIdentityDocument - ImageID
	//"placement/availability-zone", GET FROM EC2InstanceIdentityDocument - AvailabilityZone
	"hostname",
	//"instance-id", GET FROM EC2InstanceIdentityDocument - InstanceID
	//"instance-type", GET FROM EC2InstanceIdentityDocument - InstanceType
	// "local-ipv4",  GET FROM EC2InstanceIdentityDocument - PrivateIP
	"local-hostname",
	"public-hostname",
	"public-ipv4",
	"reservation-id",
	"security-groups",
	//"identity", GET FROM EC2InstanceIdentityDocument - Marshall object to JSON
	//"region", GET FROM EC2InstanceIdentityDocument - Region
	//"account", GET FROM EC2InstanceIdentityDocument - AccountID
	"iam/security-credentials",
	"mac",
	"network/interfaces/macs",
}

var DYNAMIC_PATHS = []string{
	"instance-identity/pkcs7",
}

type MetadataOptions struct {

}

type Metadata struct {
	rawProperties map[string]*json.RawMessage
	client *ec2metadata.EC2Metadata
}

func NewMetadataSource(session session.Session) *Metadata {
	// We need to assemble our client manually so we can override host and timeout if we want
	c := session.ClientConfig("ec2metadata", aws.NewConfig())
	endpoint := viper.GetString("metadata.host")
	client := ec2metadata.NewClient(*c.Config, c.Handlers, endpoint, c.SigningRegion,
		// Additional functions to modify the client
		func(client *client.Client) {
			if (viper.IsSet("metadata.version")) {
				client.APIVersion = viper.GetString("metadata.version")
			}
		}, func(client *client.Client) {
			if (viper.IsSet("metadata.timeout")) {
				client.Config.HTTPClient.Timeout = viper.GetDuration("metadata.timeout")
			}
		})

	return &Metadata{
		client: client,
		rawProperties: make(map[string]*json.RawMessage),
	}
}

func (m *Metadata) Get() {
	resc, errc := make(chan string), make(chan error)
	totalRequests := len(METADATA_PATHS) + len(DYNAMIC_PATHS)
	f := m.fetch(resc, errc)

	for _, path := range METADATA_PATHS {
		go f(path, m.client.GetMetadata)
	}

	for _, path := range DYNAMIC_PATHS {
		go f(path, m.client.GetDynamicData)
	}

	for i := 0; i < totalRequests; i++ {
		select {
		case res := <-resc:
			fmt.Println(res)
		case err := <-errc:
			fmt.Println(err)
		}
	}

}

func (m *Metadata) fetch(resc chan string, errc chan error) func(path string, method func(string) (string, error)) {
	return func(path string, method func(string) (string, error)) {
		body, err := method(path)
		if err != nil {
			errc <- err
			return
		}
		resc <- string(body)
	}
}

func (m *Metadata) getInstanceIdentityDocument() {

}

//func (m *Metadata)