//+build e2e

package test

import (
	"net/http"
	"os"
	"path/filepath"
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
	testPath, err := os.Getwd()
	c.Assert(err, check.IsNil)

	cleanup := resetEnv()
	defer cleanup()

	testConfPath := filepath.Join(testPath, "../pkg/setting/testdata")
	os.Setenv(confEnvName, testConfPath)
	routers = router.InitRouter(false, false)
}

//TearDownSuite
//SetUpTest
//TearDownTest

func (e2e *e2eSuite) Test_readinessProbe(c *check.C) {
	uri := "/readiness"
	code, _ := Get(uri, routers)

	c.Assert(code, check.Equals, http.StatusOK, "/readiness return none OK!")
}

func (e2e *e2eSuite) Test_livenessProbe(c *check.C) {
	uri := "/liveness"
	code, _ := Get(uri, routers)

	c.Assert(code, check.Equals, http.StatusOK, "/liveness return none OK!")
}

func (e2e *e2eSuite) Test_Prometheus(c *check.C) {
	uri := "/metrics"
	code, _ := Get(uri, routers)

	c.Assert(code, check.Equals, http.StatusOK, "/liveness return none OK!")
}
