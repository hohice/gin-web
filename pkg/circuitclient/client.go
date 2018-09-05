package circuitclient

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"

	"github.com/hohice/gin-web/pkg/setting"
)

func init() {
	reconfig := make(chan bool, 1)
	readConfig()
	setting.RegNotifyChannel(reconfig)
	go func() {
		select {
		case <-reconfig:
			readConfig()
		}
	}()
}

func readConfig() {
	for name, value := range setting.Configs.Circuit {
		hystrix.ConfigureCommand(name, hystrix.CommandConfig{
			Timeout:                value.Timeout,
			MaxConcurrentRequests:  value.MaxConcurrentRequests,
			RequestVolumeThreshold: value.RequestVolumeThreshold,
			SleepWindow:            value.SleepWindow,
			ErrorPercentThreshold:  value.ErrorPercentThreshold,
		})
	}
}

func DefaultGenerateName(url string) string {
	return url
}

//Get to exec http.Get
//param context context of server
//param name hystrix Circuit name
//parem url url to get
func Get(c context.Context, name, uri string) (resp *http.Response, err error) {
	done := make(chan *http.Response, 1)
	if reps, err := http.NewRequest("GET", uri, nil); err != nil {
		return nil, err
	} else {
		errChan := hystrix.GoC(c, name, func(c context.Context) error {
			if resq, err := http.DefaultClient.Do(reps); err != nil {
				return err
			} else {
				done <- resq
				return nil
			}
		},
			func(c context.Context, err error) error {
				return nil
			})

		select {
		case err := <-errChan:
			return nil, err
		case resp := <-done:
			return resp, nil
		}
	}

}

func PostJson(c context.Context, name, uri string, param map[string]string) (resp *http.Response, err error) {
	jsonByte, _ := json.Marshal(param)
	done := make(chan *http.Response, 1)
	if reps, err := http.NewRequest("POST", uri, bytes.NewReader(jsonByte)); err != nil {
		return nil, err
	} else {
		errChan := hystrix.GoC(c, name, func(c context.Context) error {
			if resq, err := http.DefaultClient.Do(reps); err != nil {
				return err
			} else {
				done <- resq
				return nil
			}
		},
			func(c context.Context, err error) error {
				return nil
			})

		select {
		case err := <-errChan:
			return nil, err
		case resp := <-done:
			return resp, nil
		}
	}

}

func PostForm(c context.Context, name, uri string, param map[string]string) (resp *http.Response, err error) {
	done := make(chan *http.Response, 1)
	if reps, err := http.NewRequest("POST", uri+ParseToStr(param), nil); err != nil {
		return nil, err
	} else {
		errChan := hystrix.GoC(c, name, func(c context.Context) error {
			if resq, err := http.DefaultClient.Do(reps); err != nil {
				return err
			} else {
				done <- resq
				return nil
			}
		},
			func(c context.Context, err error) error {
				return nil
			})

		select {
		case err := <-errChan:
			return nil, err
		case resp := <-done:
			return resp, nil
		}
	}

}

//ParseToStr parse map params to strings that can be added to url
func ParseToStr(mp map[string]string) string {
	values := ""
	for key, val := range mp {
		values += "&" + key + "=" + val
	}
	temp := values[1:]
	values = "?" + temp
	return values
}
