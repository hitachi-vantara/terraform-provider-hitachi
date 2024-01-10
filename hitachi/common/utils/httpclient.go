package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	log "github.com/romana/rlog"
)

var TIMEOUT time.Duration = 120

func IsHttpError(statusCode int) bool {
	if statusCode >= 400 && statusCode <= 599 {
		return true
	}
	return false
}

func HTTPGet(url string, headers *map[string]string, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   TIMEOUT * time.Second,
		Transport: tr,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if basicAuthentication != nil {
		/*

			usernameDecoded, err := DecodeBase64EncodedString(basicAuthentication[0].Username)
			if err != nil {
				return "", GetFormatedError(err, "DecodeBase64EncodedString:", "Failed")
			}

			passwordDecoded, err := DecodeBase64EncodedString(basicAuthentication[0].Password)
			if err != nil {
				return "", GetFormatedError(err, "DecodeBase64EncodedString:", "Failed")
			}

			req.SetBasicAuth(usernameDecoded, passwordDecoded)
		*/

		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	if headers != nil {
		for k, v := range *headers {
			strValue := fmt.Sprintf("header key=[%s], value=[%s]", k, v)
			log.Infof(strValue)
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer resp.Body.Close()

	if IsHttpError(resp.StatusCode) {
		return "", fmt.Errorf("%v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Debugf("HTTP Response: %s\n", string(body))
	return string(body), nil
}

func HTTPPost(url string, headers *map[string]string, httpBody []byte, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   TIMEOUT * time.Second,
		Transport: tr,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(httpBody))
	if err != nil {
		log.Error(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if basicAuthentication != nil {
		/*

			usernameDecoded, err := DecodeBase64EncodedString(basicAuthentication[0].Username)
			if err != nil {
				return "", GetFormatedError(err, "DecodeBase64EncodedString:", "Failed")
			}

			passwordDecoded, err := DecodeBase64EncodedString(basicAuthentication[0].Password)
			if err != nil {
				return "", GetFormatedError(err, "DecodeBase64EncodedString:", "Failed")
			}

			req.SetBasicAuth(usernameDecoded, passwordDecoded)
		*/
		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	if headers != nil {
		for k, v := range *headers {
			strValue := fmt.Sprintf("header key=[%s], value=[%s]", k, v)
			log.Infof(strValue)
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer resp.Body.Close()

	if IsHttpError(resp.StatusCode) {
		return "", fmt.Errorf("%v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Debugf("HTTP Response: %s\n", string(body))
	return string(body), nil
}

func HTTPPostWithCreds(url string, creds *map[string]string, headers *map[string]string, httpBody []byte) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   TIMEOUT * time.Second,
		Transport: tr,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(httpBody))
	if err != nil {
		log.Error(err)
		return "", err
	}

	if creds != nil {
		req.SetBasicAuth((*creds)["username"], (*creds)["password"])
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if headers != nil {
		for k, v := range *headers {
			strValue := fmt.Sprintf("header key=[%s], value=[%s]", k, v)
			log.Infof(strValue)
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer resp.Body.Close()

	if IsHttpError(resp.StatusCode) {
		return "", fmt.Errorf("%v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Debugf("HTTP Response: %s\n", string(body))
	return string(body), nil
}

func HTTPDelete(url string, headers *map[string]string, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Error(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if basicAuthentication != nil {
		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	if headers != nil {
		for k, v := range *headers {
			strValue := fmt.Sprintf("header key=[%s], value=[%s]", k, v)
			log.Infof(strValue)
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer resp.Body.Close()

	if IsHttpError(resp.StatusCode) {
		return "", fmt.Errorf("%v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Debugf("HTTP Response: %s\n", string(body))
	return string(body), nil
}

func HTTPDeleteWithBody(url string, headers *map[string]string, httpBody []byte, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   TIMEOUT * time.Second,
		Transport: tr,
	}
	var req *http.Request
	var err error
	nullValue := []byte("null")
	if reflect.DeepEqual(httpBody, nullValue) {
		req, err = http.NewRequest("DELETE", url, nil)
	} else {
		req, err = http.NewRequest("DELETE", url, bytes.NewBuffer(httpBody))
	}

	if err != nil {
		log.Error(err)
		return "", err
	}

	if headers != nil {
		for k, v := range *headers {
			strValue := fmt.Sprintf("header key=[%s], value=[%s]", k, v)
			log.Infof(strValue)
			req.Header.Add(k, v)
		}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if basicAuthentication != nil {
		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer resp.Body.Close()

	if IsHttpError(resp.StatusCode) {
		return "", fmt.Errorf("%v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Debugf("HTTP Response: %s\n", string(body))
	return string(body), nil
}

func HTTPPatch(url string, headers *map[string]string, httpBody []byte, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   TIMEOUT * time.Second,
		Transport: tr,
	}
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(httpBody))
	if err != nil {
		log.Error(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if basicAuthentication != nil {
		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	if headers != nil {
		for k, v := range *headers {
			strValue := fmt.Sprintf("header key=[%s], value=[%s]", k, v)
			log.Infof(strValue)
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return "", err
	}

	defer resp.Body.Close()

	if IsHttpError(resp.StatusCode) {
		return "", fmt.Errorf("%v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	log.Debugf("HTTP Response: %s\n", string(body))
	return string(body), nil
}
