package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/hohice/gin-web/docs"
	. "github.com/hohice/gin-web/pkg/util/log"

	"github.com/hohice/gin-web/server/middleware"
)

// @title ginS
// @version 1.0.0
// @description Gin Web API server starter.

// @contact.name hohice
// @contact.url https://github.com/hohice
// @contact.email hohice@163.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

func (server *Server) InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	//add Probe for readiness and liveness
	r.Use(middleware.ReadinessProbe, middleware.LivenessProbe)
	//enable swagger UI
	r.GET("/swagger/*any", gin.BasicAuth(gin.Accounts{server.Account}), ginSwagger.WrapHandler(swaggerFiles.Handler))

	if server.Debug {
		gin.SetMode(gin.DebugMode)
		Log.SetLevel(logrus.DebugLevel)
		r.Use(gin.LoggerWithWriter(Log.Out))
		r.Use(gin.RecoveryWithWriter(Log.Out))
	} else {
		Log.SetLevel(logrus.InfoLevel)
		//add Prometheus Metric
		p := middleware.NewPrometheus(server.ServiceName)
		p.ReqCntURLLabelMappingFn = middleware.MapURLWithParamsBackToRouteTemplate
		p.Use(r)
		//use with auth
		//p.UseWithAuth(r,gin.BasicAuth(gin.Accounts{server.Account})
	}

	//define api group
	apiv1 := r.Group("/api/v1")

	if !server.Debug {
		apiv1.Use(middleware.SpanFromHeaders("api", middleware.CPsrFunc, false), middleware.InjectToHeaders(false))
	}
	{
		podGroup := apiv1.Group("pod").Use(middleware.JWT())
		{
			podGroup.GET("/:namespace/:pod/shell/:container")
		}
	}

	return r
}
