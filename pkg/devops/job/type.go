package job

const urltag = "job"

type Status []struct {
	Status      string `json:"status"`
	PipeLineId  uint64 `json:"pipeline_id"`
	StageStatus []struct {
		Status     string `json:"status"`
		TaskStatus []struct {
			Status string `json:"status"`
			Name   string `json:"name"`
		} `json:"task_ststus"`
		Name string `json:"name"`
	}
	Id   uint64 `json:"id"`
	Logs []struct {
		Content string `json:"content"`
		Name    string `json:"name"`
	} `json:"logs"`
}
