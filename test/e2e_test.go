//+build e2e

package test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hohice/gin-web/router"
)

var confEnvName = "GINS_CONF_PATH"

func Test(t *testing.T) { check.TestingT(t) }

type e2eSuite struct {
	found       bool
	oldEnvValue string
}

var _ = check.Suite(&e2eSuite{})

var routers *gin.Engine

//SetUpSuite
func (e2e *e2eSuite) SetUpSuite() {
	testPath, err := os.Getwd()
	c.Assert(err, check.IsNil)

	/*
		cleanup := resetEnv()
		defer cleanup()
	*/
	e2e.oldEnvValue, e2e.found = os.LookupEnv()

	testConfPath := filepath.Join(testPath, "../pkg/setting/testdata")
	os.Setenv(confEnvName, testConfPath)
	routers = router.InitRouter(false, false)
}

//TearDownSuite
func (e2e *e2eSuite) TearDownSuite() {
	//os.Clearenv()
	if e2e.found {
		os.Setenv(confEnvName, e2e.oldEnvValue)
	} else {
		os.Unsetenv(confEnvName)
	}
}

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

	c.Assert(code, check.Equals, http.StatusOK, "/metrics return none OK!")
}
