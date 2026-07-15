package configuration

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"
)

type Configuration struct {
	Key       valueobject.ConfigurationKey   `json:"key"`
	Value     valueobject.ConfigurationValue `json:"value"`
	Type      ConfigType                     `json:"type"`
	Scope     ConfigScope                    `json:"scope"`
	UpdatedAt time.Time                      `json:"updated_at"`
	UpdatedBy string                         `json:"updated_by,omitempty"`
}

type ConfigType string

const (
	CfgString  ConfigType = "string"
	CfgInt     ConfigType = "int"
	CfgFloat   ConfigType = "float"
	CfgBool    ConfigType = "bool"
	CfgJSON    ConfigType = "json"
)

type ConfigScope string

const (
	ScopeGlobal  ConfigScope = "global"
	ScopeCity    ConfigScope = "city"
	ScopeRegion  ConfigScope = "region"
)

func NewConfiguration(key valueobject.ConfigurationKey, value valueobject.ConfigurationValue, cfgType ConfigType) *Configuration {
	return &Configuration{
		Key:       key,
		Value:     value,
		Type:      cfgType,
		Scope:     ScopeGlobal,
		UpdatedAt: time.Now(),
	}
}

func (c *Configuration) Update(value valueobject.ConfigurationValue) {
	c.Value = value
	c.UpdatedAt = time.Now()
}
