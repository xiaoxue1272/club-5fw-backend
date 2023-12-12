package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xiaoxue1272/club-5fw-backend/utils"
)

type Configuration struct {
	Logger   *LoggerConfiguration   `json:"logger" yaml:"logger"`
	Web      *WebConfiguration      `json:"web" yaml:"web"`
	Database *DataBaseConfiguration `json:"database" yaml:"database"`
}

type LoggerConfiguration struct {
	Formatter  *FormatterConfiguration `json:"formatter" yaml:"formatter"`
	Level      string                  `json:"level" yaml:"level"`
	MaxSize    int                     `json:"maxSize" yaml:"maxSize"`
	MaxBackups int                     `json:"maxBackups" yaml:"maxBackups"`
	MaxAge     int                     `json:"maxAge" yaml:"maxAge"`
	Compress   bool                    `json:"compress" yaml:"compress"`
	LocalTime  bool                    `json:"localTime" yaml:"localTime"`
}
type FormatterConfiguration struct {
	DisableColors   bool   `json:"disableColors" yaml:"disableColors"`
	TimestampFormat string `json:"timestampFormat" yaml:"timestampFormat"`
	DisableSorting  bool   `json:"disableSorting" yaml:"disableSorting"`
}

type RsaKeyPairConfiguration struct {
	PrivateKey string `json:"privateKey" yaml:"privateKey"`
	PublicKey  string `json:"publicKey" yaml:"publicKey"`
}

type WebConfiguration struct {
	Jwt        *JwtConfiguration        `json:"jwt" yaml:"jwt"`
	CipherJson *CipherJsonConfiguration `json:"cipherJson" yaml:"cipherJson"`
	Port       int                      `json:"port" yaml:"port"`
	Host       string                   `json:"host" yaml:"host"`
}

type CipherJsonConfiguration struct {
	Rsa *RsaKeyPairConfiguration
}

type JwtConfiguration struct {
	Rsa *RsaKeyPairConfiguration
	Alg string `json:"alg" yaml:"alg"`
}

type DataBaseConfiguration struct {
	NamingStrategy                           *NamingStrategyConfiguration `json:"namingStrategy" yaml:"namingStrategy"`
	Dns                                      string                       `json:"dns" yaml:"dns"`
	SkipDefaultTransaction                   bool                         `json:"skipDefaultTransaction" yaml:"skipDefaultTransaction"`
	FullSaveAssociations                     bool                         `json:"fullSaveAssociations" yaml:"fullSaveAssociations"`
	DryRun                                   bool                         `json:"dryRun" yaml:"dryRun"`
	PrepareStmt                              bool                         `json:"prepareStmt" yaml:"prepareStmt"`
	DisableAutomaticPing                     bool                         `json:"disableAutomaticPing" yaml:"disableAutomaticPing"`
	DisableForeignKeyConstraintWhenMigrating bool                         `json:"disableForeignKeyConstraintWhenMigrating" yaml:"disableForeignKeyConstraintWhenMigrating"`
	IgnoreRelationshipsWhenMigrating         bool                         `json:"ignoreRelationshipsWhenMigrating" yaml:"ignoreRelationshipsWhenMigrating"`
	DisableNestedTransaction                 bool                         `json:"disableNestedTransaction" yaml:"disableNestedTransaction"`
	AllowGlobalUpdate                        bool                         `json:"allowGlobalUpdate" yaml:"allowGlobalUpdate"`
	QueryFields                              bool                         `json:"queryFields" yaml:"queryFields"`
	CreateBatchSize                          int                          `json:"createBatchSize" yaml:"createBatchSize"`
	TranslateError                           bool                         `json:"translateError" yaml:"translateError"`
}

type NamingStrategyConfiguration struct {
	TablePrefix         string `json:"tablePrefix" yaml:"tablePrefix"`
	SingularTable       bool   `json:"singularTable" yaml:"singularTable"`
	NoLowerCase         bool   `json:"noLowerCase" yaml:"noLowerCase"`
	IdentifierMaxLength int    `json:"identifierMaxLength" yaml:"identifierMaxLength"`
}

var standardConfiguration = &Configuration{}

func useDefaultConfiguration() {
	standardConfiguration.Logger = defaultLoggerConfiguration()
	standardConfiguration.Web = defaultWebConfiguration()
	standardConfiguration.Database = defaultDatabaseConfiguration()
	viper.Set("logger", standardConfiguration.Logger)
	viper.Set("web", standardConfiguration.Web)
	viper.Set("database", standardConfiguration.Database)
}

func defaultLoggerConfiguration() *LoggerConfiguration {
	return &LoggerConfiguration{
		Formatter: &FormatterConfiguration{
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
	jwtRsaKey := utils.GenerateRsaKey(2048)
	jwtPrivateKeyPem, _ := utils.RsaPrivateKeyToString(jwtRsaKey)
	jwtPublicKeyPem, _ := utils.RsaPublicKeyToString(&jwtRsaKey.PublicKey)
	cipherJsonRsaKey := utils.GenerateRsaKey(4096)
	cipherJsonPrivateKeyPem, _ := utils.RsaPrivateKeyToString(cipherJsonRsaKey)
	cipherJsonPublicKeyPem, _ := utils.RsaPublicKeyToString(&cipherJsonRsaKey.PublicKey)
	return &WebConfiguration{
		Host: "127.0.0.1",
		Port: 9090,
		Jwt: &JwtConfiguration{
			Rsa: &RsaKeyPairConfiguration{
				PrivateKey: jwtPrivateKeyPem,
				PublicKey:  jwtPublicKeyPem,
			},
			Alg: "RS512",
		},
		CipherJson: &CipherJsonConfiguration{
			Rsa: &RsaKeyPairConfiguration{
				PrivateKey: cipherJsonPrivateKeyPem,
				PublicKey:  cipherJsonPublicKeyPem,
			},
		},
	}
}

func defaultDatabaseConfiguration() *DataBaseConfiguration {
	return &DataBaseConfiguration{
		Dns: "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai&allowAllFiles=true&timeout=30s",
		NamingStrategy: &NamingStrategyConfiguration{
			TablePrefix:   "club_",
			SingularTable: true,
			NoLowerCase:   false,
		},
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
		err = viper.Unmarshal(standardConfiguration)
		if err == nil {
			return
		}
	}
	useDefaultConfiguration()
	err = viper.SafeWriteConfig()
	if err != nil {
		panic(err)
	}
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
