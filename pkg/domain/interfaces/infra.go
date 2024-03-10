package interfaces

import (
	"context"
	"io"
	"net/http"

	"github.com/m-mizutani/hatchery/pkg/domain/model"
)

type CloudStorage interface {
	NewObjectWriter(ctx context.Context, bucket model.CSBucket, object model.CSObjectName) io.WriteCloser
	NewObjectReader(ctx context.Context, bucket model.CSBucket, object model.CSObjectName) (io.ReadCloser, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
