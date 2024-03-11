// Code generated from Pkl module `org.github.m_mizutani.hatchery.config`. DO NOT EDIT.
package config

type FalconDataReplicator interface {
	Action

	GetAwsRegion() string

	GetAwsAccessKeyId() string

	GetAwsSecretAccessKey() string

	GetSqsUrl() string

	GetMaxMessages() *int

	GetMaxPulls() *int
}

var _ FalconDataReplicator = (*FalconDataReplicatorImpl)(nil)

type FalconDataReplicatorImpl struct {
	AwsRegion string `pkl:"aws_region"`

	AwsAccessKeyId string `pkl:"aws_access_key_id"`

	AwsSecretAccessKey string `pkl:"aws_secret_access_key"`

	SqsUrl string `pkl:"sqs_url"`

	MaxMessages *int `pkl:"max_messages"`

	MaxPulls *int `pkl:"max_pulls"`

	Id string `pkl:"id"`

	Bucket string `pkl:"bucket"`

	Prefix *string `pkl:"prefix"`
}

func (rcv *FalconDataReplicatorImpl) GetAwsRegion() string {
	return rcv.AwsRegion
}

func (rcv *FalconDataReplicatorImpl) GetAwsAccessKeyId() string {
	return rcv.AwsAccessKeyId
}

func (rcv *FalconDataReplicatorImpl) GetAwsSecretAccessKey() string {
	return rcv.AwsSecretAccessKey
}

func (rcv *FalconDataReplicatorImpl) GetSqsUrl() string {
	return rcv.SqsUrl
}

func (rcv *FalconDataReplicatorImpl) GetMaxMessages() *int {
	return rcv.MaxMessages
}

func (rcv *FalconDataReplicatorImpl) GetMaxPulls() *int {
	return rcv.MaxPulls
}

func (rcv *FalconDataReplicatorImpl) GetId() string {
	return rcv.Id
}

func (rcv *FalconDataReplicatorImpl) GetBucket() string {
	return rcv.Bucket
}

func (rcv *FalconDataReplicatorImpl) GetPrefix() *string {
	return rcv.Prefix
}
