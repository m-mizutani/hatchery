package cs

import (
	"bytes"
	"context"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/domain/interfaces"
	"github.com/m-mizutani/hatchery/pkg/domain/types"
)

type Mock struct {
	NewObjectWriterFn func(ctx context.Context, bucket types.CSBucket, object types.CSObjectName) io.WriteCloser
	Results           []*MockResult
}

var _ interfaces.CloudStorage = &Mock{}

type MockResult struct {
	Body   Writer
	Bucket types.CSBucket
	Object types.CSObjectName
}

type Writer struct {
	bytes.Buffer
	Closed bool
}

func (x *Writer) Close() error {
	x.Closed = true
	return nil
}

func NewMock() *Mock {
	return &Mock{}
}

func (x *Mock) NewObjectWriter(ctx context.Context, bucket types.CSBucket, object types.CSObjectName) io.WriteCloser {
	if x.NewObjectWriterFn != nil {
		return x.NewObjectWriterFn(ctx, bucket, object)
	}

	var result MockResult
	x.Results = append(x.Results, &result)
	result.Bucket = bucket
	result.Object = object
	return &result.Body
}

func (x *Mock) NewObjectReader(ctx context.Context, bucket types.CSBucket, object types.CSObjectName) (io.ReadCloser, error) {
	for _, r := range x.Results {
		if r.Bucket == bucket && r.Object == object {
			return io.NopCloser(bytes.NewReader(r.Body.Bytes())), nil
		}
	}

	return nil, goerr.New("not found")
}
