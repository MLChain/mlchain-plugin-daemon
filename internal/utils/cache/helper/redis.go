package helper

import (
	"errors"
	"strings"

	"github.com/mlchain/mlchain-plugin-daemon/internal/db"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/models"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/cache"
)

var (
	ErrPluginNotFound = errors.New("plugin not found")
)

func CombinedGetPluginDeclaration(
	plugin_unique_identifier plugin_entities.PluginUniqueIdentifier,
	tenant_id string,
	runtime_type plugin_entities.PluginRuntimeType,
) (*plugin_entities.PluginDeclaration, error) {
	return cache.AutoGetWithGetter(
		strings.Join(
			[]string{
				string(runtime_type),
				plugin_unique_identifier.String(),
			},
			":",
		),
		func() (*plugin_entities.PluginDeclaration, error) {
			if runtime_type != plugin_entities.PLUGIN_RUNTIME_TYPE_REMOTE {
				declaration, err := db.GetOne[models.PluginDeclaration](
					db.Equal("plugin_unique_identifier", plugin_unique_identifier.String()),
				)
				if err == db.ErrDatabaseNotFound {
					return nil, ErrPluginNotFound
				}

				if err != nil {
					return nil, err
				}

				return &declaration.Declaration, nil
			} else {
				// try to fetch the declaration from plugin if it's remote
				plugin, err := db.GetOne[models.Plugin](
					db.Equal("plugin_unique_identifier", plugin_unique_identifier.String()),
					db.Equal("install_type", string(plugin_entities.PLUGIN_RUNTIME_TYPE_REMOTE)),
				)
				if err == db.ErrDatabaseNotFound {
					return nil, ErrPluginNotFound
				}

				if err != nil {
					return nil, err
				}

				return &plugin.Declaration, nil
			}
		},
	)
}
