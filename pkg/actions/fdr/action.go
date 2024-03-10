package fdr

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/utils"
)

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

	// Create an SQS client
	sqsClient := sqs.New(awsSession)

	// Receive messages from SQS queue
	result, err := sqsClient.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(req.SqsUrl),
		MaxNumberOfMessages: aws.Int64(10),
	})
	if err != nil {
		return goerr.Wrap(err, "failed to receive messages from SQS").With("req", req)
	}

	utils.CtxLogger(ctx).Info("FDR: received messages from SQS", "count", len(result.Messages))
	/*
		prefix := config.ToObjNamePrefix(req, utils.CtxNow(ctx))

		// Iterate over received messages
		for _, message := range result.Messages {
			// Get the S3 object key from the message
			s3ObjectKey := *message.Body

			// Download the object from S3
			s3Client := s3.New(awsSession)
			s3ObjectOutput, err := s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
				Bucket: aws.String(req.S3Bucket),
				Key:    aws.String(s3ObjectKey),
			})
			if err != nil {
				log.Printf("failed to download object from S3: %v", err)
				continue
			}

			w := clients.CloudStorage().NewObjectWriter(ctx, model.CSBucket(req.GetBucket()), objName)
			// Upload the object to Google Cloud Storage
			gcsClient, err := storage.NewService(ctx, option.WithCredentialsFile(req.GCSCredentialsFile))
			if err != nil {
				log.Printf("failed to create GCS client: %v", err)
				continue
			}

			objectName := uuid.New().String() // Generate a unique object name
			gcsObject := &storage.Object{
				Name:     objectName,
				Bucket:   req.Bucket,
				Metadata: make(map[string]string),
			}

			// Copy the object data to GCS
			_, err = gcsClient.Objects.Insert(req.Bucket, gcsObject).Media(s3ObjectOutput.Body).Do()
			if err != nil {
				log.Printf("failed to upload object to GCS: %v", err)
				continue
			}

			// Delete the message from SQS
			_, err = sqsClient.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(req.SqsUrl),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				return goerr.Wrap(err, "failed to delete message from SQS").With("req", req)
			}

			utils.CtxLogger(ctx).Info("FDR: object forwarded from S3 to GCS", "s3ObjectKey", s3ObjectKey, "gcsObjectName", objectName)
		}
	*/
	return nil
}
