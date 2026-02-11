package utils

import (
	"bytes"
	b64 "encoding/base64"

	// "encoding/json"
	"errors"
	"fmt"

	"golang.org/x/exp/slices"

	// "io/ioutil"
	// "os"
	"math"
	"regexp"
	"strconv"
	"strings"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"unicode"
)

// these funcs are not used in the code
/*
func GetOutputValue(value interface{}) interface{} {
	log := commonlog.GetLogger()

	valArray, okArray := value.([]interface{})
	if okArray {
		// log.WriteInfo("It's an array : %+v\n", valArray)
		outArray := make([]interface{}, 0)
		for _, item := range valArray {
			outMap := GetOutputValue(item)
			outArray = append(outArray, outMap)

		}
		return outArray
	} else {
		valMap, okMap := value.(map[string]interface{})
		if okMap {
			// log.WriteInfo("It's a map : %+v\n", valMap)
			outMap := map[string]interface{}{}
			for k, v := range valMap {
				outMap[ToSnakeCase(k)] = GetOutputValue(v)
			}
			return outMap
		} else {
			// log.WriteInfo("It's primitive type : %+v\n", value)
			return value
		}
	}
}

func PopulateOutput(value interface{}) []interface{} {
	lastValueInArray := make([]interface{}, 0)
	lastValue := GetOutputValue(value)
	valMap, okMap := lastValue.(map[string]interface{})
	if okMap {
		// log.WriteInfo("Last value is a map : %+v\n", lastValue)
		lastValueInArray = append(lastValueInArray, valMap)
	} else {
		// log.WriteInfo("Last value is an array : %+v\n", lastValue)
		lastValueInArray = lastValue.([]interface{})
	}

	log.WriteInfo("Output Schema Value: %s\n", ConvertToJson(lastValueInArray))
	return lastValueInArray
}

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func IsInterfaceArray(v interface{}) bool {
	log := commonlog.GetLogger()

	switch x := v.(type) {
	case []interface{}:
		log.WriteDebug("[]interface, len:", len(x))
		return true
	case interface{}:
		log.WriteDebug("interface:", x)
		return false
	default:
		log.WriteDebugf("Unsupported type: %T\n", x)
		return false
	}
}

func ConvertToJson(m interface{}) string {
	log := commonlog.GetLogger()

	b, err := json.Marshal(m)
	if err != nil {
		log.WriteError(err)
	}
	return string(b)
}

func SaveOutputToFile(resourceName string, id string, output interface{}) error {
	log := commonlog.GetLogger()

	name := "output/" + resourceName + "_" + id + ".out"
	_ = os.Mkdir("output", os.ModeDir|os.ModePerm)

	file, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		log.WriteError(err)
		return err
	}

	err = ioutil.WriteFile(name, file, 0644)
	if err != nil {
		log.WriteError(err)
		return err
	}

	return nil
}

func SaveDataToFile(dirPath, fileName string, data interface{}) error {
	log := commonlog.GetLogger()

	err := os.MkdirAll(dirPath, os.ModeDir|os.ModePerm)
	if err != nil {
		log.Error(err)
		return err
	}

	filePath := dirPath + "/" + fileName

	jsonString, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Error(err)
		return err
	}

	err = ioutil.WriteFile(filePath, jsonString, 0644)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
*/

// TransformSizeToUnit --
func TransformSizeToUnit(size uint64) string {
	if size >= TERABYTES {
		return fmt.Sprintf("%0.2f TB", float64(size)/float64(TERABYTES))
	}
	if size >= GIGABYTES {
		return fmt.Sprintf("%0.2f GB", float64(size)/float64(GIGABYTES))
	}
	if size >= MEGABYTES {
		return fmt.Sprintf("%0.2f MB", float64(size)/float64(MEGABYTES))
	}
	if size >= KILOBYTES {
		return fmt.Sprintf("%0.2f KB", float64(size)/float64(KILOBYTES))
	}
	return strconv.Itoa(int(size))
}

func ConvertSizeFromBytesToMb(size uint64) uint64 {
	return uint64(float64(size) / float64(MEGABYTES))
}

func ConvertSizeFromKbToMb(size uint64) uint64 {
	return uint64(float64(size*1024) / float64(MEGABYTES))
}

func ConvertSizeFromKbToGb(size uint64) uint64 {
	return uint64(float64(size*1024) / float64(GIGABYTES))
}

// ConvertSizeToBytes -- examples: 1GB, 2.5MB
func ConvertSizeToBytes(size string) (uint64, error) {
	var sizeInBytes uint64

	var reTotCap = regexp.MustCompile("((\\d*\\.)?\\d+)\\s*([A-Z]+)?")
	matches := reTotCap.FindStringSubmatch(size)
	snum, _ := strconv.ParseFloat(matches[1], 64)

	switch strings.ToUpper(matches[3]) {
	case "KB":
	case "K":
		sizeInBytes = uint64((snum * KILOBYTES / BLOCKSIZE)) * BLOCKSIZE
	case "MB":
	case "M":
		sizeInBytes = uint64((snum * MEGABYTES / BLOCKSIZE)) * BLOCKSIZE
	case "GB":
	case "G":
		sizeInBytes = uint64((snum * GIGABYTES / BLOCKSIZE)) * BLOCKSIZE
	case "TB":
	case "T":
		sizeInBytes = uint64((snum * TERABYTES / BLOCKSIZE)) * BLOCKSIZE
	case "":
		sizeInBytes = uint64((snum / BLOCKSIZE)) * BLOCKSIZE
	default:
		return sizeInBytes, errors.New("Cannot convert size")
	}

	return sizeInBytes, nil
}

func ParseCapacityToMiB(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, fmt.Errorf("empty capacity string")
	}

	// Require explicit unit suffix (M, G, or T)
	re := regexp.MustCompile(`^([0-9]+(\.[0-9]+)?)([MGT])$`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid format (must include unit suffix M, G, or T)")
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid numeric value")
	}

	switch matches[3] {
	case "T":
		value *= 1024 * 1024
	case "G":
		value *= 1024
	case "M":
		// already in MiB
	case "K":
		value /= 1024
	default:
		return 0, fmt.Errorf("unsupported unit %q", matches[3])
	}

	return int64(math.Round(value)), nil
}

// ConvertFloatSizeToSmartUnit converts a float64 size in GB into a smart string
// using T, G, M, or K, always using the largest unit that yields an integer.
func ConvertFloatSizeToSmartUnit(sizeGB float64) string {
	if sizeGB <= 0 {
		return "0K"
	}

	bytes := sizeGB * 1024 * 1024 * 1024 // convert GB to bytes

	const (
		kb = 1024.0
		mb = kb * 1024
		gb = mb * 1024
		tb = gb * 1024
	)

	// Helper to check if value is whole number
	isInt := func(f float64) bool {
		return math.Abs(f-math.Round(f)) < 1e-9
	}

	// Try largest possible unit that results in an integer
	tbVal := bytes / tb
	if isInt(tbVal) && tbVal >= 1 {
		return fmt.Sprintf("%dT", int64(tbVal))
	}

	gbVal := bytes / gb
	if isInt(gbVal) && gbVal >= 1 {
		return fmt.Sprintf("%dG", int64(gbVal))
	}

	mbVal := bytes / mb
	if isInt(mbVal) && mbVal >= 1 {
		return fmt.Sprintf("%dM", int64(mbVal))
	}

	kbVal := bytes / kb
	return fmt.Sprintf("%dK", int64(math.Round(kbVal)))
}

func DecodeBase64EncodedString(strvalue string) (string, error) {
	log := commonlog.GetLogger()

	var x string
	// fmt.Println(strvalue)
	uDec, err := b64.URLEncoding.DecodeString(strvalue)
	// fmt.Println(string(uDec))
	if err != nil {
		log.WriteError(err)
		return x, err
	}

	return strings.TrimSuffix(bytes.NewBuffer(uDec).String(), "\n"), nil
}

func IsParityGroupPresent(parityGroupId string, parityGroupIDs []string) bool {
	for _, pg_id := range parityGroupIDs {
		if pg_id == parityGroupId {
			return true
		}
	}
	return false
}

// GetStringSliceDiff - see demo https://go.dev/play/p/rljfryR0Ek_n
//
//	added - new elements added to the new slice
//	removed - elements removed from the old slice
//	common - elements existed in both old slice and new slice
//
// TO DO - error handling
func GetStringSliceDiff(oldSlice, newSlice []string) (added []string, removed []string, common []string, err error) {

	for _, value := range oldSlice {
		idx := slices.IndexFunc(newSlice, func(c string) bool { return strings.ToLower(c) == strings.ToLower(value) })
		if idx < 0 {
			removed = append(removed, value)
		} else {
			common = append(common, value)
		}
	}
	for _, value := range newSlice {
		idx := slices.IndexFunc(oldSlice, func(c string) bool { return strings.ToLower(c) == strings.ToLower(value) })
		if idx < 0 {
			added = append(added, value)
		}
	}
	//fmt.Printf("%v + %v - %v = %v\n", oldSlice, added, removed, newSlice)
	return added, removed, common, nil
}

// MapKeysToSlice extract keys of map as slice,
// example: https://go.dev/play/p/TSFAD0M5sf5
func MapKeysToSlice[K comparable, V any](m map[K]V) []K {
	keys := make([]K, len(m))

	i := 0
	for k := range m {
		//fmt.Printf("%v = %v\n", i, k)
		keys[i] = k
		i++
	}
	return keys
}

func IsIqn(value string) (isIqn bool) {
	index := strings.Index(strings.ToLower(value), "iqn.")
	if index == 0 {
		isIqn = true
	} else {
		isIqn = false
	}
	return isIqn
}

func IsWwn(value string) bool {
	index := strings.Index(strings.ToLower(value), "wwwn.")
	if index == 0 {
		return true
	} else {
		return false
	}
}

func RemoveDuplicateFromStringArray(input []string) (output []string) {
	bucket := make(map[string]bool)
	var result []string
	for _, str := range input {
		if _, ok := bucket[str]; !ok {
			bucket[str] = true
			result = append(result, str)
		}
	}
	return result
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func IsValidName(name string) bool {
	if len(name) > 255 {
		return false
	}
	return true
}

func ConvertInterfaceToSlice(interface_obj []interface{}) []string {
	new_slice := make([]string, 0)
	for _, v := range interface_obj {
		new_slice = append(new_slice, v.(string))
	}
	return new_slice
}

func CapitalizeFirst(s string) string {
    if s == "" {
        return s
    }
    runes := []rune(s)
    runes[0] = unicode.ToUpper(runes[0])
    return string(runes)
}

// IntToHexString converts an int to uppercase hex with 0x prefix.
func IntToHexString(v int) string {
	return fmt.Sprintf("0x%X", v)
}

// HexStringToInt converts hex (with or without 0x prefix) to int.
func HexStringToInt(s string) (int, error) {
	s = strings.TrimSpace(s)

	// Allow 0x or 0X prefixes
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}

	v, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid hex string '%s': %w", s, err)
	}

	return int(v), nil
}

// ParseLdev: accepts either *ldevID or *ldevHex
func ParseLdev(ldevID *int, ldevHex *string) (int, error) {
	// If ldevID is provided, it wins
	if ldevID != nil {
		return *ldevID, nil
	}

	// Else try hex
	if ldevHex != nil && *ldevHex != "" {
		return HexStringToInt(*ldevHex)
	}

	return 0, fmt.Errorf("either ldev_id or ldev_id_hex must be provided")
}

// helper to convert a value to a pointer
func Ptr[T any](v T) *T { return &v }

