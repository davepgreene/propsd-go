package parsers

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/davepgreene/propsd/utils"
	"reflect"
	"strings"
	"time"
)

type MetadataPropertiesCredentials struct {
	LastUpdated     time.Time `json:"lastUpdated,omitempty"`
	Type            string    `json:"type,omitempty"`
	AccessKeyId     string    `json:"accessKeyId,omitempty"`
	SecretAccessKey string    `json:"secretAccessKey,omitempty"`
	Expiration      time.Time `json:"expires,omitempty"`
	Code            string    `json:"-"`
	Token           string    `json:"-"`
}

type MetadataPropertiesInterface struct {
	VPCIPV4CIDRBlock    string `json:"vpc-ipv4-cidr-block,omitempty"`
	SubnetIPV4CIDRBlock string `json:"subnet-ipv4-cidr-block,omitempty"`
	MAC                 string `json:"mac,omitempty"`
	LocalIPV4s          string `json:"local-ipv4s,omitempty"`
	PublicIPV4s         string `json:"public-ipv4s,omitempty"`
	InterfaceID         string `json:"interface-id,omitempty"`
	VPCID               string `json:"-"`
}

type MetadataProperties struct {
	Account          string                         `json:"account,omitempty"`
	AmiID            string                         `json:"ami-id,omitempty"`
	AutoScalingGroup string                         `json:"auto-scaling-group,omitempty"`
	AvailabilityZone string                         `json:"availability-zone,omitempty"`
	Credentials      *MetadataPropertiesCredentials `json:"credentials,omitempty"`
	Hostname         string                         `json:"hostname,omitempty"`
	IAMRole          string                         `json:"iam-role,omitempty"`
	Identity         struct {
		Document string `json:"document,omitempty"`
		Pkcs7    string `json:"pkcs7,omitempty"`
	} `json:"identity,omitempty"`
	InstanceID     string                       `json:"instance-id,omitempty"`
	InstanceType   string                       `json:"instance-type,omitempty"`
	Interface      *MetadataPropertiesInterface `json:"interface,omitempty"`
	LocalHostname  string                       `json:"local-hostname,omitempty"`
	LocalIPV4      string                       `json:"local-ipv4,omitempty"`
	PublicHostname string                       `json:"public-hostname,omitempty"`
	PublicIPV4     string                       `json:"public-ipv4,omitempty"`
	Region         string                       `json:"region,omitempty"`
	ReservationID  string                       `json:"reservation-id,omitempty"`
	SecurityGroups string                       `json:"security-groups,omitempty"`
	VPCID          string                       `json:"vpc-id,omitempty"`
	//Tags map[string]string			`json:"tags"` // Rethink this as it'll need to be injected into the JSON rather than the struct
}

type MetadataParser func(string)

type Metadata struct {
	properties *MetadataProperties
	session    session.Session
	Parsers    map[string]MetadataParser
}

func NewMetadataParser(session session.Session) *Metadata {
	properties := &MetadataProperties{
		Credentials: &MetadataPropertiesCredentials{},
		Interface:   &MetadataPropertiesInterface{},
	}
	metadataClient := utils.CreateMetadataClient(session.ClientConfig("ec2metadata", aws.NewConfig()))
	parsers := map[string]MetadataParser{
		"instance-identity/document": func(body string) {
			if len(body) == 0 {
				return
			}

			properties.Identity.Document = body

			var document ec2metadata.EC2InstanceIdentityDocument
			err := json.Unmarshal([]byte(body), &document)
			if err != nil {
				return
			}

			properties.Account = document.AccountID
			properties.Region = document.Region
			properties.AvailabilityZone = document.AvailabilityZone
			properties.AmiID = document.ImageID
			properties.InstanceID = document.InstanceID
			properties.InstanceType = document.InstanceType
		},
		"hostname":                func(body string) { properties.Hostname = body },
		"local-ipv4":              func(body string) { properties.LocalIPV4 = body },
		"local-hostname":          func(body string) { properties.LocalHostname = body },
		"public-hostname":         func(body string) { properties.PublicHostname = body },
		"public-ipv4":             func(body string) { properties.PublicIPV4 = body },
		"reservation-id":          func(body string) { properties.ReservationID = body },
		"security-groups":         func(body string) { properties.SecurityGroups = body },
		"instance-identity/pkcs7": func(body string) { properties.Identity.Pkcs7 = body },
		"iam/security-credentials": func(body string) {
			if len(body) == 0 {
				return
			}

			properties.IAMRole = body

			// We need to make another request to get role data
			roleData, err := metadataClient.GetMetadata(fmt.Sprintf("iam/security-credentials/%s", body))
			if err != nil {
				return
			}

			var creds MetadataPropertiesCredentials
			err = json.Unmarshal([]byte(roleData), &creds)
			if err != nil {
				return
			}

			properties.Credentials = &creds
		},
		"network/interfaces/macs": func(body string) {
			if len(body) == 0 {
				return
			}

			i := &MetadataPropertiesInterface{}

			interfacePaths := map[string]string{
				"vpc-ipv4-cidr-block":    "VPCIPV4CIDRBlock",
				"subnet-ipv4-cidr-block": "SubnetIPV4CIDRBlock",
				"public-ipv4s":           "PublicIPV4s",
				"mac":                    "MAC",
				"local-ipv4s":            "LocalIPV4s",
				"interface-id":           "InterfaceID",
				"vpc-id":                 "VPCID",
			}

			i.MAC = strings.TrimSuffix(body, "/")

			for path, field := range interfacePaths {
				data, err := metadataClient.GetMetadata(fmt.Sprintf("network/interfaces/macs/%s/%s", i.MAC, path))

				if err != nil {
					continue
				}

				// Using reflection we can assign a value to a struct field by name.
				v := reflect.ValueOf(i).Elem().FieldByName(field)
				if v.IsValid() {
					v.SetString(data)
				}
			}

			properties.Interface = i
			properties.VPCID = i.VPCID
		},
		"auto-scaling-group": func(body string) {
			autoscalingClient := autoscaling.New(&session, aws.NewConfig().WithRegion(properties.Region))

			result, err := autoscalingClient.DescribeAutoScalingInstances(&autoscaling.DescribeAutoScalingInstancesInput{
				InstanceIds: []*string{
					aws.String(properties.InstanceID),
				},
			})
			if err != nil {
				return
			}

			// As long as we get an instance back we should trust that the ASG API's query functionality works
			if len(result.AutoScalingInstances) > 0 {
				properties.AutoScalingGroup = *result.AutoScalingInstances[0].AutoScalingGroupName
			}
		},
	}

	return &Metadata{
		session:    session,
		properties: properties,
		Parsers:    parsers,
	}
}

//
//func (m *Metadata) parsers() map[string]func(string) {
//	metadataClient := utils.CreateMetadataClient(m.session.ClientConfig("ec2metadata", aws.NewConfig()))
//	return
//}

func (m *Metadata) Properties() *MetadataProperties {
	return m.properties
}
