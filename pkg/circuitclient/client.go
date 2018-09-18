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
	reconfig := make(chan struct{}, 1)
	setting.RegNotifyChannel(reconfig)
	go func() {
		for {
			select {
			case _, ok := <-reconfig:
				{
					if !ok {
						return
					} else {
						readConfig()
					}
				}
			}
		}

	}()
}

func readConfig() {
	for name, value := range setting.Config.Circuit {
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

func Delete(c context.Context, name, uri string) (resp *http.Response, err error) {
	done := make(chan *http.Response, 1)
	if reps, err := http.NewRequest("DELETE", uri, nil); err != nil {
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

func PostByte(c context.Context, name, uri string, data []byte) (resp *http.Response, err error) {

	done := make(chan *http.Response, 1)
	if reps, err := http.NewRequest("POST", uri, bytes.NewReader(data)); err != nil {
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

func PutByte(c context.Context, name, uri string, data []byte) (resp *http.Response, err error) {

	done := make(chan *http.Response, 1)
	if reps, err := http.NewRequest("PUT", uri, bytes.NewReader(data)); err != nil {
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
