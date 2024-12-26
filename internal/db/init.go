package db

import (
	"fmt"
	"time"

	"github.com/mlchain/mlchain-plugin-daemon/internal/types/app"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/models"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initMlchainPluginDB(host string, port int, db_name string, user string, pass string, sslmode string) error {
	// create db if not exists
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, "postgres", sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	pgsqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// check if the db exists
	rows, err := pgsqlDB.Query(fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", db_name))
	if err != nil {
		return err
	}

	if !rows.Next() {
		// create database
		_, err = pgsqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", db_name))
		if err != nil {
			return err
		}
	}

	// close db
	err = pgsqlDB.Close()
	if err != nil {
		return err
	}

	// connect to the new db
	dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, db_name, sslmode)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	pgsqlDB, err = db.DB()
	if err != nil {
		return err
	}

	// check if uuid-ossp extension exists
	rows, err = pgsqlDB.Query("SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp'")
	if err != nil {
		return err
	}

	if !rows.Next() {
		// create the uuid-ossp extension
		_, err = pgsqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		if err != nil {
			return err
		}
	}

	pgsqlDB.SetConnMaxIdleTime(time.Minute * 1)
	MlchainPluginDB = db

	return nil
}

func autoMigrate() error {
	return MlchainPluginDB.AutoMigrate(
		models.Plugin{},
		models.PluginInstallation{},
		models.PluginDeclaration{},
		models.Endpoint{},
		models.ServerlessRuntime{},
		models.ToolInstallation{},
		models.AIModelInstallation{},
		models.InstallTask{},
		models.TenantStorage{},
		models.AgentStrategyInstallation{},
	)
}

func Init(config *app.Config) {
	err := initMlchainPluginDB(
		config.DBHost,
		int(config.DBPort),
		config.DBDatabase,
		config.DBUsername, config.DBPassword, config.DBSslMode,
	)

	if err != nil {
		log.Panic("failed to init mlchain plugin db: %v", err)
	}

	err = autoMigrate()
	if err != nil {
		log.Panic("failed to auto migrate: %v", err)
	}

	log.Info("mlchain plugin db initialized")
}

func Close() {
	db, err := MlchainPluginDB.DB()
	if err != nil {
		log.Error("failed to close mlchain plugin db: %v", err)
		return
	}

	err = db.Close()
	if err != nil {
		log.Error("failed to close mlchain plugin db: %v", err)
		return
	}

	log.Info("mlchain plugin db closed")
}
