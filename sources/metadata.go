package sources

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/davepgreene/propsd/api"
	log "github.com/Sirupsen/logrus"
	"github.com/davepgreene/propsd/parsers"
	"github.com/davepgreene/propsd/utils"
)

type MetadataOptions struct {
}

type MetadataChannelResponse struct {
	Path string
	Body string
}

type MetadataChannelErrorResponse struct {
	Path  string
	Error error
}

type Metadata struct {
	client *ec2metadata.EC2Metadata
	parser *parsers.Metadata
}

func NewMetadataSource(session session.Session) *Metadata {
	// We need to assemble our client manually so we can override host and timeout if we want
	c := session.ClientConfig("ec2metadata", aws.NewConfig())

	return &Metadata{
		client: utils.CreateMetadataClient(c),
		parser: parsers.NewMetadataParser(session),
	}
}

func (m *Metadata) Get() {
	resc, errc := make(chan string), make(chan string)
	paths := map[string]func(string) (string, error){
		"instance-identity/document": m.client.GetDynamicData,
		"hostname":                   m.client.GetMetadata,
		"local-ipv4":                 m.client.GetMetadata,
		"local-hostname":             m.client.GetMetadata,
		"public-hostname":            m.client.GetMetadata,
		"public-ipv4":                m.client.GetMetadata,
		"reservation-id":             m.client.GetMetadata,
		"security-groups":            m.client.GetMetadata,
		"instance-identity/pkcs7":    m.client.GetDynamicData,
		"iam/security-credentials":   m.client.GetMetadata,
		"network/interfaces/macs":    m.client.GetMetadata,
	}

	for path, fn := range paths {
		go m.fetch(resc, errc, path, fn, m.parser.Parsers[path])
	}

	for i := 0; i < len(paths); i++ {
		select {
		case res := <-resc:
			log.Debugf("Parsed data from %s", res)
		case err := <-errc:
			log.Errorf("Aws-sdk returned the following error during the metadata service request to %s", err)
		}
	}

	// We can use goroutines for all the other metadata but because ASG relies on instance region and ID we
	// have to wait until those are complete.
	//
	// If we run into performance issues here we can implement a queue with re-queue operations, spawn
	// goroutines for every queued element and then re-queue the ASG request if region and instance ID
	// haven't been retrieved yet.
	m.parser.Parsers["auto-scaling-group"]("")
}

func (m *Metadata) Poll() {
	log.Info(m.Properties())
}

func (m *Metadata) fetch(resc chan string, errc chan string, path string, method func(string) (string, error), parser func(string)) {
	body, err := method(path)
	if err != nil {
		errc <- path
		return
	}
	parser(body)
	resc <- path
}

func (m *Metadata) Properties() *parsers.MetadataProperties {
	return m.parser.Properties()
}
