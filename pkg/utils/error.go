package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/domain/model"
)

func HandleError(ctx context.Context, msg string, err error) {
	var evID *sentry.EventID

	// ErrActonFailed is for summarized error message. Each error should be reported to Sentry one by one.
	if !errors.Is(err, model.ErrActonFailed) {
		// Sending error to Sentry
		hub := sentry.CurrentHub().Clone()
		hub.ConfigureScope(func(scope *sentry.Scope) {
			if goErr := goerr.Unwrap(err); goErr != nil {
				for k, v := range goErr.Values() {
					scope.SetExtra(fmt.Sprintf("%v", k), v)
				}
			}
		})
		evID = hub.CaptureException(err)
	}

	CtxLogger(ctx).Error(msg, ErrLog(err), "sentry.EventID", evID)
}
