package escli

import (
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
)

const (
	DEFAULT_HEALTH_CHECK_INTERVAL = 10
	DEFAULT_MAX_RETRIES           = 5
)

var (
	Factory EsFactory
)

// EsFactory elastic factory
type EsFactory struct {
	HealthCheckInterval int
	MaxRetries          int
	GZip                bool
	Sniff               bool
}

// EsClient elastic client
type EsClient struct {
	*elastic.Client
}

// CreateEsCli create elastic client
func (p *EsFactory) CreateEsCli(user, pwd string, adders []string) (*EsClient, error) {
	if len(adders) == 0 {
		return nil, fmt.Errorf("addr is nil")
	}

	if p.HealthCheckInterval == 0 {
		p.HealthCheckInterval = DEFAULT_HEALTH_CHECK_INTERVAL
	}
	if p.MaxRetries == 0 {
		p.MaxRetries = DEFAULT_MAX_RETRIES
	}

	es, err := elastic.NewClient(
		elastic.SetURL(adders...),
		elastic.SetBasicAuth(user, pwd),
		elastic.SetGzip(p.GZip),
		elastic.SetSniff(p.Sniff),
		elastic.SetHealthcheckInterval(time.Duration(p.HealthCheckInterval)*time.Second),
		elastic.SetMaxRetries(p.MaxRetries),
	)
	if err != nil {
		return nil, err
	}
	return &EsClient{es}, nil
}
