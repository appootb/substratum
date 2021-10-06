package storage

import (
	es6 "github.com/elastic/go-elasticsearch/v6"
	es7 "github.com/elastic/go-elasticsearch/v7"
)

type ElasticOption func(*es6.Config, *es7.Config)

// WithoutRetry disables retrying.
func WithoutRetry() ElasticOption {
	return func(cfg6 *es6.Config, cfg7 *es7.Config) {
		if cfg6 != nil {
			cfg6.DisableRetry = true
		}
		if cfg7 != nil {
			cfg7.DisableRetry = true
		}
	}
}

// WithRetryStatus sets status codes for retry.
// Default: 502, 503, 504.
func WithRetryStatus(status ...int) ElasticOption {
	return func(cfg6 *es6.Config, cfg7 *es7.Config) {
		if cfg6 != nil {
			cfg6.RetryOnStatus = status
		}
		if cfg7 != nil {
			cfg7.RetryOnStatus = status
		}
	}
}

// WithRetryOnTimeout enables retry on timeout.
func WithRetryOnTimeout() ElasticOption {
	return func(cfg6 *es6.Config, cfg7 *es7.Config) {
		if cfg6 != nil {
			cfg6.EnableRetryOnTimeout = true
		}
		if cfg7 != nil {
			cfg7.EnableRetryOnTimeout = true
		}
	}
}

// WithElasticSearchMaxRetry sets retry times.
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
