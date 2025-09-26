package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bitwarden/sdk-go"
)

var (
	BUILD_VERSION = "0.4.1" // x-release-please-version

	DEFAULT_BWS_API_URL      = "https://api.bitwarden.eu"
	DEFAULT_BWS_IDENTITY_URL = "https://identity.bitwarden.eu"

	INPUT_BWS_ACCESS_TOKEN = "BWS_ACCESS_TOKEN"
	INPUT_BWS_API_URL      = "BWS_API_URL"
	INPUT_BWS_IDENTITY_URL = "BWS_IDENTITY_URL"
)

var versionFlag bool

func init() {
	flag.BoolVar(&versionFlag, "V", false, "show version")
	flag.BoolVar(&versionFlag, "version", false, "show version")
}

func main() {
	// Parse input parameters
	flag.Parse()

	if versionFlag {
		fmt.Println(BUILD_VERSION)
		return
	}

	secretId := flag.Arg(0)
	if secretId == "" {
		log.Fatal("secret ID is empty")
	}

	// Find secret value
	secret, err := findSecret(secretId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(secret)
}

func findSecret(secretId string) (secret string, err error) {
	// Client configuration
	apiURL := os.Getenv(INPUT_BWS_API_URL)
	if apiURL == "" {
		apiURL = DEFAULT_BWS_API_URL
	}

	identityURL := os.Getenv(INPUT_BWS_IDENTITY_URL)
	if identityURL == "" {
		identityURL = DEFAULT_BWS_IDENTITY_URL
	}

	bitwardenClient, _ := sdk.NewBitwardenClient(&apiURL, &identityURL)
	defer bitwardenClient.Close()

	// Authentication
	accessToken := os.Getenv(INPUT_BWS_ACCESS_TOKEN)
	err = bitwardenClient.AccessTokenLogin(accessToken, nil)
	if err != nil {
		return
	}

	// Read secret
	secretId, secretKey := parseSecretId(secretId)
	secretEntry, err := bitwardenClient.Secrets().Get(secretId)
	if err != nil {
		return
	}

	return secretValue(secretEntry, secretKey), err
}

func parseSecretId(name string) (string, string) {
	// Split name into entry 'id' and 'key'
	nameSplit := strings.Split(name, ":")
	if len(nameSplit) > 1 {
		return nameSplit[0], nameSplit[1]
	}
	return name, "Value"
}

func secretValue(entry *sdk.SecretResponse, key string) string {
	// https://stackoverflow.com/a/18931036
	// r := reflect.ValueOf(entry)
	// f := reflect.Indirect(r).FieldByName(key)
	// return f.String()

	switch key := strings.ToLower(key); key {
	case "key":
		return entry.Key
	case "note":
		return entry.Note
	default:
		return entry.Value
	}
}
