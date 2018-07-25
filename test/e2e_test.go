//+build e2e
package test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"gopkg.in/check.v1"

	"github.com/hohice/gin-web/router"
)

func Test(t *testing.T) { check.TestingT(t) }

type e2eSuite struct {
}

var _ = check.Suite(&e2eSuite{})

var routers *gin.Engine

//SetUpSuite
func (e2e *e2eSuite) SetUpSuite() {
	routers = router.InitRouter(false, false)
}

//TearDownSuite
//SetUpTest
//TearDownTest

func (e2e *e2eSuite) Test_readinessProbe(c *check.C) {
	uri := "/readiness"
	code, body := Get(uri, routers)

	c.Assert(code, check.Equals, http.StatusOK, "/readiness return none OK!")
}

func (e2e *e2eSuite) Test_livenessProbe(c *check.C) {
	uri := "/liveness"
	code, body := Get(uri, routers)

	c.Assert(code, check.Equals, http.StatusOK, "/liveness return none OK!")
}
