package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/hohice/gin-web/docs"

	"github.com/hohice/gin-web/server/handler/v1/config"
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
	r.Use(middleware.ReadinessProbe(), middleware.LivenessProbe())
	//enable swagger UI  /swagger/index.html
	r.GET("/swagger/*any",
		gin.BasicAuth(server.Account),
		ginSwagger.WrapHandler(swaggerFiles.Handler))

	if server.Debug {
		gin.SetMode(gin.DebugMode)
		r.Use(gin.LoggerWithWriter(os.Stdout))
		r.Use(gin.RecoveryWithWriter(os.Stderr))
	} else {
		//add Prometheus Metric
		middleware.StartPrometheusProbes(r)
		//use with auth
		//p.UseWithAuth(r,gin.BasicAuth(gin.Accounts(server.Account)
	}

	//define api group
	apiv1 := r.Group("/api/v1")

	if !server.Debug {
		apiv1.Use(middleware.SpanFromHeaders("api", middleware.CPsrFunc, false), middleware.InjectToHeaders(false))
	}
	{
		configGroup := apiv1.Group("config") //.Use(middleware.JWT())
		{
			configGroup.GET("/name/:name/version/:version", config.GetConfig)
			configGroup.DELETE("/name/:name/version/:version", config.DelConfig)
			configGroup.PUT("/", config.ModConfig)
			configGroup.POST("/", config.NewConfig)
		}
		pipeLineGroup := apiv1.Group("test") //.Use(middleware.JWT())
		{
			pipeLineGroup.GET("/start/name/:name/version/:version", config.StartTest)
		}
	}

	return r
}
