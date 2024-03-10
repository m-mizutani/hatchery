package infra

import (
	"net/http"

	"github.com/m-mizutani/hatchery/pkg/domain/interfaces"
)

type Clients struct {
	cs   interfaces.CloudStorage
	http interfaces.HTTPClient
}

type Option func(*Clients)

func New(opts ...Option) *Clients {
	c := &Clients{
		http: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Clients) CloudStorage() interfaces.CloudStorage {
	return c.cs
}

func WithCloudStorage(cs interfaces.CloudStorage) Option {
	return func(c *Clients) {
		c.cs = cs
	}
}

func (c *Clients) HTTPClient() interfaces.HTTPClient {
	return c.http
}

func WithHTTPClient(http interfaces.HTTPClient) Option {
	return func(c *Clients) {
		c.http = http
	}
}
