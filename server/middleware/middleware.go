package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hohice/gin-web/pkg/setting"
	"github.com/hohice/gin-web/server/ex"
)

func Init(conf *setting.Configs) error {
	for _, fn := range Initlist {
		if err, closeFn := fn(conf); err != nil {
			Close()
			return err
		} else {
			Closelist = append(Closelist, closeFn)
		}
	}
	return nil
}

func Close() {
	for _, fn := range Closelist {
		fn()
	}
	Closelist = []Closeble{}
}

type Closeble func()
type Register func(conf *setting.Configs) (error, Closeble)

var Initlist []Register
var Closelist []Closeble

func registerSelf(regfunc Register) {
	Initlist = append(Initlist, regfunc)
}

/*
mapURLWithParamsBackToRouteTemplate is a valid ginprometheus ReqCntURLLabelMappingFn.
For every route containing parameters (e.g. `/charts/:filename`, `/api/charts/:name/:version`, etc)
the actual parameter values will be replaced by their name, to minimize the cardinality of the
`chartmuseum_requests_total{url=..}` Prometheus counter.
*/

//add Probe for readiness and liveness
var (
	ReadinessPath = "/readiness"
	LivenessPath  = "/liveness"
)

func ReadinessProbe() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.String() == ReadinessPath {
			c.JSON(ex.ReturnOK())
			return
		}
	}
}

func LivenessProbe() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.String() == LivenessPath {
			c.JSON(ex.ReturnOK())
			return
		}
	}
}
