package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/GerardRodes/muzz-backend/internal/config"
	"github.com/GerardRodes/muzz-backend/internal/domain"
	"github.com/GerardRodes/muzz-backend/internal/httpserver"
	"github.com/GerardRodes/muzz-backend/internal/mariadb"
	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.New()

	db, err := initDBHandle(cfg)
	if err != nil {
		log.Fatalf("cannot init db handler: %s", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		} else {
			log.Println("closed db connection")
		}
	}()

	userSvc := domain.NewUserSvc(mariadb.NewUserRepo(db))

	if err := httpserver.Init(httpserver.Config{
		HTTPPort:   cfg.HTTPPort,
		UserSvc:    userSvc,
		ProfileSvc: userSvc,
	}); err != nil {
		log.Fatalf("cannot init http server: %s", err)
	}
}

func initDBHandle(cfg config.Config) (*sql.DB, error) {
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.User = cfg.DBUser
	mysqlCfg.Passwd = cfg.DBPassword
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = cfg.DBAddr
	mysqlCfg.DBName = cfg.DBName

	connector, err := mysql.NewConnector(mysqlCfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create connector: %w", err)
	}

	db := sql.OpenDB(connector)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot reach the database: %w", err)
	}

	return db, nil
}
