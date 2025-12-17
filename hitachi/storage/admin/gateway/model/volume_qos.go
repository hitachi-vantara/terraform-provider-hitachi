package admin

type VolumeQosThreshold struct {
	IsUpperIopsEnabled         bool `json:"isUpperIopsEnabled"`
	UpperIops                  int  `json:"upperIops,omitempty"`
	IsUpperTransferRateEnabled bool `json:"isUpperTransferRateEnabled"`
	UpperTransferRate          int  `json:"upperTransferRate,omitempty"`
	IsLowerIopsEnabled         bool `json:"isLowerIopsEnabled"`
	LowerIops                  int  `json:"lowerIops,omitempty"`
	IsLowerTransferRateEnabled bool `json:"isLowerTransferRateEnabled"`
	LowerTransferRate          int  `json:"lowerTransferRate,omitempty"`
	IsResponsePriorityEnabled  bool `json:"isResponsePriorityEnabled"`
	ResponsePriority           int  `json:"responsePriority,omitempty"`
	TargetResponseTime         int  `json:"targetResponseTime,omitempty"`
}

type VolumeQosAlertSetting struct {
	IsUpperAlertEnabled        bool `json:"isUpperAlertEnabled"`
	UpperAlertAllowableTime    int  `json:"upperAlertAllowableTime,omitempty"`
	IsLowerAlertEnabled        bool `json:"isLowerAlertEnabled"`
	LowerAlertAllowableTime    int  `json:"lowerAlertAllowableTime,omitempty"`
	IsResponseAlertEnabled     bool `json:"isResponseAlertEnabled"`
	ResponseAlertAllowableTime int  `json:"responseAlertAllowableTime,omitempty"`
}

type VolumeQosAlertTime struct {
	UpperAlertTime    string `json:"upperAlertTime,omitempty"`
	LowerAlertTime    string `json:"lowerAlertTime,omitempty"`
	ResponseAlertTime string `json:"responseAlertTime,omitempty"`
}

type VolumeQosSettingsRequest struct {
	Threshold *VolumeQosThreshold    `json:"threshold,omitempty"`
	Alert     *VolumeQosAlertSetting `json:"alertSetting,omitempty"`
}

type VolumeQosResponse struct {
	VolumeId     int                   `json:"volumeId,omitempty"`
	Threshold    VolumeQosThreshold    `json:"threshold,omitempty"`
	AlertSetting VolumeQosAlertSetting `json:"alertSetting,omitempty"`
	AlertTime    VolumeQosAlertTime    `json:"alertTime,omitempty"`
}
