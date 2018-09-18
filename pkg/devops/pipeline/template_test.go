package pipeline

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var confEnvName = "GINS_CONF_PATH"
var configPath = "/etc/ginS/"
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

var cleanup func() = func() func() {
	testPath, _ := os.Getwd()
	cleanup := resetEnv()

	testConfPath := filepath.Join(testPath, "/Users/hohice/MyWorkspace/go/src/github.com/hohice/gin-web/pkg/setting/testdata")
	os.Setenv(confEnvName, testConfPath)
	return cleanup
}()

func TestGetJsonInstance(t *testing.T) {
	defer cleanup()

	type args struct {
		pipeline *PipeLine
	}
	tests := []struct {
		name    string
		args    args
		want    *PipeLine
		wantErr bool
	}{
		{
			name: "suite_1",
			args: args{
				pipeline: &PipeLine{
					Name: "test_1",
				},
			},
			want: &PipeLine{
				Name: "test_1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := GetJsonInstance(tt.args.pipeline)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetJsonInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJsonInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}
