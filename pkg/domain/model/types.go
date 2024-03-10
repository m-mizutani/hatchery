package model

import "github.com/google/uuid"

type RequestID string

func NewRequestID() RequestID {
	return RequestID(uuid.NewString())
}

type CSBucket string
type CSObjectName string

type OnePasswordAPIToken string

func (t OnePasswordAPIToken) Bearer() string {
	return "Bearer " + string(t)
}
