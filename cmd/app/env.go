package app

import (
	"github.com/go-sql-driver/mysql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Файл .env не найден")
	}

	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")

	return &Config{
		MySQLHost:     mysqlHost,
		MySQLPort:     mysqlPort,
		MySQLUser:     mysqlUser,
		MySQLPassword: mysqlPassword,
		MySQLDatabase: mysqlDatabase,
	}, nil
}
func (c *Config) ToMySQLConfig() mysql.Config {
	return mysql.Config{
		User:   c.MySQLUser,
		Passwd: c.MySQLPassword,
		Net:    "tcp",
		Addr:   c.MySQLHost + ":" + c.MySQLPort,
		DBName: c.MySQLDatabase,
	}
}
