package cli

import (
	"context"
	"github.com/cloudslit/cloudslit/provider/cmd/server/cli/up"
	"github.com/cloudslit/cloudslit/provider/internal"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewServerCmd(ctx context.Context) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "up",
			Usage: "Up server provider",
			Action: func(c *cli.Context) error {
				return NewUpAction(ctx, c.String("conf"))
			},
		},
	}
}

func NewUpAction(ctx context.Context, configPath string) error {
	handle := func(ctx context.Context) (func(), error) {
		initCleanFunc, err := internal.Init(ctx,
			internal.SetConfigFile(configPath),
		)
		if err != nil {
			return nil, err
		}
		go up.NewUp().RunUp(ctx)
		return func() {
			initCleanFunc()
		}, nil
	}
	return Run(ctx, handle)
}

func Run(ctx context.Context, f func(ctx context.Context) (func(), error)) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := f(ctx)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("received signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.WithContext(ctx).Infof("shutdown!")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}