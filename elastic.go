package esconvert

import (
	"fmt"

	"esconvert/client"

	es7 "github.com/elastic/go-elasticsearch/v7"
	es8 "github.com/elastic/go-elasticsearch/v8"
)

type ES_VERSION string
type EsClientError error

var ErrorVersion EsClientError = fmt.Errorf("es version must v7 or v8")

type ConfigParam func(p *EsConfig)

const (
	V7 ES_VERSION = "v7"
	V8 ES_VERSION = "v8"
)

type EsConfig struct {
	Addresses              []string // A list of Elasticsearch nodes to use.
	Username               string   // Username for HTTP Basic Authentication.
	Password               string   // Password for HTTP Basic Authentication.
	CloudID                string   // Endpoint for the Elastic Service (https://elastic.co/cloud).
	APIKey                 string   // Base64-encoded token for authorization; if set, overrides username/password and service token.
	ServiceToken           string   // Service token for authorization; if set, overrides username/password.
	CertificateFingerprint string   // SHA256 hex fingerprint given by Elasticsearch on first launch.
	CACert                 []byte
	Above                  int
}

//es host
func WithHost(host ...string) ConfigParam {
	return func(p *EsConfig) {
		p.Addresses = host
	}
}

func WithUser(user string) ConfigParam {
	return func(p *EsConfig) {
		p.Username = user
	}
}

func WithPwd(pwd string) ConfigParam {
	return func(p *EsConfig) {
		p.Password = pwd
	}
}

func WithCloudId(cloudId string) ConfigParam {
	return func(p *EsConfig) {
		p.CloudID = cloudId
	}
}
func WithApiKey(apiKey string) ConfigParam {
	return func(p *EsConfig) {
		p.APIKey = apiKey
	}
}

func WithToken(Token string) ConfigParam {
	return func(p *EsConfig) {
		p.ServiceToken = Token
	}
}

func WithFinger(Finger string) ConfigParam {
	return func(p *EsConfig) {
		p.CertificateFingerprint = Finger
	}
}

func WithCACert(CACert []byte) ConfigParam {
	return func(p *EsConfig) {
		p.CACert = CACert
	}
}

//ignore_above
func WithIgnoreAbove(above int) ConfigParam {
	return func(p *EsConfig) {
		p.Above = above
	}
}

type MappingTool interface {
	//if the index not exists u should create mapping
	Create(index string, param interface{}) (string, error)
	//its mean put _mapping
	Put(index string, param interface{}) (string, error)
	//query _mapping
	GetMapping(index string) (string, error)
}

func NewConver(version ES_VERSION, cfg ...ConfigParam) (MappingTool, error) {
	esCfg := &EsConfig{}
	for _, c := range cfg {
		c(esCfg)
	}
	switch version {
	case V7:
		v7cfg := es7.Config{
			Addresses:              esCfg.Addresses,
			Username:               esCfg.Username,
			Password:               esCfg.Password,
			CloudID:                esCfg.CloudID,
			APIKey:                 esCfg.APIKey,
			ServiceToken:           esCfg.ServiceToken,
			CertificateFingerprint: esCfg.CertificateFingerprint,
			CACert:                 esCfg.CACert,
		}
		cli, err := es7.NewClient(v7cfg)
		if err != nil {
			return nil, err
		}
		return client.NewClient(cli, esCfg.Above), nil
	case V8:
		v8cfg := es8.Config{
			Addresses:              esCfg.Addresses,
			Username:               esCfg.Username,
			Password:               esCfg.Password,
			CloudID:                esCfg.CloudID,
			APIKey:                 esCfg.APIKey,
			ServiceToken:           esCfg.ServiceToken,
			CertificateFingerprint: esCfg.CertificateFingerprint,
			CACert:                 esCfg.CACert,
		}
		cli, err := es8.NewClient(v8cfg)
		if err != nil {
			return nil, err
		}
		return client.NewClient(cli, esCfg.Above), nil
	}
	return nil, ErrorVersion
}
