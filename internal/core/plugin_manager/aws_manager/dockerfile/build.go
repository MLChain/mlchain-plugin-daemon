package dockerfile

import (
	"fmt"

	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/constants"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/strings"
)

func handleTemplate(configuration *plugin_entities.PluginDeclaration, templateFunc func(configuration *plugin_entities.PluginDeclaration) (string, error)) (string, error) {
	if templateFunc == nil {
		return "", fmt.Errorf("template function is nil, language: %s, version: %s", configuration.Meta.Runner.Language, configuration.Meta.Runner.Version)
	}
	return templateFunc(configuration)
}

// GenerateDockerfile generates a Dockerfile for the plugin
func GenerateDockerfile(configuration *plugin_entities.PluginDeclaration) (string, error) {
	if !strings.Find(configuration.Meta.Arch, constants.AMD64) {
		return "", fmt.Errorf("unsupported architecture: %s", configuration.Meta.Arch)
	}

	switch configuration.Meta.Runner.Language {
	case constants.Python:
		return handleTemplate(configuration, pythonTemplates[configuration.Meta.Runner.Version])
	}

	return "", fmt.Errorf("unsupported language: %s", configuration.Meta.Runner.Language)
}
