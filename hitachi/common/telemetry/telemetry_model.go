package telemetry

// for sending to AWS api
type TelemetryPayload struct {
	ModuleName      string  `json:"module_name"`
	OperationName   string  `json:"operation_name"`
	OperationStatus int     `json:"operation_status"`
	ConnectionType  string  `json:"connection_type"`
	ProcessTime     float64 `json:"process_time"`
	StorageModel    string  `json:"storage_model"`
	StorageSerial   string  `json:"storage_serial"`
	StorageType     string  `json:"storage_type"`
	Site            string  `json:"site"`
}

// usages.json file
type UsagesTelemetry struct {
	ExecutionStats    map[string]ExecutionStat `json:"executionStats"`
	SanStorageSystems []SanStorageSystem       `json:"sanStorageSystems"`
	SdsBlockSystems   []SdsBlockSystem         `json:"SdsBlockSystems"`
}

// averagetime saved to json file
type ExecutionStat struct {
	Success     int     `json:"success"`
	Failure     int     `json:"failure"`
	AverageTime float64 `json:"averageTimeInSec"`
}

// san/vsp
type SanStorageSystem struct {
	StorageModel  string `json:"storageModel"`
	StorageSerial string `json:"storageSerial"`
}

// sds block
type SdsBlockSystem struct {
	ClusterAddress string `json:"clusterAddress"`
	Version        string `json:"version"`
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
