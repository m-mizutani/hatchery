// Code generated from Pkl module `org.github.m_mizutani.hatchery.config`. DO NOT EDIT.
package config

import "github.com/apple/pkl-go/pkl"

type OnePassword interface {
	Action

	GetApiToken() string

	GetDuration() *pkl.Duration

	GetLimit() int
}

var _ OnePassword = (*OnePasswordImpl)(nil)

type OnePasswordImpl struct {
	ApiToken string `pkl:"api_token"`

	Duration *pkl.Duration `pkl:"duration"`

	Limit int `pkl:"limit"`

	Id string `pkl:"id"`

	Bucket string `pkl:"bucket"`

	Prefix *string `pkl:"prefix"`
}

func (rcv *OnePasswordImpl) GetApiToken() string {
	return rcv.ApiToken
}

func (rcv *OnePasswordImpl) GetDuration() *pkl.Duration {
	return rcv.Duration
}

func (rcv *OnePasswordImpl) GetLimit() int {
	return rcv.Limit
}

func (rcv *OnePasswordImpl) GetId() string {
	return rcv.Id
}

func (rcv *OnePasswordImpl) GetBucket() string {
	return rcv.Bucket
}

func (rcv *OnePasswordImpl) GetPrefix() *string {
	return rcv.Prefix
}
