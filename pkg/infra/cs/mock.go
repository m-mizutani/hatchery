package cs

import (
	"bytes"
	"context"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/domain/interfaces"
	"github.com/m-mizutani/hatchery/pkg/domain/model"
)

type Mock struct {
	data map[string]*mockWriter
}

var _ interfaces.CloudStorage = (*Mock)(nil)

func NewMock() *Mock {
	return &Mock{
		data: map[string]*mockWriter{},
	}
}

type mockWriter struct {
	bytes.Buffer
}

func (m *mockWriter) Close() error {
	return nil
}

// NewObjectReader implements interfaces.CloudStorage.
func (m *Mock) NewObjectReader(ctx context.Context, bucket model.CSBucket, object model.CSObjectName) (io.ReadCloser, error) {
	buf, ok := m.data[string(bucket)+"/"+string(object)]
	if !ok {
		return nil, goerr.New("no such object")
	}
	return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// NewObjectWriter implements interfaces.CloudStorage.
func (m *Mock) NewObjectWriter(ctx context.Context, bucket model.CSBucket, object model.CSObjectName) io.WriteCloser {
	buf := &mockWriter{}
	m.data[string(bucket)+"/"+string(object)] = buf
	return buf
}
