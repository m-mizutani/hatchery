package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
)

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

func LogObjNamePrefix(action config.Action, now time.Time) CSObjectName {
	objPrefix := now.Format("logs/2006/01/02/15/")
	if prefix := action.GetPrefix(); prefix != nil {
		objPrefix = *prefix + objPrefix
	}

	return CSObjectName(objPrefix)
}
