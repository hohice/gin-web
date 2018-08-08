package setting

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/hohice/gin-web/pkg/util/log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var confEnvName = "GINS_CONF_PATH"

var configPath = "/etc/ginS/"

var DefaultWalmHome = filepath.Join(HomeDir(), ".ginS")

var Config Configs

var regNotifyChannel []*chan bool

type Configs struct {
	Service string `mapstructure:"service"`
	Home    string `mapstructure:"home"`
	Debug   bool   `mapstructure:"debug"`

	Http struct {
		HTTPPort     int           `mapstructure:"port"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
	} `mapstructure:"http"`

	Secret struct {
		Tls       bool   `mapstructure:"tls"`
		TlsVerify bool   `mapstructure:"tls-verify"`
		TlsKey    string `mapstructure:"tls-key"`
		TlsCert   string `mapstructure:"tls-cert"`
		TlsCaCert string `mapstructure:"tls-ca-cert"`
	} `mapstructure:"secret"`

	Helm struct {
		TillerConnectionTimeout time.Duration `mapstructure:"tiller_time_out"`
		TillerHost              string        `mapstructure:"tillerHost"`
	} `mapstructure:"helm"`

	Repo struct {
		Name string `mapstructure:"name"`
		URL  string `mapstructure:"url"`
	} `mapstructure:"repo"`

	Kube struct {
		MasterHost string `mapstructure:"master_host"`
		Context    string `mapstructure:"config"`
		Config     string `mapstructure:"context"`
	} `mapstructure:"kube"`

	Trace struct {
		ZipkinUrl string `mapstructure:"zipkin_url"`
	} `mapstructure:"trace"`

	Auth struct {
		Enable      bool   `mapstructure:"enable"`
		JwtSecret   string `mapstructure:"jwtsecret"`
		TokenLookup string `mapstructure:"tokenlookup"`
		AuthScheme  string `mapstructure:"authscheme"`
	} `mapstructure:"auth"`
}

// Init sets values from the environment.
func Init() {
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigName("conf")
	vp.SetDefault("home", DefaultWalmHome)
	vp.SetDefault("http.port", 8000)
	if str, have := getEnv(); have {
		configPath = str
	}
	vp.AddConfigPath(configPath)
	if err := vp.ReadInConfig(); err != nil {
		Log.Fatalf("Read config file faild! %s\n", err.Error())
	}
	if err := vp.Unmarshal(&Config); err != nil {
		Log.Fatalf("Unmarshal config file faild! %s\n", err.Error())
	}

	vp.OnConfigChange(func(in fsnotify.Event) {
		if err := vp.Unmarshal(&Config); err != nil {
			Log.Warnf("Unmarshal config file faild when update config! %s\n", err.Error())
		}
		for _, pchan := range regNotifyChannel {
			*pchan <- true
		}
	})
	defer vp.WatchConfig()

}

func RegNotifyChannel(channel *chan bool) {
	regNotifyChannel = append(regNotifyChannel, channel)
}

func Close() {
	for _, pchan := range regNotifyChannel {
		close(*pchan)
	}
}

func getEnv() (string, bool) {
	if str := os.Getenv(confEnvName); len(str) > 0 {
		return str, true
	} else {
		return str, false
	}
}
