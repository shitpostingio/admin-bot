package config

import (
	"log"

	"github.com/spf13/viper"

	"github.com/shitpostingio/admin-bot/config/structs"
)

const (
	defaultFileSizeThreshold  = 20971520 //20MB
	defaultDatabaseAddress    = "localhost"
	defaultDatabasePort       = 3306
	defaultDocumentStoreHosts = "localhost:27017"
	defaultSocketPath         = "/tmp/log.socket"
)

// Load reads a configuration file and returns its config instance
func Load(path string, useWebhook bool) (cfg structs.Config, err error) {

	setDefaultValuesForOptionalFields()

	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	err = CheckMandatoryFields(false, cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if useWebhook {
		checkWebhookConfig(&cfg.Webhook)
	}

	err = viper.WriteConfig()
	return
}

func setDefaultValuesForOptionalFields() {
	viper.SetDefault("fpserver.filesizethreshold", defaultFileSizeThreshold)
	viper.SetDefault("log.socketpath", defaultSocketPath)
	viper.SetDefault("database.address", defaultDatabaseAddress)
	viper.SetDefault("database.port", defaultDatabasePort)
	viper.SetDefault("documentstore.hosts", []string{defaultDocumentStoreHosts})
}
