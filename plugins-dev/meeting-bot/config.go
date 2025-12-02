package main

import (
	"fmt"
	"strings"
	"sync"
)

// Configuration holds the plugin configuration
type Configuration struct {
	OpenAIAPIKey      string
	EnableDebugLogging bool
}

type pluginConfig struct {
	configuration *Configuration
	mu            sync.RWMutex
}

var config = &pluginConfig{
	configuration: &Configuration{},
}

// OnConfigurationChange is called when the plugin configuration changes
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(Configuration)

	// Load the configuration from Mattermost server
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return fmt.Errorf("failed to load plugin configuration: %w", err)
	}

	// Validate configuration
	if err := validateConfiguration(configuration); err != nil {
		p.API.LogWarn("Configuration validation warning", "error", err.Error())
		// Don't return error - allow plugin to run with invalid config
		// Admin will see warning in logs
	}

	// Update configuration under lock
	config.mu.Lock()
	config.configuration = configuration
	config.mu.Unlock()

	p.API.LogInfo("Configuration reloaded")
	return nil
}

// validateConfiguration validates the plugin configuration
func validateConfiguration(cfg *Configuration) error {
	if cfg.OpenAIAPIKey == "" {
		return fmt.Errorf("OpenAI API key is not configured")
	}

	// Validate API key format (OpenAI keys start with "sk-")
	if !strings.HasPrefix(cfg.OpenAIAPIKey, "sk-") {
		return fmt.Errorf("OpenAI API key format appears invalid (should start with 'sk-')")
	}

	return nil
}

// getOpenAIKey retrieves the OpenAI API key from configuration
// Returns error if key is not configured
func (p *Plugin) getOpenAIKey() (string, error) {
	config.mu.RLock()
	defer config.mu.RUnlock()

	if config.configuration.OpenAIAPIKey == "" {
		return "", fmt.Errorf("OpenAI API key is not configured. Please set it in System Console > Plugins > Meeting Bot > Settings")
	}

	return config.configuration.OpenAIAPIKey, nil
}

// isDebugLoggingEnabled checks if debug logging is enabled
func (p *Plugin) isDebugLoggingEnabled() bool {
	config.mu.RLock()
	defer config.mu.RUnlock()

	return config.configuration.EnableDebugLogging
}

