package setting

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var confEnvName = "GINS_CONF_PATH"

var configPath = "/etc/ginS/"

var DefaultWalmHome = filepath.Join(HomeDir(), ".ginS")

var Config Configs

var regNotifyChannel []chan bool

type Configs struct {
	Service   string `mapstructure:"service"`
	Home      string `mapstructure:"home"`
	Debug     bool   `mapstructure:"debug"`
	Logformat string `mapstructure:"logformat"`

	Http struct {
		HTTPPort     int           `mapstructure:"port"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
	} `mapstructure:"http"`

	Secret struct {
		Account   map[string]string `mapstructure:"account"`
		Tls       bool              `mapstructure:"tls"`
		TlsVerify bool              `mapstructure:"tls-verify"`
		TlsKey    string            `mapstructure:"tls-key"`
		TlsCert   string            `mapstructure:"tls-cert"`
		TlsCaCert string            `mapstructure:"tls-ca-cert"`
	} `mapstructure:"secret"`

	Database struct {
		Enable   bool   `mapstructure:"enable"`
		Dirver   string `mapstructure:"mysql"`
		Username string `mapstructure:"root"`
		Password string `mapstructure:"password"`
		Dbname   string `mapstructure:"dbname"`
	} `mapstructure:"database"`

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

	Limit struct {
		AddrMap     map[string]string `mapstructure:"addr_map"`
		DefaultRate int               `mapstructure:"default_rate"`
	} `mapstructure:"limit"`

	Circuit map[string]struct {
		Timeout                int `mapstructure:"timeout"`
		MaxConcurrentRequests  int `mapstructure:"max_concurrent_requests"`
		RequestVolumeThreshold int `mapstructure:"request_volume_threshold"`
		SleepWindow            int `mapstructure:"sleep_window"`
		ErrorPercentThreshold  int `mapstructure:"error_percent_threshold"`
	} `mapstructure:"circuit"`
}

// Init sets values from the environment.
func init() { //Init
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigName("config")
	vp.SetDefault("home", DefaultWalmHome)
	vp.SetDefault("http.port", 8000)

	ReadConfigPath(vp)

}

func ReadConfigPath(vp *viper.Viper) {
	if str, have := getEnv(); have {
		configPath = str
	}
	vp.AddConfigPath(configPath)
	if err := vp.ReadInConfig(); err != nil {
		panic("Read config file faild! " + err.Error())
	}
	if err := vp.Unmarshal(&Config); err != nil {
		panic("Unmarshal config file faild! " + err.Error())
	}

	vp.OnConfigChange(func(in fsnotify.Event) {
		if err := vp.Unmarshal(&Config); err != nil {
			panic("Unmarshal config file faild when update config!" + err.Error())
		}
		for _, pchan := range regNotifyChannel {
			pchan <- true
		}
	})
	defer vp.WatchConfig()
}

func RegNotifyChannel(channel chan bool) {
	regNotifyChannel = append(regNotifyChannel, channel)
}

func Close() {
	for _, pchan := range regNotifyChannel {
		close(pchan)
	}
}

func getEnv() (string, bool) {
	if str := os.Getenv(confEnvName); len(str) > 0 {
		return str, true
	} else {
		return str, false
	}
}
