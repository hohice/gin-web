package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/hohice/gin-web/pkg/util/logger"
)

var log = logger.DefaultLogger

func GetJsonInstance(pipeline *PipeLine) ([]byte, *PipeLine, error) {
	t := template.New("jsontemp")
	var err error
	if t, err = t.Parse(jsonTemp); err != nil {
		return []byte{}, nil, err
	}

	var w = new(bytes.Buffer)
	if err = t.ExecuteTemplate(w, "jsontemp", pipeline); err != nil {
		return []byte{}, nil, err
	} else {
		if err = json.Unmarshal(w.Bytes(), pipeline); err != nil {
			return []byte{}, nil, err
		} else {
			log.Debugw(fmt.Sprintf("%v", pipeline), "template", "pipeline")
			return w.Bytes(), pipeline, nil
		}
	}

}

var jsonTemp = `
{{ define "jsontemp" }}
	{
	  "stages": [
		{
		  "next_stage_name": "syntax-check",
		  "tasks": [
			{
			  "action": [
				{
				  "content": "flask milkcow syntax-check {{.Name}}",
				  "type": "string"
				}
			  ],
			  "name": "parameter-check"
			}
		  ],
		  "name": "parameter-check",
		  "checkouts": [
			{
			  "credential": "{{.UserId}}",
			  "project_url": "string",
			  "branch": "{{.Branch}}"
			}
		  ],
		  "is_head": true,
		  "environment": "172.16.1.99/postcommit/milkcow:v1",
		  "is_tail": false,
		  "parameter": [
			{
			  "name": "string",
			  "value": "string"
			}
		  ]
		},
		{
		  "next_stage_name": "generator-chart",
		  "tasks": [
			{
			  "action": [
				{
				  "content": "flask milkcow parameter-check  {{.Name}}",
				  "type": "string"
				}
			  ],
			  "name": "syntax-check"
			}
		  ],
		  "name": "syntax-check",
		  "checkouts": [
			{
			  "credential": "{{.UserId}}",
			  "project_url": "string",
			  "branch": "{{.Branch}}"
			}
		  ],
		  "is_head": false,
		  "environment": "172.16.1.99/postcommit/milkcow:v1",
		  "is_tail": false,
		  "parameter": [
			{
			  "name": "string",
			  "value": "string"
			}
		  ]
		},
		{
		  "next_stage_name": "",
		  "tasks": [
			{
			  "action": [
				{
				  "content": "flask milkcow generator-chart  {{.Name}} -o ./",
				  "type": "sh"
				}
			  ],
			  "name": "generator-chart"
			}
		  ],
		  "name": "generator-chart",
		  "checkouts": [
			{
			  "credential": "{{.UserId}}",
			  "project_url": "string",
			  "branch": "{{.Branch}}"
			}
		  ],
		  "is_head": false,
		  "environment": "172.16.1.99/postcommit/milkcow:v1",
		  "is_tail": true,
		  "parameter": [
			{
			  "name": "string",
			  "value": "string"
			}
		  ]
		}
	  ],
	  "user_id": {{.UserId}},
	  "name": "{{.Name}}",
	  "trigger": "manual",
	  "branch": "{{.Branch}}",
	  "project_id": {{.ProjectId}},
	  "id": {{.Id}}
	}
	{{ end }}`
