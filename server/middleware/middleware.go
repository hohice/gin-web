package middleware

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/hohice/gin-web/pkg/setting"
	"github.com/hohice/gin-web/server/ex"
)

func Init() error {
	conf := &setting.Config
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
func MapURLWithParamsBackToRouteTemplate(c *gin.Context) string {
	url := c.Request.URL.String()
	for _, p := range c.Params {
		re := regexp.MustCompile(fmt.Sprintf(`(^.*?)/\b%s\b(.*$)`, regexp.QuoteMeta(p.Value)))
		url = re.ReplaceAllString(url, fmt.Sprintf(`$1/:%s$2`, p.Key))
	}
	return url
}

//add Probe for readiness and liveness
var (
	ReadinessPath = "/readiness"
	LivenessProbe = "/liveness"
)

func ReadinessProbe() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.String() == p.ReadinessPath {
			c.JSON(ex.ReturnOK())
			return
		}
	}
}

func LivenessProbe()gin.HandlerFunc {
	func LivenessProbe(c *gin.Context) {
		if c.Request.URL.String() == p.LivenessProbe {
			c.JSON(ex.ReturnOK())
			return
		}
	}
}

