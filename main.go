package main

import (
	"github.com/orandin/lumberjackrus"
	logger "github.com/sirupsen/logrus"
	"github.com/xiaoxue1272/club-5fw-backend/config"
	"github.com/xiaoxue1272/club-5fw-backend/web"
)

func configureLogger(logCfg *config.LoggerConfiguration) {
	level, err := logger.ParseLevel(logCfg.Level)
	if err != nil {
		panic("Unknown logger level, please ensure that the configured level is correct")
	}
	logger.SetLevel(level)
	textFormatter := &logger.TextFormatter{
		DisableColors:   logCfg.Formatter.DisableColors,
		TimestampFormat: logCfg.Formatter.TimestampFormat,
		DisableSorting:  logCfg.Formatter.DisableSorting,
		FieldMap: logger.FieldMap{
			"FieldKeyTime":  "@timestamp",
			"FieldKeyLevel": "@level",
			"FieldKeyMsg":   "@message",
		},
	}
	logger.SetFormatter(textFormatter)
	hook, err := lumberjackrus.NewHook(
		configureLogFile(logCfg, "log/general.log"),
		level,
		textFormatter,
		&lumberjackrus.LogFileOpts{
			logger.WarnLevel:  configureLogFile(logCfg, "log/warn.log"),
			logger.ErrorLevel: configureLogFile(logCfg, "log/error.log"),
		},
	)
	if err != nil {
		panic(err)
	}
	logger.AddHook(hook)
}

func configureLogFile(logCfg *config.LoggerConfiguration, fileName string) *lumberjackrus.LogFile {
	return &lumberjackrus.LogFile{
		Filename:   fileName,
		LocalTime:  logCfg.LocalTime,
		MaxAge:     logCfg.MaxAge,
		MaxBackups: logCfg.MaxBackups,
		MaxSize:    logCfg.MaxSize,
	}
}

func main() {
	configuration := config.GetConfiguration()
	if config.InitConfigErr != nil {
		logger.Warnf("Loading configuration file failed. \nError: %v\n", config.InitConfigErr)
		logger.Infoln("Using default builtin configurations.")
	}
	configureLogger(configuration.Logger)
	web.StartWebServer(configuration.Web)
}
