package database

import (
	"fmt"

	"go-gin-gorm-mysql/internal/core/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// Database global variable database
	Database = &gorm.DB{}
)

// MySQLConfig config mysql
type MySQLConfig struct {
	Database *gorm.DB
}

// InitConnection open initialize a new db connection.
func InitConnection(cf *config.Configs) error {
	dns := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cf.Database.MySQL.Username,
		cf.Database.MySQL.Password,
		cf.Database.MySQL.Protocol,
		cf.Database.MySQL.Host,
		cf.Database.MySQL.Port,
		cf.Database.MySQL.DatabaseName,
	)

	database, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		logrus.Errorf("[InitConnection] failed to connect to the database error: %s", err)
		return err
	}

	sqlDB, err := database.DB()
	if err != nil {
		logrus.Errorf("[InitConnection] set up to connect to the database error: %s", err)
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		logrus.Errorf("[InitConnection] ping error: %s", err)
		return err
	}

	return nil
}
