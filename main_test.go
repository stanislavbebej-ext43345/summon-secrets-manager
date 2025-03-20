package main

import (
	"os"
	"testing"

	"github.com/bitwarden/sdk-go"
)

const (
	SECRET_ID = "acd2d25f-1fd2-4604-9a6e-b2a600f71a31"
)

func TestMainFunc(t *testing.T) {
	// Redirect standard out to null
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
		os.Args = os.Args[:len(os.Args)-1]
	}()
	os.Stdout = os.NewFile(0, os.DevNull)

	os.Args = append(os.Args, SECRET_ID)
	main()
}

func TestVersionFlag(t *testing.T) {
	// Redirect standard out to null
	stdout := os.Stdout
	defer func() {
		os.Stdout = stdout
		os.Args = os.Args[:len(os.Args)-1]
	}()
	os.Stdout = os.NewFile(0, os.DevNull)

	os.Args = append(os.Args, "-V")
	main()
}

func TestFindSecret(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"", "API error: Invalid command value: UUID parsing failed: invalid length: expected length 32 for simple format, found 0"},
		{SECRET_ID + ":undefined", ""},
		{SECRET_ID + ":VALUE", ""},
		{SECRET_ID, ""},
	}

	for _, test := range tests {
		_, got := findSecret(test.input)
		if got != nil && got.Error() != test.want {
			t.Errorf("findSecret(%q) = %q, wanted %q", test.input, got, test.want)
		}
	}
}

func TestEnvFindSecret(t *testing.T) {
	var tests = []struct {
		input             string
		want              string
		envBwsApiUrl      string
		envBwsIdentityUrl string
	}{
		{SECRET_ID, "", "", ""},
		{SECRET_ID, "", DEFAULT_BWS_API_URL, DEFAULT_BWS_IDENTITY_URL},
		{SECRET_ID, "API error: builder error", "x", "x"},
		{SECRET_ID, "API error: Received error message from server: [400 Bad Request] {\"error\":\"invalid_client\"}", DEFAULT_BWS_API_URL, "https://identity.bitwarden.com"},
		{SECRET_ID, "API error: Received error message from server: [401 Unauthorized] ", "https://api.bitwarden.com", DEFAULT_BWS_IDENTITY_URL},
	}

	for _, test := range tests {
		os.Setenv(INPUT_BWS_API_URL, test.envBwsApiUrl)
		os.Setenv(INPUT_BWS_IDENTITY_URL, test.envBwsIdentityUrl)

		_, got := findSecret(test.input)
		if got != nil && got.Error() != test.want {
			t.Errorf("findSecret(%q) = %q, wanted %q", test.input, got, test.want)
		}
	}
}

func TestParseSecretId(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"a:b", "a-b"},
		{"a", "a-Value"},
	}

	for _, test := range tests {
		id, key := parseSecretId(test.input)
		if got := id + "-" + key; got != test.want {
			t.Errorf("parseSecretId(%q) = %q, wanted %q", test.input, got, test.want)
		}
	}
}

func TestSecretValue(t *testing.T) {
	const (
		ENTRY_ID    = SECRET_ID
		ENTRY_KEY   = "user123"
		ENTRY_VALUE = "default123"
		ENTRY_NOTE  = "description"
	)

	var (
		entry = sdk.SecretResponse{
			ID:    ENTRY_ID,
			Key:   ENTRY_KEY,
			Note:  ENTRY_NOTE,
			Value: ENTRY_VALUE,
		}

		tests = []struct {
			input string
			want  string
		}{
			{"id", ENTRY_VALUE},
			{"key", ENTRY_KEY},
			{"note", ENTRY_NOTE},
			{"value", ENTRY_VALUE},
			{"VALUE", ENTRY_VALUE},
		}
	)

	for _, test := range tests {
		got := secretValue(&entry, test.input)
		if got != test.want {
			t.Errorf("secretValue(%q) = %q, wanted %q", test.input, got, test.want)
		}
	}
}
