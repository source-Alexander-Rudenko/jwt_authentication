package main

import (
	"jwt/cmd/app"
	"jwt/internal/repo"
	"log"
)

func main() {

	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("Ощибка загрузки конфигурации: %v", err)
	}
	mysqlCfg := cfg.ToMySQLConfig()
	db, err := repo.NewMySQLAStorage(mysqlCfg)
	if err != nil {
		log.Fatalf("Ошибка при подключении к бд: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка при пинге бд %v", err)
	}

	server := app.NewAPIServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal("Unable to run server", err)
	}
}
