package configparser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"sync"
)

func resolvePath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	if !strings.Contains(path, "~") {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %v", err)
	}

	homeDir := usr.HomeDir
	if homeDir == "" {
		return "", fmt.Errorf("home directory not found for current user")
	}

	resolvedPath := strings.Replace(path, "~", homeDir, 1)
	return resolvedPath, nil
}

func resolvePathsInMap(configMap map[string]interface{}) error {
	for key, value := range configMap {
		switch v := value.(type) {
		case string:
			if strings.Contains(v, "~") {
				resolvedPath, err := resolvePath(v)
				if err != nil {
					return err
				}
				configMap[key] = resolvedPath
			}
		case map[string]interface{}:
			if err := resolvePathsInMap(v); err != nil {
				return err
			}
		}
	}
	return nil
}

var (
	configInstance     map[string]interface{}
	configInstanceOnce sync.Once
)

func loadConfig() error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("error parsing JSON data: %v", err)
	}

	if err := resolvePathsInMap(config); err != nil {
		return err
	}

	configInstance = config
	return nil
}

func GetConfigInstance() map[string]interface{} {
	configInstanceOnce.Do(func() {
		if err := loadConfig(); err != nil {
			log.Fatalf("Error getting config instance: %v", err)
		}
	})
	return configInstance
}
