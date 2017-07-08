package sources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	//"github.com/davepgreene/propsd/api"
	"github.com/davepgreene/propsd/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/davepgreene/propsd/parsers"
)

type MetadataOptions struct {

}

type MetadataChannelResponse struct {
	Path string
	Body string
}

type MetadataChannelErrorResponse struct {
	Path string
	Error error
}

type Metadata struct {
	rawProperties map[string]string
	client *ec2metadata.EC2Metadata
	parser *parsers.Metadata
}

func NewMetadataSource(session session.Session) *Metadata {
	// We need to assemble our client manually so we can override host and timeout if we want
	c := session.ClientConfig("ec2metadata", aws.NewConfig())

	return &Metadata{
		client: utils.CreateMetadataClient(c),
		rawProperties: make(map[string]string),
		parser: parsers.NewMetadataParser(session),
	}
}

func (m *Metadata) Get() {
	resc, errc := make(chan MetadataChannelResponse), make(chan MetadataChannelErrorResponse)
	paths := map[string]string {
		"instance-identity/document": "GetDynamicData",
		"hostname": "GetMetadata",
		"local-ipv4": "GetMetadata",
		"local-hostname": "GetMetadata",
		"public-hostname": "GetMetadata",
		"public-ipv4": "GetMetadata",
		"reservation-id": "GetMetadata",
		"security-groups": "GetMetadata",
		"instance-identity/pkcs7": "GetDynamicData",
		"iam/security-credentials": "GetMetadata",
		"network/interfaces/macs": "GetMetadata",
	}

	for path, function := range paths {
		fn := utils.GetMethod(m.client, function).(func(string) (string, error))
		go m.fetch(resc, errc)(path, fn)
	}

	for i := 0; i < len(paths); i++ {
		select {
		case res := <-resc:
			m.rawProperties[res.Path] = res.Body
		case err := <-errc:
			log.Errorf("Aws-sdk returned the following error during the metadata service request to %s", err.Path)
		}
	}

	m.rawProperties["auto-scaling-group"] = ""
}

func (m *Metadata) Poll() {
	m.parser.Parse(m.rawProperties)
}

func (m *Metadata) fetch(resc chan MetadataChannelResponse, errc chan MetadataChannelErrorResponse) func(string, func(string) (string, error)) {
	return func(path string, method func(string) (string, error)) {
		body, err := method(path)
		if err != nil {
			errc <- MetadataChannelErrorResponse{path,err}
			return
		}
		resc <- MetadataChannelResponse{path,body}
	}
}

func (m *Metadata) Properties() *parsers.MetadataProperties {
	return m.parser.Properties()
}
