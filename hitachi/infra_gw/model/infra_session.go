package infra_gw

type TaskResponse struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Data    struct {
		TaskId       string `json:"taskId"`
		Name         string `json:"name"`
		Initiator    string `json:"initiator"`
		ResourceType string `json:"resourceType"`
		Status       string `json:"status"`
		Events       []struct {
			Description       string `json:"description"`
			Time              string `json:"time"`
			RecommendedAction string `json:"recommendedAction"`
		} `json:"events"`
	} `json:"data"`
	StartTime            string `json:"startTime"`
	EndTime              string `json:"endTime"`
	AdditionalAttributes []struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"additionalAttributes"`
}

type Response struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Data    struct {
		TaskId     string `json:"taskId"`
		ResourceId string `json:"resourceId"`
		State      string `json:"state"`
	} `json:"data"`
}

type BasicResponse struct {
	
		TaskId     string `json:"taskId"`
		ResourceId string `json:"resourceId"`
		State      string `json:"state"`

}