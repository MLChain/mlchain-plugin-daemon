package plugin_manager

import (
	"fmt"
	"time"

	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_manager/aws_manager"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_manager/basic_manager"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_manager/positive_manager"
	"github.com/mlchain/mlchain-plugin-daemon/internal/db"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/models"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/cache"
)

const (
	PLUGIN_SERVERLESS_CACHE_KEY = "serverless:runtime:%s"
)

func (p *PluginManager) getServerlessRuntimeCacheKey(
	identity plugin_entities.PluginUniqueIdentifier,
) string {
	return fmt.Sprintf(PLUGIN_SERVERLESS_CACHE_KEY, identity.String())
}

func (p *PluginManager) getServerlessPluginRuntime(
	identity plugin_entities.PluginUniqueIdentifier,
) (plugin_entities.PluginLifetime, error) {
	model, err := p.getServerlessPluginRuntimeModel(identity)
	if err != nil {
		return nil, err
	}

	declaration := model.Declaration

	// init runtime entity
	runtimeEntity := plugin_entities.PluginRuntime{
		Config: declaration,
	}
	runtimeEntity.InitState()

	// convert to plugin runtime
	pluginRuntime := aws_manager.AWSPluginRuntime{
		PositivePluginRuntime: positive_manager.PositivePluginRuntime{
			BasicPluginRuntime: basic_manager.NewBasicPluginRuntime(p.mediaBucket),
			InnerChecksum:      model.Checksum,
		},
		PluginRuntime: runtimeEntity,
		LambdaURL:     model.FunctionURL,
		LambdaName:    model.FunctionName,
	}

	if err := pluginRuntime.InitEnvironment(); err != nil {
		return nil, err
	}

	return &pluginRuntime, nil
}

func (p *PluginManager) getServerlessPluginRuntimeModel(
	identity plugin_entities.PluginUniqueIdentifier,
) (*models.ServerlessRuntime, error) {
	// check if plugin is a serverless runtime
	runtime, err := cache.Get[models.ServerlessRuntime](
		p.getServerlessRuntimeCacheKey(identity),
	)
	if err != nil && err != cache.ErrNotFound {
		return nil, fmt.Errorf("unexpected error occurred during fetch serverless runtime cache: %v", err)
	}

	if err == cache.ErrNotFound {
		runtimeModel, err := db.GetOne[models.ServerlessRuntime](
			db.Equal("plugin_unique_identifier", identity.String()),
		)

		if err == db.ErrDatabaseNotFound {
			return nil, fmt.Errorf("plugin serverless runtime not found: %s", identity.String())
		}

		if err != nil {
			return nil, fmt.Errorf("failed to load serverless runtime from db: %v", err)
		}

		cache.Store(p.getServerlessRuntimeCacheKey(identity), runtimeModel, time.Minute*30)
		runtime = &runtimeModel
	} else if err != nil {
		return nil, fmt.Errorf("failed to load serverless runtime from cache: %v", err)
	}

	return runtime, nil
}
