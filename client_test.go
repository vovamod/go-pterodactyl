package pterodactyl_test

import (
	"context"
	"github.com/vovamod/go-pterodactyl"
	"testing"
)

func TestNewClientSuccess(t *testing.T) {

	t.Parallel()

	testCases := []struct {
		name          string
		apiKey        string
		keyType       pterodactyl.KeyType
		expectedError bool
	}{
		{
			name:          "Valid ApplicationAPI Key",
			apiKey:        "ptla_abc123", // A dummy key with the correct prefix
			keyType:       pterodactyl.ApplicationKey,
			expectedError: false,
		},
		{
			name:          "Valid ClientAPI Key",
			apiKey:        "ptlc_def456", // A dummy key with the correct prefix
			keyType:       pterodactyl.ClientKey,
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := pterodactyl.NewClient("https://fake-panel.com", tc.apiKey, tc.keyType)

			if tc.expectedError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}

				// On success, we should also check that the client is not nil and its services are initialized.
				if client == nil {
					t.Fatal("expected client to be non-nil on success")
				}
				if client.ApplicationAPI == nil {
					t.Error("expected ApplicationAPI service to be initialized")
				}
				if client.ClientAPI == nil {
					t.Error("expected ClientAPI service to be initialized")
				}
			}
		})
	}
}

// TestNewClient_InvalidKeyFormat checks for errors when API keys have incorrect prefixes.
func TestNewClient_InvalidKeyFormat(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		apiKey        string
		keyType       pterodactyl.KeyType
		expectedError bool
	}{
		{
			name:          "Invalid ApplicationAPI Key",
			apiKey:        "ptlc_wrongprefix", // A client key prefix used for an application client
			keyType:       pterodactyl.ApplicationKey,
			expectedError: true,
		},
		{
			name:          "Invalid ClientAPI Key",
			apiKey:        "ptla_wrongprefix", // An application key prefix used for a client client
			keyType:       pterodactyl.ClientKey,
			expectedError: true,
		},
		{
			name:          "ApplicationAPI Key without prefix",
			apiKey:        "noprefix",
			keyType:       pterodactyl.ApplicationKey,
			expectedError: true,
		},
		{
			name:          "ClientAPI Key without prefix",
			apiKey:        "noprefix",
			keyType:       pterodactyl.ClientKey,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := pterodactyl.NewClient("https://fake-panel.com", tc.apiKey, tc.keyType)

			if !tc.expectedError {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected an error for invalid key format but got none")
				}
			}
		})
	}
}

// TestNewClient_InvalidURL demonstrates testing for a malformed base URL.
// While this case might be less common, it's good practice to ensure robustness.
func TestNewClient_InvalidURL(t *testing.T) {
	t.Parallel()

	// This is a special character that will cause url.Parse to fail in NewRequest
	malformedURL := "::not a valid url"

	client, err := pterodactyl.NewClient(malformedURL, "ptlc_dummykey", pterodactyl.ClientKey)
	if err != nil {
		// The error doesn't happen in NewClient itself, but on the first request.
		// Let's test that by trying to make a request.
		// Since we can't access client.NewRequest directly (it's unexported in our test package),
		// we test a public method that uses it.
		// We'll need to update the client.go NewClient to check the URL at creation time.

		// Let's go back to client.go and improve it first.
		t.Skip("Skipping test: NewClient should validate the baseURL upon creation.")
	}

	_, listErr := client.ClientAPI.ListPermissions(context.Background())
	if listErr == nil {
		t.Errorf("expected an error from a request with a malformed baseURL, but got none")
	}
}
