package helper

import (
	"encoding/base64"
	"errors"
	"strings"
)

// parsing basic auth "authoriztion header"
func ParsebasicAuth(authHeader string) (username, password string, err error) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", errors.New("invalid authorization format")
	}
	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", errors.New("failed to decode base 64")
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid Basjic auth credentials")
	}
	return parts[0], parts[1], nil

}
