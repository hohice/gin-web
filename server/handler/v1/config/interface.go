package config

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hohice/gin-web/pkg/devops/project"
	"github.com/hohice/gin-web/server/ex"
	"github.com/hohice/gin-web/server/handler/util"
)

// GetConfig godoc
// @Tags config
// @Description Get Application Config
// @OperationId GetConfig
// @Accept  json
// @Produce  json
// @Param   name     path    string     true      "name of the config"
// @Param   version     path    string     true      "version of the config"
// @Success 200 {object} config.ConfigType	"ok"
// @Failure 400 {object} ex.Response "Invalid Name supplied!"
// @Failure 404 {object} ex.Response "Instance not found"
// @Failure 405 {object} ex.Response "Invalid input"
// @Failure 500 {object} ex.Response "Server Error"
// @Router "/config/name/{name}/version/{version}" [get]
func GetConfig(c *gin.Context) { //TerminalResponse
	if values, err := util.GetPathParams(c, []string{"name", "version"}); err != nil {
		c.JSON(ex.ReturnBadRequest())
	} else {
		name, version := values["name"], values["version"]
		if configInfo, _, err := getConfigInfo(name, version); err != nil {
			c.JSON(ex.ReturnConfigNotExistError())
		} else {
			c.JSON(http.StatusOK, configInfo.Config)
		}
	}
}

// DelConfig godoc
// @Tags config
// @Description Delete Application Config
// @OperationId DelConfig
// @Accept  json
// @Produce  json
// @Param   name     path    string     true      "name of the config"
// @Param   version     path    string     true      "version of the config"
// @Success 200 {object} ex.Response	"ok"
// @Failure 400 {object} ex.Response "Invalid Name supplied!"
// @Failure 404 {object} ex.Response "Instance not found"
// @Failure 405 {object} ex.Response "Invalid input"
// @Failure 500 {object} ex.Response "Server Error"
// @Router "/config/name/{name}/version/{version}" [delete]
func DelConfig(c *gin.Context) {
	if values, err := util.GetPathParams(c, []string{"name", "version"}); err != nil {
		c.JSON(ex.ReturnBadRequest())
	} else {
		name, version := values["name"], values["version"]
		//purge, _ := strconv.ParseBool(c.Query("purge"))

		if configInfo, found, err := getConfigInfo(name, version); err != nil {
			c.JSON(ex.ReturnInternalServerError(err))
		} else {
			if found {
				if err := delAll(c, configInfo); err != nil {
					c.JSON(ex.ReturnInternalServerError(err))
				}
			}
			c.JSON(ex.ReturnOK())
		}
	}
}

// ModConfig godoc
// @Tags config
// @Description Modify Application Config
// @OperationId ModConfig
// @Accept  json
// @Produce  json
// @Param   config     body    ConfigType     true      "data of the config"
// @Success 200 {object} ex.Response	"ok"
// @Failure 400 {object} ex.Response "Invalid Name supplied!"
// @Failure 404 {object} ex.Response "Instance not found"
// @Failure 405 {object} ex.Response "Invalid input"
// @Failure 500 {object} ex.Response "Server Error"
// @Router "/config" [put]
func ModConfig(c *gin.Context) {
	configType := &ConfigType{}
	if err := c.BindJSON(configType); err != nil {
		c.JSON(ex.ReturnBadRequest())
	}
	if configInfo, found, err := getConfigInfo(configType.Name, configType.Version); err != nil {
		c.JSON(ex.ReturnInternalServerError(err))
	} else {
		if found {
			if err := setConfigInfo(configInfo); err != nil {
				c.JSON(ex.ReturnInternalServerError(err))
			} else {
				userid := getUserId(c)
				doneChan, errChan := project.ModProject(configType.Name, configType.Version, configType.Context, userid, configInfo.ProjectInfo.Id)
				select {
				case <-doneChan:
					go startPipeLine(c, configInfo)
				case err := <-errChan:
					c.JSON(ex.ReturnInternalServerError(err))
				}
			}
		} else {
			NewConfig(c)
		}
	}
}

// NewConfig godoc
// @Tags config
// @Description Modify Application Config
// @OperationId NewConfig
// @Accept  json
// @Produce  json
// @Param   config     body    ConfigType     true      "data of the config"
// @Success 200 {object} ex.Response	"ok"
// @Failure 400 {object} ex.Response "Invalid Name supplied!"
// @Failure 404 {object} ex.Response "Instance not found"
// @Failure 405 {object} ex.Response "Invalid input"
// @Failure 500 {object} ex.Response "Server Error"
// @Router "/config" [post]
func NewConfig(c *gin.Context) {
	configtype := ConfigType{}
	if err := c.BindJSON(&configtype); err != nil {
		c.JSON(ex.ReturnBadRequest())
	}
	userid := getUserId(c)

	projectInfoChan, errChan := project.NewProject(configtype.Name, configtype.Version, configtype.Context, userid)

	select {
	case projectInfo := <-projectInfoChan:
		{
			configInfo := &ConfigInfo{
				Config:      configtype,
				ProjectInfo: projectInfo,
			}
			go addConfig(c, configInfo)
		}
	case err := <-errChan:
		c.JSON(ex.ReturnInternalServerError(err))
	}
}

// StartTest godoc
// @Tags application
// @Description Modify Application Config
// @OperationId NewConfig
// @Accept  json
// @Produce  json
// @Param   name     path    string     true      "name of the config"
// @Param   version     path    string     true      "version of the config"
// @Success 200 {object} ex.Response	"ok"
// @Failure 400 {object} ex.Response "Invalid Name supplied!"
// @Failure 404 {object} ex.Response "Instance not found"
// @Failure 405 {object} ex.Response "Invalid input"
// @Failure 500 {object} ex.Response "Server Error"
// @Router "/application/build/name/{name}/version/{version}" [post]
func StartTest(c *gin.Context) {
	if values, err := util.GetPathParams(c, []string{"name", "version"}); err != nil {
		c.JSON(ex.ReturnBadRequest())
	} else {
		name, version := values["name"], values["version"]
		if configInfo, _, err := getConfigInfo(name, version); err != nil {
			c.JSON(ex.ReturnInternalServerError(err))
		} else {
			go startJob(c, configInfo)
		}
	}
}
