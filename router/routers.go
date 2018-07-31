package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/hohice/gin-web/docs"
	. "github.com/hohice/gin-web/pkg/util/log"

	"github.com/hohice/gin-web/router/ex"
	"github.com/hohice/gin-web/router/middleware"
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

func InitRouter(oauth, Debug bool) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	if Debug {
		gin.SetMode(gin.DebugMode)
		Log.SetLevel(logrus.DebugLevel)
		r.Use(gin.LoggerWithWriter(Log.Out))
		r.Use(gin.RecoveryWithWriter(Log.Out))
	} else {
		Log.SetLevel(logrus.InfoLevel)
		//add Prometheus Metric
		p := middleware.NewPrometheus("Walm")
		p.Use(r)
	}

	//enable swagger UI
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"swag": "password",
	}))
	authorized.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//add Probe for readiness and liveness
	r.GET("/readiness", readinessProbe)
	r.GET("/liveness", livenessProbe)

	//define api group
	apiv1 := r.Group("/api/v1")
	if oauth {
		apiv1.Use(middleware.JWT())
	}
	if !Debug && middleware.Tracer != nil {
		//add opentracing
		apiv1.Use(middleware.SpanFromHeaders(middleware.Tracer, "ginS", middleware.CPsr, false), middleware.InjectToHeaders(middleware.Tracer, false))
	}
	{
		podGroup := apiv1.Group("pod")
		{
			podGroup.GET("/:namespace/:pod/shell/:container")
		}
	}

	return r
}

func readinessProbe(c *gin.Context) {
	c.JSON(ex.ReturnOK())
}

func livenessProbe(c *gin.Context) {
	c.JSON(ex.ReturnOK())
}
