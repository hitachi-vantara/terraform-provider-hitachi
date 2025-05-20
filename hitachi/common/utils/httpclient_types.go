package utils

import (
	"net/http"
	"encoding/base64"
)

// HttpBasicAuthentication is used for authentication
type HttpBasicAuthentication struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Define the SetAuthHeaders method to add Authorization headers
func (auth *HttpBasicAuthentication) SetAuthHeaders(req *http.Request) {
	// Set the basic authentication header
	// Basic auth is typically in the form of "Authorization: Basic <base64-encoded-username-password>"
	// Base64 encode username and password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth.Username + ":" + auth.Password))
	req.Header.Set("Authorization", "Basic "+encoded)
}