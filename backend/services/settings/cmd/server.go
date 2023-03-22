package cmd

import (
	"context"
	"net/http"

	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/log"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/pkg/service/rpc"
	"github.com/ONLYOFFICE/onlyoffice-pipedrive/services/settings/web"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

func Server() *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    "starts a new rpc server instance",
		Category: "server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config_path",
				Usage:   "sets custom configuration path",
				Aliases: []string{"config", "conf", "c"},
			},
			&cli.StringFlag{
				Name:    "environment",
				Usage:   "sets servers environment (development, testing, production)",
				Aliases: []string{"env", "e"},
			},
		},
		Action: func(c *cli.Context) error {
			var (
				CONFIG_PATH = c.String("config_path")
				// ENVIRONMENT = c.String("environment")
			)

			fx.New(
				pkg.Bootstrap(CONFIG_PATH),
				fx.Provide(rpc.NewService),
				fx.Provide(web.NewDocserverRPCServer),
				fx.Invoke(func(lifecycle fx.Lifecycle, service micro.Service, repl *http.Server, logger log.Logger) {
					lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							go repl.ListenAndServe()
							go service.Run()
							return nil
						},
						OnStop: func(ctx context.Context) error {
							g, gCtx := errgroup.WithContext(ctx)
							g.Go(func() error {
								return repl.Shutdown(gCtx)
							})
							return g.Wait()
						},
					})
				}),
				fx.NopLogger,
			).Run()

			return nil
		},
	}
}