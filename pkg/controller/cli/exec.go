package cli

import (
	"context"
	"log/slog"
	"sync"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/hatchery/pkg/actions/fdr"
	"github.com/m-mizutani/hatchery/pkg/actions/one_password"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/domain/model"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/infra/cs"
	"github.com/m-mizutani/hatchery/pkg/utils"
	"github.com/urfave/cli/v2"
)

func cmdExec(rt *runtime) *cli.Command {
	var (
		actionIDs cli.StringSlice
		allAction bool
		dryRun    bool
	)

	return &cli.Command{
		Name:      "exec",
		Aliases:   []string{"e"},
		Usage:     "Execute actions",
		UsageText: `hatchery [global options] exec [command options]`,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "id",
				Aliases:     []string{"i"},
				Usage:       "Action ID",
				EnvVars:     []string{"HATCHERY_EXEC_ID"},
				Destination: &actionIDs,
			},
			&cli.BoolFlag{
				Name:        "all",
				Aliases:     []string{"a"},
				Usage:       "Execute all actions",
				EnvVars:     []string{"HATCHERY_EXEC_ALL"},
				Destination: &allAction,
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Aliases:     []string{"d"},
				Usage:       "Dry run",
				EnvVars:     []string{"HATCHERY_EXEC_DRY_RUN"},
				Destination: &dryRun,
			},
		},
		Action: func(c *cli.Context) error {
			reqID, ctx := utils.CtxRequestID(c.Context)
			ctx = utils.CtxLoggerWith(ctx, slog.Any("request_id", reqID))

			if allAction && len(actionIDs.Value()) > 0 {
				return goerr.Wrap(model.ErrInvalidOption, "both --all and --id are specified")
			}

			ids := actionIDs.Value()
			if allAction {
				for _, action := range rt.config.Actions {
					ids = append(ids, action.GetId())
				}
			} else if len(ids) == 0 {
				return goerr.Wrap(model.ErrInvalidOption, "either --all or --id is required")
			}

			if dryRun {
				return execDryRun(rt, ids)
			}

			csClient, err := cs.New(ctx)
			if err != nil {
				return err
			}

			clients := infra.New(infra.WithCloudStorage(csClient))

			errCh := make(chan error, len(ids))
			var wg sync.WaitGroup

			for _, id := range ids {
				action := rt.config.LookupAction(id)
				if action == nil {
					return goerr.Wrap(model.ErrInvalidOption, "action not found", "id", id)
				}

				wg.Add(1)
				go func(action config.Action) {
					defer wg.Done()
					if err := execute(ctx, clients, action); err != nil {
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
				return goerr.Wrap(model.ErrActonFailed, "failed to execute actions").With("errors", errs)
			}

			return nil
		},
	}
}

func execute(ctx context.Context, clients *infra.Clients, action config.Action) error {
	switch v := action.(type) {
	case *config.OnePasswordImpl:
		return one_password.Exec(ctx, clients, v)
	case *config.FalconDataReplicatorImpl:
		return fdr.Exec(ctx, clients, v)
	default:
		return goerr.Wrap(model.ErrAssertFailed, "unknown action type").With("action", action)
	}
}

func execDryRun(rt *runtime, actionIDs []string) error {
	var attrs []slog.Attr
	for _, id := range actionIDs {
		action := rt.config.LookupAction(id)
		if action == nil {
			return goerr.Wrap(model.ErrInvalidOption, "action not found", "id", id)
		}
		attrs = append(attrs, actionToAttr(action))
	}

	utils.Logger().Info("DryRun", "actions", attrs)
	return nil
}

func actionToAttr(action config.Action) slog.Attr {
	switch v := action.(type) {
	case *config.OnePasswordImpl:
		return slog.Any(v.Id, *v)
	case *config.FalconDataReplicatorImpl:
		return slog.Any(v.Id, *v)
	default:
		return slog.Any(action.GetId(), action)
	}
}
