package setting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/check.v1"
)

type settingSuit struct{}

var _ = check.Suite(&settingSuit{})

func (ss *settingSuit) Test(c *check.C) {
	testPath, err := os.Getwd()
	c.Assert(err, check.IsNil)

	cleanup := resetEnv()
	defer cleanup()

	testConfPath := filepath.Join(testPath, "/conf")
	fmt.Println(testConfPath)
	os.Setenv(confEnvName, testConfPath)

	c.Assert(Config.Http.HTTPPort, check.Equals, 9000)
	c.Assert(Config.Home, check.Equals, os.Getenv("HOME")+".walm")
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
