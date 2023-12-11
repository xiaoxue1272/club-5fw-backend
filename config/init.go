package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Configuration struct {
	Logger   *LoggerConfiguration   `mapstructure:"logger,sqursh" json:"logger" yaml:"logger" `
	Web      *WebConfiguration      `mapstructure:"web,sqursh" json:"web" yaml:"web"`
	Database *DataBaseConfiguration `mapstructure:"database,sqursh" json:"database" yaml:"database"`
}

type LoggerConfiguration struct {
	Level      string                        `mapstructure:"level" json:"level" yaml:"level"`
	Formatter  *LoggerFormatterConfiguration `mapstructure:"formatter,sqursh" json:"formatter" yaml:"formatter"`
	MaxSize    int                           `mapstructure:"maxSize" json:"maxSize" yaml:"maxSize"`
	MaxBackups int                           `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`
	MaxAge     int                           `mapstructure:"maxAge" json:"maxAge" yaml:"maxAge"`
	Compress   bool                          `mapstructure:"compress" json:"compress" yaml:"compress"`
	LocalTime  bool                          `mapstructure:"localTime" json:"localTime" yaml:"localTime"`
}
type LoggerFormatterConfiguration struct {
	DisableColors   bool   `mapstructure:"disableColors" json:"disableColors" yaml:"disableColors"`
	TimestampFormat string `mapstructure:"timestampFormat" json:"timestampFormat" yaml:"timestampFormat"`
	DisableSorting  bool   `mapstructure:"disableSorting" json:"disableSorting" yaml:"disableSorting"`
}

type WebConfiguration struct {
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Host string `mapstructure:"host" json:"host" yaml:"host"`
}

type DataBaseConfiguration struct {
	Dns string `mapstructure:"dns" json:"dns" yaml:"dns"`
}

var standardConfiguration *Configuration = nil

var InitConfigErr error = nil

func useDefaultConfiguration() {
	standardConfiguration = &Configuration{
		Logger:   defaultLoggerConfiguration(),
		Web:      defaultWebConfiguration(),
		Database: defaultDatabaseConfiguration(),
	}
	viper.Set("logger", standardConfiguration.Logger)
	viper.Set("web", standardConfiguration.Web)
}

func defaultLoggerConfiguration() *LoggerConfiguration {
	return &LoggerConfiguration{
		Formatter: &LoggerFormatterConfiguration{
			DisableColors:   true,
			DisableSorting:  true,
			TimestampFormat: "2006-01-02 15:01:05",
		},
		Compress:   false,
		LocalTime:  false,
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     1,
		Level:      "debug",
	}
}

func defaultWebConfiguration() *WebConfiguration {
	return &WebConfiguration{
		Host: "127.0.0.1",
		Port: 9090,
	}
}

func defaultDatabaseConfiguration() *DataBaseConfiguration {
	return &DataBaseConfiguration{
		Dns: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai&allowAllFiles=true&timeout=30s",
	}
}

func init() {
	configDir := getConfigDir()
	configType := getConfigType()
	configName := getConfigName()
	pflag.Parse()
	viper.AddConfigPath(*configDir)
	viper.SetConfigType(*configType)
	viper.SetConfigName(*configName)
	viper.SetConfigPermissions(0644)
	err := viper.ReadInConfig()
	if err == nil {
		err = viper.Unmarshal(&standardConfiguration)
		if err == nil {
			return
		}
	}
	useDefaultConfiguration()
	err = viper.SafeWriteConfig()
	InitConfigErr = err
}

func getConfigDir() *string {
	return pflag.String("configDir", ".", "the application configuration dictionary path")
}

func getConfigType() *string {
	return pflag.String("configType", "yaml", "configuration file type")
}

func getConfigName() *string {
	return pflag.String("configName", "config", "configuration file name")
}

func GetConfiguration() *Configuration {
	return standardConfiguration
}
