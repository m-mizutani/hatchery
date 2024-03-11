package fdr

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/domain/interfaces"
	"github.com/m-mizutani/hatchery/pkg/domain/model"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/utils"
)

type fdrMessage struct {
	Bucket     string `json:"bucket"`
	Cid        string `json:"cid"`
	FileCount  int64  `json:"fileCount"`
	Files      []file `json:"files"`
	PathPrefix string `json:"pathPrefix"`
	Timestamp  int64  `json:"timestamp"`
	TotalSize  int64  `json:"totalSize"`
}

type file struct {
	Checksum string `json:"checksum"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
}

type fdrClients struct {
	infra *infra.Clients
	sqs   interfaces.SQS
	s3    interfaces.S3
}

func Exec(ctx context.Context, clients *infra.Clients, req *config.FalconDataReplicatorImpl) error {
	// Create an AWS session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(req.AwsRegion),
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     req.AwsAccessKeyId,
				SecretAccessKey: req.AwsSecretAccessKey,
			},
		}),
	})
	if err != nil {
		return goerr.Wrap(err, "failed to create AWS session").With("req", req)
	}

	// Create AWS service clients
	sqsClient := clients.NewSQS(awsSession)
	s3Client := clients.NewS3(awsSession)
	prefix := config.LogObjNamePrefix(req, utils.CtxNow(ctx))

	// Receive messages from SQS queue
	input := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(req.SqsUrl),
	}
	if req.MaxMessages != nil {
		input.MaxNumberOfMessages = aws.Int64(int64(*req.MaxMessages))
	}

	for i := 0; ; i++ {
		if req.MaxPulls != nil && i >= *req.MaxPulls {
			break
		}

		c := &fdrClients{infra: clients, sqs: sqsClient, s3: s3Client}
		if err := copy(ctx, c, input, model.CSBucket(req.Bucket), prefix); err != nil {
			if err == errNoMoreMessage {
				break
			}
			return err
		}
	}

	return nil
}

var (
	errNoMoreMessage = errors.New("no more message")
)

func copy(ctx context.Context, clients *fdrClients, input *sqs.ReceiveMessageInput, bucket model.CSBucket, prefix model.CSObjectName) error {
	result, err := clients.sqs.ReceiveMessageWithContext(ctx, input)
	if err != nil {
		return goerr.Wrap(err, "failed to receive messages from SQS").With("input", input)
	}
	if len(result.Messages) == 0 {
		return errNoMoreMessage
	}

	// Iterate over received messages
	for _, message := range result.Messages {
		// Get the S3 object key from the message
		var msg fdrMessage
		if err := json.Unmarshal([]byte(*message.Body), &msg); err != nil {
			return goerr.Wrap(err, "failed to unmarshal message").With("message", *message.Body)
		}

		for _, file := range msg.Files {
			// Download the object from S3
			s3Input := &s3.GetObjectInput{
				Bucket: aws.String(msg.Bucket),
				Key:    aws.String(file.Path),
			}
			s3Obj, err := clients.s3.GetObjectWithContext(ctx, s3Input)
			if err != nil {
				return goerr.Wrap(err, "failed to download object from S3").With("msg", msg)
			}
			defer utils.SafeClose(s3Obj.Body)

			csObj := prefix + model.CSObjectName(file.Path)
			w := clients.infra.CloudStorage().NewObjectWriter(ctx, bucket, csObj)

			if _, err := io.Copy(w, s3Obj.Body); err != nil {
				return goerr.Wrap(err, "failed to write object to GCS").With("msg", msg)
			}
			if err := w.Close(); err != nil {
				return goerr.Wrap(err, "failed to close object writer").With("msg", msg)
			}

			utils.CtxLogger(ctx).Info("FDR: object forwarded from S3 to GCS", "s3", s3Input, "gcsObj", csObj)
		}

		// Delete the message from SQS
		_, err = clients.sqs.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      input.QueueUrl,
			ReceiptHandle: message.ReceiptHandle,
		})
		if err != nil {
			return goerr.Wrap(err, "failed to delete message from SQS")
		}
	}

	return nil
}
