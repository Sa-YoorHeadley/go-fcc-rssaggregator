package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Extracts API Key from headers
//  Example
// Authorization: ApiKey (*****)
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == ""{
		return "", errors.New("no authentication info found")
	}

	authHeader := strings.Split(val, " ")
	
	if(len(authHeader) != 2){
		return "", errors.New("malformed authentication header")
	}

	if(authHeader[0] != "ApiKey"){
		return "", errors.New("malformed authentication header")
	}

	return authHeader[1], nil
}
