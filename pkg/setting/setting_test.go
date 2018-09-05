package setting

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type settingSuit struct {
	cleanup func()
}

var _ = check.Suite(&settingSuit{
	cleanup: func() func() {
		testPath, _ := os.Getwd()
		cleanup := resetEnv()

		testConfPath := filepath.Join(testPath, "/testdata")
		os.Setenv(confEnvName, testConfPath)
		return cleanup
	}(),
})

/*
func (ss *settingSuit) SetUpSuite(c *check.C) {
	testPath, _ := os.Getwd()
	ss.cleanup = resetEnv()

	testConfPath := filepath.Join(testPath, "/testdata")
	os.Setenv(confEnvName, testConfPath)
}
*/

func (ss *settingSuit) TearDownSuite(c *check.C) {
	ss.cleanup()
}

func (ss *settingSuit) Test_ReadConfigPath(c *check.C) {
	vp := viper.New()
	ReadConfigPath(vp)
	c.Assert(Config.Http.HTTPPort, check.Equals, 9000)
	c.Assert(Config.Home, check.Equals, os.Getenv("HOME")+"/.ginS")
	c.Assert(Config.Trace.ZipkinUrl, check.Equals, "http://zipkin:9411/api/v1/spans")
	c.Assert(Config.Limit.AddrMap["server1"], check.Equals, "host1:port1")
	c.Assert(Config.Circuit["default"].SleepWindow, check.Equals, 5000)
	c.Assert(Config.Circuit["url1"].SleepWindow, check.Equals, 1000)
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
