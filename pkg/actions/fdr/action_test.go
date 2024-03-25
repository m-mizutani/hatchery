package fdr_test

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/m-mizutani/gt"
	"github.com/m-mizutani/hatchery/pkg/actions/fdr"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/domain/interfaces"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/infra/cs"
	"github.com/m-mizutani/hatchery/pkg/utils"
)

type mockSQS struct {
	FnDeleteMessage  func(ctx context.Context, input *sqs.DeleteMessageInput, opts ...request.Option) (*sqs.DeleteMessageOutput, error)
	FnReceiveMessage func(ctx context.Context, input *sqs.ReceiveMessageInput, opts ...request.Option) (*sqs.ReceiveMessageOutput, error)
	messages         []*sqs.ReceiveMessageOutput
}

// DeleteMessageWithContext implements interfaces.SQS.
func (m *mockSQS) DeleteMessageWithContext(ctx context.Context, input *sqs.DeleteMessageInput, opts ...request.Option) (*sqs.DeleteMessageOutput, error) {
	return m.FnDeleteMessage(ctx, input, opts...)
}

// ReceiveMessageWithContext implements interfaces.SQS.
func (m *mockSQS) ReceiveMessageWithContext(ctx context.Context, input *sqs.ReceiveMessageInput, opts ...request.Option) (*sqs.ReceiveMessageOutput, error) {
	if m.FnReceiveMessage != nil {
		return m.FnReceiveMessage(ctx, input, opts...)
	}

	if len(m.messages) == 0 {
		return &sqs.ReceiveMessageOutput{
			Messages: []*sqs.Message{},
		}, nil
	}

	msg := m.messages[0]
	m.messages = m.messages[1:]
	return msg, nil
}

var _ interfaces.SQS = &mockSQS{}

type mockS3 struct {
	DataSet [][]byte
}

// GetObjectWithContext implements interfaces.S3.
func (m *mockS3) GetObjectWithContext(ctx context.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error) {
	if len(m.DataSet) == 0 {
		return nil, errors.New("no data")
	}

	data := m.DataSet[0]
	m.DataSet = m.DataSet[1:]
	return &s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader(data)),
	}, nil
}

var _ interfaces.S3 = &mockS3{}

//go:embed testdata/body.json
var bodyJSON string

func TestFalconDataReplicator(t *testing.T) {
	var calledRecv, calledDelete int

	mockCS := cs.NewMock()
	mockSQS := &mockSQS{
		FnDeleteMessage: func(ctx context.Context, input *sqs.DeleteMessageInput, opts ...request.Option) (*sqs.DeleteMessageOutput, error) {
			calledDelete++
			gt.Equal(t, "test-sqs-url", *input.QueueUrl)
			gt.Equal(t, "test-receipt-handle", *input.ReceiptHandle)
			return nil, nil
		},
		FnReceiveMessage: func(ctx context.Context, input *sqs.ReceiveMessageInput, opts ...request.Option) (*sqs.ReceiveMessageOutput, error) {
			calledRecv++
			gt.Equal(t, "test-sqs-url", *input.QueueUrl)

			if calledRecv > 1 {
				return &sqs.ReceiveMessageOutput{
					Messages: []*sqs.Message{},
				}, nil
			}
			return &sqs.ReceiveMessageOutput{
				Messages: []*sqs.Message{
					{
						Body:          &bodyJSON,
						ReceiptHandle: aws.String("test-receipt-handle"),
					},
				},
			}, nil
		},
	}
	mockS3 := &mockS3{
		DataSet: [][]byte{
			[]byte("test-data-1"),
			[]byte("test-data-2"),
		},
	}
	clients := infra.New(
		infra.WithCloudStorage(mockCS),
		infra.WithNewSQS(func(s *session.Session) interfaces.SQS { return mockSQS }),
		infra.WithNewS3(func(s *session.Session) interfaces.S3 { return mockS3 }),
	)

	now := time.Date(2021, 9, 1, 2, 3, 0, 0, time.UTC)
	ctx := utils.CtxWithNow(context.Background(), func() time.Time { return now })
	gt.NoError(t, fdr.Exec(ctx, clients, &config.FalconDataReplicatorImpl{
		AwsRegion:          "us-west-2",
		Bucket:             "test-bucket",
		AwsAccessKeyId:     "test-access-key",
		AwsSecretAccessKey: "test-secret",
		SqsUrl:             "test-sqs-url",
	}))
	gt.V(t, calledDelete).Equal(1)
	gt.V(t, calledRecv).Equal(2)
	gt.A(t, mockCS.Results).Length(2).
		At(0, func(t testing.TB, v *cs.MockResult) {
			gt.Equal(t, v.Bucket, "test-bucket")
			gt.Equal(t, v.Object, "logs/2021/09/01/02/dAnpZeYcYD1J1B00-9f25c8f9/data/C246521D-D19E-43DD-9EB9-4EEE07F53D5A/part-00000.gz")
		}).
		At(1, func(t testing.TB, v *cs.MockResult) {
			gt.Equal(t, v.Bucket, "test-bucket")
			gt.Equal(t, v.Object, "logs/2021/09/01/02/dAnpZeYcYD1J1B00-9f25c8f9/data/C246521D-D19E-43DD-9EB9-4EEE07F53D5A/part-00001.gz")
		})
}
