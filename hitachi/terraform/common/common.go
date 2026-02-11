package terraform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"
	// mc "terraform-provider-hitachi/hitachi/messagecatalog"
	sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
	vssbmodel "terraform-provider-hitachi/hitachi/storage/vosb/reconciler/model"
	// sanprov "terraform-provider-hitachi/hitachi/storage/san/provisioner"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// file not used

func GetSanSettingsFromFile(serialNumber string) (*sanmodel.StorageDeviceSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	filePath := utils.SAN_SETTINGS_DIR + "/" + serialNumber
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	data := sanmodel.StorageSettingsAndInfo{}
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	log.WriteDebug("storageSetting: %+v", data.Settings)

	return &data.Settings, nil
}

func GetVssbSettingsFromFile(vssbAddr string) (*vssbmodel.StorageDeviceSettings, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	filePath := utils.VOSB_SETTINGS_DIR + "/" + vssbAddr
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	data := vssbmodel.StorageSettingsAndInfo{}
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.WriteError(err)
		return nil, err
	}

	log.WriteDebug("vssb storageSetting: %+v", data.Settings)

	return &data.Settings, nil
}

func ParseHostGroupFromID(id string) (portID string, hostGroupNumber int, hostGroupName string, err error) {
    if id == "" {
        // No ID yet — treat as "no stored values"
        return "", 0, "", nil
    }

    parts := strings.Split(id, ",")
    if len(parts) != 3 {
        return "", 0, "", fmt.Errorf("invalid ID format: %s", id)
    }

    portID = parts[0]

    hostGroupNumber, err = strconv.Atoi(parts[1])
    if err != nil {
        return "", 0, "", fmt.Errorf("invalid hostGroupNumber in ID: %s", parts[1])
    }

    hostGroupName = parts[2]

    return portID, hostGroupNumber, hostGroupName, nil
}

// ExtractLdevFields reads integer and hex LDEV fields from Terraform schema,
// normalizes them via utils.ParseLdev, and returns a final integer LDEV pointer.
func ExtractLdevFields(d *schema.ResourceData, idField, hexField string) (*int, error) {
	// Read integer LDEV field
	var ldevID *int
	if v, ok := d.GetOk(idField); ok {
		val := v.(int)
		ldevID = &val
	}

	// Read hex LDEV field
	var ldevHex *string
	if v, ok := d.GetOk(hexField); ok {
		str := v.(string)
		ldevHex = &str
	}

	// If neither field exists, do nothing
	if ldevID == nil && ldevHex == nil {
		return nil, nil
	}

	// Normalize using your existing utility
	finalLdev, err := utils.ParseLdev(ldevID, ldevHex)
	if err != nil {
		return nil, err
	}

	return &finalLdev, nil
}

// ExtractLdevFromMap reads integer and hex LDEV fields from a TF map
// and returns the parsed final LDEV pointer.
func ExtractLdevFromMap(m map[string]interface{}, idField, hexField string) (*int, error) {
	var ldevID *int
	var ldevHex *string

	if v, ok := m[idField]; ok && v != nil {
		val := v.(int)
        // Terraform sets zero-value ints for absent fields; only treat as valid if user set it
        if val != 0 && val != -1 {
		    ldevID = &val
        }
	}

	if v, ok := m[hexField]; ok && v != nil {
		str := v.(string)
        if str != "" {
		    ldevHex = &str
        }
	}

	// Neither provided → skip
	if ldevID == nil && ldevHex == nil {
		return nil, nil
	}

	// Enforce mutual exclusivity (already handled by schema but safe)
	finalLdev, err := utils.ParseLdev(ldevID, ldevHex)
	if err != nil {
		return nil, err
	}

	return &finalLdev, nil
}

// GetStringPointer returns a pointer to the string if it's set in the HCL
func GetStringPointer(d *schema.ResourceData, key string) *string {
	if v, ok := d.GetOkExists(key); ok {
		val := v.(string)
		return &val
	}
	return nil
}

// GetIntPointer returns a pointer to the int if it's set in the HCL
func GetIntPointer(d *schema.ResourceData, key string) *int {
	if v, ok := d.GetOkExists(key); ok {
		val := v.(int)
		return &val
	}
	return nil
}

// GetBoolPointer returns a pointer to the bool if it's set in the HCL
func GetBoolPointer(d *schema.ResourceData, key string) *bool {
	// We use GetOkExists for bools because a 'false' value is 
	// technically a zero-value and GetOk might skip it.
	if v, ok := d.GetOkExists(key); ok {
		val := v.(bool)
		return &val
	}
	return nil
}
