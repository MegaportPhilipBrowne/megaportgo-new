package megaport

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/megaport/megaportgo/shared"
	"github.com/stretchr/testify/assert"
)

var accessKey string
var secretKey string

var logger *DefaultLogger

var megaportClient *Client

const (
	MEGAPORTURL = "https://api-staging.megaport.com/"
)

func TestMain(m *testing.M) {
	logger = NewDefaultLogger()
	logger.SetLevel(DebugLevel)

	accessKey = os.Getenv("MEGAPORT_ACCESS_KEY")
	secretKey = os.Getenv("MEGAPORT_SECRET_KEY")

	logLevel := os.Getenv("LOG_LEVEL")

	fmt.Println(logLevel)
	if logLevel != "" {
		logger.SetLevel(StringToLogLevel(logLevel))
	}

	httpClient := NewHttpClient()
	baseURL, err := url.Parse(MEGAPORTURL)
	if err != nil {
		log.Fatalf("invalid base URL: %s", MEGAPORTURL)
	}
	megaportClient = NewClient(httpClient, baseURL)
	os.Exit(m.Run())
}

func TestLoginOauth(t *testing.T) {
	if accessKey == "" {
		logger.Error("MEGAPORT_ACCESS_KEY environment variable not set.")
		os.Exit(1)
	}

	if secretKey == "" {
		logger.Error("MEGAPORT_SECRET_KEY environment variable not set.")
		os.Exit(1)
	}

	ctx := context.Background()
	token, loginErr := megaportClient.AuthenticationService.LoginOauth(ctx, accessKey, secretKey)
	if loginErr != nil {
		logger.Error("login error", "error", loginErr.Error())
	}
	assert.NoError(t, loginErr)

	// Session Token is not empty
	assert.NotEmpty(t, token)
	// SessionToken is a valid guid
	assert.NotNil(t, shared.IsGuid(token))

	logger.Info("", "token", token)
	megaportClient.SessionToken = token
}
