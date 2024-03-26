package usecase

import (
	"context"
	"log/slog"
	"sync"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/actions/fdr"
	"github.com/m-mizutani/hatchery/pkg/actions/one_password"
	"github.com/m-mizutani/hatchery/pkg/actions/slack"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/domain/model"
	"github.com/m-mizutani/hatchery/pkg/domain/types"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/utils"
)

type executeConfig struct {
	dryRun bool
	execFn func(context.Context, *infra.Clients, config.Action) error
}

type ExecuteOption func(*executeConfig)

func WithDryRun() ExecuteOption {
	return func(c *executeConfig) {
		c.dryRun = true
	}
}

// WithExecFn is an option to specify a function to execute an action. This is used for testing.
func WithExecFn(fn func(context.Context, *infra.Clients, config.Action) error) ExecuteOption {
	return func(c *executeConfig) {
		c.execFn = fn
	}
}

func Execute(ctx context.Context, clients *infra.Clients, actions []config.Action, selector *model.Selector, options ...ExecuteOption) error {
	_, ctx = utils.CtxRequestID(ctx)
	cfg := executeConfig{
		execFn: executeAction,
	}
	for _, opt := range options {
		opt(&cfg)
	}

	var executable []config.Action
	for _, action := range actions {
		if !selector.Contains(action) {
			continue
		}

		executable = append(executable, action)
	}

	var attrs []any
	for _, action := range executable {
		attrs = append(attrs, actionToAttr(action))
	}
	utils.CtxLogger(ctx).Info("Start execution", slog.Group("actions", attrs...))

	wg := sync.WaitGroup{}
	errCh := make(chan error, len(executable))

	for _, action := range executable {
		wg.Add(1)
		go func(action config.Action) {
			defer wg.Done()

			attr := actionToAttr(action)
			if cfg.dryRun {
				utils.CtxLogger(ctx).Info("Dry run", attr)
				return
			}

			utils.CtxLogger(ctx).Info("Start action", attr)
			if err := cfg.execFn(ctx, clients, action); err != nil {
				utils.HandleError(ctx, "failed to execute action", err)
				errCh <- err
			}
		}(action)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		// This is a case that multiple actions are executed and some of them are failed. This error will not be reported to Sentry, just logging.
		return goerr.Wrap(types.ErrActonFailed, "failed to execute actions").With("errors", errs)
	}

	return nil
}

func executeAction(ctx context.Context, clients *infra.Clients, action config.Action) error {
	switch v := action.(type) {
	case *config.OnePasswordImpl:
		return one_password.Exec(ctx, clients, v)
	case *config.FalconDataReplicatorImpl:
		return fdr.Exec(ctx, clients, v)
	case *config.SlackImpl:
		return slack.Exec(ctx, clients, v)
	default:
		return goerr.Wrap(types.ErrAssertFailed, "unknown action type").With("action", action)
	}
}

func actionToAttr(action config.Action) slog.Attr {
	switch v := action.(type) {
	case *config.OnePasswordImpl:
		return slog.Group(action.GetId(),
			slog.String("type", "OnePassword"),
			slog.Any("config", *v),
		)

	case *config.FalconDataReplicatorImpl:
		return slog.Group(action.GetId(),
			slog.String("type", "FalconDataReplicator"),
			slog.Any("config", *v),
		)

	default:
		return slog.Group(action.GetId(),
			slog.String("type", "unknown"),
			slog.Any("config", action),
		)
	}
}
