package admin

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/singleflight"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	"terraform-provider-hitachi/hitachi/common/utils"
)

var tokenGroup singleflight.Group // prevent duplicate token refreshes

func GetUrl(ip string, urlPath string) string {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	url := fmt.Sprintf("https://%s/ConfigurationManager/simple/v1/%s", ip, urlPath)
	log.WriteDebug("TFDebug|url: %s", url)
	return url
}

// token timeout is 300 sec
func GetAuthTokenNoCache(mgmtIP, username, password string) (string, error) {
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
		err := CheckHttpErrorResponse(resJSONString, err)
		return "", err
	}

	log.WriteDebug("TFDebug|resJSONString: %s", resJSONString)

	type ResponseSession struct {
		// Token     string `json:"token"`
		SessionId string `json:"sessionId"`
	}

	var responseSession ResponseSession

	err2 := json.Unmarshal([]byte(resJSONString), &responseSession)
	if err2 != nil {
		log.WriteDebug("TFError| error in Unmarshal call, err: %v", err2)
		return "", fmt.Errorf("failed to unmarshal json response: %+v", err2)
	}

	log.WriteInfo("Successfully obtained new session token.")
	return responseSession.SessionId, nil
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

	// Prevent duplicate refresh across goroutines
	v, err, _ := tokenGroup.Do(cachekey, func() (interface{}, error) {
		token, err := GetAuthTokenNoCache(mgmtIP, username, password)
		if err != nil {
			return "", err
		}

		// Add small random jitter (avoid simultaneous expiry refresh)
		jitter := time.Duration(rand.Intn(30)) * time.Second
		cacheStorage.Set(cachekey, token, cacheStorageDuration+jitter)

		log.WriteDebug("TFDebug|Stored new token in cache for %v (+%v jitter)", cacheStorageDuration, jitter)
		return token, nil
	})

	if err != nil {
		return "", err
	}
	return v.(string), nil
}

func GetAuthTokenHeader(mgmtIP, username, password string) (headers map[string]string, err error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	token, err := GetAuthToken(mgmtIP, username, password)
	if err != nil {
		return nil, err
	}

	headers = map[string]string{}
	headers["Authorization"] = fmt.Sprintf("Session %s", token)
	return headers, nil
}
