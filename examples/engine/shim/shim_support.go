package main

import (
	"os"

	_ "github.com/qingcloudhx/core/data/expression/script"
	"github.com/qingcloudhx/core/engine"
	"github.com/qingcloudhx/core/support/log"
)

var (
	cfgJson       string
	cfgCompressed bool
)

func init() {
	log.SetLogLevel(log.RootLogger(), log.ErrorLevel)

	cfg, err := engine.LoadAppConfig(cfgJson, cfgCompressed)
	if err != nil {
		log.RootLogger().Errorf("Failed to create engine: %s", err.Error())
		os.Exit(1)
	}

	_, err = engine.New(cfg, engine.DirectRunner)
	if err != nil {
		log.RootLogger().Errorf("Failed to create engine: %s", err.Error())
		os.Exit(1)
	}
}
