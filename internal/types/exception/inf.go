package exception

import "github.com/mlchain/mlchain-plugin-daemon/internal/types/entities"

type PluginDaemonError interface {
	error

	ToResponse() *entities.Response
}
