package cmd

import (
	"fmt"

	"github.com/cocktailcloud/acloud-alarm-collector/application/api"
	"github.com/cocktailcloud/acloud-alarm-collector/application/config"
	"github.com/cocktailcloud/acloud-alarm-collector/application/database"
	"github.com/cocktailcloud/acloud-alarm-collector/application/kafka"
	adapter "github.com/cocktailcloud/acloud-alarm-collector/application/logger"
	"github.com/cocktailcloud/acloud-alarm-collector/application/route"
	"github.com/cocktailcloud/acloud-alarm-collector/application/scheduler"
	"github.com/cocktailcloud/acloud-alarm-collector/application/service"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/server"
)

var gatewayServer *server.Instance

// var mainQueueScheduler scheduler.MainQueueScheduler
var reloadMemoryScheduler scheduler.ReloadMemoryScheduler
var orphanScheduler scheduler.OrphanScheduler
var retryScheduler scheduler.RetryScheduler

func Execute() {
	gatewayServer.Start()
	// mainQueueScheduler.Run()
	reloadMemoryScheduler.Run()
	orphanScheduler.Run()
	retryScheduler.Run()
}

func Close() {
	logger.Info("close process start")
	gatewayServer.DBShutdownIfExists()
	gatewayServer.Shutdown()
	// mainQueueScheduler.Stop()
	reloadMemoryScheduler.Stop()
	orphanScheduler.Stop()
	retryScheduler.Stop()
	logger.Info("close process finished ################")
}

func init() {
	fmt.Println("##### Start Metric Collector  - Cocktail Version 4.8.1.0 #####")

	context := &config.Context{}

	conf, err := config.Load()
	if err != nil {
		panic(err)
	}
	context.Config = conf

	err = adapter.NewLogger(conf)
	if err != nil {
		logger.Fatalf("Could not instantiate log %ss", err.Error())
	}

	// create server instance & initialize the server
	gatewayServer = server.NewInstance(conf.ServerConfig)
	gatewayServer.Init()

	// database connection - how to switch stat, control
	dbAdapter, err := database.NewAdapter(conf.DBConfig)
	if err != nil {
		logger.Error("Could not generate database object: ", err)
		panic(err)
	}
	dbAdapter.OrmMapping()
	context.DBAdapter = dbAdapter

	// Queue 생성 - main, retry
	context.Init()

	service, err := service.New(context)
	if err != nil {
		panic(err)
	}

	api, _ := api.New(conf, service)

	route.SetRoutes(api, gatewayServer)

	go kafka.New(conf.KafkaConfig, service)

	// mainQueueScheduler, err = scheduler.NewMainQueueJob(conf, service)
	// if err != nil {
	// 	panic(err)
	// }
	go service.LoopMainQueueProcess()

	reloadMemoryScheduler, err = scheduler.NewReloadMemoryJob(context, service)
	if err != nil {
		panic(err)
	}

	orphanScheduler, err = scheduler.NewOrphanScheduler(context, service)
	if err != nil {
		panic(err)
	}

	retryScheduler, err = scheduler.NewRetryScheduler(context, service)
	if err != nil {
		panic(err)
	}
}
