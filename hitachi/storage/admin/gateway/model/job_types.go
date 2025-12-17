package admin

// GatewayError .
type GatewayError struct {
	ErrorSource  string `json:"errorSource"`
	Message      string `json:"message"`
	Cause        string `json:"cause"`
	Solution     string `json:"solution"`
	SolutionType string `json:"solutionType"`
	MessageID    string `json:"messageId"`
	ErrorCode    struct {
		Ssb2 string `json:"SSB2"`
		Ssb1 string `json:"SSB1"`
	} `json:"errorCode"`
	DetailCode string `json:"detailCode"`
}

type JobResponse struct {
	JobID             int               `json:"id"`                          // make it similar to other JobResponse
	Progress          string            `json:"progress"`
	Status            string            `json:"status,omitempty"`            // Nullable
	AffectedResources []string          `json:"affectedResources,omitempty"` // Nullable
	ErrorResource     string            `json:"errorResource,omitempty"`     // Nullable
	ErrorCode         ErrorCode         `json:"errorCode,omitempty"`         // Nullable
	ErrorMessage      string            `json:"errorMessage,omitempty"`      // Nullable
	OperationDetails  []OperationDetail `json:"operationDetails,omitempty"`  // Nullable
	Error             GatewayError      `json:"error"`                       // added to make it similar to other JobResponse
	MultipleJobs      []int      		`json:"multipleJobs"`                // added to support multiple jobs
}

type ErrorCode struct {
	SSB1 string `json:"SSB1"`
	SSB2 string `json:"SSB2"`
}

type OperationDetail struct {
	OperationType string  `json:"operationType"`        // CREATE, UPDATE, DELETE
	ResourceType  string  `json:"resourceType"`         // CommandStatus, Pool, Port, Server, Snapshot, etc.
	ResourceID    string  `json:"resourceId,omitempty"` // Nullable
}

type AsyncCommandStatus struct {
	StatusResource string `json:"statusResource,omitempty"`
}
