package vssbstorage

import (
	"time"
)

type ChangeUserPasswordReq struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type UserGroup struct {
	UserGroupId       string `json:"userGroupId"`
	UserGroupObjectId string `json:"userGroupObjectId"`
}

type Privilege struct {
	Scope    string   `json:"scope"`
	RoleNames []string `json:"roleNames"`
}

// Response: User
type User struct {
	UserId                 string       `json:"userId"`
	UserObjectId           string       `json:"userObjectId"`
	PasswordExpirationTime time.Time    `json:"passwordExpirationTime"`
	IsEnabled              bool         `json:"isEnabled"`
	UserGroups            []UserGroup  `json:"userGroups"`
	IsBuiltIn             bool         `json:"isBuiltIn"`
	Authentication        string       `json:"authentication"`
	RoleNames             []string     `json:"roleNames"`
	IsEnabledConsoleLogin interface{}  `json:"isEnabledConsoleLogin"` // Can be nil
	VpsId                  string       `json:"vpsId"`
	Privileges            []Privilege  `json:"privileges"`
}
