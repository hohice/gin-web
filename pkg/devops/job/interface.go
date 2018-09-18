package job

import (
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

func GetPipelineStatus(jobId uint64) (<-chan Status, <-chan error) {
	data := make(chan Status, 1)
	errchan := make(chan error, 1)
	go func() {
		if url, err := devops.GetUrlByTag(urltag); err != nil {
			errchan <- err
		} else {
			url = fmt.Sprintf("%s?id=%d", url, jobId)
			c, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if resp, err := circuitclient.Get(c, urltag, url); err != nil {
				errchan <- err
			} else {
				defer resp.Body.Close()
				if body, err := ioutil.ReadAll(resp.Body); err != nil {
					errchan <- err
				} else {
					status := Status{}
					if err := json.Unmarshal(body, &status); err != nil {
						errchan <- err
					} else {
						data <- status
					}
				}
			}
		}
	}()
	return data, errchan
}

func DelPipelineJob(jobId uint64) (<-chan bool, <-chan error) {
	done := make(chan bool, 1)
	errchan := make(chan error, 1)
	go func() {
		if url, err := devops.GetUrlByTag(urltag); err != nil {
			errchan <- err
		} else {
			url = fmt.Sprintf("%s?id=%d", url, jobId)
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

func NewPipelineJob(pipelineId uint64) (<-chan uint64, <-chan error) {
	id := make(chan uint64, 1)
	errchan := make(chan error, 1)

	go func() {
		if url, err := devops.GetUrlByTag(urltag); err != nil {
			errchan <- err
		} else {
			url = fmt.Sprintf("%s?pipeline_id=%d", url, pipelineId)
			c, _ := context.WithTimeout(context.Background(), 5*time.Second)
			if resp, err := circuitclient.PostByte(c, urltag, url, []byte{}); err != nil {
				errchan <- err
			} else {
				defer resp.Body.Close()
				if body, err := ioutil.ReadAll(resp.Body); err != nil {
					errchan <- err
				} else {
					pid, n := binary.Uvarint(body)
					if n == 0 {
						errchan <- errors.New("parse project id failed")
					} else {
						id <- pid
					}
				}
			}
		}
	}()

	return id, errchan
}
