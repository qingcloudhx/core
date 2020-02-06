package app

import (
	"github.com/qingcloudhx/core/action"
	"github.com/qingcloudhx/core/app/resource"
	"github.com/qingcloudhx/core/data"
	"github.com/qingcloudhx/core/data/schema"
	"github.com/qingcloudhx/core/support/log"
	"github.com/qingcloudhx/core/trigger"
)

// Def is the configuration for the App
type Config struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
	AppModel    string `json:"appModel"`

	Imports    []string               `json:"imports,omitempty"`
	Properties []*data.Attribute      `json:"properties,omitempty"`
	Channels   []string               `json:"channels,omitempty"`
	Triggers   []*trigger.Config      `json:"triggers"`
	Resources  []*resource.Config     `json:"resources,omitempty"`
	Actions    []*action.Config       `json:"actions,omitempty"`
	Schemas    map[string]*schema.Def `json:"schemas,omitempty"`
	Log        *log.Option            `json:"log,omitempty"`
}

func LoadConfig() *Config {
	return nil
}
