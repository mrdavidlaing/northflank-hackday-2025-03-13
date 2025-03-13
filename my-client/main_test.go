package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Masterminds/semver/v3"
)

// captureLogOutput captures log output during tests
func captureLogOutput(t *testing.T) *bytes.Buffer {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() {
		log.SetOutput(os.Stderr)
	})
	return &buf
}

func TestVersionParsing(t *testing.T) {
	testCases := []struct {
		name            string
		serverVersion   string
		versionRange    string
		expectCompatible bool
	}{
		{
			name:            "Compatible version with v prefix",
			serverVersion:   "v0.1.1",
			versionRange:    ">=0.1.0",
			expectCompatible: true,
		},
		{
			name:            "Compatible version with dev suffix",
			serverVersion:   "v0.1.1-dev",
			versionRange:    ">=0.1.0",
			expectCompatible: true,
		},
		{
			name:            "Incompatible version",
			serverVersion:   "v0.1.1-dev",
			versionRange:    ">=0.2.0",
			expectCompatible: false,
		},
		{
			name:            "Version range with upper bound",
			serverVersion:   "v0.1.1",
			versionRange:    ">=0.1.0 <0.2.0",
			expectCompatible: true,
		},
		{
			name:            "Version outside upper bound",
			serverVersion:   "v0.2.0",
			versionRange:    ">=0.1.0 <0.2.0",
			expectCompatible: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test server that returns the specified version
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(ServerInfo{Version: tc.serverVersion})
			}))
			defer server.Close()

			// Set up environment
			os.Setenv(envServerURL, server.URL)
			os.Setenv(envSupportedVersions, tc.versionRange)
			
			// Parse the constraint
			constraint, err := semver.NewConstraint(tc.versionRange)
			if err != nil {
				t.Fatalf("Failed to parse constraint: %v", err)
			}

			// Get server info
			info, err := getServerInfo(server.URL)
			if err != nil {
				t.Fatalf("Failed to get server info: %v", err)
			}

			// Parse version
			versionStr := info.Version
			versionStr = strings.TrimPrefix(versionStr, "v")
			versionStr = strings.Split(versionStr, "-")[0]
			
			serverVersion, err := semver.NewVersion(versionStr)
			if err != nil {
				t.Fatalf("Failed to parse version: %v", err)
			}

			// Check compatibility
			isCompatible := constraint.Check(serverVersion)
			if isCompatible != tc.expectCompatible {
				t.Errorf("Expected compatibility to be %v, but got %v for version %s with range %s",
					tc.expectCompatible, isCompatible, tc.serverVersion, tc.versionRange)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	// Test with environment variable set
	const testKey = "TEST_ENV_VAR"
	const testValue = "test_value"
	const defaultValue = "default_value"

	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	result := getEnv(testKey, defaultValue)
	if result != testValue {
		t.Errorf("getEnv returned %s, expected %s", result, testValue)
	}

	// Test with environment variable not set
	const unsetKey = "UNSET_ENV_VAR"
	result = getEnv(unsetKey, defaultValue)
	if result != defaultValue {
		t.Errorf("getEnv returned %s, expected %s", result, defaultValue)
	}
}

func TestGetEnvInt(t *testing.T) {
	// Capture log output
	logBuf := captureLogOutput(t)

	// Test with valid integer
	const testKey = "TEST_INT_VAR"
	const testValue = "42"
	const defaultValue = 10

	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	result := getEnvInt(testKey, defaultValue)
	if result != 42 {
		t.Errorf("getEnvInt returned %d, expected %d", result, 42)
	}

	// Test with invalid integer
	const invalidKey = "INVALID_INT_VAR"
	os.Setenv(invalidKey, "not_an_int")
	defer os.Unsetenv(invalidKey)

	result = getEnvInt(invalidKey, defaultValue)
	if result != defaultValue {
		t.Errorf("getEnvInt returned %d, expected %d", result, defaultValue)
	}

	// Verify the warning message
	logOutput := logBuf.String()
	expectedWarning := "Warning: Invalid value for INVALID_INT_VAR, using default 10"
	if !strings.Contains(logOutput, expectedWarning) {
		t.Errorf("Expected log output to contain '%s', but got '%s'", expectedWarning, logOutput)
	}

	// Test with environment variable not set
	const unsetKey = "UNSET_INT_VAR"
	result = getEnvInt(unsetKey, defaultValue)
	if result != defaultValue {
		t.Errorf("getEnvInt returned %d, expected %d", result, defaultValue)
	}
} 