package main

import (
	"context"
	"os"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/opensourceways/repo-file-cache/routers"
	"github.com/opensourceways/robot-gitee-plugin-lib/interrupts"

	"github.com/opensourceways/repo-file-cache/config"
	"github.com/opensourceways/repo-file-cache/dbmodels"
	"github.com/opensourceways/repo-file-cache/mongodb"
)

func main() {
	configFile, err := beego.AppConfig.String("appconf")
	if err != nil {
		fatal(err)
	}

	if err := config.InitAppConfig(configFile); err != nil {
		fatal(err)
	}
	cfg := config.AppConfig

	clears := []func(){}
	setClear := func(f func()) {
		clears = append(clears, f)
	}
	clear := func() {
		for _, f := range clears {
			f()
		}
	}

	if err := startMongoService(&cfg.Mongodb, setClear); err != nil {
		clear()
		fatal(err)
	}

	run(clear)
}

func fatal(err error) {
	logs.Error(err)
	os.Exit(1)
}

func startMongoService(cfg *config.MongodbConfig, setClear func(func())) error {
	c, err := mongodb.Initialize(cfg)
	if err != nil {
		return err
	}
	dbmodels.RegisterDB(c)

	setClear(func() {
		logs.Info("closing mongodb ...")
		if err := dbmodels.GetDB().Close(); err != nil {
			logs.Error(err)
		}
	})
	return nil
}

func run(clear func()) {
	defer interrupts.WaitForGracefulShutdown()

	interrupts.OnInterrupt(func() {
		shutdown()
		if clear != nil {
			clear()
		}
	})

	beego.Run()
}

func shutdown() {
	logs.Info("server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := beego.BeeApp.Server.Shutdown(ctx); err != nil {
		logs.Error("error to shut down server, err:", err.Error())
	}
	cancel()
}
