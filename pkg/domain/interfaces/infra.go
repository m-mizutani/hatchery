package interfaces

import (
	"context"
	"io"
	"net/http"

	"github.com/m-mizutani/hatchery/pkg/domain/types"
)

type CloudStorage interface {
	NewObjectWriter(ctx context.Context, bucket types.CSBucket, object types.CSObjectName) io.WriteCloser
	NewObjectReader(ctx context.Context, bucket types.CSBucket, object types.CSObjectName) (io.ReadCloser, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
