package sanstorage

import (
	"encoding/json"
	"fmt"

	// "strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
	// sanmodel "terraform-provider-hitachi/hitachi/storage/san/model"
)

func GetUrl(ip string, urlPath string) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	url := fmt.Sprintf("https://%s/ConfigurationManager/v1/%s", ip, urlPath)
	log.WriteDebug("TFDebug|url: %s", url)
	return url
}

// token timeout is 300 sec
func GetAuthTokenNoCache(mgmtIP, username, password string) (string, error) {
	// curl -k -v -H "Accept:application/json" -H "Content-Type:application/json" -u user:passw  -X POST
	// https://mgmtIP/ConfigurationManager/v1/objects/sessions/ -d ""
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	creds := map[string]string{
		"username": username,
		"password": password,
	}

	url := GetUrl(mgmtIP, "objects/sessions/")
	resJSONString, err := utils.HTTPPostWithCreds(url, &creds, nil, nil) // no additional headers, no request body
	if err != nil {
		log.WriteError(err)
		log.WriteDebug("TFError| error in HTTPPostWithCreds call, err: %v", err)
		return "", err
	}

	log.WriteDebug("TFDebug|resJSONString: %s", resJSONString)

	type ResponseSession struct {
		Token     string `json:"token"`
		SessionId int    `json:"sessionId"`
	}

	var responseSession ResponseSession

	err2 := json.Unmarshal([]byte(resJSONString), &responseSession)
	if err2 != nil {
		log.WriteDebug("TFError| error in Unmarshal call, err: %v", err2)
		return "", fmt.Errorf("failed to unmarshal json response: %+v", err2)
	}

	return responseSession.Token, nil
}

// using cache
func GetAuthToken(mgmtIP, username, password string) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	cachekey := fmt.Sprintf("%s_%s_%s_%s", mgmtIP, username, password, "Token")

	// get from cache
	value, found := cacheStorage.Get(cachekey)
	log.WriteDebug("TFDebug|CACHEKEY FOUND: %v\n", found)

	if found {
		return value.(string), nil
	}

	token, err := GetAuthTokenNoCache(mgmtIP, username, password)
	if err != nil {
		log.WriteDebug("TFError| error in GetAuthTokenNoCache call, err: %v", err)
		return "", err
	}

	// store to cache
	// token last 300 secs, but store only in cache for 270 secs
	cacheStorage.Set(cachekey, token, 270*time.Second)

	return token, nil
}

func GetAuthTokenHeader(mgmtIP, username, password string) (headers map[string]string, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	token, err := GetAuthToken(mgmtIP, username, password)
	if err != nil {
		log.WriteDebug("TFError| error in GetAuthToken call, err: %v", err)
		return nil, err
	}

	headers = map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Session %s", token)
	return headers, nil
}
