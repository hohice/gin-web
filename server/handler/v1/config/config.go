package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hohice/gin-web/pkg/devops/job"
	"github.com/hohice/gin-web/pkg/devops/pipeline"
	"github.com/hohice/gin-web/pkg/devops/project"
	"github.com/hohice/gin-web/pkg/store"
	"github.com/hohice/gin-web/server/ex"
)

var storeMethod = store.GetMethodInstance()

func getUserId(c *gin.Context) uint64 {
	return 0
}

func startPipeLine(c *gin.Context, config *ConfigInfo) {
	uid := getUserId(c)
	if uid <= 0 {
		return
	}
	pidChan, errChan := pipeline.NewDefaultPipeLine(uid, config.ProjectInfo.Id, config.Config.Name, config.Config.Version)
	select {
	case pid := <-pidChan:
		{
			config.PipeLineId = pid
			go addConfig(c, config)
		}
	case err := <-errChan:
		{
			c.JSON(ex.ReturnInternalServerError(err))
		}
	}
}

func startJob(c *gin.Context, config *ConfigInfo) {
	uid := getUserId(c)
	if uid <= 0 {
		return
	}
	jobChan, errChan := job.NewPipelineJob(config.PipeLineId)
	select {
	case jobId := <-jobChan:
		{
			config.Jobid = jobId
			go addConfig(c, config)
		}
	case err := <-errChan:
		{
			c.JSON(ex.ReturnInternalServerError(err))
		}
	}
}

//delAll can delete project pipeline and job ,if failed return error info
func delAll(c *gin.Context, configInfo *ConfigInfo) error {
	doneProjChan, errProjChan := project.DelProject(configInfo.ProjectInfo.Id)
	donePipelineChan, errPipelineChan := pipeline.DelPipeLine(configInfo.PipeLineId)
	doneJobChan, errJobChan := job.DelPipelineJob(configInfo.Jobid)
	result := 1

listen:
	for {
		select {
		case err := <-errProjChan:
			{
				return err
			}
		case err := <-errPipelineChan:
			{
				return err
			}
		case err := <-errJobChan:
			{
				return err
			}
		case <-doneProjChan:
			{
				result <<= 1
			}
		case <-donePipelineChan:
			{
				result <<= 1
			}
		case <-doneJobChan:
			{
				result >>= 3
				break listen
			}
		}
	}
	if result > 0 {
		switch result {
		case 1 << 1:
			{
				return errors.New("delete pipeline failed")
			}
		case 1 << 2:
			{
				return errors.New("delete job of pipeline failed")
			}
		}

	}
	{
		key := fmt.Sprintf("%s_%s", configInfo.Config.Name, configInfo.Config.Version)
		if err := storeMethod.Del(key); err != nil {
			return err
		}
	}

	return nil
}

//addConfig add config with [#setConfigInfo] and delete config when set faild
func addConfig(c *gin.Context, configInfo *ConfigInfo) {
	if err := setConfigInfo(configInfo); err != nil {
		delAll(c, configInfo)
		c.JSON(ex.ReturnInternalServerError(err))
	} else {
		c.JSON(ex.ReturnOK())
	}
}

//------------------------------------------- store ------------------------------//
//getConfigInfo get config with key made by param
func getConfigInfo(name, version string) (*ConfigInfo, bool, error) {
	pConfigInfo := &ConfigInfo{}
	key := fmt.Sprintf("%s_%s", name, version)
	if data, err := storeMethod.Get(key); err != nil {
		return nil, false, err
	} else {
		if err := json.Unmarshal([]byte(data), pConfigInfo); err != nil {
			return nil, false, err
		}
		if pConfigInfo.Config.Name == name && pConfigInfo.Config.Version == version {
			return pConfigInfo, true, nil
		} else {
			return nil, false, nil
		}
	}
}

//setConfigInfo set config with key made by param
func setConfigInfo(configinfo *ConfigInfo) error {

	key := fmt.Sprintf("%s_%s", configinfo.Config.Name, configinfo.Config.Version)
	if data, err := json.Marshal(configinfo); err != nil {
		return err
	} else {
		if err := storeMethod.Set(key, data); err != nil {
			return err
		} else {
			return nil
		}
	}

}
