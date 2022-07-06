package server

import (
	"github.com/cloudslit/cloudslit/fullnode/pkg/confer"
	"github.com/cloudslit/cloudslit/fullnode/pkg/logger"
	"github.com/cloudslit/cloudslit/fullnode/pkg/mysql"
	"github.com/cloudslit/cloudslit/fullnode/pkg/redis"
	"github.com/cloudslit/cloudslit/fullnode/pkg/web3/eth"
	"github.com/cloudslit/cloudslit/fullnode/pkg/web3/w3s"
	"github.com/urfave/cli"
)

func InitService(c *cli.Context) (err error) {
	if err = confer.Init(c.String("c")); err != nil {
		return
	}
	cfg := confer.GlobalConfig()
	logger.Init(&logger.Config{
		Level:       logger.LogLevel(),
		Filename:    logger.LogFile(),
		SendToFile:  logger.SendLogToFile(),
		Development: confer.ConfigEnvIsDev(),
	})
	if err = redis.Init(&cfg.Redis); err != nil {
		logger.Errorf(nil, "redis init error : %v", err)
		return
	}
	if err = mysql.Init(&cfg.Mysql); err != nil {
		logger.Errorf(nil, "mysql init error : %v", err)
		return
	}
	if err = w3s.Init(&cfg.Web3); err != nil {
		logger.Errorf(nil, "w3s init error : %v", err)
		return
	}
	if err = eth.Init(&cfg.Web3); err != nil {
		logger.Errorf(nil, "eth init error : %v", err)
		return
	}
	if confer.GlobalConfig().Web3.Register == "true" {
		if err = runETH(); err != nil {
			logger.Errorf(nil, "runETH error : %v", err)
			return
		}
	}
	// 判断是否开启P2P
	if confer.GlobalConfig().P2P.Enable {
		return runP2P()
	}
	return
}
