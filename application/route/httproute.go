package route

import (
	"github.com/cocktailcloud/acloud-alarm-collector/application/api"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/server"
)

func SetRoutes(api *api.API, server *server.Instance) {
	setHttpRoutes(api, server)
}

func setHttpRoutes(api *api.API, server *server.Instance) {
	v1 := server.HTTPServer.Group("/monitoring-api/v1")
	v1.GET("/check/readiness", api.CheckReadiness)
	v1.GET("/check/liveness", api.CheckLiveness)

	v2 := server.HTTPServer.Group("/collector-api/v2")
	v2.PUT("/alerts", api.ReceiveAlarms)
	v2.POST("/alerts", api.ReceiveAlarms)
}
