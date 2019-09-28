package engine

import (
	"encoding/json"
	"github.com/qingcloudhx/core/data/schema"
	"github.com/qingcloudhx/core/support/logEx"
	"io/ioutil"
	"os"

	"github.com/qingcloudhx/core/app"
	"github.com/qingcloudhx/core/engine/secret"
	"github.com/qingcloudhx/core/support"
)

var appName, appVersion string

func init() {
	if IsSchemaSupportEnabled() {
		schema.Enable()

		if !IsSchemaValidationEnabled() {
			schema.DisableValidation()
		}
	}
}

// Returns name of the application
func GetAppName() string {
	return appName
}

// Returns version of the application
func GetAppVersion() string {
	return appVersion
}

func LoadAppConfig(flogoJson string, compressed bool) (*app.Config, error) {

	var jsonBytes []byte

	if flogoJson == "" {

		// a json string wasn't provided, so lets lookup the file in path
		configPath := GetFlogoConfigPath()

		flogo, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}

		jsonBytes, err = ioutil.ReadAll(flogo)
		if err != nil {
			return nil, err
		}
	} else {

		if compressed {
			var err error
			jsonBytes, err = support.DecodeAndUnzip(flogoJson)
			if err != nil {
				return nil, err
			}
		} else {
			jsonBytes = []byte(flogoJson)
		}
	}

	updated, err := secret.PreProcessConfig(jsonBytes)
	if err != nil {
		return nil, err
	}

	appConfig := &app.Config{}
	err = json.Unmarshal(updated, &appConfig)
	if err != nil {
		return nil, err
	}

	appName = appConfig.Name
	appVersion = appConfig.Version
	logEx.Init(appName)
	return appConfig, nil
}
