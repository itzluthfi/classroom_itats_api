package helper

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)
//tes
type databaseConnection struct {
}

type DatabaseConnection interface {
	Connect() (*gorm.DB, error)
}

func NewDatabaseConnection() *databaseConnection {
	return &databaseConnection{}
}

func (c *databaseConnection) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", viper.GetString("DB_HOST"), viper.GetString("DB_USERNAME"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_NAME"), viper.GetString("DB_PORT"))

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}

	return conn, nil
}
