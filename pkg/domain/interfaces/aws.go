package interfaces

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type NewSQS func(*session.Session) SQS
type NewS3 func(*session.Session) S3

func DefaultNewSQS(s *session.Session) SQS {
	return sqs.New(s)
}

func DefaultNewS3(s *session.Session) S3 {
	return s3.New(s)
}

type SQS interface {
	ReceiveMessageWithContext(ctx context.Context, input *sqs.ReceiveMessageInput, opts ...request.Option) (*sqs.ReceiveMessageOutput, error)
	DeleteMessageWithContext(ctx context.Context, input *sqs.DeleteMessageInput, opts ...request.Option) (*sqs.DeleteMessageOutput, error)
}

type S3 interface {
	GetObjectWithContext(ctx context.Context, input *s3.GetObjectInput, opts ...request.Option) (*s3.GetObjectOutput, error)
}
