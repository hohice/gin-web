package setting

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type settingSuit struct{}

var _ = check.Suite(&settingSuit{})

func (ss *settingSuit) Test_Init(c *check.C) {
	testPath, err := os.Getwd()
	c.Assert(err, check.IsNil)

	cleanup := resetEnv()
	defer cleanup()

	testConfPath := filepath.Join(testPath, "/testdata")
	os.Setenv(confEnvName, testConfPath)
	c.Log(testConfPath)
	Init()
	c.Assert(Config.Http.HTTPPort, check.Equals, 9000)
	c.Assert(Config.Home, check.Equals, os.Getenv("HOME")+"/.ginS")
	c.Assert(Config.Trace.ZipkinUrl, check.Equals, "http://zipkin:9411/api/v1/spans")
}

var envMap = map[string]string{
	confEnvName: configPath,
}

func resetEnv() func() {
	origEnv := os.Environ()

	// ensure any local envvars do not hose us
	for _, e := range envMap {
		os.Unsetenv(e)
	}

	return func() {
		for _, pair := range origEnv {
			kv := strings.SplitN(pair, "=", 2)
			os.Setenv(kv[0], kv[1])
		}
	}
}
