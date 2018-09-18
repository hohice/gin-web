package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hohice/gin-web/pkg/circuitclient"
	"github.com/hohice/gin-web/pkg/devops"
)

func DelProject(projectId uint64) (<-chan bool, <-chan error) {
	done := make(chan bool, 1)
	errchan := make(chan error, 1)
	go func() {
		if url, err := devops.GetUrlByTag(urltag); err != nil {
			errchan <- err
		} else {
			url = fmt.Sprintf("%s?id=%d", url, projectId)
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

func NewProject(name, branch, contexts string, userId uint64) (<-chan ProjectInfo, <-chan error) {
	info := make(chan ProjectInfo, 1)
	errchan := make(chan error, 1)
	go func() {
		proj := Project{
			Name: name,
			Files: []File{
				File{
					Content:       contexts,
					Path:          name,
					Branch:        branch,
					CommitMessage: "",
				},
			},
			NamespaceId: userId,
		}
		if jsonByte, err := json.Marshal(proj); err != nil {
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
						errchan <- errors.New("create project failed")
					} else {
						if body, err := ioutil.ReadAll(resp.Body); err != nil {
							errchan <- err
						} else {
							pinfo := ProjectInfo{
								Name: name,
							}
							if err := json.Unmarshal(body, &pinfo); err != nil {
								errchan <- err
							} else {
								info <- pinfo
							}
						}
					}
				}
			}
		}

	}()
	return info, errchan
}

func ModProject(name, branch, contexts string, userId, projectId uint64) (<-chan bool, <-chan error) {
	done := make(chan bool, 1)
	errchan := make(chan error, 1)

	go func() {
		proj := Project{
			Name: name,
			Files: []File{
				File{
					Content:       contexts,
					Path:          name,
					Branch:        branch,
					CommitMessage: "",
				},
			},
			Id:          projectId,
			NamespaceId: userId,
		}
		if jsonByte, err := json.Marshal(proj); err != nil {
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
						errchan <- errors.New("modify project failed")
					} else {
						done <- true
					}
				}
			}
		}

	}()
	return done, errchan
}
