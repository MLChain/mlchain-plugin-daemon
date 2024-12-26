package model_entities

import "github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"

type GetModelSchemasResponse struct {
	ModelSchema *plugin_entities.ModelDeclaration `json:"model_schema" validate:"omitempty"`
}
