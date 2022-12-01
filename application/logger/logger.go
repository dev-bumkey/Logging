package config

import (
	"github.com/cocktailcloud/acloud-alarm-collector/application/config"
	"github.com/cocktailcloud/acloud-monitoring-common/v2/logger"
)

func NewLogger(conf *config.Config) error {
	return logger.NewLogger(logger.Config{
		EnableConsole: true,
		ConsoleLevel:  logger.GetLogLevel(conf.LoggingLevel),
		// ConsoleJSONFormat: true,
		EnableFile:     conf.LoggingFileUse,
		FileLevel:      logger.GetLogLevel(conf.LoggingLevel),
		FileJSONFormat: false,
		FileLocation:   conf.LoggingFilePath,
	}, logger.InstanceZapLogger)
}
