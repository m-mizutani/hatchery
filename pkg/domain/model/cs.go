package model

import (
	"context"
	"fmt"
	"time"

	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/domain/types"
	"github.com/m-mizutani/hatchery/pkg/utils"
)

func LogObjNamePrefix(action config.Action, now time.Time) types.CSObjectName {
	objPrefix := now.Format("logs/2006/01/02/15/")
	if prefix := action.GetPrefix(); prefix != nil {
		objPrefix = *prefix + objPrefix
	}

	return types.CSObjectName(objPrefix)
}

func DefaultLogObjectName(ctx context.Context, action config.Action, now time.Time, seq int) types.CSObjectName {
	reqID, _ := utils.CtxRequestID(ctx)
	objName := types.CSObjectName(
		fmt.Sprintf("%s-%s-%08d.json.gz", now.Format("20060102T150304"), reqID, seq),
	)

	return LogObjNamePrefix(action, now) + objName
}
