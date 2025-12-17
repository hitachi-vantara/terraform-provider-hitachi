package messagecatalog

// GetMessage .
func GetMessage(id interface{}) string {
	str := MessageCatalog[id]
	return str
}

var MessageCatalog = map[interface{}]string{
	Default1: "%v",

	// Messages for GetStorageAdminInfo operation
	INFO_GET_STORAGE_ADMIN_INFO_BEGIN: "Getting storage admin information",
	INFO_GET_STORAGE_ADMIN_INFO_END:   "Sucessfully getting storage admin information",
	ERR_GET_STORAGE_ADMIN_INFO_FAILED: "Failed to get storage admin information",
}
