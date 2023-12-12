package db

import (
	"github.com/sirupsen/logrus"
	"github.com/xiaoxue1272/club-5fw-backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

var dialect gorm.Dialector

var gormConfig *gorm.Config

func Init(dbConfig *config.DataBaseConfiguration) {
	dialect = mysql.Open(dbConfig.Dns)
	gormConfig = &gorm.Config{
		SkipDefaultTransaction: dbConfig.SkipDefaultTransaction,
		Logger:                 logger.New(logrus.StandardLogger(), logger.Config{}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:         dbConfig.NamingStrategy.TablePrefix,
			SingularTable:       dbConfig.NamingStrategy.SingularTable,
			NoLowerCase:         dbConfig.NamingStrategy.NoLowerCase,
			IdentifierMaxLength: dbConfig.NamingStrategy.IdentifierMaxLength,
		},
		FullSaveAssociations:                     dbConfig.FullSaveAssociations,
		DryRun:                                   dbConfig.DryRun,
		PrepareStmt:                              dbConfig.PrepareStmt,
		DisableAutomaticPing:                     dbConfig.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: dbConfig.DisableForeignKeyConstraintWhenMigrating,
		IgnoreRelationshipsWhenMigrating:         dbConfig.IgnoreRelationshipsWhenMigrating,
		DisableNestedTransaction:                 dbConfig.DisableNestedTransaction,
		AllowGlobalUpdate:                        dbConfig.AllowGlobalUpdate,
		QueryFields:                              dbConfig.QueryFields,
		CreateBatchSize:                          dbConfig.CreateBatchSize,
		TranslateError:                           dbConfig.TranslateError,
	}
}

func connect() {
	db, err := gorm.Open(dialect, gormConfig)
	if err != nil {
		logrus.Panicf("Failed to connect database\nError %v\n", err)
	}
	DB = db
}
