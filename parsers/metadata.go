package parsers

import "time"

type MetadataPropertiesIdentity struct {
	Document string
	Pkcs7 string
}

type MetadataPropertiesCredentials struct {
	LastUpdated time.Time
	Type string
	AccessKeyId string
	SecretAccessKey string
	Expires time.Time
}

type MetadataPropertiesInterface struct {
	VPCIPV4CIDRBlock string
	SubnetIPV4CIDRBlock string
	MAC string
	LocalIPV4s string
	InterfaceID string
}

type MetadataProperties struct {
	AmiID string
	AvailabilityZone string
	Hostname string
	InstanceID string
	InstanceType string
	LocalIPV4 string
	LocalHostname string
	PublicHostname string
	PublicIPV4 string
	ReservationID string
	SecurityGroups string
	Identity *MetadataPropertiesIdentity
	Account string
	Region string
	IAMRole string
	Credentials *MetadataPropertiesCredentials
	Interface *MetadataPropertiesInterface
	VPCID string
	AutoScalingGroup string
	Tags map[string]string
}

type Metadata struct {
	properties *MetadataProperties
}
