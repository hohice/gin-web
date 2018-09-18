package devops

import (
	"errors"

	"github.com/hohice/gin-web/pkg/setting"
)

type DevopsConf struct {
	Conf setting.Configs
	Urls map[string]string
}

var devops DevopsConf

type UrlTag struct {
	name string
	path string
}

var UrlTags = []UrlTag{
	UrlTag{
		name: "pipeline",
		path: "pipeline",
	},
	UrlTag{
		name: "job",
		path: "pipeline/job",
	},
	UrlTag{
		name: "projects",
		path: "projects",
	},
}

func init() {
	configChan := make(chan struct{})
	setting.RegNotifyChannel(configChan)
	go func() {
		for {
			select {
			case _, ok := <-configChan:
				{
					if !ok {
						return
					} else {
						devops := DevopsConf{
							Conf: setting.Config,
							Urls: map[string]string{},
						}

						for _, urlTag := range UrlTags {
							devops.Urls[urlTag.name] = devops.Conf.Devops.Url + "/" + urlTag.path
						}
					}
				}
			}
		}
	}()

}

func GetUrlByTag(tag string) (string, error) {
	if url, ok := devops.Urls[tag]; !ok {
		return "", errors.New("the Requst Url not exist!")
	} else {
		return url, nil
	}
}

func GetDefaultTemplateFileName() string {
	return devops.Conf.Devops.DefaultTemplateFile
}