package middleware

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/hohice/gin-web/pkg/setting"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

var service string

func init() {
	registerSelf(func(conf *setting.Configs) (error, Closeble) {
		service = conf.Service
		return nil, func() {}
	})
}

func StartPrometheusProbes(engine *gin.Engine) {
	p := ginprometheus.NewPrometheus(service)
	p.ReqCntURLLabelMappingFn = mapURLWithParamsBackToRouteTemplate
	p.Use(engine)
}

func mapURLWithParamsBackToRouteTemplate(c *gin.Context) string {
	url := c.Request.URL.String()
	for _, p := range c.Params {
		re := regexp.MustCompile(fmt.Sprintf(`(^.*?)/\b%s\b(.*$)`, regexp.QuoteMeta(p.Value)))
		url = re.ReplaceAllString(url, fmt.Sprintf(`$1/:%s$2`, p.Key))
	}
	return url
}
