package pipeline

const urltag = "pipeline"

type PipeLine struct {
	//Status    string `json:"status"`
	ProjectId uint64  `json:"project_id"`
	UserId    uint64  `json:"user_id"`
	Name      string  `json:"name"`
	Trigger   string  `json:"trigger"`
	Periodic  string  `json:"periodic"`
	Branch    string  `json:"branch"`
	Stages    []Stage `json:"stages"`
	Id        uint64  `json:"id"`
} //`json:"stages"`

type Checkout struct {
	Credential string `json:"credential"`
	ProjectUrl string `json:"project_url"`
	Branch     string `json:"branch"`
}

type Parameters struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Actions struct {
	Content string `json:"content"`
	Type    string `json:"type"`
}

type Task struct {
	Action []Actions `json:"action"`
	Name   string    `json:"name"`
}

type Stage struct {
	Environment   string `json:"environment"`
	NextStageName string `json:"next_stage_name"`
	Tasks         []Task `json:"tasks"`

	Checkouts []Checkout   `json:"checkouts"`
	Name      string       `json:"name"`
	IsHead    bool         `json:"is_head"`
	IsTail    bool         `json:"is_tail"`
	Parameter []Parameters `json:"paremeter"`
}
