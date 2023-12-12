package main

import (
	"github.com/orandin/lumberjackrus"
	logger "github.com/sirupsen/logrus"
	"github.com/xiaoxue1272/club-5fw-backend/config"
	"github.com/xiaoxue1272/club-5fw-backend/db"
	"github.com/xiaoxue1272/club-5fw-backend/web"
)

func configureLogger(logCfg *config.LoggerConfiguration) {
	level, err := logger.ParseLevel(logCfg.Level)
	if err != nil {
		panic(err)
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
	var hook *lumberjackrus.Hook
	hook, err = lumberjackrus.NewHook(
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
	configureLogger(configuration.Logger)
	web.Init(configuration.Web)
	db.Init(configuration.Database)
	web.StartWebServer()
}
