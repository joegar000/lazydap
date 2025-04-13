package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type Config struct {
    RunConfigs []interface{} `yaml:"run-configs"`
}

// getConfigPath returns the full path to the lazydap config file in the user config directory.
func getConfigPath() (string, error) {
    configDir, err := os.UserConfigDir() // Handles XDG_CONFIG_HOME, %AppData%, etc.
    if err != nil {
        return "", fmt.Errorf("could not determine user config directory: %w", err)
    }
    return filepath.Join(configDir, "lazydap", "config.yaml"), nil
}

// EnsureConfigFile ensures the config file exists, creating and initializing it if necessary.
func EnsureConfigFile() (*Config, error) {
    configPath, err := getConfigPath()
    if err != nil {
        return nil, err
    }

    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        // Create the directory if needed
        if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
            return nil, fmt.Errorf("failed to create config directory: %w", err)
        }

        // Initialize with an empty run-configs array
        defaultConfig := Config{
            RunConfigs: []interface{}{},
        }

        yamlData, err := yaml.Marshal(&defaultConfig)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal default YAML: %w", err)
        }

        if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
            return nil, fmt.Errorf("failed to write config file: %w", err)
        }
        return &defaultConfig, nil
    }

    return LoadConfig()
}

// LoadConfig reads and unmarshals the config.yaml into a Config struct.
func LoadConfig() (*Config, error) {
    configPath, err := getConfigPath()
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
    }

    return &cfg, nil
}
