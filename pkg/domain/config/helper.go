package config

import (
	"time"

	"github.com/m-mizutani/hatchery/pkg/domain/model"
)

func (x Config) LookupAction(id string) Action {
	for _, action := range x.Actions {
		if action.GetId() == id {
			return action
		}
	}

	return nil
}

func LogObjNamePrefix(action Action, now time.Time) model.CSObjectName {
	objPrefix := now.Format("logs/2006/01/02/15/04/")
	if prefix := action.GetPrefix(); prefix != nil {
		objPrefix = *prefix + objPrefix
	}

	return model.CSObjectName(objPrefix)
}
