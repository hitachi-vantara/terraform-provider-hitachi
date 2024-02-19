package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"

	// "log"
	b64 "encoding/base64"

	log "github.com/romana/rlog"
)

func GetOutputValue(value interface{}) interface{} {
	valArray, okArray := value.([]interface{})
	if okArray {
		// log.Infof("It's an array : %+v\n", valArray)
		outArray := make([]interface{}, 0)
		for _, item := range valArray {
			outMap := GetOutputValue(item)
			outArray = append(outArray, outMap)

		}
		return outArray
	} else {
		valMap, okMap := value.(map[string]interface{})
		if okMap {
			// log.Infof("It's a map : %+v\n", valMap)
			outMap := map[string]interface{}{}
			for k, v := range valMap {
				outMap[ToSnakeCase(k)] = GetOutputValue(v)
			}
			return outMap
		} else {
			// log.Infof("It's primitive type : %+v\n", value)
			return value
		}
	}
}

func PopulateOutput(value interface{}) []interface{} {
	lastValueInArray := make([]interface{}, 0)
	lastValue := GetOutputValue(value)
	valMap, okMap := lastValue.(map[string]interface{})
	if okMap {
		// log.Infof("Last value is a map : %+v\n", lastValue)
		lastValueInArray = append(lastValueInArray, valMap)
	} else {
		// log.Infof("Last value is an array : %+v\n", lastValue)
		lastValueInArray = lastValue.([]interface{})
	}

	log.Infof("Output Schema Value: %s\n", ConvertToJson(lastValueInArray))
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
	switch x := v.(type) {
	case []interface{}:
		log.Debug("[]interface, len:", len(x))
		return true
	case interface{}:
		log.Debug("interface:", x)
		return false
	default:
		log.Debugf("Unsupported type: %T\n", x)
		return false
	}
}

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

func ConvertToJson(m interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		log.Error(err)
	}
	return string(b)
}

func SaveOutputToFile(resourceName string, id string, output interface{}) error {
	name := "output/" + resourceName + "_" + id + ".out"
	_ = os.Mkdir("output", os.ModeDir|os.ModePerm)

	file, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		log.Error(err)
		return err
	}

	err = ioutil.WriteFile(name, file, 0644)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func DecodeBase64EncodedString(strvalue string) (string, error) {
	var x string
	// fmt.Println(strvalue)
	uDec, err := b64.URLEncoding.DecodeString(strvalue)
	// fmt.Println(string(uDec))
	if err != nil {
		log.Error(err)
		return x, err
	}

	return strings.TrimSuffix(bytes.NewBuffer(uDec).String(), "\n"), nil
}

func SaveDataToFile(dirPath, fileName string, data interface{}) error {
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

func IsValidPortName(name string) bool {
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

func MakeResponse(resp http.Response) (string, error) {

	var ErrorModel *PorcelainError
	var BadRequestError *BadRequestError
	var message string
	body, bodyErr := io.ReadAll(resp.Body)

	if IsHttpError(resp.StatusCode) {

		if bodyErr == nil {
			err := json.Unmarshal(body, &ErrorModel)
			if err != nil {
				log.Debug(err)
				return "", fmt.Errorf("%v", resp.Status)
			}

			if ErrorModel.Error.Message != "" {
				message = ErrorModel.Error.Message
			} else if ErrorModel.Message != "" {
				message = ErrorModel.Message
			} else {

				err := json.Unmarshal(body, &BadRequestError)
				if err != nil {
					log.Debug(err)
					return "", fmt.Errorf("%v", resp.Status)
				}
				message = fmt.Sprintf("%s: %s", BadRequestError.Title, BadRequestError.Detail)
			}

			return "", fmt.Errorf("%s", message)
		}
		return "", fmt.Errorf("%v", resp.Status)

	}
	if bodyErr != nil {
		log.Error(bodyErr)
		return "", bodyErr
	}

	log.Debugf("HTTP Response: %s\n", string(body))
	return string(body), nil

}
