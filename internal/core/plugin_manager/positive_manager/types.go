package positive_manager

import (
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_manager/basic_manager"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/plugin_packager/decoder"
)

type PositivePluginRuntime struct {
	basic_manager.BasicPluginRuntime

	WorkingPath string
	// plugin decoder used to manage the plugin
	Decoder decoder.PluginDecoder

	InnerChecksum string
}
