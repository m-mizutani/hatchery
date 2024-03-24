package cli

import (
	"github.com/m-mizutani/hatchery/pkg/domain/model"
	"github.com/m-mizutani/hatchery/pkg/infra"
	"github.com/m-mizutani/hatchery/pkg/infra/cs"
	"github.com/m-mizutani/hatchery/pkg/usecase"
	"github.com/m-mizutani/hatchery/pkg/utils"
	"github.com/urfave/cli/v2"
)

func cmdExec(rt *runtime) *cli.Command {
	var (
		actionIDs  cli.StringSlice
		actionTags cli.StringSlice
		allAction  bool
		dryRun     bool
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
			&cli.StringSliceFlag{
				Name:        "tag",
				Aliases:     []string{"t"},
				Usage:       "Action tag",
				EnvVars:     []string{"HATCHERY_EXEC_TAG"},
				Destination: &actionTags,
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
			_, ctx := utils.CtxRequestID(c.Context)

			selector := &model.Selector{
				IDs:  actionIDs.Value(),
				Tags: actionTags.Value(),
				All:  allAction,
			}
			if err := selector.Validate(); err != nil {
				return err
			}

			var options []usecase.ExecuteOption
			if dryRun {
				options = append(options, usecase.WithDryRun())
			}

			csClient, err := cs.New(ctx)
			if err != nil {
				return err
			}

			clients := infra.New(infra.WithCloudStorage(csClient))

			if err := usecase.Execute(ctx, clients, rt.config.Actions, selector, options...); err != nil {
				return err
			}

			return nil
		},
	}
}
