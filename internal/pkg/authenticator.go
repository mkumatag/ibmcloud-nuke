package pkg

import (
	"fmt"
	"os"

	"github.com/IBM/go-sdk-core/v5/core"
)

const (
	// APIKeyEnv is environment which contains the IBMCloud API Key
	APIKeyEnv = "IBMCLOUD_API_KEY"
)

// GetAuthenticator Returns an authenticator for IBM Cloud
func GetAuthenticator() (core.Authenticator, error) {
	apiKey := os.Getenv(APIKeyEnv)
	if apiKey == "" {
		return nil, fmt.Errorf("please set %s environment, it cannot be empty", APIKeyEnv)
	}

	auth := &core.IamAuthenticator{
		ApiKey: apiKey,
	}
	return auth, nil
}
