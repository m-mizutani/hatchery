package infra

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/m-mizutani/hatchery/pkg/domain/interfaces"
)

type Clients struct {
	cs     interfaces.CloudStorage
	http   interfaces.HTTPClient
	newS3  interfaces.NewS3
	newSQS interfaces.NewSQS
}

type Option func(*Clients)

func New(opts ...Option) *Clients {
	c := &Clients{
		http:   http.DefaultClient,
		newS3:  interfaces.DefaultNewS3,
		newSQS: interfaces.DefaultNewSQS,
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

func (c *Clients) NewS3(s *session.Session) interfaces.S3 {
	return c.newS3(s)
}

func WithNewS3(newS3 interfaces.NewS3) Option {
	return func(c *Clients) {
		c.newS3 = newS3
	}
}

func (c *Clients) NewSQS(s *session.Session) interfaces.SQS {
	return c.newSQS(s)
}

func WithNewSQS(newSQS interfaces.NewSQS) Option {
	return func(c *Clients) {
		c.newSQS = newSQS
	}
}
