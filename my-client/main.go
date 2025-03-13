package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	defaultServerURL      = "http://localhost:8080/info"
	defaultVersionRange   = ">=0.1.0"
	defaultPollInterval   = 10 * time.Second
	envServerURL          = "SERVER_URL"
	envSupportedVersions  = "SUPPORTED_VERSIONS"
	envPollInterval       = "POLL_INTERVAL_SECONDS"
)

// ServerInfo represents the response from the server's /info endpoint
type ServerInfo struct {
	Version string `json:"version"`
}

func main() {
	// Get configuration from environment variables
	serverURL := getEnv(envServerURL, defaultServerURL)
	supportedVersionsRange := getEnv(envSupportedVersions, defaultVersionRange)
	pollIntervalSeconds := getEnvInt(envPollInterval, int(defaultPollInterval.Seconds()))
	pollInterval := time.Duration(pollIntervalSeconds) * time.Second

	// Parse the version constraint
	constraint, err := semver.NewConstraint(supportedVersionsRange)
	if err != nil {
		log.Fatalf("Invalid version range '%s': %v", supportedVersionsRange, err)
	}

	log.Printf("Starting client with configuration:")
	log.Printf("  Server URL: %s", serverURL)
	log.Printf("  Supported versions: %s", supportedVersionsRange)
	log.Printf("  Poll interval: %s", pollInterval)

	// Start polling loop
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// Poll immediately on startup
	checkServerVersion(serverURL, constraint)

	// Then poll on the ticker interval
	for range ticker.C {
		checkServerVersion(serverURL, constraint)
	}
}

// checkServerVersion polls the server and checks if its version is compatible
func checkServerVersion(serverURL string, constraint *semver.Constraints) {
	// Get server info
	info, err := getServerInfo(serverURL)
	if err != nil {
		log.Printf("Error getting server info: %v", err)
		return
	}

	// Parse the server version - remove 'v' prefix if present
	versionStr := info.Version
	versionStr = strings.TrimPrefix(versionStr, "v")
	
	// Handle pre-release versions by removing the -dev suffix for compatibility checking
	versionStr = strings.Split(versionStr, "-")[0]
	
	serverVersion, err := semver.NewVersion(versionStr)
	if err != nil {
		log.Printf("Server returned invalid version '%s': %v", info.Version, err)
		return
	}

	// Check if the server version satisfies our constraint
	if constraint.Check(serverVersion) {
		log.Printf("Server version %s is compatible ✅", info.Version)
	} else {
		log.Printf("Server version %s is NOT compatible ❌", info.Version)
	}
}

// getServerInfo fetches the server info from the /info endpoint
func getServerInfo(serverURL string) (*ServerInfo, error) {
	// Make the HTTP request
	resp, err := http.Get(serverURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status code %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON
	var info ServerInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &info, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvInt gets an environment variable as an integer or returns a default value
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue := 0
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		log.Printf("Warning: Invalid value for %s, using default %d", key, defaultValue)
		return defaultValue
	}

	return intValue
} 