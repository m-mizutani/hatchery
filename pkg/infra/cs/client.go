package cs

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/domain/interfaces"
	"github.com/m-mizutani/hatchery/pkg/domain/types"
	"google.golang.org/api/option"
)

type Client struct {
	client *storage.Client
}

var _ interfaces.CloudStorage = (*Client)(nil)

func New(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	c, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create Google Cloud Storage client")
	}

	return &Client{
		client: c,
	}, nil
}

// NewObjectReader implements interfaces.CloudStorage.
func (c *Client) NewObjectReader(ctx context.Context, bucket types.CSBucket, object types.CSObjectName) (io.ReadCloser, error) {
	r, err := c.client.Bucket(string(bucket)).Object(string(object)).NewReader(ctx)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create object reader")
	}

	return r, nil
}

// NewObjectWriter implements interfaces.CloudStorage.
func (c *Client) NewObjectWriter(ctx context.Context, bucket types.CSBucket, object types.CSObjectName) io.WriteCloser {
	return c.client.Bucket(string(bucket)).Object(string(object)).NewWriter(ctx)
}
