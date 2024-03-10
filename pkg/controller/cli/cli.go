package cli

import (
	"context"

	"github.com/m-mizutani/hatchery/pkg/controller/cli/flags"
	"github.com/m-mizutani/hatchery/pkg/domain/config"
	"github.com/m-mizutani/hatchery/pkg/utils"
	"github.com/urfave/cli/v2"
)

type runtime struct {
	config config.Config
}

func Run(argv []string) error {
	var (
		rt         runtime
		configPath string
		logger     flags.Logger
		sentry     flags.Sentry
	)

	app := cli.App{
		Name:  "hatchery",
		Usage: "Hatchery is a tool to import SaaS data and logs into object storage",
		Flags: mergeFlags([]cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Path to configuration file",
				Destination: &configPath,
				Required:    true,
			},
		}, logger.Flags(), sentry.Flags()),

		Before: func(ctx *cli.Context) error {
			cfg, err := config.LoadFromPath(ctx.Context, configPath)
			if err != nil {
				return err
			}
			rt.config = *cfg

			logger, err := logger.Configure()
			if err != nil {
				return err
			}
			utils.SetLogger(logger)

			return nil
		},
		Commands: []*cli.Command{
			cmdExec(&rt),
		},
	}

	if err := app.Run(argv); err != nil {
		ctx := context.Background()
		utils.HandleError(ctx, "cli failed", err)
		return err
	}

	return nil
}

func mergeFlags(flagsList ...[]cli.Flag) []cli.Flag {
	var ret []cli.Flag
	for _, flags := range flagsList {
		ret = append(ret, flags...)
	}
	return ret
}
