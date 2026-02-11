package messagecatalog

// MessageID .
type MessageID uint64

const sidx = 2000

const (
	// Default1 .
	Default1 MessageID = iota + sidx

	// Messages for GetStorageAdminInfo operation
	INFO_GET_STORAGE_ADMIN_INFO_BEGIN
	INFO_GET_STORAGE_ADMIN_INFO_END
	ERR_GET_STORAGE_ADMIN_INFO_FAILED
)

var enumStrings = map[MessageID]string{
	Default1: "Default1",

	// Messages for GetStorageAdminInfo operation
	INFO_GET_STORAGE_ADMIN_INFO_BEGIN: "INFO_GET_STORAGE_ADMIN_INFO_BEGIN",
	INFO_GET_STORAGE_ADMIN_INFO_END:   "INFO_GET_STORAGE_ADMIN_INFO_END",
	ERR_GET_STORAGE_ADMIN_INFO_FAILED: "ERR_GET_STORAGE_ADMIN_INFO_FAILED",
}

func (s MessageID) String() string { return enumStrings[s] }

// GetEnumString .
func GetEnumString(m interface{}) string {
	if m, ok := m.(MessageID); ok {
		return m.String()
	}

	return "UNKNOWN"
}
