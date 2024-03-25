// Code generated from Pkl module `org.github.m_mizutani.hatchery.config`. DO NOT EDIT.
package config

import "github.com/apple/pkl-go/pkl"

type Slack interface {
	Action

	GetAccessToken() string

	GetDuration() *pkl.Duration

	GetLimit() int

	GetMaxPages() *int
}

var _ Slack = (*SlackImpl)(nil)

type SlackImpl struct {
	AccessToken string `pkl:"access_token"`

	Duration *pkl.Duration `pkl:"duration"`

	Limit int `pkl:"limit"`

	MaxPages *int `pkl:"max_pages"`

	Id string `pkl:"id"`

	Tags *[]string `pkl:"tags"`

	Bucket string `pkl:"bucket"`

	Prefix *string `pkl:"prefix"`
}

func (rcv *SlackImpl) GetAccessToken() string {
	return rcv.AccessToken
}

func (rcv *SlackImpl) GetDuration() *pkl.Duration {
	return rcv.Duration
}

func (rcv *SlackImpl) GetLimit() int {
	return rcv.Limit
}

func (rcv *SlackImpl) GetMaxPages() *int {
	return rcv.MaxPages
}

func (rcv *SlackImpl) GetId() string {
	return rcv.Id
}

func (rcv *SlackImpl) GetTags() *[]string {
	return rcv.Tags
}

func (rcv *SlackImpl) GetBucket() string {
	return rcv.Bucket
}

func (rcv *SlackImpl) GetPrefix() *string {
	return rcv.Prefix
}
