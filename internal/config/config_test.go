package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	assert := assert.New(t)

	config, err := LoadConfigFromFile("./testdata/test_config.yaml")

	assert.Nil(err)
	assert.Equal(config.App.Port, 3000)
	assert.Equal(config.Logging.LogLevel, "INFO")
	assert.Equal(config.Logging.File.LogFilePath, "/app/logs/comiket_backend.log")
	assert.Equal(config.Db.Postgres.Host, "comiket-db")
	assert.Equal(config.Db.Postgres.Port, 5432)
	assert.Equal(config.Db.Postgres.DatabaseName, "comiket")
	assert.Equal(config.Db.Postgres.Username, "username")
	assert.Equal(config.Db.Postgres.Password, "password")

}

func TestLoadConfigFromFileNoLoggingFile(t *testing.T) {
	assert := assert.New(t)

	config, err := LoadConfigFromFile("./testdata/test_config_no_log_file.yaml")

	assert.Nil(err)
	assert.Equal(config.App.Port, 3000)
	assert.Equal(config.Logging.LogLevel, "INFO")
	assert.Equal(config.Db.Postgres.Host, "comiket-db")
	assert.Equal(config.Db.Postgres.Port, 5432)
	assert.Equal(config.Db.Postgres.DatabaseName, "comiket")
	assert.Equal(config.Db.Postgres.Username, "username")
	assert.Equal(config.Db.Postgres.Password, "password")

}
