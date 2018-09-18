package config

import (
	"github.com/hohice/gin-web/pkg/devops/project"
)

type ConfigType struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Context string `json:"context"`
}

type ConfigInfo struct {
	Config      ConfigType          `json:"config"`
	ProjectInfo project.ProjectInfo `json:"project_info"`
	PipeLineId  uint64              `json:"pipeline_id"`
	Jobid       uint64              `json:"job_id"`
}
