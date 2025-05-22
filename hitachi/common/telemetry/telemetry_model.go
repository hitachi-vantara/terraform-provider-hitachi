package telemetry

// for sending to AWS api
type TelemetryPayload struct {
	ModuleName      string  `json:"moduleName"`
	OperationName   string  `json:"operationName"`
	OperationStatus int     `json:"operationStatus"`
	ConnectionType  string  `json:"connectionType"`
	ProcessTime     float64 `json:"processTime"`
	StorageModel    string  `json:"storageModel"`
	StorageSerial   string  `json:"storageSerial"`
	StorageType     string  `json:"storageType"`
	Site            string  `json:"site"`
}

// averagetime saved to json file
type ExecutionJsonFileStats struct {
	Success     int     `json:"success"`
	Failure     int     `json:"failure"`
	AverageTime float64 `json:"averageTimeInSec"`
}

// intermediate info
type MethodTaskExecution struct {
	TerraformTfType       string  `json:"terraformTfType"`       // ex: datasource:hitachi_vsp_hostgroup
	TerraformResourceName string  `json:"terraformResourceName"` // ex: san.datasource.DataSourceStorageHostGroup
	GatewayMethodName     string  `json:"gatewayMethodName"`     // gateway func name
	Status                string  `json:"status"`                // success or failure
	ElapsedTime           float64 `json:"elapsedTime"`           // in sec
	ConnectionType        string  `json:"connectionType"`        // san or vosb
	StorageModel          string  `json:"storageModel"`
	StorageSerial         string  `json:"storageSerial"`
	StorageType           string  `json:"storageType"` // vsp or sds_block
	SiteId                string  `json:"siteId"`
}

// User Consent data
type UserConsent struct {
	ConsentHistory      []ConsentHistory `json:"consent_history"`
	SiteID              string           `json:"site_id"`
	Time                string           `json:"time"`
	UserConsentAccepted bool             `json:"user_consent_accepted"`
}

type ConsentHistory struct {
	Time                string `json:"time"`
	UserConsentAccepted bool   `json:"user_consent_accepted"`
}
