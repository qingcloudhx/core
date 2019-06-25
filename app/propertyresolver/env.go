package propertyresolver

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/qingcloudhx/core/data/property"
	"github.com/qingcloudhx/core/engine"
	"github.com/qingcloudhx/core/support/log"
)

const EnvAppPropertyEnvConfigKey = "FLOGO_APP_PROPS_ENV"

type PropertyMappings struct {
	Mappings map[string]string `json:"mappings"`
}

func init() {

	logger := log.RootLogger()
	mappingsEnv := getEnvValue()

	if mappingsEnv != "" {

		resolver := &EnvVariableValueResolver{}

		if strings.EqualFold(mappingsEnv, "auto") {
			resolver.autoMapping = true
		} else {

			var mappings *PropertyMappings
			e := json.Unmarshal([]byte(mappingsEnv), &mappings)
			if e != nil {
				logger.Errorf("Can not parse value set to '%s' due to error - '%v'", EnvAppPropertyEnvConfigKey, e)
				panic("")
			}
		}
		_ = property.RegisterPropertyResolver(resolver)
	}
}

func getEnvValue() string {
	key := os.Getenv(EnvAppPropertyEnvConfigKey)
	if len(key) > 0 {
		return key
	}
	return ""
}

// Resolve property value from environment variable
type EnvVariableValueResolver struct {
	autoMapping bool
	mappings    PropertyMappings
}

func (resolver *EnvVariableValueResolver) Name() string {
	return engine.PropertyResolverEnv
}

func (resolver *EnvVariableValueResolver) LookupValue(key string) (interface{}, bool) {

	if resolver.autoMapping {
		value, exists := os.LookupEnv(key) // first try with the name of the property as is
		if exists {
			return value, exists
		}

		// Replace dot with underscore e.g. a.b would be a_b
		key = strings.Replace(key, ".", "_", -1)
		value, exists = os.LookupEnv(key)
		if exists {
			return value, exists
		}

		// Try upper case form e.g. a.b would be A_B
		key = strings.ToUpper(key)
		value, exists = os.LookupEnv(key) // if not found try with the canonical form
		return value, exists

	} else {
		// Lookup based on mapping defined
		keyMapping, ok := resolver.mappings.Mappings[key]
		if ok {
			return os.LookupEnv(keyMapping)
		}
	}

	return nil, false
}
