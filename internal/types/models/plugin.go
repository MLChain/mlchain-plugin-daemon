package models

import (
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/manifest_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
)

type Plugin struct {
	Model
	// PluginUniqueIdentifier is a unique identifier for the plugin, it contains version and checksum
	PluginUniqueIdentifier string `json:"plugin_unique_identifier" gorm:"index;size:255"`
	// PluginID is the id of the plugin, only plugin name is considered
	PluginID     string                             `json:"id" gorm:"index;size:255"`
	Refers       int                                `json:"refers" gorm:"default:0"`
	InstallType  plugin_entities.PluginRuntimeType  `json:"install_type" gorm:"size:127;index"`
	ManifestType manifest_entities.DifyManifestType `json:"manifest_type" gorm:"size:127"`
	Declaration  plugin_entities.PluginDeclaration  `json:"declaration" gorm:"serializer:json;type:text;size:65535"`
}

type ServerlessRuntimeType string

const (
	SERVERLESS_RUNTIME_TYPE_AWS_LAMBDA ServerlessRuntimeType = "aws_lambda"
)

type ServerlessRuntime struct {
	Model
	PluginUniqueIdentifier string                            `json:"plugin_unique_identifier" gorm:"size:255;unique"`
	FunctionURL            string                            `json:"function_url" gorm:"size:255"`
	FunctionName           string                            `json:"function_name" gorm:"size:127"`
	Type                   ServerlessRuntimeType             `json:"type" gorm:"size:127"`
	Declaration            plugin_entities.PluginDeclaration `json:"declaration" gorm:"serializer:json;type:text;size:65535"`
	Checksum               string                            `json:"checksum" gorm:"size:127;index"`
}

type PluginDeclaration struct {
	Model
	PluginUniqueIdentifier string                            `json:"plugin_unique_identifier" gorm:"size:255;unique"`
	PluginID               string                            `json:"plugin_id" gorm:"size:255;index"`
	Declaration            plugin_entities.PluginDeclaration `json:"declaration" gorm:"serializer:json;type:text;size:65535"`
}