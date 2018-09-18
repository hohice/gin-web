package pipeline

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hohice/gin-web/pkg/circuitclient"
	"github.com/hohice/gin-web/pkg/devops"
)

const DefaultPipeLineTemp = "init_chart"

func GetPipeLine(pipeLineId int) (<-chan PipeLine, <-chan error) {
	data := make(chan PipeLine, 1)
	errchan := make(chan error, 1)
	go func() {
		if url, err := devops.GetUrlByTag(urltag); err != nil {
			errchan <- err
		} else {
			url = fmt.Sprintf("%s?id=%d", url, pipeLineId)
			c, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if resp, err := circuitclient.Get(c, urltag, url); err != nil {
				errchan <- err
			} else {
				defer resp.Body.Close()
				if body, err := ioutil.ReadAll(resp.Body); err != nil {
					errchan <- err
				} else {
					pipeline := PipeLine{}
					if err := json.Unmarshal(body, &pipeline); err != nil {
						errchan <- err
					} else {
						data <- pipeline
					}
				}
			}
		}
	}()
	return data, errchan
}

func DelPipeLine(pipeLineId uint64) (<-chan bool, <-chan error) {
	done := make(chan bool, 1)
	errchan := make(chan error, 1)
	go func() {
		if url, err := devops.GetUrlByTag(urltag); err != nil {
			errchan <- err
		} else {
			url = fmt.Sprintf("%s?id=%d", url, pipeLineId)
			c, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if resp, err := circuitclient.Delete(c, urltag, url); err != nil {
				errchan <- err
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					done <- false
				} else {
					done <- true
				}
			}
		}
	}()
	return done, errchan
}

func NewDefaultPipeLine(userId, projectId uint64, name, branch string) (<-chan uint64, <-chan error) {
	pipeline := &PipeLine{
		Name:      name,
		Branch:    branch,
		ProjectId: projectId,
		UserId:    userId,
	}

	if jsonByte, _, err := GetJsonInstance(pipeline); err != nil {
		pipeLineId := make(chan uint64, 1)
		errchan := make(chan error, 1)
		errchan <- err
		return pipeLineId, errchan
	} else {
		pipeLineId, errchan := newPipeLineByBytes(jsonByte)
		return pipeLineId, errchan
	}
}

//NewPipeLine create pipeline by []byte drictery and return pipeline id when success
func newPipeLineByBytes(jsonByte []byte) (<-chan uint64, <-chan error) {
	pipeLineId := make(chan uint64, 1)
	errchan := make(chan error, 1)

	go func() {
		if url, err := devops.GetUrlByTag(urltag); err != nil {
			errchan <- err
		} else {
			c, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if resp, err := circuitclient.PostByte(c, urltag, url, jsonByte); err != nil {
				errchan <- err
			} else {
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					errchan <- errors.New("create pipeline failed")
				} else {
					if body, err := ioutil.ReadAll(resp.Body); err != nil {
						errchan <- err
					} else {
						bytesBuffer := bytes.NewBuffer(body)
						var tmp uint64
						binary.Read(bytesBuffer, binary.BigEndian, &tmp)
						pipeLineId <- tmp
					}
				}
			}
		}

	}()

	return pipeLineId, errchan
}

//NewPipeLine create pipeline and return pipeline id when success
func newPipeLine(pipeline *PipeLine) (<-chan uint64, <-chan error) {
	pipeLineId := make(chan uint64, 1)
	errchan := make(chan error, 1)

	go func() {
		if jsonByte, err := json.Marshal(pipeline); err != nil {
			errchan <- err
		} else {
			if url, err := devops.GetUrlByTag(urltag); err != nil {
				errchan <- err
			} else {
				c, _ := context.WithTimeout(context.Background(), 5*time.Second)
				if resp, err := circuitclient.PostByte(c, urltag, url, jsonByte); err != nil {
					errchan <- err
				} else {
					defer resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						errchan <- errors.New("create pipeline failed")
					} else {
						if body, err := ioutil.ReadAll(resp.Body); err != nil {
							errchan <- err
						} else {
							bytesBuffer := bytes.NewBuffer(body)
							var tmp uint64
							binary.Read(bytesBuffer, binary.BigEndian, &tmp)
							pipeLineId <- tmp
						}
					}
				}
			}
		}

	}()

	return pipeLineId, errchan
}

func ModPipeLine(pipeline PipeLine) (<-chan bool, <-chan error) {
	done := make(chan bool, 1)
	errchan := make(chan error, 1)
	go func() {
		if jsonByte, err := json.Marshal(pipeline); err != nil {
			errchan <- err
		} else {
			if url, err := devops.GetUrlByTag(urltag); err != nil {
				errchan <- err
			} else {
				c, _ := context.WithTimeout(context.Background(), 5*time.Second)
				if resp, err := circuitclient.PutByte(c, urltag, url, jsonByte); err != nil {
					errchan <- err
				} else {
					defer resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						errchan <- errors.New("create pipeline failed")
					} else {
						done <- true
					}
				}
			}
		}

	}()
	return done, errchan
}
