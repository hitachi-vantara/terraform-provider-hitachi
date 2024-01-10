package vssbstorage

import (
	"time"
)

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

// JobResponse .
type JobResponse struct {
	JobID         string    `json:"jobId"`
	Self          string    `json:"self"`
	UserID        string    `json:"userId"`
	Status        string    `json:"status"`
	State         string    `json:"state"`
	CreatedTime   time.Time `json:"createdTime"`
	UpdatedTime   time.Time `json:"updatedTime"`
	CompletedTime time.Time `json:"completedTime"`
	Request       struct {
		RequestURL    string `json:"requestUrl"`
		RequestMethod string `json:"requestMethod"`
		RequestBody   string `json:"requestBody"`
	} `json:"request"`
	Error             GatewayError `json:"error"`
	AffectedResources []string     `json:"affectedResources"`
}
