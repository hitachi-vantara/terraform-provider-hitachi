package infra_gw

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Roles    []Role `json:"roles"`
}

type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserWithDetails struct {
	Path    string   `json:"path"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type RoleWithDetails struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Data    []Role `json:"data"`
}
type UserData struct {
	TotalUserCount int    `json:"totalUserCount"`
	Users          []User `json:"users"`
}
