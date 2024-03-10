package config

import (
	"time"
)

func (x Config) LookupAction(id string) Action {
	for _, action := range x.Actions {
		if action.GetId() == id {
			return action
		}
	}

	return nil
}

func ToObjNamePrefix(action Action, now time.Time) string {
	objPrefix := now.Format("logs/2006/01/02/15/20060102T150405")
	if prefix := action.GetPrefix(); prefix != nil {
		objPrefix = *prefix + objPrefix
	}

	return objPrefix
}
