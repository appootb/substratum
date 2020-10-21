package storage

import (
	es6 "github.com/elastic/go-elasticsearch/v6"
	es7 "github.com/elastic/go-elasticsearch/v7"
)

type ElasticOption func(*es6.Config, *es7.Config)

func WithElasticSearchMaxRetry(retries int) ElasticOption {
	return func(cfg6 *es6.Config, cfg7 *es7.Config) {
		if cfg6 != nil {
			cfg6.MaxRetries = retries
		}
		if cfg7 != nil {
			cfg7.MaxRetries = retries
		}
	}
}
