package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	stdot "github.com/opentracing/opentracing-go"

	_ "github.com/hohice/gin-web/docs"
	. "github.com/hohice/gin-web/pkg/util/log"

	"github.com/hohice/gin-web/router/ex"
	"github.com/hohice/gin-web/router/middleware"
)

// @title Walm
// @version 1.0.0
// @description Warp application lifecycle manager.

// @contact.name bing.han
// @contact.url http://transwarp.io
// @contact.email bing.han@transwarp.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /walm/api/v1

func InitRouter(oauth, runmode bool) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	//r.Use(gin.RecoveryWithWriter(Log.Out))

	if runmode {
		gin.SetMode(gin.DebugMode)
		Log.SetLevel(logrus.DebugLevel)
		r.Use(gin.LoggerWithWriter(Log.Out))
	} else {
		Log.SetLevel(logrus.InfoLevel)
		//add Prometheus Metric
		p := middleware.NewPrometheus("Walm")
		p.Use(r)
	}

	//enable swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//add Probe for readiness and liveness
	r.GET("/readiness", readinessProbe)
	r.GET("/liveness", livenessProbe)

	//define api group
	apiv1 := r.Group("/api/v1")
	if oauth {
		apiv1.Use(middleware.JWT())
	}
	if !runmode && middleware.Tracer != nil {
		//add opentracing
		psr := func(spancontext stdot.SpanContext) stdot.StartSpanOption {
			return stdot.ChildOf(spancontext)
		}
		apiv1.Use(middleware.SpanFromHeaders(middleware.Tracer, "Walm", psr, false), middleware.InjectToHeaders(middleware.Tracer, false))
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
