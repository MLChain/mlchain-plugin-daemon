package persistence

import (
	"github.com/mlchain/mlchain-plugin-daemon/internal/oss"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/app"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/log"
)

var (
	persistence *Persistence
)

func InitPersistence(oss oss.OSS, config *app.Config) {
	persistence = &Persistence{
		storage:        NewWrapper(oss, config.PersistenceStoragePath),
		maxStorageSize: config.PersistenceStorageMaxSize,
	}

	log.Info("Persistence initialized")
}

func GetPersistence() *Persistence {
	return persistence
}
