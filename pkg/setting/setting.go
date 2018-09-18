package setting

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var confEnvName = "GINS_CONF_PATH"

var DefaultConfPath = "/Users/hohice/MyWorkspace/go/src/github.com/hohice/gin-web/pkg/setting/testdata"

var DefaultWalmHome = filepath.Join(HomeDir(), ".ginS")

var Config Configs

var configPath string

var regNotifyChannel []chan struct{}

type Configs struct {
	Service string `mapstructure:"service"`
	Home    string `mapstructure:"home"`
	Debug   bool   `mapstructure:"debug"`
	Log     struct {
		Logformat string `mapstructure:"logformat"`
		LogPath   string `mapstructure:"logpath"`
	} `mapstructure:"log"`

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
		Enable      bool   `mapstructure:"enable"`
		Dirver      string `mapstructure:"mysql"`
		Username    string `mapstructure:"root"`
		Password    string `mapstructure:"password"`
		Host        string `mapstructure:"host"`
		Dbname      string `mapstructure:"dbname"`
		MaxOpenConn int    `mapstructure:"max_open_conn"`
		MaxIdleConn int    `mapstructure:"max_idle_conn"`
		MaxLifeTime int    `mapstructure:"max_life_time"`
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

	Devops struct {
		Url                 string `mapstructure:"url"`
		DefaultTemplateFile string `mapstructure:"default_template_file"`
	} `mapstructure:"devops"`

	Store struct {
		Bases map[string]struct {
			Type      string `mapstructure:"type"`
			Enable    bool   `mapstructure:"enable"`
			BasePath  string `mapstructure:"base_path"`
			IndexPath string `mapstructure:"index_path"`
			ValuePath string `mapstructure:"value_path"`
		} `mapstructure:"bases"`
	} `mapstructure:"store"`
}

// Init sets values from the environment.
func Init() error { //Init
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigName("config")
	vp.SetDefault("http.port", 8000)
	vp.SetDefault("home", DefaultWalmHome)

	vp.AddConfigPath(DefaultConfPath)
	if str, have := getEnv(); have {
		cp := str
		vp.AddConfigPath(cp)
	}
	if len(configPath) > 0 {
		vp.AddConfigPath(configPath)
	}

	return ReadConfigPath(vp)
}

func ReadConfigPath(vp *viper.Viper) error {
	if err := vp.ReadInConfig(); err != nil {
		return err
	}
	if err := vp.Unmarshal(&Config); err != nil {
		return err
	}

	if err := SyncNotify(); err != nil {
		return err
	}

	go notify()

	vp.OnConfigChange(func(in fsnotify.Event) {
		if err := vp.Unmarshal(&Config); err != nil {
			return
		}
		SyncNotify()
		go notify()
	})
	defer vp.WatchConfig()
	return nil
}

func notify() {
	for _, pchan := range regNotifyChannel {
		pchan <- struct{}{}
	}
}

func SyncNotify() error {
	for _, fun := range SyncNotifyFuncs {
		if err := fun(); err != nil {
			return err
		}
	}
	return nil
}

type SyncNotifyFunc func() error

var SyncNotifyFuncs []SyncNotifyFunc

func RegSyncNotify(snf SyncNotifyFunc) {
	SyncNotifyFuncs = append(SyncNotifyFuncs, snf)
}

func RegNotifyChannel(channel chan struct{}) {
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

func AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&configPath, "config-path", "", "path of config file,config file name: config.yaml")
}
