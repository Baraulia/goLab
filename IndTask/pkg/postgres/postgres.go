package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
)

type PostgresDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
	logger   logging.Logger
}

func NewPostgresDB(database PostgresDB) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		database.Host, database.Port, database.Username, database.DBName, database.Password, database.SSLMode))
	if err != nil {
		database.logger.Panicf("Database open error:%s", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		database.logger.Errorf("DB ping error:%s", err)
		return nil, err
	}
	return db, nil
}
