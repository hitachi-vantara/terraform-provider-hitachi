package messagecatalog

// MessageID .
type MessageID uint64

const sidx = 2000

const (
	// Default1 .
	Default1 MessageID = iota + sidx
	ERR_GET_LUN_1
)

var enumStrings = map[interface{}]string{
	Default1:      "Default1",
	ERR_GET_LUN_1: "ERR_GET_LUN_1",
}

func (s MessageID) String() string { return enumStrings[s] }

// GetEnumString .
func GetEnumString(m interface{}) string {
	if m, ok := m.(MessageID); ok {
		return m.String()
	}

	return "UNKNOWN"
}
