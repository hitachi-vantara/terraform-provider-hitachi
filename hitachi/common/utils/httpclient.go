package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	config "terraform-provider-hitachi/hitachi/common/config"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
)

var (
	sharedClient      *http.Client
	sharedClientOnce  sync.Once
	noKeepAliveClient *http.Client
	noKeepAliveOnce   sync.Once
)

// SharedClient returns a singleton, reusable HTTP client
func SharedClient() *http.Client {
	sharedClientOnce.Do(func() {
		log := commonlog.GetLogger()
		log.WriteEnter()
		defer log.WriteExit()

		apiTimeout := config.DEFAULT_API_TIMEOUT
		if config.ConfigData != nil && config.ConfigData.APITimeout > 0 {
			apiTimeout = config.ConfigData.APITimeout
		}

		timeout := time.Duration(apiTimeout) * time.Second
		sharedClient = newHTTPClient(timeout, true)

		log.WriteDebug("API Execution Timeout: %v", timeout)
	})

	return sharedClient
}

func NoKeepAliveClient() *http.Client {
	noKeepAliveOnce.Do(func() {
		log := commonlog.GetLogger()
		log.WriteEnter()
		defer log.WriteExit()

		apiTimeout := config.DEFAULT_API_TIMEOUT
		if config.ConfigData != nil && config.ConfigData.APITimeout > 0 {
			apiTimeout = config.ConfigData.APITimeout
		}

		timeout := time.Duration(apiTimeout) * time.Second

		tr := &http.Transport{
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives:     true,
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   0,
			IdleConnTimeout:       1 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		noKeepAliveClient = &http.Client{
			Timeout:   timeout,
			Transport: tr,
		}

		log.WriteDebug("API Execution Timeout: %v", timeout)
	})

	return noKeepAliveClient
}

// Helper function to build a configured client
func newHTTPClient(timeout time.Duration, skipTLSVerify bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: skipTLSVerify},
		MaxIdleConns:          200, // Reuses connections across requests
		MaxIdleConnsPerHost:   200,
		IdleConnTimeout:       90 * time.Second, // Keeps unused connections alive
		DisableKeepAlives:     false,             // Enables connection reuse
		TLSHandshakeTimeout:   10 * time.Second,  // Gives up TLS setup after 10s
		ExpectContinueTimeout: 1 * time.Second,   // Waits 1s for "100 Continue"
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}
}

func IsHttpError(statusCode int) bool {
	if statusCode >= 400 && statusCode <= 599 {
		return true
	}
	return false
}

// executeRequestWithRetry handles the Do() call and retries on EOF
func executeRequestWithRetry(client *http.Client, req *http.Request) (*http.Response, error) {
    log := commonlog.GetLogger()
    resp, err := client.Do(req)

    // Retry once if we hit an EOF (common in concurrent Hitachi API calls due to keep-alive races)
    if err != nil && (errors.Is(err, io.EOF) || strings.Contains(err.Error(), "EOF")) {
        log.WriteInfo("Transient EOF detected during %s, retrying with fresh connection...", req.Method)
        req.Close = true // Force fresh connection for the retry
        return client.Do(req)
    }
    return resp, err
}

func HTTPGet(url string, headers *map[string]string, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if len(basicAuthentication) > 0 && basicAuthentication[0] != nil {
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
			req.Header.Set(k, v)
		}
	}

	logRequest(req, nil)

	resp, err := executeRequestWithRetry(SharedClient(), req)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	log.WriteDebug("HTTP Response: %s\n", string(body))

	if IsHttpError(resp.StatusCode) {
		io.Copy(io.Discard, resp.Body) // keep-alive safe
		return string(body), fmt.Errorf("%v", resp.Status)
	}

	return string(body), nil
}

func HTTPPost(url string, headers *map[string]string, httpBody []byte, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(httpBody))
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if len(basicAuthentication) > 0 && basicAuthentication[0] != nil {
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
			req.Header.Set(k, v)
		}
	}

	logRequest(req, httpBody)

	resp, err := executeRequestWithRetry(SharedClient(), req)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	log.WriteDebug("HTTP Response: %s\n", string(body))

	if IsHttpError(resp.StatusCode) {
		io.Copy(io.Discard, resp.Body) // keep-alive safe
		return string(body), fmt.Errorf("%v", resp.Status)
	}

	return string(body), nil
}

// this is used for session api call, it should use NoKeepAliveClient
// because of UCT-588. Disable keep-alive for session api
func HTTPPostWithCreds(url string, creds *map[string]string, headers *map[string]string, httpBody []byte) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(httpBody))
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	if creds != nil {
		req.SetBasicAuth((*creds)["username"], (*creds)["password"])
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}

	logRequest(req, httpBody)

	resp, err := NoKeepAliveClient().Do(req)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	if IsHttpError(resp.StatusCode) {
		io.Copy(io.Discard, resp.Body) // keep-alive safe
		return string(body), fmt.Errorf("%v", resp.Status)
	}

	log.WriteDebug("HTTP Response: %s\n", string(body))
	return string(body), nil
}

func HTTPPostNoBody(url string, headers *map[string]string, httpBody []byte, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var bodyReader io.Reader
	if httpBody != nil {
		bodyReader = bytes.NewBuffer(httpBody)
	} else {
		bodyReader = bytes.NewReader([]byte{}) // empty body
	}

	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	if httpBody != nil {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.ContentLength = 0
	}

	// Add auth
	if len(basicAuthentication) > 0 && basicAuthentication[0] != nil {
		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}

	logRequest(req, httpBody)

	resp, err := executeRequestWithRetry(SharedClient(), req)
	if err != nil {
		log.WriteError(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	log.WriteDebug("HTTP Response: %s\n", string(body))

	if IsHttpError(resp.StatusCode) {
		io.Copy(io.Discard, resp.Body) // keep-alive safe
		return string(body), fmt.Errorf("%v", resp.Status)
	}

	return string(body), nil
}

func HTTPDelete(url string, headers *map[string]string, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if len(basicAuthentication) > 0 && basicAuthentication[0] != nil {
		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}

	logRequest(req, nil)

    resp, err := executeRequestWithRetry(SharedClient(), req)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	if IsHttpError(resp.StatusCode) {
		io.Copy(io.Discard, resp.Body) // keep-alive safe
		return string(body), fmt.Errorf("%v", resp.Status)
	}

	log.WriteDebug("HTTP Response: %s\n", string(body))
	return string(body), nil
}

func HTTPDeleteWithBody(url string, headers *map[string]string, httpBody []byte, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	var req *http.Request
	var err error
	nullValue := []byte("null")
	if reflect.DeepEqual(httpBody, nullValue) {
		req, err = http.NewRequest("DELETE", url, nil)
	} else {
		req, err = http.NewRequest("DELETE", url, bytes.NewBuffer(httpBody))
	}

	if err != nil {
		log.WriteError(err)
		return "", err
	}

	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if len(basicAuthentication) > 0 && basicAuthentication[0] != nil {
		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	logRequest(req, httpBody)

	resp, err := executeRequestWithRetry(SharedClient(), req)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	if IsHttpError(resp.StatusCode) {
		io.Copy(io.Discard, resp.Body) // keep-alive safe
		return string(body), fmt.Errorf("%v", resp.Status)
	}

	log.WriteDebug("HTTP Response: %s\n", string(body))
	return string(body), nil
}

func HTTPPatch(url string, headers *map[string]string, httpBody []byte, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(httpBody))
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Add auth
	if len(basicAuthentication) > 0 && basicAuthentication[0] != nil {
		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}

	logRequest(req, httpBody)

	resp, err := executeRequestWithRetry(SharedClient(), req)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	if IsHttpError(resp.StatusCode) {
		io.Copy(io.Discard, resp.Body) // keep-alive safe
		return string(body), fmt.Errorf("%v", resp.Status)
	}

	log.WriteDebug("HTTP Response: %s\n", string(body))
	return string(body), nil
}

// HTTPDownloadFile downloads a file via GET and saves it to the specified directory.
// It uses basic authentication and optionally reads the filename from the Content-Disposition header.
func HTTPDownloadFile(url string, toFilePath string, headers *map[string]string, basicAuthentication ...*HttpBasicAuthentication) (string, error) {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	// Basic auth
	if len(basicAuthentication) > 0 && basicAuthentication[0] != nil {
		req.SetBasicAuth(basicAuthentication[0].Username, basicAuthentication[0].Password)
	}

	// Optional headers
	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}

	logRequest(req, nil)

	resp, err := executeRequestWithRetry(SharedClient(), req)
	if err != nil {
		log.WriteError(err)
		return "", err
	}
	defer resp.Body.Close()

	if IsHttpError(resp.StatusCode) {
		bodyBytes, _ := io.ReadAll(resp.Body)
		io.Copy(io.Discard, resp.Body) // ensure keep-alive connection can be reused
		log.WriteError("HTTP error: %s", string(bodyBytes))
		return string(bodyBytes), fmt.Errorf("%v", resp.Status)
	}

	// Try to extract filename from Content-Disposition header
	filename := "downloaded_file"
	if cd := resp.Header.Get("Content-Disposition"); strings.Contains(cd, "filename=") {
		parts := strings.Split(cd, "filename=")
		if len(parts) > 1 {
			filename = strings.Trim(parts[1], `"`)
		}
	}

	// Determine the final file path
	var finalPath string
	if toFilePath == "" {
		// Behave like curl -O
		finalPath = filename
	} else {
		info, err := os.Stat(toFilePath)
		if err == nil && info.IsDir() {
			// It's an existing directory: secure join
			finalPath, err = secureJoin(toFilePath, filename)
			if err != nil {
				log.WriteError(err)
				return "", err
			}
		} else if strings.HasSuffix(toFilePath, string(os.PathSeparator)) {
			// Ends with / or \ but does not exist yet: treat as directory
			if err := os.MkdirAll(toFilePath, 0755); err != nil {
				log.WriteError(err)
				return "", err
			}
			finalPath, err = secureJoin(toFilePath, filename)
			if err != nil {
				log.WriteError(err)
				return "", err
			}
		} else {
			// It's a full file path
			finalPath = toFilePath
		}
	}

	// Write to file
	outFile, err := os.Create(finalPath)
	if err != nil {
		log.WriteError(err)
		return "", err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	log.WriteInfo("File downloaded successfully: %s", finalPath)
	return finalPath, nil
}

func logRequest(req *http.Request, reqBodyInBytes []byte) {
	log := commonlog.GetLogger()

	// bodyBytes := []byte{}
	// if req.Body != nil {
	// 	var err error
	// 	bodyBytes, err = io.ReadAll(req.Body)
	// 	if err != nil {
	// 		log.WriteError("Error reading request body: %v", err)
	// 		return
	// 	}

	// 	// Since we've read the body, we need to reset it so that the server can read it too
	// 	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	// }

	maskedReq := ""
	if reqBodyInBytes != nil {
		maskedReq, _ = MaskSensitiveData(string(reqBodyInBytes))
	}

	log.WriteDebug("Method: %s", req.Method)
	log.WriteDebug("URL: %s", req.URL.String())
	log.WriteDebug("Headers: %v", sanitizeHeaders(req.Header))
	// log.WriteDebug("Body: %s", string(bodyBytes))
	log.WriteDebug("Body: %s", string(maskedReq))
	log.WriteDebug("ContentLength: %d", req.ContentLength)
	log.WriteDebug("Host: %s", req.Host)
	log.WriteDebug("RequestURI: %s", req.RequestURI)
}

func HTTPPostForm(
	url string, headers *map[string]string, httpBody []byte,
	form bytes.Buffer,
	basicAuthentication ...*HttpBasicAuthentication,
) (string, error) {

	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(httpBody))
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")

	// Add auth
	if len(basicAuthentication) > 0 && basicAuthentication[0] != nil {
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
			req.Header.Set(k, v)
		}
	} else {
		// the above is doing Add, so we need to set the content type here
		req.Header.Set("Content-Type", "application/json")
	}

	logRequest(req, httpBody)

	resp, err := executeRequestWithRetry(SharedClient(), req)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WriteError(err)
		return "", err
	}

	log.WriteDebug("HTTP Response: %s\n", string(body))

	if IsHttpError(resp.StatusCode) {
		io.Copy(io.Discard, resp.Body) // keep-alive safe
		return string(body), fmt.Errorf("%v", resp.Status)
	}

	return string(body), nil
}

// secureJoin ensures finalPath is inside baseDir to prevent path traversal
func secureJoin(baseDir, filename string) (string, error) {
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}

	targetPath := filepath.Join(baseAbs, filename)
	targetAbs, err := filepath.Abs(targetPath)
	if err != nil {
		return "", err
	}

	// Ensure the target path is within the base directory
	if !strings.HasPrefix(targetAbs, baseAbs+string(os.PathSeparator)) && targetAbs != baseAbs {
		return "", errors.New("path traversal detected in filename")
	}

	return targetAbs, nil
}

// mask headers in logs
func sanitizeHeaders(h http.Header) http.Header {
	clone := h.Clone()
	for key := range clone {
		if strings.EqualFold(key, "Authorization") ||
			strings.Contains(strings.ToLower(key), "password") {
			clone.Set(key, "***MASKED***")
		}
	}
	return clone
}
