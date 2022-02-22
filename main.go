package main

import (
	"wxcrm/pkg/common"
	"github.com/bamboer/log"
	"wxcrm/pkg/router"
)

var (
	Logger *log.Logger
	cfg    *common.Opts
)

func main() {
	cfg = common.NewOpts()
	Logger = common.NewLogger(cfg.LogFile, cfg.LogLevel)
	router.Run(cfg, Logger)
}
